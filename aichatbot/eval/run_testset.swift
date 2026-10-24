// 用測試題庫（eval/testset.json）評估模型，自動評分並把結果存到 eval/results/
// 執行：swift eval/run_testset.swift                                  （本機 Apple 模型）
//       MODEL=haiku ANTHROPIC_API_KEY=[REDACTED] swift eval/run_testset.swift  （Claude Haiku 4.5，會產生 API 費用）
//       ROUNDS=3 swift eval/run_testset.swift                          （每題跑幾輪，預設 2）
//       MODEL=replay ANSWERS=eval/results/xxx.jsonl swift eval/run_testset.swift  （只評分已存好的回答，例如用訂閱的 Haiku 子代理跑出來的）
// 店家問題的問法與 chat.swift 的 A 版本相同：系統提示 = shop.txt，訊息前加【回覆語言】，輸出 answer + status
import Foundation
import FoundationModels
import NaturalLanguage

let env = ProcessInfo.processInfo.environment
let modelName = env["MODEL"] ?? "local"
let rounds = modelName == "replay" ? 1 : (Int(env["ROUNDS"] ?? "") ?? 2)
let dir = URL(fileURLWithPath: #filePath).deletingLastPathComponent()
let shop = (try? String(contentsOf: dir.appendingPathComponent("shop.txt"), encoding: .utf8)) ?? ""
let chatSource = (try? String(contentsOf: dir.appendingPathComponent("chat.swift"), encoding: .utf8)) ?? ""

// MARK: - 題庫

struct TestCase: Decodable {
    let id, cat, q, intent, expect: String
    let all, any, none: [String]?
    let lang: String?
}
struct TestSet: Decodable { let cases: [TestCase] }
let cases = try! JSONDecoder().decode(TestSet.self, from: Data(contentsOf: dir.appendingPathComponent("testset.json"))).cases

// MARK: - 語言

func detectLanguage(_ text: String) -> String {
    // 先看文字系統：有假名就是日文、有漢字就是中文
    // （「可以用LINE Pay嗎」這種夾英文品牌名的句子，語言辨識常誤判成英文）
    let scalars = text.unicodeScalars
    if scalars.contains(where: { (0x3040...0x30FF).contains($0.value) }) { return "ja" }
    if scalars.contains(where: { $0.properties.isIdeographic }) { return "zh" }
    let r = NLLanguageRecognizer()
    r.languageConstraints = [.traditionalChinese, .simplifiedChinese, .english, .japanese]
    r.processString(text)
    guard let best = r.dominantLanguage else { return "zh" }
    return best.rawValue.hasPrefix("zh") ? "zh" : String(best.rawValue.prefix(2))
}
let languageNames = ["zh": "Traditional Chinese", "en": "English", "ja": "Japanese"]
func buildPrompt(_ q: String) -> String {
    "【回覆語言】請用 \(languageNames[detectLanguage(q)] ?? "Traditional Chinese") 回答。\n\n顧客訊息：\(q)"
}

// MARK: - 意圖類別（從 chat.swift 讀，保持一致）

let intentRegex = try! NSRegularExpression(pattern: #"\("([a-z_]+)", "([^"]+)"\),"#)
let intents: [(String, String)] = intentRegex.matches(in: chatSource, range: NSRange(chatSource.startIndex..., in: chatSource)).compactMap { m in
    guard let a = Range(m.range(at: 1), in: chatSource), let b = Range(m.range(at: 2), in: chatSource) else { return nil }
    return (String(chatSource[a]), String(chatSource[b]))
}
let intentDescription = intents.map { "\($0.0)=\($0.1)" }.joined(separator: "；")

let answerDescription = "依店家資料回答顧客（需要時套用折扣、付款門檻等規則計算）；簡潔、不要用「你好」開頭、不要提到資料裡沒有的人名"
let statusDescription = "answered=上面的回答有店家資料根據；not_in_data=店家資料完全沒有相關資訊；unclear=看不懂顧客在問什麼（太短、沒頭沒尾、指代不明如「他是誰」）"
let statuses = ["answered", "not_in_data", "unclear"]

// MARK: - 模型

struct Answer { var status: String; var answer: String; var inputTokens = 0; var outputTokens = 0 }

protocol Responder {
    func answer(_ q: String) async throws -> Answer
    func classify(_ q: String) async throws -> String
}

struct LocalModel: Responder {
    let answerSchema = try! GenerationSchema(root: DynamicGenerationSchema(name: "Answer", properties: [
        .init(name: "answer", description: answerDescription, schema: DynamicGenerationSchema(type: String.self)),
        .init(name: "status", description: statusDescription, schema: DynamicGenerationSchema(name: "AnswerStatus", anyOf: statuses)),
    ]), dependencies: [])
    let intentSchema = try! GenerationSchema(root: DynamicGenerationSchema(name: "Intent", properties: [
        .init(name: "intent", description: intentDescription,
              schema: DynamicGenerationSchema(name: "IntentValue", anyOf: intents.map(\.0)))
    ]), dependencies: [])

    func answer(_ q: String) async throws -> Answer {
        let c = try await LanguageModelSession(instructions: shop)
            .respond(to: buildPrompt(q), schema: answerSchema, options: GenerationOptions(maximumResponseTokens: 200)).content
        return Answer(status: (try? c.value(String.self, forProperty: "status")) ?? "?",
                      answer: (try? c.value(String.self, forProperty: "answer")) ?? "")
    }
    func classify(_ q: String) async throws -> String {
        let c = try await LanguageModelSession(instructions: "判斷顧客訊息的意圖。顧客可能使用任何語言。").respond(to: q, schema: intentSchema).content
        return (try? c.value(String.self, forProperty: "intent")) ?? "?"
    }
}

struct HaikuModel: Responder {
    let apiKey: [REDACTED]
    let model = "claude-haiku-4-5"

    func call(system: String, user: String, schema: [String: Any]) async throws -> ([String: Any], Int, Int) {
        var req = URLRequest(url: URL(string: "https://example.invalid/v1/messages")!)
        req.httpMethod = "POST"
        req.setValue(apiKey, forHTTPHeaderField: "x-api-key")
        req.setValue("2023-06-01", forHTTPHeaderField: "anthropic-version")
        req.setValue("application/json", forHTTPHeaderField: "content-type")
        req.httpBody = try JSONSerialization.data(withJSONObject: [
            "model": model, "max_tokens": 512, "system": system,
            "messages": [["role": "user", "content": user]],
            "output_config": ["format": ["type": "json_schema", "schema": schema]],
        ])
        let (data, resp) = try await URLSession.shared.data(for: req)
        guard (resp as? HTTPURLResponse)?.statusCode == 200,
              let json = try JSONSerialization.jsonObject(with: data) as? [String: Any] else {
            throw NSError(domain: "haiku", code: 1, userInfo: [NSLocalizedDescriptionKey: String(data: data, encoding: .utf8) ?? "?"])
        }
        let text = ((json["content"] as? [[String: Any]])?.first { $0["type"] as? String == "text" }?["text"] as? String) ?? "{}"
        let usage = json["usage"] as? [String: Any]
        let parsed = (try? JSONSerialization.jsonObject(with: Data(text.utf8))) as? [String: Any] ?? [:]
        return (parsed, usage?["input_tokens"] as? Int ?? 0, usage?["output_tokens"] as? Int ?? 0)
    }

    func answer(_ q: String) async throws -> Answer {
        let (r, i, o) = try await call(system: shop, user: buildPrompt(q), schema: [
            "type": "object",
            "properties": ["answer": ["type": "string", "description": answerDescription],
                           "status": ["type": "string", "enum": statuses, "description": statusDescription]],
            "required": ["answer", "status"], "additionalProperties": false,
        ])
        return Answer(status: r["status"] as? String ?? "?", answer: r["answer"] as? String ?? "", inputTokens: i, outputTokens: o)
    }
    func classify(_ q: String) async throws -> String {
        let (r, _, _) = try await call(system: "判斷顧客訊息的意圖。顧客可能使用任何語言。", user: q, schema: [
            "type": "object",
            "properties": ["intent": ["type": "string", "enum": intents.map(\.0), "description": intentDescription]],
            "required": ["intent"], "additionalProperties": false,
        ])
        return r["intent"] as? String ?? "?"
    }
}

// 讀取已存好的回答（每行 {"id","answer","status"}），只做評分，不呼叫模型
struct ReplayModel: Responder {
    let answers: [String: Answer]
    init(path: String) {
        var map: [String: Answer] = [:]
        for line in ((try? String(contentsOfFile: path, encoding: .utf8)) ?? "").split(separator: "\n") {
            guard let obj = try? JSONSerialization.jsonObject(with: Data(line.utf8)) as? [String: Any],
                  let id = obj["id"] as? String else { continue }
            map[id] = Answer(status: obj["status"] as? String ?? "?", answer: obj["answer"] as? String ?? "")
        }
        answers = map
    }
    func answer(_ q: String) async throws -> Answer {
        guard let c = cases.first(where: { $0.q == q }), let a = answers[c.id] else { throw NSError(domain: "replay", code: 1, userInfo: [NSLocalizedDescriptionKey: "缺少這題的回答"]) }
        return a
    }
    func classify(_ q: String) async throws -> String { throw NSError(domain: "replay", code: 2) }
}

let responder: Responder
let displayName: String
if modelName == "replay" {
    guard let path = env["ANSWERS"] else { print("❌ MODEL=replay 需要 ANSWERS=檔案路徑"); exit(1) }
    responder = ReplayModel(path: path); displayName = "重播：\(URL(fileURLWithPath: path).lastPathComponent)"
} else if modelName == "haiku" {
    guard let key = env["ANTHROPIC_API_KEY"], !key.isEmpty else { print("❌ MODEL=haiku 需要 ANTHROPIC_API_KEY"); exit(1) }
    responder = HaikuModel(apiKey: [REDACTED]; displayName = "Claude Haiku 4.5"
} else {
    responder = LocalModel(); displayName = "本機 \(SystemLanguageModel.default.variant.displayName)"
}

// MARK: - 評分

func judge(_ c: TestCase, _ a: Answer) -> (Bool, String) {
    let answered = a.status == "answered" && !a.answer.trimmingCharacters(in: .whitespaces).isEmpty
    func checks() -> (Bool, String) {
        if let all = c.all, let miss = all.first(where: { !a.answer.contains($0) }) { return (false, "缺少「\(miss)」") }
        if let any = c.any, !any.isEmpty, !any.contains(where: a.answer.contains) { return (false, "沒有提到 \(any.joined(separator: "/"))") }
        if let none = c.none, let bad = none.first(where: a.answer.contains) { return (false, "不該出現「\(bad)」") }
        if let lang = c.lang, detectLanguage(a.answer) != lang { return (false, "語言應為 \(lang)，實際是 \(detectLanguage(a.answer))") }
        return (true, "")
    }
    switch c.expect {
    case "answer": return answered ? checks() : (false, "應該回答，卻判成 \(a.status)")
    case "no_answer": return answered ? (false, "應該不回答，卻回答了") : (true, "")
    default: return answered ? checks() : (true, "")
    }
}

// MARK: - 執行

struct Record: Codable { let id, cat, q, status, answer, reason: String; let pass: Bool; let round: Int }
var records: [Record] = []
var intentOK = 0, intentTotal = 0
var inTokens = 0, outTokens = 0, errors = 0
var times: [Double] = []

print("模型：\(displayName)　題數：\(cases.count)　每題 \(rounds) 輪\n")
for round in 1...rounds {
    for c in cases {
        let t = Date()
        do {
            let a = try await responder.answer(c.q)
            times.append(Date().timeIntervalSince(t))
            inTokens += a.inputTokens; outTokens += a.outputTokens
            let (pass, reason) = judge(c, a)
            records.append(Record(id: c.id, cat: c.cat, q: c.q, status: a.status,
                                  answer: a.answer.replacingOccurrences(of: "\n", with: " "), reason: reason, pass: pass, round: round))
        } catch {
            errors += 1
            records.append(Record(id: c.id, cat: c.cat, q: c.q, status: "error", answer: "", reason: String(describing: error).prefix(120).description, pass: false, round: round))
        }
    }
    print("第 \(round) 輪完成")
}
for c in cases where c.intent != "*" && modelName != "replay" {
    intentTotal += 1
    if (try? await responder.classify(c.q)) == c.intent { intentOK += 1 }
}

// MARK: - 報告

let categories = cases.map(\.cat).reduce(into: [String]()) { if !$0.contains($1) { $0.append($1) } }
print("\n========== 各類別 ==========")
for cat in categories {
    let rs = records.filter { $0.cat == cat }
    print(String(format: "%@　%d/%d", cat.padding(toLength: 8, withPad: "　", startingAt: 0), rs.filter(\.pass).count, rs.count))
}
let passed = records.filter(\.pass).count
print("\n========== 總結 ==========")
print(String(format: "店家問題：%d/%d（%.0f%%）", passed, records.count, Double(passed) / Double(records.count) * 100))
print("意圖分類：\(intentOK)/\(intentTotal)")
let unstable = Set(records.filter(\.pass).map(\.id)).intersection(records.filter { !$0.pass }.map(\.id))
print("不穩定（有幾輪對、有幾輪錯）：\(unstable.count) 題")
print(String(format: "平均每次回答：%.2f 秒　錯誤：%d 次", times.isEmpty ? 0 : times.reduce(0, +) / Double(times.count), errors))
if inTokens > 0 {
    let cost = Double(inTokens) / 1_000_000 * 1 + Double(outTokens) / 1_000_000 * 5
    print(String(format: "token：輸入 %d、輸出 %d　本次費用約 US$%.4f　平均每題 US$%.5f", inTokens, outTokens, cost, cost / Double(records.count)))
}

print("\n========== 答錯的題目（第 1 輪）==========")
for r in records where !r.pass && r.round == 1 {
    print("❌ [\(r.id)] \(r.q)\n   → [\(r.status)] \(r.answer.isEmpty ? "（無）" : r.answer)\n   原因：\(r.reason)")
}

let resultsDir = dir.appendingPathComponent("results")
try? FileManager.default.createDirectory(at: resultsDir, withIntermediateDirectories: true)
let outURL = resultsDir.appendingPathComponent("\(modelName)-\(Int(Date().timeIntervalSince1970)).json")
let encoder = JSONEncoder(); encoder.outputFormatting = [.prettyPrinted, .withoutEscapingSlashes]
try? encoder.encode(records).write(to: outURL)
print("\n詳細結果：\(outURL.path)")
