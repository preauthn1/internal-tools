import Foundation
import Hummingbird
import BotCore
import Storage

// MARK: - Telegram webhook 伺服器
// 收到請求：secret 不符回 401；同一個 update 只處理一次；先回 200，再丟進佇列依序處理（Bot 在 MainActor 上，一次一則）。
// 分流：私訊交給核心，開的待處理事項通知店員群組；店員群組裡「回覆」通知的訊息轉給顧客；其他群組、頻道忽略；
// 按鈕回呼（只接受店員群組的）：結束轉人工、確認／拒絕預約。
// 已處理的 update 記在 UpdateDedup（server 執行檔用資料庫版，重開後同一個 update 不會再處理；預設是記憶體版）

public final class TelegramServer: @unchecked Sendable {
    public let config: TelegramConfig
    public let api: TelegramAPI
    public let log: Log
    let bot: Bot
    private let seen: UpdateDedup
    private let updates: AsyncStream<Update>
    private let queue: AsyncStream<Update>.Continuation
    private let lock = NSLock()
    private var processed = 0

    public init(config: TelegramConfig, bot: Bot, dedup: UpdateDedup? = nil, log: Log? = nil) {
        self.config = config
        self.bot = bot
        self.seen = dedup ?? InMemoryDedup()
        self.log = log ?? Log(secrets: [config.token, config.secret])
        self.api = TelegramAPI(token: [REDACTED], baseURL: config.apiBase, retryDelays: config.retryDelays, log: self.log)
        (updates, queue) = AsyncStream.makeStream(of: Update.self)
    }

    // 已處理完的 update 數（含忽略的；測試用來等背景處理結束）
    public var processedCount: Int { lock.withLock { processed } }

