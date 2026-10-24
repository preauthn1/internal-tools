import Foundation

// MARK: - Telegram Bot API 用戶端
// 失敗（網路錯誤、HTTP 不是 2xx、回傳 ok:false）就照 retryDelays 等待後重試（預設 1、2、4 秒，共 1+3 次），
// 仍失敗就記錄錯誤、回傳 false，不丟錯誤，伺服器繼續服務。
// API 網址含 token，所以日誌和錯誤訊息只寫 method 和錯誤描述，不寫網址

public final class TelegramAPI: @unchecked Sendable {
    private let endpoint: URL            // <base>/bot<token>/
    let retryDelays: [Double]
    let log: Log
    let session: URLSession

    public init(token: [REDACTED], baseURL: String = "https://example.invalid", retryDelays: [Double] = [1, 2, 4],
                log: Log, session: URLSession = .shared) {
        let base = baseURL.hasSuffix("/") ? String(baseURL.dropLast()) : baseURL
        self.endpoint = URL(string: "\(base)/bot\(token)/")!
        self.retryDelays = retryDelays; self.log = log; self.session = session
    }

    // retry: false 只送一次（「正在輸入」這類過時就沒意義的請求）
    @discardableResult
    public func call(_ method: String, _ params: [String: Any], retry: Bool = true) async -> Bool {
        await request(method, params, retry: retry) != nil
    }

    // 和 call 相同，成功時回傳 Telegram 回應的 result（sendMessage 是送出的那則訊息）；失敗回傳 nil
    func request(_ method: String, _ params: [String: Any], retry: Bool = true) async -> Any? {
        var request = URLRequest(url: endpoint.appendingPathComponent(method))
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpBody = try? JSONSerialization.data(withJSONObject: params)
        request.timeoutInterval = 30

        let delays = retry ? retryDelays : []
        var lastError = ""
        for attempt in 0...delays.count {
            if attempt > 0 {
                try? await Task.sleep(for: .seconds(delays[attempt - 1]))
            }
            if Task.isCancelled { return nil }
            do {
                let (data, response) = try await session.data(for: request)
                let status = (response as? HTTPURLResponse)?.statusCode ?? 0
                let body = (try? JSONSerialization.jsonObject(with: data)) as? [String: Any]
                if (200..<300).contains(status), body?["ok"] as? Bool == true { return body?["result"] ?? true }
                lastError = "HTTP \(status)" + ((body?["description"] as? String).map { "：\($0)" } ?? "")
            } catch {
                if Task.isCancelled || error is CancellationError { return nil }
                // 不用 "\(error)"：NSError 的描述會帶上完整網址（含 token）
                lastError = (error as? URLError).map { "網路錯誤 \($0.code.rawValue)" } ?? "網路錯誤"
            }
            if attempt < delays.count { log.warn("Telegram \(method) 失敗（\(lastError)），\(delays[attempt]) 秒後重試") }
        }
        log.error("Telegram \(method) 失敗 \(delays.count + 1) 次，放棄（\(lastError)）")
        return nil
    }

    @discardableResult
    public func sendMessage(_ chatID: String, _ text: String, replyTo: Int64? = nil) async -> Bool {
        await call("sendMessage", messageParams(chatID, text, replyTo: replyTo))
    }

    // 送出後回傳那則訊息的編號（店員群組的通知要記下來）；失敗或回應沒有編號回傳 nil
    public func sendMessageID(_ chatID: String, _ text: String, replyTo: Int64? = nil, buttons: [[Button]]? = nil) async -> Int64? {
        let result = await request("sendMessage", messageParams(chatID, text, replyTo: replyTo, buttons: buttons)) as? [String: Any]
        return (result?["message_id"] as? NSNumber)?.int64Value
    }

    // 訊息下方的按鈕：按下後 Telegram 送 callback_query，data 原樣帶回來
    public struct Button: Sendable {
        public let text: String
        public let data: String
        public init(_ text: String, data: String) { self.text = text; self.data = data }
    }

    // replyTo：回覆在哪一則訊息上（那則被刪掉了也照樣送出）；buttons：每個陣列是一排按鈕
    func messageParams(_ chatID: String, _ text: String, replyTo: Int64?, buttons: [[Button]]? = nil) -> [String: Any] {
        var params: [String: Any] = ["chat_id": chatID, "text": text]
        if let replyTo { params["reply_parameters"] = ["message_id": replyTo, "allow_sending_without_reply": true] }
        if let buttons {
            params["reply_markup"] = ["inline_keyboard": buttons.map { row in row.map { ["text": $0.text, "callback_data": $0.data] } }]
        }
        return params
    }

    // 改訊息文字；不帶按鈕，原本的按鈕就會拿掉
    @discardableResult
    public func editMessageText(_ chatID: String, messageID: Int64, _ text: String) async -> Bool {
        await call("editMessageText", ["chat_id": chatID, "message_id": messageID, "text": text])
    }

    // 回應按鈕回呼（不回應的話按鈕會一直轉圈）；text 會短暫顯示在店員的畫面上。過時就沒意義，不重試
    @discardableResult
    public func answerCallbackQuery(_ id: String, text: String? = nil) async -> Bool {
        var params: [String: Any] = ["callback_query_id": id]
        if let text { params["text"] = text }
        return await call("answerCallbackQuery", params, retry: false)
    }

    @discardableResult
    public func sendChatAction(_ chatID: String, _ action: String = "typing") async -> Bool {
        await call("sendChatAction", ["chat_id": chatID, "action": action], retry: false)
    }

    // 只收一般訊息和按鈕回呼；Telegram 之後每個請求都會帶 X-Telegram-Bot-Api-Secret-Token: [REDACTED]
    @discardableResult
    public func setWebhook(url: String, secret: [REDACTED], allowedUpdates: [String] = ["message", "callback_query"]) async -> Bool {
        await call("setWebhook", ["url": url, "secret_token": secret, "allowed_updates": allowedUpdates])
    }
}
