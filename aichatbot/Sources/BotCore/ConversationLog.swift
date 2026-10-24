import Foundation

// MARK: - 對話紀錄：每則訊息一筆（顧客、機器人、店員），機器人那筆附上為了產生這則回覆做的每一次 AI 判斷
// 正式系統存在資料庫（Sources/Storage），預設和測試用記憶體版

public enum LogDirection: String, Codable, Sendable { case customer, bot, staff }

// 一次 AI 呼叫：題目名稱、送進去的文字、各欄位的回傳值、花了幾毫秒、是否失敗。不存 instructions（很長，同一版程式都一樣）。
// 逾時在模型層也是回傳 nil，這裡分不出失敗還是逾時，都記成 failed
public struct AITrace: Codable, Equatable, Sendable {
    public var name: String
    public var text: String
    public var values: [String: FieldValue]
    public var ms: Int
    public var failed: Bool
}

public struct LogEntry: Equatable {
    public var time: Date
    public var customer: CustomerKey
    public var chatID: String
    public var direction: LogDirection
    public var text: String
    public var replyKey: String?     // 機器人、店員的回覆種類（silent = 轉人工中不回覆）
    public var traces: [AITrace]
    public init(time: Date, customer: CustomerKey, chatID: String, direction: LogDirection, text: String,
                replyKey: String? = nil, traces: [AITrace] = []) {
        self.time = time; self.customer = customer; self.chatID = chatID; self.direction = direction
        self.text = text; self.replyKey = replyKey; self.traces = traces
    }
}

@MainActor public protocol ConversationLogStore: AnyObject {
    func append(_ entry: LogEntry)
    // 這位顧客最近 limit 則（舊的在前），略過沒有內容的（轉人工中 AI 不回覆的那筆）
    func recent(_ key: CustomerKey, limit: Int) -> [LogEntry]
}

@MainActor public final class InMemoryConversationLog: ConversationLogStore {
    public private(set) var entries: [LogEntry] = []
    public init() {}
    public func append(_ entry: LogEntry) { entries.append(entry) }
    public func recent(_ key: CustomerKey, limit: Int) -> [LogEntry] {
        Array(entries.filter { $0.customer == key && !$0.text.isEmpty }.suffix(limit))
    }
}

// 包住模型，記下每一次呼叫（處理一則訊息時建一個，處理完把 traces 寫進紀錄）
final class TracingModel: ChatModel, @unchecked Sendable {
    let base: ChatModel
    private let lock = NSLock()
    private var _traces: [AITrace] = []
    init(base: ChatModel) { self.base = base }

    var traces: [AITrace] { lock.withLock { _traces } }

    func respond(_ request: ModelRequest) async -> ModelReply? {
        let start = Date()
        let reply = await base.respond(request)
        let trace = AITrace(name: request.name, text: request.text, values: reply?.values ?? [:],
                            ms: Int(Date().timeIntervalSince(start) * 1000), failed: reply == nil)
        lock.withLock { _traces.append(trace) }
        return reply
    }
}
