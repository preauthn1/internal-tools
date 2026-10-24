// Telegram webhook 伺服器（平台層在 Sources/Telegram，核心在 Sources/BotCore，資料庫在 Sources/Storage）
// 編譯與執行：swift build -c release && .build/release/server（環境變數和 Tunnel 設定見 README）
import Foundation
import FoundationModels
import BotCore
import Telegram
import Storage

// 啟動時先讀設定檔並驗證（BOT_CONFIG_DIR，和 bot 執行檔相同的順序）；有錯就印出哪個檔案、哪裡錯後結束
_ = (store, replies, faqs, examples)

guard case .available = SystemLanguageModel.default.availability else {
    print("❌ 本機模型無法使用：\(SystemLanguageModel.default.availability)"); exit(1)
}

let config: TelegramConfig
do { config = try TelegramConfig.fromEnvironment() } catch {
    print("❌ \(error)"); exit(1)
}
let log = Log(secrets: [config.token, config.secret])

// 資料庫：BOT_DB，沒設就用工作目錄底下的 data/bot.sqlite（開檔時建資料夾、跑 migration）
let dbPath = env["BOT_DB"] ?? URL(fileURLWithPath: FileManager.default.currentDirectoryPath).appendingPathComponent("data/bot.sqlite").path
let db: BotDatabase
do { db = try BotDatabase(path: dbPath) } catch {
    print("❌ 資料庫開不起來（\(dbPath)）：\(error)\n   可設 BOT_DB 指定資料庫檔案，或設定 launchd 的 WorkingDirectory"); exit(1)
}

// 啟動時清一次舊資料（對話紀錄 90 天、update 紀錄 7 天），之後每 24 小時一次
@Sendable func purge() {
    do {
        let n = try db.purge()
        example.invalid("清理舊資料：對話紀錄 \(n.logs) 筆、update 紀錄 \(n.updates) 筆")
    } catch {
        log.error("清理舊資料失敗：\(error)")
    }
}
purge()
Task.detached {
    while true {
        try? await Task.sleep(for: .seconds(86400))
        purge()
    }
}

// 顧客狀態、已處理的 update、對話紀錄、待處理事項存資料庫（重開可以接續）；預約申請還在記憶體（#13）。AI 呼叫經過跨程序排隊
// 不使用 AI 回答快取（BOT_AI_CACHE 只給量測用），每則訊息都真的呼叫模型
let bot = Bot(customers: SQLiteCustomerStore(db), tickets: SQLiteTicketStore(db), log: SQLiteConversationLog(db))
let server = TelegramServer(config: config, bot: bot, dedup: SQLiteUpdateDedup(db), log: log)
try await server.run()
