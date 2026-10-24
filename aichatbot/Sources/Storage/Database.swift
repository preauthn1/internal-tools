import Foundation
import GRDB
import BotCore

// MARK: - 資料庫（SQLite）：顧客狀態、已處理的 Telegram update、對話紀錄、待處理事項、預約申請
// server 執行檔用檔案（BOT_DB），終端機版 bot 用記憶體資料庫

public final class BotDatabase: Sendable {
    public let queue: DatabaseQueue

    // 對話紀錄保留 90 天、update 紀錄保留 7 天（purge 刪掉更舊的）
    public static let logRetentionDays = 90
    public static let updateRetentionDays = 7

    // 上層目錄不存在就建立；開檔後跑 migration
    public init(path: String) throws {
        try FileManager.default.createDirectory(at: URL(fileURLWithPath: path).deletingLastPathComponent(), withIntermediateDirectories: true)
        queue = try DatabaseQueue(path: path)
        try Self.migrator.migrate(queue)
    }

    // 記憶體資料庫：關掉就消失，不寫入正式資料
    public init() throws {
        queue = try DatabaseQueue()
        try Self.migrator.migrate(queue)
    }

    static var migrator: DatabaseMigrator {
        var m = DatabaseMigrator()
        m.registerMigration("v1") { db in
            // 顧客狀態：Customer 整個存成 JSON
            try db.create(table: "customer") { t in
                t.column("platform", .text).notNull()
                t.column("userID", .text).notNull()
                t.column("data", .text).notNull()
                t.column("updatedAt", .datetime).notNull()
                t.primaryKey(["platform", "userID"])
            }
            try db.create(table: "processed_update") { t in
                t.primaryKey("updateID", .integer)
                t.column("processedAt", .datetime).notNull()
            }
            // 對話紀錄：traces 是 AI 判斷的 JSON 陣列（顧客、店員那筆是空陣列）
            try db.create(table: "conversation_log") { t in
                t.autoIncrementedPrimaryKey("id")
                t.column("time", .datetime).notNull().indexed()
                t.column("platform", .text).notNull()
                t.column("userID", .text).notNull()
                t.column("chatID", .text).notNull()
                t.column("direction", .text).notNull()
                t.column("text", .text).notNull()
                t.column("replyKey", .text)
                t.column("traces", .text).notNull()
            }
            try db.create(index: "conversation_log_customer", on: "conversation_log", columns: ["platform", "userID", "time"])
        }
        m.registerMigration("v2") { db in
            // 待處理事項：groupMessageID 是店員群組裡那則通知的編號（店員「回覆」通知時用它找回來）
            try db.create(table: "ticket") { t in
                t.autoIncrementedPrimaryKey("id")
                t.column("kind", .text).notNull()
                t.column("platform", .text).notNull()
                t.column("userID", .text).notNull()
                t.column("chatID", .text).notNull()
                t.column("original", .text).notNull()
                t.column("text", .text).notNull()
                t.column("status", .text).notNull()          // open、replied
                t.column("relatedTo", .integer)               // 緊急通知：對應的原本那張
                t.column("groupMessageID", .integer).indexed()
                t.column("createdAt", .datetime).notNull()
            }
        }
        m.registerMigration("v3") { db in
            // 預約申請：編號和開出來的待處理事項相同；handledBy 是確認或拒絕的店員
            try db.create(table: "booking") { t in
                t.primaryKey("id", .integer)
                t.column("platform", .text).notNull()
                t.column("userID", .text).notNull()
                t.column("chatID", .text).notNull()
                t.column("date", .text)
                t.column("time", .text)
                t.column("people", .text)
                t.column("name", .text)
                t.column("phone", .text)
                t.column("summary", .text).notNull()
                t.column("status", .text).notNull()          // pending、confirmed、rejected
                t.column("handledBy", .text)
                t.column("handledAt", .datetime)
                t.column("createdAt", .datetime).notNull()
            }
        }
        return m
    }

    // 刪掉超過保留天數的對話紀錄和 update 紀錄，回傳各刪了幾筆
    @discardableResult public func purge(now: Date = Date()) throws -> (logs: Int, updates: Int) {
        let day = 86400.0
        return try queue.write { db in
            try db.execute(sql: "DELETE FROM conversation_log WHERE time < ?",
                           arguments: [now.addingTimeInterval(-Double(Self.logRetentionDays) * day)])
            let logs = db.changesCount
            try db.execute(sql: "DELETE FROM processed_update WHERE processedAt < ?",
                           arguments: [now.addingTimeInterval(-Double(Self.updateRetentionDays) * day)])
            return (logs, db.changesCount)
        }
    }
}

