// Meeting minutes with on-device Apple Foundation Models, tuned for its small (~3B) model and 8K context:
// each request does one job on one chunk (proofread, or take notes), then a final request merges the
// chunk notes. Prints the same JSON shape haiku.py gets from Claude, so afm.py can render it identically.
//
//   afm <transcript.md> [glossary.txt]
import Foundation
import FoundationModels

struct Line {
    let number: Int
    let seconds: Int?
    let text: String
}

/// Measured: without this, one of two runs on the 24-minute talk came back entirely in English.
let languageRule = "\n所有輸出一律使用繁體中文（英文術語保持英文）。"

func stderr(_ s: String) {
    FileHandle.standardError.write((s + "\n").data(using: .utf8)!)
}

// MARK: - Schemas (DynamicGenerationSchema because the @Generable macro ships only with Xcode)

let string = DynamicGenerationSchema(type: String.self)
let int = DynamicGenerationSchema(type: Int.self)

func array(_ item: DynamicGenerationSchema, max: Int? = nil) -> DynamicGenerationSchema {
    DynamicGenerationSchema(arrayOf: item, maximumElements: max)
}

func object(_ name: String, _ props: [(String, String?, DynamicGenerationSchema)]) -> DynamicGenerationSchema {
    DynamicGenerationSchema(name: name, properties: props.map { .init(name: $0.0, description: $0.1, schema: $0.2) })
}

let actionItem = object("ActionItem", [
    ("task", nil, string),
    ("owner", "負責人；逐字稿沒提到就填「未指定」", string),
    ("due", "期限；逐字稿沒提到就填「未指定」", string),
])

let topic = object("Topic", [
    ("title", "主題名稱，不超過 15 字", string),
    ("points", "這個主題的重點，用自己的話濃縮，每項不超過 30 字", array(string, max: 5)),
])

let fixSchema = try! GenerationSchema(root: object("Fixes", [
    ("corrections", nil, array(object("Correction", [
        ("line", "行號", int),
        ("original", "該行裡原封不動出現的錯誤字串，越短越好", string),
        ("corrected", "正確的詞", string),
        ("reason", "依據，不超過 15 字", string),
    ]))),
    ("uncertain", "看不懂、或和上下文無關的片段", array(object("Uncertain", [
        ("line", "行號", int),
        ("text", "該行裡原封不動出現的字串", string),
    ]))),
]), dependencies: [])

let notesSchema = try! GenerationSchema(root: object("Notes", [
    ("topics", "這段依講述順序的主題，每個話題都要列到", array(topic, max: 6)),
    ("decisions", "會議中明確做出的決定；演講裡的建議或舉的案例不算；沒有就是空陣列", array(string)),
    ("actionItems", "會議中明確指派給某人的工作；演講裡的建議不算；沒有就是空陣列", array(actionItem)),
]), dependencies: [])

let summarySchema = try! GenerationSchema(root: object("Summary", [
    ("title", "一句話的主題，不超過 25 字", string),
    ("topics", "依講述順序、涵蓋開頭到結尾的主題；最後一段的內容也要列到", array(topic, max: 10)),
    ("decisions", "會議中明確做出的決定；沒有就是空陣列", array(string)),
    ("actionItems", "會議中明確指派的工作；沒有就是空陣列", array(actionItem)),
]), dependencies: [])

// MARK: - Prompts

func fixInstructions(_ glossary: [String]) -> String {
    """
    你是語音辨識逐字稿的校對員。使用者給你一段逐字稿，每行開頭是行號（例如 L12:）。
    找出辨識錯誤：同音或近音錯字、被寫成中文或拼錯的英文術語。
    - 只有在上下文能確定正確詞語時才列在 corrections；不確定就不要列
    - original 必須是該行裡原封不動出現的字串，只包含要改的詞
    - 看不懂或和上下文無關的片段列在 uncertain
    - 沒有錯誤就回傳空陣列
    例：「L3: 我們把服務低破到 staging」→ corrections: line 3、original「低破」、corrected「deploy」
    """ + languageRule + (glossary.isEmpty ? "" : "\n講者常用的詞彙，發音相近的詞很可能是這些：" + glossary.joined(separator: "、"))
}

