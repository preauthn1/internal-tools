import Foundation
import FoundationModels

// MARK: - Apple 本機模型（模型層的第一個實作；排隊由 QueuedModel 負責）

// 有逾時的呼叫，結果分得出是逾時還是 op 丟出錯誤（錯誤留型別名稱和訊息，Error 不一定是 Sendable）；模型服務的 log 要寫出是哪一種（#59）
public enum TimedResult<T: Sendable>: Sendable { case value(T), timedOut, failed(type: String, message: String) }

// 提示詞超過本機模型上下文（8192 tokens，輸入加輸出共用）時 .failed 的 type；模型服務看到它就在 502 帶 reason context_exceeded（#65）
public let contextExceededType = "GenerationError.exceededContextWindowSize"
// 被 Apple 的安全過濾擋下（May contain unsafe content）時 .failed 的 type；模型服務看到它就在 502 帶 reason guardrail（#68）
public let guardrailType = "GenerationError.guardrailViolation"

// 安全過濾擋下的錯誤：macOS 26 是 GenerationError，macOS 27 之後可能改丟 LanguageModelError
func isGuardrailViolation(_ error: any Error) -> Bool {
    if case LanguageModelSession.GenerationError.guardrailViolation = error { return true }
    if #available(macOS 27, *), case LanguageModelError.guardrailViolation = error { return true }
    return false
}

public func runWithTimeout<T: Sendable>(_ seconds: Double, _ op: @escaping @Sendable () async throws -> T) async -> TimedResult<T> {
    await withTaskGroup(of: TimedResult<T>.self) { group in
        group.addTask {
            do { return .value(try await op()) } catch {
                if case LanguageModelSession.GenerationError.exceededContextWindowSize = error {
                    return .failed(type: contextExceededType, message: String(describing: error))
                }
                if isGuardrailViolation(error) { return .failed(type: guardrailType, message: String(describing: error)) }
                return .failed(type: String(describing: type(of: error)), message: String(describing: error))
            }
        }
        group.addTask { try? await Task.sleep(for: .seconds(seconds)); return .timedOut }
        let first = await group.next() ?? .timedOut
        group.cancelAll()
        return first
    }
}

// 本機模型一次回答用掉的 token 數（模型服務回給 Go 版記進 AI 費用頁，#60）；輸入含指示、提示和輸出格式
public struct ModelUsage: Sendable, Encodable {
    public let inputTokens: Int
    public let outputTokens: Int
    // 模型服務回應 JSON 的欄位名稱（Go 端照這個讀）
    enum CodingKeys: String, CodingKey {
        case inputTokens = "input_tokens", outputTokens = "output_tokens"
    }
}

public struct AppleModel: ChatModel {
    public init() {}

    // 模型版本：AI 回答快取（AICache.swift）的 key 包含這個字串。改了下面的生成參數（取樣、長度、逾時）要一起改，舊的快取才會失效
    public static let version = "apple-fm/greedy/max150/timeout15/permissive-translate"
    public static let timeout: Double = 15

    // 只改寫既有內容的請求（翻譯）：安全過濾常把正常的句子當成不安全擋下（#68），改用寬鬆的過濾。
    // 寬鬆過濾只對純文字回答有效（guided generation 照樣擋），所以這些請求只能有一個文字欄位，回答直接放進那個欄位
    static let permissiveRequests: Set<String> = ["Translate", "TranslateCustomer"]
    public static func usesPermissiveGuardrails(_ name: String) -> Bool { permissiveRequests.contains(name) }

    // greedy：每次都選機率最高的字，同一句話每次結果相同（預設的隨機抽樣讓同一題時對時錯，分類 7～9 題、單題 9 題不穩定；改了之後都是 0）
    public func respond(_ request: ModelRequest) async -> ModelReply? {
        if case .value(let answer) = await respondDetailed(request) { return answer.reply }
        return nil
    }

