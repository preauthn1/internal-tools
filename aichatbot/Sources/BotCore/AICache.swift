import Foundation
import CryptoKit

// MARK: - AI 回答快取：量測和驗收重跑題組時，同一個請求直接用存好的回答（AI 用 greedy，同一個請求每次回答都一樣）
// BOT_AI_CACHE=1 才開（bot 執行檔的量測和對話模式）；沒開時行為完全不變，Telegram 伺服器不使用
// 目錄：BOT_AI_CACHE_DIR，預設 ~/Library/Caches/aichatbot/ai；命名空間 BOT_AI_CACHE_NS=dev（預設）或 holdout，驗證題組分開存
// 檔名是請求的 SHA256，內容只有各欄位的值，不存明文輸入。清快取：python3 eval/clear_ai_cache.py 或直接刪資料夾

public final class CachedModel: ChatModel, @unchecked Sendable {
    let base: ChatModel
    let dir: URL
    let modelVersion: String
    let systemVersion: String
    private let state = NSLock()
    private var _hits = 0, _misses = 0

    // 包在 QueuedModel 外面：命中時不用排隊
    public init(base: ChatModel, dir: URL, modelVersion: String = AppleModel.version,
                systemVersion: String = ProcessInfo.processInfo.operatingSystemVersionString) {
        self.base = base; self.dir = dir; self.modelVersion = modelVersion; self.systemVersion = systemVersion
    }

    public var hits: Int { state.withLock { _hits } }
    public var misses: Int { state.withLock { _misses } }
    // 量測結果附上這一行（eval/run_flows.py、eval/snapshot.py、loop/score.py 會讀）
    public var summary: String { "AI 快取：命中 \(hits) 次、實際呼叫 \(misses) 次" }

    public func respond(_ request: ModelRequest) async -> ModelReply? {
        let file = path(Self.key(request, modelVersion: modelVersion, systemVersion: systemVersion))
        if let data = try? Data(contentsOf: file), let values = try? JSONDecoder().decode([String: FieldValue].self, from: data) {
            state.withLock { _hits += 1 }
            return ModelReply(values)
        }
        state.withLock { _misses += 1 }
        // 失敗或逾時（nil）不存，下次重新呼叫
        guard let reply = await base.respond(request) else { return nil }
        if let data = try? JSONEncoder().encode(reply.values) {
            try? FileManager.default.createDirectory(at: file.deletingLastPathComponent(), withIntermediateDirectories: true)
            try? data.write(to: file, options: .atomic)
        }
        return reply
    }

    func path(_ key: String) -> URL {
        dir.appendingPathComponent(String(key.prefix(2))).appendingPathComponent(key + ".json")
    }

    // 送給模型的全部內容＋模型版本＋系統版本（系統更新可能換模型），編成 JSON 陣列再算 SHA256，欄位串接不會互相碰撞
    public static func key(_ r: ModelRequest, modelVersion: String, systemVersion: String) -> String {
        let fields: [[String]] = r.fields.map { f in
            switch f {
            case .choice(let n, let d, let options): return ["choice", n, d] + options
            case .text(let n, let d): return ["text", n, d]
            case .integer(let n, let d): return ["integer", n, d]
            }
        }
        let parts: [Any] = ["v1", modelVersion, systemVersion, r.name, r.instructions, r.text, fields]
        let data = try! JSONSerialization.data(withJSONObject: parts, options: [.sortedKeys])
        return SHA256.hash(data: data).map { String(format: "%02x", $0) }.joined()
    }

    public static let defaultDir = FileManager.default.urls(for: .cachesDirectory, in: .userDomainMask)[0]
        .appendingPathComponent("aichatbot/ai")

    // BOT_AI_CACHE=1 時包一層快取，否則原樣回傳 base；命名空間只接受 dev、holdout
    public static func fromEnvironment(base: ChatModel, vars: [String: String] = env) -> ChatModel {
        guard vars["BOT_AI_CACHE"] == "1" else { return base }
        let ns = vars["BOT_AI_CACHE_NS"] ?? "dev"
        guard ["dev", "holdout"].contains(ns) else {
            FileHandle.standardError.write("❌ BOT_AI_CACHE_NS 只能是 dev 或 holdout（現在是 \(ns)）\n".data(using: .utf8)!)
            exit(1)
        }
        let root = vars["BOT_AI_CACHE_DIR"].map { URL(fileURLWithPath: $0) } ?? defaultDir
        return CachedModel(base: base, dir: root.appendingPathComponent(ns))
    }
}