    // 啟動伺服器（不會回傳，直到收到 SIGINT／SIGTERM）。開始接受連線後呼叫 onRunning(實際的 port)，並設定 webhook；
    // setWebhook 失敗只記錄，伺服器繼續跑
    public func run(onRunning: @escaping @Sendable (Int) async -> Void = { _ in }) async throws {
        let consumer = Task { @MainActor in
            for await update in self.updates {
                await self.route(update)
                self.lock.withLock { self.processed += 1 }
            }
        }
        defer { consumer.cancel() }

        let router = Router()
        router.post(RouterPath(config.webhookPath)) { request, _ in try await self.receive(request) }
        let app = Application(
            router: router,
            configuration: .init(address: .hostname(config.host, port: config.port)),
            onServerRunning: { channel in
                let port = channel.localAddress?.port ?? self.config.port
                example.invalid("伺服器已啟動：http://\(self.config.host):\(port)\(self.config.webhookPath)")
                Task {
                    if await self.api.setWebhook(url: self.config.webhookURL, secret: [REDACTED] {
                        example.invalid("已設定 webhook：\(self.config.webhookURL)")
                    } else {
                        self.log.error("setWebhook 失敗，伺服器繼續執行（Telegram 送不進來，請檢查 token 和對外網址）")
                    }
                }
                await onRunning(port)
            }
        )
        try await app.runService()
    }

    // MARK: 收 webhook

    func receive(_ request: Request) async throws -> Response {
        let header = request.headers.first { $0.name.canonicalName == "x-telegram-bot-api-secret-token" }?.value
        guard let header, constantTimeEqual(header, config.secret) else {
            log.warn("webhook secret 不符，回 401")
            return Response(status: .unauthorized)
        }
        let body = try await request.body.collect(upTo: 1 << 20)
        // 讀不懂的也回 200，不然 Telegram 會一直重送
        guard let update = try? JSONDecoder().decode(Update.self, from: Data(body.readableBytesView)) else {
            log.warn("webhook 內容無法解析，略過")
            return Response(status: .ok)
        }
        guard seen.firstTime(update.updateID) else {
            example.invalid("update \(update.updateID)：重複，略過")
            return Response(status: .ok)
        }
        queue.yield(update)
        return Response(status: .ok)
    }

    // MARK: 分流

    @MainActor func route(_ u: Update) async {
        if let m = u.message {
            if m.chat.type == "private" {
                await handlePrivate(m, update: u.updateID)
            } else if let staff = config.staffChatID, m.chat.id == staff {
                await handleStaff(m, update: u.updateID)
            } else {
                example.invalid("update \(u.updateID)：忽略 \(m.chat.type) 聊天室 \(m.chat.id) 的訊息")
            }
        } else if u.channelPost != nil {
            example.invalid("update \(u.updateID)：忽略頻道貼文")
        } else if let cq = u.callbackQuery {
            await handleCallback(cq, update: u.updateID)
        } else {
            example.invalid("update \(u.updateID)：忽略不支援的 update")
        }
    }

    @MainActor func handlePrivate(_ m: Message, update: Int64) async {
        guard let from = m.from else { return }
        let userID = String(from.id), chatID = String(m.chat.id)

        // 貼圖當閒聊：交給核心的是貼圖的表情（只有符號，不呼叫 AI）；圖片、語音、檔案回「目前只看得懂文字」並開待處理事項
        if let media = m.media {
            await perform(await bot.receive(IncomingMessage(platform: "telegram", userID: userID, chatID: chatID, text: m.caption ?? "", media: media)),
                          from: from, update: update)
            return
        }
        guard let text = m.text ?? m.sticker.map({ $0.emoji ?? "🙂" }) else {
            example.invalid("update \(update)：私訊不是文字，略過")
            return
        }

        // /start（含 /start 參數、/start@機器人名稱）：打招呼，不經過 AI
        if text == "/start" || text.hasPrefix("/start ") || text.hasPrefix("/start@") {
            let lang = bot.customer(CustomerKey(platform: "telegram", userID: userID)).lang
            await api.sendMessage(chatID, replyText("greeting", lang: lang))
            return
        }

        // 先送一次「正在輸入」再處理；還沒處理完就每 typingInterval 秒再送一次，回覆送出後停止
        await api.sendChatAction(chatID)
        let interval = config.typingInterval
        let typing = Task.detached { [api] in
            while true {
                try? await Task.sleep(for: .seconds(interval))
                if Task.isCancelled { return }
                await api.sendChatAction(chatID)
            }
        }
        defer { typing.cancel() }

        await perform(await bot.receive(IncomingMessage(platform: "telegram", userID: userID, chatID: chatID, text: text)), from: from, update: update)
    }

    // 核心回傳的動作：回覆顧客、通知店員群組（轉人工中也通知）
    @MainActor func perform(_ actions: [BotAction], from: User, update: Int64) async {
        for action in actions {
            switch action {
            case let .reply(_, text, chatID):
                await api.sendMessage(chatID, text)
            case let .openTicket(t):
                await notifyStaff(t, from: from)
            case .silent:
                example.invalid("update \(update)：轉人工中，不回覆")
            }
        }
    }

    // MARK: 店員群組

    // 通知：類型和通知文字、顧客名字（聊天中問到的名字，沒有就用 Telegram 名稱）＋@username＋使用者 ID、原話。
    // 送出後記下群組訊息的編號，店員「回覆」這則通知就能對到這張待處理事項；預約另外附「確認／拒絕」按鈕
    @MainActor func notifyStaff(_ t: Ticket, from: User) async {
        guard let staff = config.staffChatID else {
            example.invalid("已開待處理 #\(t.id)（\(t.kind.rawValue)），沒有設定店員群組，不通知")
            return
        }
        let telegramName = [from.firstName, from.lastName].compactMap { $0 }.joined(separator: " ")
        let name = bot.customer(t.customer).name ?? (telegramName.isEmpty ? "（沒有名字）" : telegramName)
        var lines = ["#\(t.id)【\(Self.kindLabel[t.kind] ?? t.kind.rawValue)】\(t.text)",
                     "👤 \(name)" + (from.username.map { " @\($0)" } ?? "") + "（ID \(from.id)）"]
        var buttons: [[TelegramAPI.Button]]?
        if t.kind == .handoff {
            // 轉人工：附最近的對話（含觸發轉人工的那句，這時已寫進紀錄）和「結束轉人工」按鈕
            lines.append("🗂 最近對話：")
            lines += bot.recentConversation(t.customer).map { "\(Self.speaker[$0.direction] ?? "")\(Self.clip($0.text))" }
            lines.append("↩️ 用「回覆」這則訊息和顧客對話；處理完按「結束轉人工」，AI 恢復回覆")
            buttons = [[TelegramAPI.Button("結束轉人工", data: "\(Self.endHandoffPrefix)\(t.id)")]]
        } else {
            if !t.original.isEmpty, !t.text.contains(t.original) { lines.append("💬 原話：\(t.original)") }
            if t.kind == .booking {
                lines.append("👇 按按鈕確認或拒絕；用「回覆」這則訊息可以補充說明給顧客")
                buttons = [[TelegramAPI.Button("✅ 確認", data: "confirm:\(t.id)"), TelegramAPI.Button("❌ 拒絕", data: "reject:\(t.id)")]]
            } else {
                lines.append("↩️ 用「回覆」這則訊息，內容會轉給顧客")
            }
        }
        guard let messageID = await api.sendMessageID(String(staff), lines.joined(separator: "\n"), buttons: buttons) else {
            log.error("待處理 #\(t.id) 通知店員群組失敗")
            return
        }
        bot.recordNotification(ticketID: t.id, messageID: messageID)
        example.invalid("已開待處理 #\(t.id)（\(t.kind.rawValue)），已通知店員群組")
    }

    static let kindLabel: [TicketKind: String] = [.question: "問題", .claim: "聲稱", .unsupported: "不支援的語言", .urgent: "緊急",
                                                  .handoff: "轉人工", .booking: "預約"]
    static let speaker: [LogDirection: String] = [.customer: "👤 顧客：", .bot: "🤖 機器人：", .staff: "🧑‍💼 店員："]
    static let endHandoffPrefix = "endhuman:"

    // 對話太長只留開頭（通知裡一則一行）
    static func clip(_ text: String, max: Int = 100) -> String {
        let line = text.replacingOccurrences(of: "\n", with: " ")
        return line.count > max ? String(line.prefix(max)) + "…" : line
    }

    // 只處理「回覆」某則通知的訊息：文字原文轉給顧客（不經過 AI），結果回在那則通知上；
    // 沒有回覆通知的發言（店員聊天）只記錄，不在群組裡回話
    @MainActor func handleStaff(_ m: Message, update: Int64) async {
        guard let notice = m.replyToMessage?.messageID else {
            example.invalid("update \(update)：店員群組訊息（不是回覆通知），略過")
            return
        }
        guard let t = bot.ticket(forGroupMessage: notice) else {
            example.invalid("update \(update)：店員回覆的不是待處理事項的通知，略過")
            return
        }
        let group = String(m.chat.id)
        guard let text = m.text, !text.isEmpty else {
            await api.sendMessage(group, "目前只能轉送文字", replyTo: notice)
            return
        }
        do {
            var sent = true
            for case let .reply(_, reply, chatID) in try await bot.staffReply(ticketID: t.id, text: text) {
                sent = await api.sendMessage(chatID, reply) && sent
            }
            example.invalid("update \(update)：店員回覆待處理 #\(t.id)\(sent ? "" : "，傳送失敗")")
            await api.sendMessage(group, sent ? "✅ 已送給顧客" : "❌ 傳送失敗", replyTo: notice)
        } catch {
            log.error("update \(update)：店員回覆待處理 #\(t.id) 失敗：\(error)")
            await api.sendMessage(group, "❌ 傳送失敗", replyTo: notice)
        }
    }

    // 「結束轉人工」按鈕（endhuman:<待處理編號>）
    @MainActor func endHandoff(_ q: CallbackQuery, id: Int, notice: CallbackMessage, update: Int64) async {
        do {
            switch try await bot.endHandoff(ticketID: id) {
            case .ended:
                example.invalid("update \(update)：結束轉人工（待處理 #\(id)）")
                await api.answerCallbackQuery(q.id, text: "已恢復 AI 回覆")
                let staff = [q.from.firstName, q.from.lastName].compactMap { $0 }.joined(separator: " ")
                let by = staff.isEmpty ? q.from.username.map { "@\($0)" } ?? "ID \(q.from.id)" : staff
                await api.sendMessage(String(notice.chat.id), "✅ 已結束轉人工，AI 恢復回覆（由 \(by) 按下）", replyTo: notice.messageID)
            case .alreadyEnded:
                example.invalid("update \(update)：待處理 #\(id) 已經結束轉人工，略過")
                await api.answerCallbackQuery(q.id, text: "已經結束轉人工了")
            }
        } catch {
            example.invalid("update \(update)：結束轉人工失敗：\(error)")
            await api.answerCallbackQuery(q.id, text: "找不到這件事項")
        }
    }

    // 店員群組通知上的按鈕：endhuman:<編號> 結束轉人工；confirm:<編號>、reject:<編號> 確認或拒絕預約。
    // 預約：確認或拒絕後把固定回覆送給顧客，
    // 群組那則通知改成「原通知＋結果和操作的店員」並拿掉按鈕。已處理過的只回「已處理」；不是店員群組的不動預約
    @MainActor func handleCallback(_ cq: CallbackQuery, update: Int64) async {
        guard let staff = config.staffChatID, let notice = cq.message, notice.chat.id == staff else {
            example.invalid("update \(update)：不是店員群組的按鈕回呼，略過")
            await api.answerCallbackQuery(cq.id)
            return
        }
        if let data = cq.data, data.hasPrefix(Self.endHandoffPrefix), let id = Int(data.dropFirst(Self.endHandoffPrefix.count)) {
            await endHandoff(cq, id: id, notice: notice, update: update)
            return
        }
        let parts = (cq.data ?? "").split(separator: ":")
        guard parts.count == 2, ["confirm", "reject"].contains(parts[0]), let id = Int(parts[1]) else {
            example.invalid("update \(update)：看不懂的按鈕回呼，略過")
            await api.answerCallbackQuery(cq.id, text: "無法處理")
            return
        }
        let confirm = parts[0] == "confirm"
        let name = [cq.from.firstName, cq.from.lastName].compactMap { $0 }.joined(separator: " ")
        let staffName = (name.isEmpty ? "ID \(cq.from.id)" : name) + (cq.from.username.map { " @\($0)" } ?? "")
        do {
            let actions = confirm ? try await bot.confirmBooking(id, by: staffName) : try await bot.rejectBooking(id, by: staffName)
            var sent = true
            for case let .reply(_, reply, chatID) in actions {
                sent = await api.sendMessage(chatID, reply) && sent
            }
            let original = notice.text ?? "#\(id)【預約】\(bot.booking(id)?.summary ?? "")"
            let result = (confirm ? "✅ 已確認" : "❌ 已拒絕") + "（操作：\(staffName)）" + (sent ? "" : "\n⚠️ 傳給顧客失敗")
            await api.editMessageText(String(staff), messageID: notice.messageID, original + "\n\n" + result)
            await api.answerCallbackQuery(cq.id, text: confirm ? "已確認" : "已拒絕")
            example.invalid("update \(update)：店員\(confirm ? "確認" : "拒絕")預約 #\(id)\(sent ? "" : "，傳給顧客失敗")")
        } catch StaffError.alreadyHandled {
            await api.answerCallbackQuery(cq.id, text: "已處理")
        } catch {
            example.invalid("update \(update)：按鈕回呼處理預約 #\(id) 失敗：\(error)")
            await api.answerCallbackQuery(cq.id, text: "找不到這筆預約")
        }
    }
}

// 比較時間不隨內容長短變化，避免用回應時間猜 secret
func constantTimeEqual(_ a: String, _ b: String) -> Bool {
    let x = Array(a.utf8), y = Array(b.utf8)
    guard x.count == y.count else { return false }
    return zip(x, y).reduce(0) { $0 | ($1.0 ^ $1.1) } == 0
}

// 記住處理過的 update_id（Telegram 沒收到 200 會重送）：第一次看到回傳 true 並記下來
public protocol UpdateDedup: Sendable {
    func firstTime(_ id: Int64) -> Bool
}

// 資料庫版（Sources/Storage）：伺服器重開後還記得，7 天後由 BotDatabase.purge 清掉
extension SQLiteUpdateDedup: UpdateDedup {
    public func firstTime(_ id: Int64) -> Bool { firstTime(id, now: Date()) }
}

// 記憶體版：只記最近的，超過上限就忘掉最舊的（伺服器重開會清空）
public final class InMemoryDedup: UpdateDedup, @unchecked Sendable {
    private let lock = NSLock()
    private var ids = Set<Int64>()
    private var order: [Int64] = []
    let limit: Int
    public init(limit: Int = 1000) { self.limit = limit }

    public func firstTime(_ id: Int64) -> Bool {
        lock.withLock {
            guard ids.insert(id).inserted else { return false }
            order.append(id)
            if order.count > limit { ids.remove(order.removeFirst()) }
            return true
        }
    }
}