    // 和 respond 一樣，但失敗時分得出是逾時還是本機模型丟出錯誤（模型服務寫 log 用），成功時另外附 token 數（算不出來是 nil）
    public func respondDetailed(_ request: ModelRequest) async -> TimedResult<(reply: ModelReply, usage: ModelUsage?)> {
        if Self.usesPermissiveGuardrails(request.name), request.fields.count == 1, case .text(let field, _) = request.fields[0] {
            return await respondPermissive(request, field: field)
        }
        let schema = Self.schema(request)
        let result = await runWithTimeout(Self.timeout) {
            let session = LanguageModelSession(instructions: request.instructions)
            return (session, try await session.respond(to: request.text, schema: schema, options: GenerationOptions(sampling: .greedy, maximumResponseTokens: 150)))
        }
        switch result {
        // token 數在逾時之外算：已經有的回答不會因為算 token 變成逾時
        case .value(let (session, response)): return .value((Self.reply(response.content, request.fields), await Self.usage(response, session)))
        case .timedOut: return .timedOut
        case .failed(let type, let message): return .failed(type: type, message: message)
        }
    }

    // 寬鬆過濾的請求：不帶 schema，純文字的回答放進唯一的文字欄位 field
    func respondPermissive(_ request: ModelRequest, field: String) async -> TimedResult<(reply: ModelReply, usage: ModelUsage?)> {
        let result = await runWithTimeout(Self.timeout) {
            let model = SystemLanguageModel(guardrails: .permissiveContentTransformations)
            let session = LanguageModelSession(model: model, instructions: request.instructions)
            return (session, try await session.respond(to: request.text, options: GenerationOptions(sampling: .greedy, maximumResponseTokens: 150)))
        }
        switch result {
        case .value(let (session, response)): return .value((ModelReply([field: .text(response.content)]), await Self.usage(response, session)))
        case .timedOut: return .timedOut
        case .failed(let type, let message): return .failed(type: type, message: message)
        }
    }

    // macOS 27 的回應自己帶用量（不多花時間）；沒有或是 0 時用 tokenCount 算這次 session 的 transcript：指示＋提示（含輸出格式）是輸入，回答是輸出。
    // 算失敗不影響回答，只是不附 token 數
    static func usage<Content>(_ response: LanguageModelSession.Response<Content>, _ session: LanguageModelSession) async -> ModelUsage? {
        if #available(macOS 27, *), response.usage.input.totalTokenCount > 0 {
            return ModelUsage(inputTokens: response.usage.input.totalTokenCount, outputTokens: response.usage.output.totalTokenCount)
        }
        guard #available(macOS 26.4, *) else { return nil }
        let entries = Array(session.transcript)
        let input = entries.filter { if case .response = $0 { false } else { true } }
        let output = entries.filter { if case .response = $0 { true } else { false } }
        guard let inputTokens = try? await SystemLanguageModel.default.tokenCount(for: input),
              let outputTokens = try? await SystemLanguageModel.default.tokenCount(for: output) else { return nil }
        return ModelUsage(inputTokens: inputTokens, outputTokens: outputTokens)
    }

    static func reply(_ content: GeneratedContent, _ fields: [Field]) -> ModelReply {
        // 每個欄位各自讀，讀不到的略過（其他欄位照用）
        var values: [String: FieldValue] = [:]
        for f in fields {
            switch f {
            case .integer(let n, _): if let v = try? content.value(Int.self, forProperty: n) { values[n] = .integer(v) }
            case .choice(let n, _, _), .text(let n, _): if let v = try? content.value(String.self, forProperty: n) { values[n] = .text(v) }
            }
        }
        return ModelReply(values)
    }

    // 欄位定義轉成 guided generation 的 schema：名稱、選項的型別名稱（欄位名＋Value）、欄位順序都和改成模型層之前相同，產生的內容才會一樣
    static func schema(_ request: ModelRequest) -> GenerationSchema {
        let properties = request.fields.map { f -> DynamicGenerationSchema.Property in
            switch f {
            case .choice(let n, let d, let options):
                return .init(name: n, description: d, schema: DynamicGenerationSchema(name: n + "Value", anyOf: options))
            case .text(let n, let d):
                return .init(name: n, description: d, schema: DynamicGenerationSchema(type: String.self))
            case .integer(let n, let d):
                return .init(name: n, description: d, schema: DynamicGenerationSchema(type: Int.self))
            }
        }
        return try! GenerationSchema(root: DynamicGenerationSchema(name: request.name, properties: properties), dependencies: [])
    }
}