let notesInstructions = """
你負責整理逐字稿的重點。使用者給你一段逐字稿，每行開頭是行號（例如 L12:）。
- 只能使用逐字稿裡的資訊，不要推論或補充
- 看不懂或明顯是辨識錯誤的句子直接略過
- 英文術語保持英文
""" + languageRule

let mergeInstructions = """
使用者會給你同一場錄音依順序分段整理的重點。請合併成最終的會議記錄：
- 主題依講述順序排列，必須涵蓋所有段落，合併重複的主題
- 只能使用這些重點裡的資訊，不要推論或補充
- 標題是一句話，不超過 25 字，不要把所有主題名稱串起來
""" + languageRule

// MARK: - Glossary candidates (code finds look-alikes, AFM only confirms)

/// Small models are poor at finding errors on their own but fine at a yes/no question about one candidate.
let verifyInstructions = """
你協助檢查語音辨識的錯誤。使用者會給你一段逐字稿，以及其中一個詞的兩種寫法 A 和 B（A 是辨識結果，B 是發音相近的專有名詞）。
先用一句話分析哪一種放進句子裡意思通順、符合上下文，再回答講者實際說的是 A 還是 B。
"""

struct Candidate {
    let line: Int
    let original: String
    let corrected: String
}

func editDistance(_ a: [Character], _ b: [Character]) -> Int {
    var prev = Array(0...b.count)
    for (i, ca) in a.enumerated() {
        var cur = [i + 1]
        for (j, cb) in b.enumerated() {
            cur.append(min(prev[j + 1] + 1, cur[j] + 1, prev[j] + (ca == cb ? 0 : 1)))
        }
        prev = cur
    }
    return prev[b.count]
}

func letters(_ s: String) -> [Character] {
    Array(s.lowercased().filter { $0.isLetter && $0.isASCII })
}

