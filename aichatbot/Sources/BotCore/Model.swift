import Foundation

// MARK: - 模型層：給顧客的話、指示和欄位定義，回傳各欄位的值（失敗或逾時回傳 nil）
// 業務程式只透過 ChatModel 呼叫 AI，不直接碰 FoundationModels；實作有 Apple 本機模型（AI.swift）和照腳本回答的假模型

public enum Field: Sendable {
    case choice(String, String, [String])   // 名稱、說明、選項
    case text(String, String)               // 名稱、說明
    case integer(String, String)            // 名稱、說明

    public var name: String {
        switch self { case .choice(let n, _, _), .text(let n, _), .integer(let n, _): return n }
    }
}

// 模型服務（Sources/model-server）收到的 JSON schema 轉回欄位定義，再交給 AppleModel 產生（和 Swift 版的生成結果完全一樣）。
// 只認 Go 版 model.Schema 產生的格式：string＋enum 是選擇題、string 是文字、integer 是整數；
// 順序照 required（properties 解析成字典後順序會亂）。看不懂就回傳 nil
public func fieldsFromJSONSchema(_ schema: Any) -> [Field]? {
    guard let root = schema as? [String: Any], let properties = root["properties"] as? [String: Any],
          let required = root["required"] as? [String] else { return nil }
    var fields: [Field] = []
    for name in required {
        guard let p = properties[name] as? [String: Any], let type = p["type"] as? String else { return nil }
        let description = p["description"] as? String ?? ""
        switch (type, p["enum"]) {
        case ("string", let options as [String]): fields.append(.choice(name, description, options))
        case ("string", nil): fields.append(.text(name, description))
        case ("integer", nil): fields.append(.integer(name, description))
        default: return nil
        }
    }
    return fields
}

public struct ModelRequest: Sendable {
    public let name: String          // 題目名稱（Intent、Pick、Booking…），也是 schema 的根名稱
    public let text: String          // 顧客的話
    public let instructions: String
    public let fields: [Field]       // 欄位依序產生，順序會影響結果
    public init(name: String, text: String, instructions: String, fields: [Field]) {
        self.name = name; self.text = text; self.instructions = instructions; self.fields = fields
    }
}

public enum FieldValue: Sendable, Equatable {
    case text(String)
    case integer(Int)
}

// 存進對話紀錄時直接寫成 JSON 的字串或數字
extension FieldValue: Codable {
    public init(from decoder: Decoder) throws {
        let c = try decoder.singleValueContainer()
        if let n = try? c.decode(Int.self) { self = .integer(n) } else { self = .text(try c.decode(String.self)) }
    }
    public func encode(to encoder: Encoder) throws {
        var c = encoder.singleValueContainer()
        switch self { case .text(let s): try c.encode(s); case .integer(let n): try c.encode(n) }
    }
}

// 讀不到的欄位就是 nil（和原本每個欄位各自 try? 的行為相同）
public struct ModelReply: Sendable {
    public var values: [String: FieldValue]
    public init(_ values: [String: FieldValue]) { self.values = values }
    public func string(_ name: String) -> String? {
        if case .text(let s)? = values[name] { return s }
        return nil
    }
    public func int(_ name: String) -> Int? {
        if case .integer(let n)? = values[name] { return n }
        return nil
    }
}

public protocol ChatModel: Sendable {
    func respond(_ request: ModelRequest) async -> ModelReply?
}

// MARK: - 假模型：依題目名稱照腳本依序回答，腳本沒寫到的題目回傳 nil（模擬 AI 失敗）。測試用，不必等本機模型

public final class ScriptedModel: ChatModel, @unchecked Sendable {
    private let lock = NSLock()
    private var script: [String: [[String: FieldValue]]]
    private var _requests: [ModelRequest] = []
    let delay: Double   // 每次回答前等幾秒（排隊測試用來模擬慢模型）

    public init(_ script: [String: [[String: FieldValue]]] = [:], delay: Double = 0) {
        self.script = script; self.delay = delay
    }

    // 收到的請求（依呼叫順序）
    public var requests: [ModelRequest] { lock.withLock { _requests } }

    public func respond(_ request: ModelRequest) async -> ModelReply? {
        let answer: [String: FieldValue]? = lock.withLock {
            _requests.append(request)
            guard var queue = script[request.name], !queue.isEmpty else { return nil }
            let first = queue.removeFirst()
            script[request.name] = queue
            return first
        }
        if delay > 0 { try? await Task.sleep(for: .seconds(delay)) }
        return answer.map(ModelReply.init)
    }
}
