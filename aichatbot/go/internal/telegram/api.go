package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// Client 是 Telegram Bot API 用戶端。
// 失敗（網路錯誤、HTTP 不是 2xx、回傳 ok:false）就照 RetryDelays 等待後重試（預設 1、2、4 秒，共 1+3 次），
// 仍失敗就記錄錯誤並回傳錯誤。API 網址含 token，日誌和錯誤訊息都先把 token 換成 ***
type Client struct {
	token       string
	endpoint    string // <base>/bot<token>/
	RetryDelays []time.Duration
	HTTP        *http.Client
	Log         *slog.Logger
}

func NewClient(token, baseURL string, log *slog.Logger) *Client {
	return &Client{
		token:       [REDACTED],
		endpoint:    strings.TrimSuffix(baseURL, "/") + "/bot" + token + "/",
		RetryDelays: []time.Duration{time.Second, 2 * time.Second, 4 * time.Second},
		HTTP:        &http.Client{Timeout: 30 * time.Second},
		Log:         log,
	}
}

type apiResponse struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	Description string          `json:"description"`
}

// Call 呼叫 method，成功時回傳 Telegram 回應的 result。retry 為 false 只送一次（「正在輸入」這類過時就沒意義的請求）
func (c *Client) Call(ctx context.Context, method string, params any, retry bool) (json.RawMessage, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	var delays []time.Duration
	if retry {
		delays = c.RetryDelays
	}
	var lastErr error
	for attempt := 0; attempt <= len(delays); attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delays[attempt-1]):
			}
		}
		result, err := c.once(ctx, method, body)
		if err == nil {
			return result, nil
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		lastErr = err
		if attempt < len(delays) {
			c.Log.Warn("Telegram API 失敗，稍後重試", "method", method, "error", err.Error(), "retry_in", delays[attempt].String())
		}
	}
	c.Log.Error("Telegram API 失敗，放棄", "method", method, "attempts", len(delays)+1, "error", lastErr.Error())
	return nil, fmt.Errorf("Telegram %s 失敗：%w", method, lastErr)
}

func (c *Client) once(ctx context.Context, method string, body []byte) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+method, bytes.NewReader(body))
	if err != nil {
		return nil, errors.New(redact(err.Error(), c.token))
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		// *url.Error 的訊息帶完整網址（含 token）
		return nil, errors.New("網路錯誤：" + redact(err.Error(), c.token))
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, errors.New("讀取回應失敗：" + redact(err.Error(), c.token))
	}
	var r apiResponse
	_ = json.Unmarshal(data, &r)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || !r.OK {
		msg := fmt.Sprintf("HTTP %d", resp.StatusCode)
		if r.Description != "" {
			msg += "：" + r.Description
		}
		return nil, errors.New(redact(msg, c.token))
	}
	if len(r.Result) == 0 {
		return json.RawMessage("true"), nil
	}
	return r.Result, nil
}

func (c *Client) SendMessage(ctx context.Context, chatID int64, text string) error {
	_, err := c.Call(ctx, "sendMessage", map[string]any{"chat_id": chatID, "text": text}, true)
	return err
}

// Button 是訊息下方的按鈕；Data 是按下時回呼帶回來的 callback_data
type Button struct {
	Text string `json:"text"`
	Data string `json:"callback_data"`
}

// SendOptions：ReplyTo 不是 0 時回覆那則訊息；Buttons 是一列一列的按鈕
type SendOptions struct {
	ReplyTo int64
	Buttons [][]Button
}

// SendMessageID 送訊息並回傳訊息編號（員工群組的通知要記下來，員工「回覆」時才找得回來）
func (c *Client) SendMessageID(ctx context.Context, chatID int64, text string, o SendOptions) (int64, error) {
	params := map[string]any{"chat_id": chatID, "text": text}
	if o.ReplyTo != 0 {
		// 原訊息不見了（被刪掉）照樣送出
		params["reply_parameters"] = map[string]any{"message_id": o.ReplyTo, "allow_sending_without_reply": true}
	}
	if len(o.Buttons) > 0 {
		params["reply_markup"] = map[string]any{"inline_keyboard": o.Buttons}
	}
	return c.callMessageID(ctx, "sendMessage", params)
}

// CopyMessage 把 fromChat 的訊息 messageID 複製一份到 toChat（不顯示「轉寄自」），回傳新訊息的編號
func (c *Client) CopyMessage(ctx context.Context, toChat, fromChat, messageID, replyTo int64) (int64, error) {
	params := map[string]any{"chat_id": toChat, "from_chat_id": fromChat, "message_id": messageID}
	if replyTo != 0 {
		params["reply_parameters"] = map[string]any{"message_id": replyTo, "allow_sending_without_reply": true}
	}
	return c.callMessageID(ctx, "copyMessage", params)
}

func (c *Client) callMessageID(ctx context.Context, method string, params map[string]any) (int64, error) {
	result, err := c.Call(ctx, method, params, true)
	if err != nil {
		return 0, err
	}
	var m struct {
		MessageID int64 `json:"message_id"`
	}
	if err := json.Unmarshal(result, &m); err != nil {
		return 0, fmt.Errorf("Telegram %s 回應看不懂：%w", method, err)
	}
	return m.MessageID, nil
}

// EditMessageText 改掉一則訊息的文字（不帶按鈕，原本的按鈕會拿掉）
func (c *Client) EditMessageText(ctx context.Context, chatID, messageID int64, text string) error {
	_, err := c.Call(ctx, "editMessageText", map[string]any{"chat_id": chatID, "message_id": messageID, "text": text}, true)
	return err
}

// AnswerCallbackQuery 回應按鈕回呼（text 是按的人畫面上跳出的提示，可以空白）；過時就沒意義，不重試
func (c *Client) AnswerCallbackQuery(ctx context.Context, id, text string) error {
	params := map[string]any{"callback_query_id": id}
	if text != "" {
		params["text"] = text
	}
	_, err := c.Call(ctx, "answerCallbackQuery", params, false)
	return err
}

// SendChatAction 送「正在輸入」；過時就沒意義，不重試
func (c *Client) SendChatAction(ctx context.Context, chatID int64) error {
	_, err := c.Call(ctx, "sendChatAction", map[string]any{"chat_id": chatID, "action": "typing"}, false)
	return err
}

// SetWebhook：只收一般訊息和按鈕回呼；Telegram 之後每個請求都會帶 X-Telegram-Bot-Api-Secret-Token: [REDACTED]
func (c *Client) SetWebhook(ctx context.Context, url, secret string) error {
	_, err := c.Call(ctx, "setWebhook", map[string]any{
		"url": url, "secret_token": secret, "allowed_updates": []string{"message", "callback_query"},
	}, true)
	return err
}
