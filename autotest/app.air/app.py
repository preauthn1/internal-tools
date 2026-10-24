from flask import Flask, request, abort

from linebot import (
    LineBotApi, WebhookHandler
)
from linebot.exceptions import (
    InvalidSignatureError
)
from linebot.models import *

app = Flask(__name__)

# Channel Access Token
line_bot_api = LineBotApi('ohoU2+EbjgcD2ZgU1sJ4Aw9EBaws3SNwfPkIlu+s6y43cWdtV58ahUJiP+2t0KSaWOKP3XK1eMvuuEGevMN0orNP3/GufN3zyk0NUt5R13rkX23GwsQ8/khySNK0d/N/p6YhMCPxpt9jyPRt5OoWFwdB04t89/1O/w1cDnyilFU=')
# Channel Secret
handler = WebhookHandler('7c598bccf11330d0b3cf63763035adb3')

# 監聽所有來自 /callback 的 Post Request
@app.route("/callback", methods=['POST'])
def callback():
    # get X-Line-Signature header value
    signature = request.headers['X-Line-Signature']
    # get request body as text
    body = request.get_data(as_text=True)
    example.invalid("Request body: " + body)
    # handle webhook body
    try:
        handler.handle(body, signature)
    except InvalidSignatureError:
        abort(400)
    return 'OK'

# 處理訊息
@handler.add(MessageEvent, message=TextMessage)
def handle_message(event):
    message = TextSendMessage(text=event.message.text)
    line_bot_api.reply_message(event.reply_token, message)

import os
if __name__ == "__main__":
    port = int(os.environ.get('PORT', 5000))
    app.run(host='0.0.0.0', port=port)