// MARK: - 顧客狀態

@MainActor public final class SQLiteCustomerStore: CustomerStore {
    let db: BotDatabase
    public init(_ db: BotDatabase) { self.db = db }

    // 沒有資料（或讀不出來）就是新顧客
    public func load(_ key: CustomerKey) -> Customer {
        let json = try? db.queue.read { db in
            try String.fetchOne(db, sql: "SELECT data FROM customer WHERE platform = ? AND userID = ?", arguments: [key.platform, key.userID])
        }
        guard let json, let c = try? JSONDecoder().decode(Customer.self, from: Data(json.utf8)) else { return Customer() }
        return c
    }

    public func save(_ key: CustomerKey, _ customer: Customer) {
        guard let data = try? JSONEncoder().encode(customer) else { return }
        try? db.queue.write { db in
            try db.execute(sql: """
                INSERT INTO customer (platform, userID, data, updatedAt) VALUES (?, ?, ?, ?)
                ON CONFLICT(platform, userID) DO UPDATE SET data = excluded.data, updatedAt = excluded.updatedAt
                """, arguments: [key.platform, key.userID, String(decoding: data, as: UTF8.self), Date()])
        }
    }
}

// MARK: - 待處理事項

// 待處理事項由資料庫自動編號；預約申請的編號和開出來的待處理事項相同
@MainActor public final class SQLiteTicketStore: TicketStore {
    let db: BotDatabase
    public init(_ db: BotDatabase) { self.db = db }

    // 依編號順序讀出全部（已回覆的不清，店員還能回覆已回覆過的通知）
    public var tickets: [Ticket] {
        ((try? db.queue.read { db in try Row.fetchAll(db, sql: "SELECT * FROM ticket ORDER BY id") }) ?? []).map(Self.ticket)
    }

    public func add(kind: TicketKind, customer: CustomerKey, chatID: String, original: String, text: String, relatedTo: Int?) -> Ticket {
        let id = (try? db.queue.write { db -> Int in
            try db.execute(sql: """
                INSERT INTO ticket (kind, platform, userID, chatID, original, text, status, relatedTo, createdAt)
                VALUES (?, ?, ?, ?, ?, ?, 'open', ?, ?)
                """, arguments: [kind.rawValue, customer.platform, customer.userID, chatID, original, text, relatedTo, Date()])
            return Int(db.lastInsertedRowID)
        }) ?? 0
        return Ticket(id: id, kind: kind, customer: customer, chatID: chatID, original: original, text: text, relatedTo: relatedTo)
    }

    // 只有狀態會變
    public func update(_ ticket: Ticket) {
        try? db.queue.write { db in
            try db.execute(sql: "UPDATE ticket SET status = ? WHERE id = ?", arguments: [ticket.done ? "replied" : "open", ticket.id])
        }
    }

    public func setGroupMessage(ticketID: Int, messageID: Int64) {
        try? db.queue.write { db in
            try db.execute(sql: "UPDATE ticket SET groupMessageID = ? WHERE id = ?", arguments: [messageID, ticketID])
        }
    }

    public func ticket(groupMessageID: Int64) -> Ticket? {
        let row = try? db.queue.read { db in try Row.fetchOne(db, sql: "SELECT * FROM ticket WHERE groupMessageID = ?", arguments: [groupMessageID]) }
        return row.flatMap { $0 }.map(Self.ticket)
    }

    // MARK: 預約申請

    static let bookingFields = ["date", "time", "people", "name", "phone"]

    public var bookings: [BookingRequest] {
        ((try? db.queue.read { db in try Row.fetchAll(db, sql: "SELECT * FROM booking ORDER BY id") }) ?? []).map(Self.booking)
    }

    public func booking(_ id: Int) -> BookingRequest? {
        let row = try? db.queue.read { db in try Row.fetchOne(db, sql: "SELECT * FROM booking WHERE id = ?", arguments: [id]) }
        return row.flatMap { $0 }.map(Self.booking)
    }

    public func add(_ booking: BookingRequest) {
        let b = booking, f = Self.bookingFields.map { b.fields[$0] }
        try? db.queue.write { db in
            try db.execute(sql: """
                INSERT INTO booking (id, platform, userID, chatID, date, time, people, name, phone, summary, status, handledBy, handledAt, createdAt)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                """, arguments: [b.id, b.customer.platform, b.customer.userID, b.chatID, f[0], f[1], f[2], f[3], f[4],
                                 b.summary, b.status.rawValue, b.handledBy, b.handledAt, Date()])
        }
    }

