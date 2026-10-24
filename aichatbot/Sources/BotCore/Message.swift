import Foundation

// MARK: - 訊息進、動作出（平台無關：LINE、終端機都換成這個格式再交給 Bot）

public struct IncomingMessage {
    public var platform: String      // line、terminal…
    public var userID: String
    public var chatID: String        // 回覆送到哪個聊天室（一對一時通常就是 userID）
    public var text: String          // 非文字訊息時是說明文字（caption），沒有就是空字串
    public var media: String?        // 非文字訊息：圖片、語音、檔案…（給店員看的名稱）；nil = 文字訊息
    public var time: Date
    public init(platform: String, userID: String, chatID: String? = nil, text: String, media: String? = nil, time: Date = Date()) {
        self.platform = platform; self.userID = userID; self.chatID = chatID ?? userID; self.text = text; self.media = media; self.time = time
    }
    public var customer: CustomerKey { CustomerKey(platform: platform, userID: userID) }
}

// 顧客狀態依「平台＋使用者 ID」分開保存
public struct CustomerKey: Hashable, CustomStringConvertible {
    public var platform: String
    public var userID: String
    public init(platform: String, userID: String) { self.platform = platform; self.userID = userID }
    public var description: String { "\(platform):\(userID)" }
}

public enum BotAction: Equatable {
    case reply(key: String, text: String, chatID: String)   // 回覆顧客；key 是回覆種類（replies.json 的模板名稱、answer…）
    case openTicket(Ticket)                                 // 開一筆待處理事項給店員
    case silent                                             // 不回覆（已轉人工）
}

public enum TicketKind: String {
    case question     // 知識庫沒有的問題，轉店員
    case booking      // 預約申請，等店員確認或拒絕
    case urgent       // 顧客追問多次
    case unsupported  // 不支援的語言
    case claim        // 顧客聲稱／更正
    case handoff      // 轉人工
}

public struct Ticket: Equatable {
    public let id: Int
    public let kind: TicketKind
    public let customer: CustomerKey
    public let chatID: String
    public let original: String      // 顧客的原話
    public let text: String          // 給店員看的通知文字
    public let relatedTo: Int?       // 緊急通知：對應的原本那張
    public var done = false          // 店員已回覆或處理
    public init(id: Int, kind: TicketKind, customer: CustomerKey, chatID: String, original: String, text: String,
                relatedTo: Int? = nil, done: Bool = false) {
        self.id = id; self.kind = kind; self.customer = customer; self.chatID = chatID
        self.original = original; self.text = text; self.relatedTo = relatedTo; self.done = done
    }
}

public enum BookingStatus: String { case pending, confirmed, rejected }

// 預約申請：編號和開出來的待處理事項相同；fields 是 date、time、people、name、phone
public struct BookingRequest: Equatable {
    public let id: Int
    public let customer: CustomerKey
    public let chatID: String
    public let summary: String
    public let fields: [String: String]
    public var status = BookingStatus.pending
    public var handledBy: String?    // 確認或拒絕的店員（顯示名稱）
    public var handledAt: Date?
    public init(id: Int, customer: CustomerKey, chatID: String, summary: String, fields: [String: String],
                status: BookingStatus = .pending, handledBy: String? = nil, handledAt: Date? = nil) {
        self.id = id; self.customer = customer; self.chatID = chatID; self.summary = summary; self.fields = fields
        self.status = status; self.handledBy = handledBy; self.handledAt = handledAt
    }
}

public enum StaffError: Error, CustomStringConvertible {
    case notFound(Int)
    case alreadyHandled(Int)
    public var description: String {
        switch self {
        case .notFound(let id): return "找不到編號 #\(id)"
        case .alreadyHandled(let id): return "#\(id) 已經處理過了"
        }
    }
}