/// Latin spans of one or two words, e.g. "radius", "Sir manager".
func latinSpans(_ text: String) -> [String] {
    let words = text.matches(of: #/[A-Za-z][A-Za-z\-]*/#).map { String($0.output) }
    var spans = words
    for i in words.indices.dropLast() { spans.append(words[i] + " " + words[i + 1]) }
    return spans.filter { text.contains($0) }
}

func candidates(_ lines: [Line], _ glossary: [String]) -> (caseFixes: [Candidate], toVerify: [Candidate]) {
    let latin = glossary.filter { $0.contains(where: { $0.isASCII && $0.isLetter }) }
    let chinese = glossary.filter { !$0.contains(where: { $0.isASCII }) && $0.count >= 2 }
    var caseFixes: [Candidate] = []
    var toVerify: [Candidate] = []
    for line in lines {
        for span in Set(latinSpans(line.text)) {
            let a = letters(span)
            guard a.count >= 3 else { continue }
            if let exact = latin.first(where: { letters($0) == a }) {
                if exact != span, !latin.contains(span) { caseFixes.append(Candidate(line: line.number, original: span, corrected: exact)) }
                continue
            }
            // Measured on script-01: radius→Redis is 0.67, "Sir manager"→cert-manager 0.73, HyPath→Hyperf 0.5.
            let best = latin.map { ($0, 1 - Double(editDistance(a, letters($0))) / Double(max(a.count, letters($0).count))) }
                .max { $0.1 < $1.1 }
            if let (term, similarity) = best, similarity >= 0.5 {
                toVerify.append(Candidate(line: line.number, original: span, corrected: term))
            }
        }
        // Chinese homophones: same length, at least one shared character, not already the term (e.g. 供單→工單).
        let chars = Array(line.text)
        for term in chinese where !line.text.contains(term) {
            let n = term.count
            guard chars.count >= n else { continue }
            for start in 0...(chars.count - n) {
                let window = chars[start..<start + n]
                guard window.allSatisfy({ !$0.isASCII && !$0.isPunctuation && !$0.isWhitespace }) else { continue }
                let shared = zip(window, term).filter { $0 == $1 }.count
                if shared >= 1, shared < n {
                    toVerify.append(Candidate(line: line.number, original: String(window), corrected: term))
                }
            }
        }
    }
    // Keep only the best candidate per (line, original).
    var seen = Set<String>()
    toVerify = toVerify.filter { seen.insert("\($0.line)|\($0.original)").inserted }
    return (caseFixes, toVerify)
}

/// Small models are poor at finding errors on their own, and a yes/no question made AFM answer "no" to every
/// candidate (22/22 on script-01). A forced choice between the two spellings, with a one-line analysis
/// generated first, is a fairer question.
func verify(_ c: Candidate, _ lines: [Line]) async -> Bool {
    let context = lines.filter { abs($0.number - c.line) <= 2 }.map(\.text).joined(separator: "\n")
    let prompt = "逐字稿：\n\(context)\n\nA：「\(c.original)」\nB：「\(c.corrected)」"
    let schema = try! GenerationSchema(root: object("Verdict", [
        ("analysis", "一句話分析，不超過 40 字", string),
        ("answer", "講者實際說的是 A 還是 B", DynamicGenerationSchema(name: "Choice", anyOf: ["A", "B"])),
    ]), dependencies: [])
    guard case .ok(let content) = await ask(verifyInstructions, prompt, schema) else { return false }
    let verdict = (try? content.value(String.self, forProperty: "answer")) == "B"
    if ProcessInfo.processInfo.environment["AFM_DEBUG"] != nil {
        stderr("  候選 L\(c.line)「\(c.original)」→「\(c.corrected)」：\(verdict)")
    }
    return verdict
}

// MARK: - Kind (a separate, focused question; inside the merge AFM called a talk a meeting on both runs)

let kindInstructions = """
使用者會給你一段錄音逐字稿的開頭。判斷這是哪一種錄音：
- meeting：多個人輪流發言、互相討論、分配工作
- talk：一個人對聽眾講話，例如演講、課程、分享會
- memo：一個人自己口述筆記
先用一句話分析，再回答。
"""

func kindOf(_ lines: [Line]) async -> String {
    var sample = ""
    for line in lines where sample.count < 1500 { sample += line.text + "\n" }
    let schema = try! GenerationSchema(root: object("Kind", [
        ("analysis", "一句話分析，不超過 40 字", string),
        ("kind", nil, DynamicGenerationSchema(name: "KindChoice", anyOf: ["meeting", "talk", "memo"])),
    ]), dependencies: [])
    guard case .ok(let content) = await ask(kindInstructions, sample, schema) else { return "other" }
    return (try? content.value(String.self, forProperty: "kind")) ?? "other"
}

/// Paragraphs from pauses instead of from AFM, which marked a new paragraph every one or two lines.
func pauseParagraphs(_ lines: [Line]) -> [Int] {
    var starts = [1]
    var lastStart = 0
    for (i, line) in lines.enumerated().dropFirst() {
        let gap = (line.seconds ?? 0) - (lines[i - 1].seconds ?? 0)
        if (gap >= 4 && i - lastStart >= 4) || i - lastStart >= 12 {
            starts.append(line.number)
            lastStart = i
        }
    }
    return starts
}

// MARK: - Requests

let model = SystemLanguageModel.default
let options = GenerationOptions(temperature: 0.1)

enum Outcome {
    case ok(GeneratedContent)
    case tooLong
    case failed(String)
}

func ask(_ instructions: String, _ prompt: String, _ schema: GenerationSchema) async -> Outcome {
    let session = LanguageModelSession(instructions: instructions)
    do {
        return .ok(try await session.respond(to: prompt, schema: schema, options: options).content)
    } catch LanguageModelSession.GenerationError.exceededContextWindowSize {
        return .tooLong
    } catch {
        return .failed("\(error)")
    }
}

/// Runs one request per chunk; a chunk that overflows the context window is split in half and retried.
func perChunk(_ chunk: [Line], _ instructions: String, _ schema: GenerationSchema, _ label: String) async -> [GeneratedContent] {
    let prompt = chunk.map { "L\($0.number): \($0.text)" }.joined(separator: "\n")
    switch await ask(instructions, prompt, schema) {
    case .ok(let content):
        return [content]
    case .tooLong where chunk.count > 1:
        stderr("  \(label) L\(chunk[0].number)：超過 context，切半重試")
        let half = chunk.count / 2
        return await perChunk(Array(chunk[..<half]), instructions, schema, label)
            + perChunk(Array(chunk[half...]), instructions, schema, label)
    case .tooLong:
        stderr("  \(label) L\(chunk[0].number)：單行就超過 context，略過")
        return []
    case .failed(let error):
        stderr("  \(label) L\(chunk[0].number)：失敗，略過（\(error)）")
        return []
    }
}

/// Packs lines into chunks of at most `budget` tokens, measured with the model's own tokenizer.
func chunks(_ lines: [Line], budget: Int) async throws -> [[Line]] {
    var result: [[Line]] = []
    var current: [Line] = []
    var used = 0
    for line in lines {
        let cost = try await model.tokenCount(for: "L\(line.number): \(line.text)\n")
        if !current.isEmpty, used + cost > budget {
            result.append(current)
            current = []
            used = 0
        }
        current.append(line)
        used += cost
    }
    if !current.isEmpty { result.append(current) }
    return result
}

// MARK: - Decoding into haiku.py's JSON shape

func strings(_ c: GeneratedContent, _ key: String) throws -> [String] {
    try c.value([String].self, forProperty: key)
}

func items(_ c: GeneratedContent, _ key: String) throws -> [GeneratedContent] {
    try c.value([GeneratedContent].self, forProperty: key)
}

func topics(_ c: GeneratedContent) throws -> [[String: Any]] {
    try items(c, "topics").map { ["title": try $0.value(String.self, forProperty: "title"), "points": try strings($0, "points")] }
}

func actionItems(_ c: GeneratedContent) throws -> [[String: Any]] {
    try items(c, "actionItems").map {
        ["task": try $0.value(String.self, forProperty: "task"),
         "owner": try $0.value(String.self, forProperty: "owner"),
         "due": try $0.value(String.self, forProperty: "due")]
    }
}

@main
struct AFM {
    static func main() async {
        let args = CommandLine.arguments
        guard args.count >= 2 else {
            stderr("用法: afm <逐字稿.md> [glossary.txt]")
            exit(1)
        }
        do {
            try await run(path: args[1], glossaryPath: args.count > 2 ? args[2] : nil)
        } catch {
            stderr("錯誤：\(error)")
            exit(1)
        }
    }

    static func run(path: String, glossaryPath: String?) async throws {
        if case .unavailable(let reason) = model.availability {
            throw NSError(domain: "afm", code: 1, userInfo: [NSLocalizedDescriptionKey: "AFM 不可用：\(reason)"])
        }
        let raw = try String(contentsOfFile: path, encoding: .utf8).split(separator: "\n").map(String.init)
        let stamp = #/^\[\d+:\d+\]\s*/#
        let timed = raw.filter { $0.contains(stamp) }
        let lines = (timed.isEmpty ? raw.filter { !$0.trimmingCharacters(in: .whitespaces).isEmpty } : timed)
            .enumerated().map { item -> Line in
                let m = item.element.firstMatch(of: #/^\[(\d+):(\d+)\]/#)
                var seconds: Int?
                if let m, let mm = Int(m.1), let ss = Int(m.2) { seconds = mm * 60 + ss }
                return Line(number: item.offset + 1, seconds: seconds, text: item.element.replacing(stamp, with: ""))
            }
        let glossary = try glossaryPath.map {
            try String(contentsOfFile: $0, encoding: .utf8).split(separator: "\n").map(String.init)
                .filter { !$0.isEmpty && !$0.hasPrefix("#") }
        } ?? []

        // Input budget: the 8K window must also hold the instructions (~400 tokens) and the response,
        // which for proofreading or notes is at most about the size of the input.
        // 3/8 overflowed by one token on the 24-minute talk (8193 > 8192), so leave more headroom.
        let budget = model.contextSize / 4
        let parts = try await chunks(lines, budget: budget)
        stderr("AFM：\(lines.count) 行，context \(model.contextSize) tokens，切成 \(parts.count) 段（每段 ≤ \(budget) tokens）")
        let clock = ContinuousClock.now

        var corrections: [[String: Any]] = []
        var uncertain: [[String: Any]] = []
        var notes: [String] = []
        for (i, part) in parts.enumerated() {
            stderr("AFM：第 \(i + 1)/\(parts.count) 段（L\(part[0].number)–L\(part.last!.number)）")
            for fix in await perChunk(part, fixInstructions(glossary), fixSchema, "校對") {
                for c in try items(fix, "corrections") {
                    corrections.append([
                        "line": try c.value(Int.self, forProperty: "line"),
                        "original": try c.value(String.self, forProperty: "original"),
                        "corrected": try c.value(String.self, forProperty: "corrected"),
                        "reason": try c.value(String.self, forProperty: "reason"),
                    ])
                }
                for u in try items(fix, "uncertain") {
                    uncertain.append(["line": try u.value(Int.self, forProperty: "line"), "text": try u.value(String.self, forProperty: "text")])
                }
            }
            for note in await perChunk(part, notesInstructions, notesSchema, "整理") {
                var text = try topics(note).map { t in
                    "主題：\(t["title"]!)\n" + (t["points"] as! [String]).map { "- \($0)" }.joined(separator: "\n")
                }.joined(separator: "\n")
                let decisions = try strings(note, "decisions")
                if !decisions.isEmpty { text += "\n決定：\n" + decisions.map { "- \($0)" }.joined(separator: "\n") }
                let actions = try actionItems(note)
                if !actions.isEmpty {
                    text += "\n指派：\n" + actions.map { "- \($0["task"]!)（\($0["owner"]!)，\($0["due"]!)）" }.joined(separator: "\n")
                }
                notes.append("【第 \(notes.count + 1) 段】\n" + text)
            }
        }

        let found = candidates(lines, glossary)
        stderr("AFM：大小寫修正 \(found.caseFixes.count) 筆")
        for c in found.caseFixes {
            corrections.append(["line": c.line, "original": c.original, "corrected": c.corrected, "reason": "詞彙表大小寫"])
        }
        // Not used: on script-01 AFM confirmed 4 of 22 look-alike candidates, all wrong (e.g. 效投注→總投注),
        // and rejected every real error (radius→Redis, Sir manager→cert-manager, 供單→工單). Kept for re-testing
        // when the OS model changes: AFM_VERIFY=1.
        if ProcessInfo.processInfo.environment["AFM_VERIFY"] != nil {
            for c in found.toVerify where await verify(c, lines) {
                corrections.append(["line": c.line, "original": c.original, "corrected": c.corrected, "reason": "與詞彙表發音相近"])
            }
        }

        stderr("AFM：合併 \(notes.count) 段重點")
        var summary: [String: Any] = ["kind": "other", "title": "（AFM 合併失敗）", "topics": [], "decisions": [], "action_items": []]
        var kind = await kindOf(lines)
        // AFM called all three recorded talks "memo". Talks keep addressing the audience (大家/各位: 19, 24 and
        // 40 times in the three talks, 2 in the script-01 stand-up), a personal memo doesn't.
        let audience = lines.reduce(0) { $0 + $1.text.matches(of: #/大家|各位/#).count }
        if kind == "memo", audience >= 5 { kind = "talk" }
        stderr("AFM：類型 \(kind)（大家/各位 \(audience) 次）")
        if case .ok(let merged) = await ask(mergeInstructions, notes.joined(separator: "\n\n"), summarySchema) {
            summary = [
                "kind": kind,
                "title": try merged.value(String.self, forProperty: "title"),
                "topics": try topics(merged),
                // A small model still lists a talk's advice or anecdotes as decisions; enforce the rule in code.
                "decisions": kind == "talk" ? [] : try strings(merged, "decisions"),
                "action_items": kind == "talk" ? [] : try actionItems(merged),
            ]
        }

        let elapsed = ContinuousClock.now - clock
        stderr("AFM：完成，耗時 \(elapsed.components.seconds) 秒，修正 \(corrections.count) 筆、不確定 \(uncertain.count) 筆")
        let result: [String: Any] = [
            "corrections": corrections,
            "uncertain": uncertain,
            "paragraphs": pauseParagraphs(lines),
            "summary": summary,
            "elapsed_seconds": Double(elapsed.components.seconds),
        ]
        let data = try JSONSerialization.data(withJSONObject: result, options: [.prettyPrinted, .sortedKeys])
        print(String(data: data, encoding: .utf8)!)
    }
}
