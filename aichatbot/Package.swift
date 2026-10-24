// swift-tools-version:6.2
// 只需要 Command Line Tools：swift build -c release && .build/release/bot；測試：swift build && .build/debug/bot-tests && .build/debug/telegram-tests
import PackageDescription

let package = Package(
    name: "aichatbot",
    platforms: [.macOS(.v26)],   // FoundationModels 需要 macOS 26 以上
    products: [
        .executable(name: "bot", targets: ["bot"]),
        .executable(name: "server", targets: ["server"]),
        .executable(name: "model-server", targets: ["model-server"]),
    ],
    dependencies: [
        .package(url: "https://example.invalid/hummingbird-project/hummingbird.git", from: "2.0.0"),
        .package(url: "https://example.invalid/groue/GRDB.swift.git", from: "7.0.0"),
    ],
    targets: [
        // 核心：分類、挑 FAQ、預約、分派（自動改進迴圈只能改這裡和 eval/ 的三個資料檔）
        .target(name: "BotCore", path: "Sources/BotCore"),
        // 資料庫（SQLite，GRDB）：顧客狀態、已處理的 update、對話紀錄、待處理事項。BotCore 只定介面、不依賴 GRDB
        .target(name: "Storage", dependencies: ["BotCore", .product(name: "GRDB", package: "GRDB.swift")], path: "Sources/Storage"),
        // 執行檔：終端機對話和 BOT_INTENT_EVAL、BOT_TESTSET 量測模式（終端機對話用記憶體資料庫）
        .executableTarget(name: "bot", dependencies: ["BotCore", "Storage"], path: "Sources/bot"),
        // Telegram 平台層：webhook 路由、update 分流、Telegram API 用戶端（含重試）、「正在輸入」
        .target(name: "Telegram", dependencies: ["BotCore", "Storage", .product(name: "Hummingbird", package: "hummingbird")],
                path: "Sources/Telegram"),
        // 執行檔：Telegram webhook 伺服器（設定從環境變數讀，見 README）
        .executableTarget(name: "server", dependencies: ["Telegram", "Storage"], path: "Sources/server"),
        // 執行檔：模型服務，只包 Apple 本機模型的 HTTP 服務，給 Go 版呼叫（見 README「模型層」）
        .executableTarget(name: "model-server", dependencies: ["BotCore", .product(name: "Hummingbird", package: "hummingbird")],
                          path: "Sources/model-server"),
        // 對話層測試：店員操作、顧客狀態分開、AI 排隊、資料庫（會呼叫 AI 的路徑用假模型，不需要本機模型）。
        // Command Line Tools 沒有 XCTest、也載不到 Swift Testing 的巨集，所以寫成執行檔
        .executableTarget(name: "bot-tests", dependencies: ["BotCore", "Storage"], path: "Tests/bot-tests"),
        // Webhook 層測試：對伺服器送真的 HTTP 請求，Telegram API 換成本機的假端點（分開寫，bot-tests 不必依賴 Hummingbird）
        .executableTarget(name: "telegram-tests", dependencies: ["Telegram", "BotCore", "Storage", .product(name: "Hummingbird", package: "hummingbird")],
                          path: "Tests/telegram-tests"),
    ],
    swiftLanguageModes: [.v5]   // 全域狀態只剩設定檔（store、replies、faqs…），沿用原本的寫法，不套 Swift 6 嚴格併發檢查
)
