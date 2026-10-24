import Foundation

// MARK: - Telegram 送來的 update（只解碼用得到的欄位，其他欄位一律忽略）
// 欄位說明：https://example.invalid/bots/api#update

public struct Update: Decodable, Sendable {
    public let updateID: Int64
    public let message: Message?
    public let channelPost: Message?        // 頻道貼文：忽略
    public let callbackQuery: CallbackQuery? // 按鈕回呼：店員群組通知上的確認／拒絕預約

    enum CodingKeys: String, CodingKey {
        case updateID = "update_id", message, channelPost = "channel_post", callbackQuery = "callback_query"
    }
}

public struct Message: Decodable, Sendable {
    public let messageID: Int64
    public let chat: Chat
    public let from: User?        // 頻道貼文沒有 from
    public let text: String?      // 貼圖、照片沒有 text
    public let caption: String?   // 照片、檔案附的說明文字
    public let replyToMessage: MessageRef?   // 這則是「回覆」哪一則訊息（店員回覆通知）
    public let sticker: Sticker?
    // 非文字訊息：只看有沒有，內容不解
    public let photo: Present?
    public let video: Present?
    public let videoNote: Present?
    public let animation: Present?
    public let voice: Present?
    public let audio: Present?
    public let document: Present?

    enum CodingKeys: String, CodingKey {
        case messageID = "message_id", chat, from, text, caption, replyToMessage = "reply_to_message", sticker
        case photo, video, videoNote = "video_note", animation, voice, audio, document
    }

    // 給店員看的種類名稱（貼圖不算，當成閒聊）；文字或不認得的種類回傳 nil
    public var media: String? {
        if photo != nil { return "圖片" }
        if video != nil || videoNote != nil || animation != nil { return "影片" }
        if voice != nil || audio != nil { return "語音" }
        if document != nil { return "檔案" }
        return nil
    }
}

// 被回覆的訊息：只需要它的編號
public struct MessageRef: Decodable, Sendable {
    public let messageID: Int64
    enum CodingKeys: String, CodingKey { case messageID = "message_id" }
}

public struct Sticker: Decodable, Sendable {
    public let emoji: String?
}

// 只判斷欄位有沒有出現（照片是陣列、其他是物件，都不解內容）
public struct Present: Decodable, Sendable {
    public init(from decoder: Decoder) throws {}
}

public struct Chat: Decodable, Sendable {
    public let id: Int64
    public let type: String       // private、group、supergroup、channel
}

public struct User: Decodable, Sendable {
    public let id: Int64
    public let firstName: String?
    public let lastName: String?
    public let username: String?

    enum CodingKeys: String, CodingKey { case id, firstName = "first_name", lastName = "last_name", username }
}

// 按鈕回呼：message 是按鈕所在的那則訊息（太舊的訊息 Telegram 不會附上）
public struct CallbackQuery: Decodable, Sendable {
    public let id: String
    public let from: User
    public let data: String?
    public let message: CallbackMessage?   // 按鈕所在的那則訊息（太舊的訊息 Telegram 不會帶）
}

// 按鈕所在的訊息：編號、聊天室、文字
public struct CallbackMessage: Decodable, Sendable {
    public let messageID: Int64
    public let chat: Chat
    public let text: String?
    enum CodingKeys: String, CodingKey { case messageID = "message_id", chat, text }
}