    // 只有狀態和處理的店員會變
    public func update(_ booking: BookingRequest) {
        try? db.queue.write { db in
            try db.execute(sql: "UPDATE booking SET status = ?, handledBy = ?, handledAt = ? WHERE id = ?",
                           arguments: [booking.status.rawValue, booking.handledBy, booking.handledAt, booking.id])
        }
    }

    static func ticket(_ r: Row) -> Ticket {
        let kind: String = r["kind"], status: String = r["status"]
        return Ticket(id: r["id"], kind: TicketKind(rawValue: kind) ?? .question, customer: CustomerKey(platform: r["platform"], userID: r["userID"]),
                      chatID: r["chatID"], original: r["original"], text: r["text"], relatedTo: r["relatedTo"], done: status == "replied")
    }

    // 欄位沒有值的不放進 fields（和記憶體版一樣）
    static func booking(_ r: Row) -> BookingRequest {
        let status: String = r["status"]
        var fields: [String: String] = [:]
        for key in bookingFields { if let v: String = r[key] { fields[key] = v } }
        return BookingRequest(id: r["id"], customer: CustomerKey(platform: r["platform"], userID: r["userID"]), chatID: r["chatID"],
                              summary: r["summary"], fields: fields, status: BookingStatus(rawValue: status) ?? .pending,
                              handledBy: r["handledBy"], handledAt: r["handledAt"])
    }
}

// MARK: - 已處理的 update（Telegram 層的 UpdateDedup 實作在 Sources/Telegram）

public final class SQLiteUpdateDedup: Sendable {
    let db: BotDatabase
    public init(_ db: BotDatabase) { self.db = db }

    // 第一次看到回傳 true 並記下來；寫不進資料庫時當成第一次（寧可多處理一次，不要漏掉訊息）
    public func firstTime(_ id: Int64, now: Date = Date()) -> Bool {
        (try? db.queue.write { db in
            try db.execute(sql: "INSERT OR IGNORE INTO processed_update (updateID, processedAt) VALUES (?, ?)", arguments: [id, now])
            return db.changesCount == 1
        }) ?? true
    }
}

// MARK: - 對話紀錄

@MainActor public final class SQLiteConversationLog: ConversationLogStore {
    let db: BotDatabase
    public init(_ db: BotDatabase) { self.db = db }

    public func append(_ e: LogEntry) {
        let traces = (try? JSONEncoder().encode(e.traces)).map { String(decoding: $0, as: UTF8.self) } ?? "[]"
        try? db.queue.write { db in
            try db.execute(sql: """
                INSERT INTO conversation_log (time, platform, userID, chatID, direction, text, replyKey, traces)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                """, arguments: [e.time, e.customer.platform, e.customer.userID, e.chatID, e.direction.rawValue, e.text, e.replyKey, traces])
        }
    }

    // 依寫入順序讀出（指定顧客時只讀那位）
    public func entries(for key: CustomerKey? = nil) -> [LogEntry] {
        let rows: [Row] = (try? db.queue.read { db in
            if let key {
                return try Row.fetchAll(db, sql: "SELECT * FROM conversation_log WHERE platform = ? AND userID = ? ORDER BY id",
                                        arguments: [key.platform, key.userID])
            }
            return try Row.fetchAll(db, sql: "SELECT * FROM conversation_log ORDER BY id")
        }) ?? []
        return rows.map(Self.entry)
    }

    // 最近 limit 則（舊的在前），略過沒有內容的（轉人工中 AI 不回覆的那筆）
    public func recent(_ key: CustomerKey, limit: Int) -> [LogEntry] {
        let rows: [Row] = (try? db.queue.read { db in
            try Row.fetchAll(db, sql: "SELECT * FROM conversation_log WHERE platform = ? AND userID = ? AND text != '' ORDER BY id DESC LIMIT ?",
                             arguments: [key.platform, key.userID, limit])
        }) ?? []
        return rows.reversed().map(Self.entry)
    }

    static func entry(_ r: Row) -> LogEntry {
        let traces: String = r["traces"], direction: String = r["direction"]
        return LogEntry(time: r["time"], customer: CustomerKey(platform: r["platform"], userID: r["userID"]), chatID: r["chatID"],
                        direction: LogDirection(rawValue: direction) ?? .bot, text: r["text"], replyKey: r["replyKey"],
                        traces: (try? JSONDecoder().decode([AITrace].self, from: Data(traces.utf8))) ?? [])
    }
}
