import Foundation

// MARK: - 伺服器設定（正式執行從環境變數讀；token、secret 不進 git、不寫進日誌）

public struct TelegramConfig: Sendable {
    public var token: [REDACTED]
    public var secret: [REDACTED]              // setWebhook 的 secret_token，Telegram 每個請求都會帶回來
    public var publicURL: String           // 對外網址（Cloudflare Tunnel），webhook 是 <publicURL>/telegram/webhook
    public var staffChatID: Int64?         // 店員群組：收待處理事項的通知，店員回覆通知就轉給顧客
    public var apiBase = "https://example.invalid"
    public var host = "0.0.0.0"          // 只聽本機，對外由 cloudflared 轉進來
    public var port = 8080                 // 0 = 隨便挑一個空的（測試用）
    public var retryDelays: [Double] = [1, 2, 4]
    public var typingInterval: Double = 5  // 還沒處理完，每隔幾秒再送一次「正在輸入」

    public init(token: [REDACTED], secret: [REDACTED], publicURL: String, staffChatID: Int64? = nil) {
        self.token = [REDACTED]; self.secret = [REDACTED]; self.staffChatID = staffChatID
        self.publicURL = publicURL.hasSuffix("/") ? String(publicURL.dropLast()) : publicURL
    }

    public var webhookPath: String { "/telegram/webhook" }
    public var webhookURL: String { publicURL + webhookPath }

    public enum ConfigError: Error, CustomStringConvertible {
        case missing([String])
        case invalid(String, String)
        // 只說是哪個變數，不印值
        public var description: String {
            switch self {
            case .missing(let names): return "缺少環境變數：\(names.joined(separator: "、"))"
            case .invalid(let name, let why): return "環境變數 \(name) 不正確：\(why)"
            }
        }
    }

    // TELEGRAM_BOT_TOKEN、TELEGRAM_WEBHOOK_SECRET、TELEGRAM_PUBLIC_URL 必填；
    // TELEGRAM_STAFF_CHAT_ID、TELEGRAM_API_BASE（測試用）、PORT 選填
    public static func fromEnvironment(_ env: [String: String] = ProcessInfo.processInfo.environment) throws -> TelegramConfig {
        let required = ["TELEGRAM_BOT_TOKEN", "TELEGRAM_WEBHOOK_SECRET", "TELEGRAM_PUBLIC_URL"]
        let missing = required.filter { (env[$0] ?? "").isEmpty }
        if !missing.isEmpty { throw ConfigError.missing(missing) }
        let secret = [REDACTED]"TELEGRAM_WEBHOOK_SECRET"]!
        // Telegram 的規定：1～256 個字，只能用 A-Z、a-z、0-9、_、-
        let allowed = CharacterSet(charactersIn: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-")
        guard secret.count <= 256, secret.unicodeScalars.allSatisfy(allowed.contains) else {
            throw ConfigError.invalid("TELEGRAM_WEBHOOK_SECRET", "只能用英數字、_、-，最多 256 個字")
        }
        let publicURL = env["TELEGRAM_PUBLIC_URL"]!
        guard publicURL.hasPrefix("https://") else { throw ConfigError.invalid("TELEGRAM_PUBLIC_URL", "要用 https:// 開頭") }
        var staff: Int64?
        if let s = env["TELEGRAM_STAFF_CHAT_ID"], !s.isEmpty {
            guard let id = Int64(s) else { throw ConfigError.invalid("TELEGRAM_STAFF_CHAT_ID", "要是數字（群組通常是負數）") }
            staff = id
        }
        var config = TelegramConfig(token: [REDACTED]"TELEGRAM_BOT_TOKEN"]!, secret: [REDACTED], publicURL: publicURL, staffChatID: staff)
        if let base = env["TELEGRAM_API_BASE"], !base.isEmpty { config.apiBase = base }
        if let p = env["PORT"], !p.isEmpty {
            guard let port = Int(p) else { throw ConfigError.invalid("PORT", "要是數字") }
            config.port = port
        }
        return config
    }
}
