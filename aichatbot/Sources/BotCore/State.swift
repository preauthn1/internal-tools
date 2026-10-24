import Foundation

// MARK: - 對話狀態（依「平台＋使用者 ID」分開保存；這裡是介面和記憶體版，資料庫版在 Sources/Storage）

public struct Customer: Equatable, Codable {
    public var name: String?
    public var lang = "zh"
    public var booking: [String: String] = [:]
    public var askingField: String?          // 預約中：正在問哪個欄位（nil = 沒有在預約）
    public var openTicket = false            // 有事情已轉給店員、還沒有回覆
    public var followUps = 0                 // 這件事被追問了幾次
    public var offTopicStreak = 0            // 連續無關或聲稱的次數
    public var human = false                 // 已轉人工，AI 不再回覆
    public init() {}

    // 資料庫存 JSON：少了的欄位用預設值（之後加欄位，舊資料照樣讀得出來）
    public init(from decoder: Decoder) throws {
        let c = try decoder.container(keyedBy: CodingKeys.self)
        name = try c.decodeIfPresent(String.self, forKey: .name)
        lang = try c.decodeIfPresent(String.self, forKey: .lang) ?? lang
        booking = try c.decodeIfPresent([String: String].self, forKey: .booking) ?? booking
        askingField = try c.decodeIfPresent(String.self, forKey: .askingField)
        openTicket = try c.decodeIfPresent(Bool.self, forKey: .openTicket) ?? openTicket
        followUps = try c.decodeIfPresent(Int.self, forKey: .followUps) ?? followUps
        offTopicStreak = try c.decodeIfPresent(Int.self, forKey: .offTopicStreak) ?? offTopicStreak
        human = try c.decodeIfPresent(Bool.self, forKey: .human) ?? human
    }
}

@MainActor public protocol CustomerStore: AnyObject {
    func load(_ key: CustomerKey) -> Customer
    func save(_ key: CustomerKey, _ customer: Customer)
}

@MainActor public final class InMemoryCustomerStore: CustomerStore {
    var customers: [CustomerKey: Customer] = [:]
    public init() {}
    public func load(_ key: CustomerKey) -> Customer { customers[key] ?? Customer() }
    public func save(_ key: CustomerKey, _ customer: Customer) { customers[key] = customer }
}

// 待處理事項和預約申請（整間店共用，換顧客不清空）；資料庫版在 Sources/Storage。
// 通知送到店員群組後記下群組訊息的編號，店員「回覆」那則通知時用它找回待處理事項
@MainActor public protocol TicketStore: AnyObject {
    var tickets: [Ticket] { get }
    var bookings: [BookingRequest] { get }
    func add(kind: TicketKind, customer: CustomerKey, chatID: String, original: String, text: String, relatedTo: Int?) -> Ticket
    func update(_ ticket: Ticket)
    func setGroupMessage(ticketID: Int, messageID: Int64)
    func ticket(groupMessageID: Int64) -> Ticket?
    func add(_ booking: BookingRequest)
    func update(_ booking: BookingRequest)
    func booking(_ id: Int) -> BookingRequest?
}

extension TicketStore {
    // 沒有對應的原本那張（測試常用的寫法）
    @discardableResult
    public func add(kind: TicketKind, customer: CustomerKey, chatID: String, original: String = "", text: String) -> Ticket {
        add(kind: kind, customer: customer, chatID: chatID, original: original, text: text, relatedTo: nil)
    }
    public func ticket(_ id: Int) -> Ticket? { tickets.first { $0.id == id } }
    public func booking(_ id: Int) -> BookingRequest? { bookings.first { $0.id == id } }
}

@MainActor public final class InMemoryTicketStore: TicketStore {
    public private(set) var tickets: [Ticket] = []
    public private(set) var bookings: [BookingRequest] = []
    var groupMessages: [Int64: Int] = [:]   // 群組訊息編號 → 待處理事項編號
    public init() {}
    public func add(kind: TicketKind, customer: CustomerKey, chatID: String, original: String, text: String, relatedTo: Int?) -> Ticket {
        let t = Ticket(id: tickets.count + 1, kind: kind, customer: customer, chatID: chatID, original: original, text: text, relatedTo: relatedTo)
        tickets.append(t)
        return t
    }
    public func update(_ ticket: Ticket) {
        if let i = tickets.firstIndex(where: { $0.id == ticket.id }) { tickets[i] = ticket }
    }
    public func setGroupMessage(ticketID: Int, messageID: Int64) { groupMessages[messageID] = ticketID }
    public func ticket(groupMessageID: Int64) -> Ticket? { groupMessages[groupMessageID].flatMap { ticket($0) } }
    public func add(_ booking: BookingRequest) { bookings.append(booking) }
    public func update(_ booking: BookingRequest) {
        if let i = bookings.firstIndex(where: { $0.id == booking.id }) { bookings[i] = booking }
    }
}

// MARK: - 處理一則訊息時的脈絡：這位顧客的狀態，和這次產生的動作
// handle、bookingStep、classify、answerStoreQuestion、R 都是它的方法，讀寫的 customer 就是這位顧客

@MainActor final class Turn {
    var customer: Customer
    let key: CustomerKey
    let chatID: String
    let now: Date
    let tickets: TicketStore
    let model: ChatModel
    var original = ""                 // 這次顧客的原話（開待處理事項時一起存）
    var actions: [BotAction] = []

    init(_ customer: Customer, key: CustomerKey, chatID: String, now: Date = Date(), tickets: TicketStore, model: ChatModel) {
        self.customer = customer; self.key = key; self.chatID = chatID; self.now = now; self.tickets = tickets; self.model = model
    }

    // AI 選擇題：失敗或逾時回傳 nil，呼叫的地方各自用預設值
    func ask(_ name: String, _ text: String, instructions: String, fields: [Field]) async -> ModelReply? {
        await model.respond(ModelRequest(name: name, text: text, instructions: instructions, fields: fields))
    }

    // 開待處理事項；同時記下「有事在等店員」（追問的次數從這裡開始算）
    @discardableResult func openTicket(_ kind: TicketKind, _ item: String, waiting: Bool = true, relatedTo: Int? = nil) -> Ticket {
        let t = tickets.add(kind: kind, customer: key, chatID: chatID, original: original, text: item, relatedTo: relatedTo)
        actions.append(.openTicket(t))
        if waiting, !customer.openTicket { customer.openTicket = true; customer.followUps = 0 }
        return t
    }
}
