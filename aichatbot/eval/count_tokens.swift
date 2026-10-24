// 用 Anthropic 的 token 計數 API 量出：如果改用 Claude Haiku 4.5，每則訊息 / 每段對話 / 1 萬段對話要多少錢
// 執行：ANTHROPIC_API_KEY=[REDACTED] swift eval/count_tokens.swift
// token 計數 API 免費（有速率限制），不會產生模型費用
import Foundation

let model = "claude-haiku-4-5"
let inputPrice = 1.0 / 1_000_000   // US$ / token
let outputPrice = 5.0 / 1_000_000
let usdToTwd = 32.0                 // 粗估匯率

guard let apiKey = [REDACTED]"ANTHROPIC_API_KEY"], !apiKey.isEmpty else {
    print("❌ 請先設定環境變數 ANTHROPIC_API_KEY"); exit(1)
}

let dir = URL(fileURLWithPath: #filePath).deletingLastPathComponent()
func read(_ name: String) -> String { (try? String(contentsOf: dir.appendingPathComponent(name), encoding: .utf8)) ?? "" }
let shop = read("shop.txt")
let faqText = read("faq.txt")
let chatSource = read("chat.swift")

// 從 chat.swift 取出意圖分類的說明，跟測試工具保持一致
let intentRegex = try! NSRegularExpression(pattern: #"\("([a-z_]+)", "([^"]+)"\),"#)
let intentLines = intentRegex.matches(in: chatSource, range: NSRange(chatSource.startIndex..., in: chatSource)).compactMap { m -> String? in
    guard let a = Range(m.range(at: 1), in: chatSource), let b = Range(m.range(at: 2), in: chatSource) else { return nil }
    return "\(chatSource[a])=\(chatSource[b])"
}
let classifySystem = "判斷顧客訊息的意圖，只輸出一個代號。顧客可能使用任何語言。\n" + intentLines.joined(separator: "\n")

let faqList = faqText.split(separator: "\n").filter { !$0.hasPrefix("#") && $0.contains(" | ") }.map { line -> String in
    let cols = line.components(separatedBy: " | ")
    return "\(cols[0])：\(cols[1])"
}.joined(separator: "\n")
let faqSystem = "從 FAQ 清單挑出能回答顧客問題的一題，只輸出代號；沒有就輸出 none。\nFAQ 清單：\n" + faqList

struct Msg: Codable { let role: String; let content: String }

func countTokens(system: String?, messages: [Msg]) async throws -> Int {
    var req = URLRequest(url: URL(string: "https://example.invalid/v1/messages/count_tokens")!)
    req.httpMethod = "POST"
    req.setValue(apiKey, forHTTPHeaderField: "x-api-key")
    req.setValue("2023-06-01", forHTTPHeaderField: "anthropic-version")
    req.setValue("application/json", forHTTPHeaderField: "content-type")
    var body: [String: Any] = ["model": model, "messages": messages.map { ["role": $0.role, "content": $0.content] }]
    if let system { body["system"] = system }
    req.httpBody = try JSONSerialization.data(withJSONObject: body)
    let (data, resp) = try await URLSession.shared.data(for: req)
    guard (resp as? HTTPURLResponse)?.statusCode == 200,
          let json = try JSONSerialization.jsonObject(with: data) as? [String: Any],
          let n = json["input_tokens"] as? Int else {
        throw NSError(domain: "count_tokens", code: 1, userInfo: [NSLocalizedDescriptionKey: String(data: data, encoding: .utf8) ?? "?"])
    }
    return n
}

// 顧客訊息與客服回覆樣本（取自測試對話）
let questions = ["抹茶拿鐵多少錢", "你們幾點開？禮拜一有開嗎", "100塊可以刷卡嗎", "包廠要提前多久阿 最少要花多少",
                 "我想自取買三杯紅玉，總共多少？", "有停車位嗎？", "Do you deliver?", "抹茶ラテはいくらですか",
                 "我想預約 10/20 下午兩點", "有幾個老闆"]
let replies = ["抹茶拿鐵一杯 150 元～", "我們週二至週日 11:00–20:00 營業，週一公休喔！",
               "100 元不適用信用卡，信用卡需滿 200 元才可使用，可以用現金或 LINE Pay 喔！",
               "包場需提前 7 天預約，低消 3000 元，最多可容納 20 人～",
               "三杯日月潭紅玉紅茶共 360 元，店內自取滿 300 元打 9 折，總共 324 元喔！",
               "這個問題我幫您轉給店員確認喔！", "Yes! You can find us on both Uber Eats and foodpanda~",
               "抹茶ラテは一杯 150元です～", "請問幾位呢？", "這個問題我幫您轉給店員確認喔！"]

func avg(_ xs: [Int]) -> Double { Double(xs.reduce(0, +)) / Double(xs.count) }
func usd(_ x: Double) -> String { String(format: "$%.5f", x) }
func money(_ x: Double) -> String { String(format: "US$%.2f（約台幣 %.0f 元）", x, x * usdToTwd) }

do {
    print("模型：\(model)　（輸入 $1 / 輸出 $5 每百萬 token）\n")

    // 1. 回答店家問題（A：系統提示含完整店家資料）
    var answerIn: [Int] = []
    for q in questions { answerIn.append(try await countTokens(system: shop, messages: [Msg(role: "user", content: q)])) }

    // 2. 意圖分類
    var classifyIn: [Int] = []
    for q in questions { classifyIn.append(try await countTokens(system: classifySystem, messages: [Msg(role: "user", content: q)])) }

    // 3. 挑 FAQ（B）
    var faqIn: [Int] = []
    for q in questions { faqIn.append(try await countTokens(system: faqSystem, messages: [Msg(role: "user", content: q)])) }

    // 4. 回覆長度（輸出 token）：把回覆當內容計數，扣掉空訊息的基本開銷
    let base = try await countTokens(system: nil, messages: [Msg(role: "user", content: "。")])
    var replyOut: [Int] = []
    for r in replies { replyOut.append(max(1, try await countTokens(system: nil, messages: [Msg(role: "user", content: r)]) - base + 1)) }

    // 5. 多輪對話：第 5 輪時要連同前 4 輪一起送
    var history: [Msg] = []
    var turnIn: [Int] = []
    for i in 0..<5 {
        history.append(Msg(role: "user", content: questions[i]))
        turnIn.append(try await countTokens(system: shop, messages: history))
        history.append(Msg(role: "assistant", content: replies[i]))
    }

    let aIn = avg(answerIn), cIn = avg(classifyIn), fIn = avg(faqIn), out = avg(replyOut)
    print("【實測 token 數】")
    print(String(format: "  系統提示（店家資料）＋一句問題：平均 %.0f token", aIn))
    print(String(format: "  意圖分類提示＋一句問題：　　　　平均 %.0f token", cIn))
    print(String(format: "  FAQ 挑選提示＋一句問題：　　　　平均 %.0f token", fIn))
    print(String(format: "  一則回覆（輸出）：　　　　　　　平均 %.0f token", out))
    print("  多輪對話每一輪的輸入：\(turnIn.map(String.init).joined(separator: " → ")) token\n")

    let answerCost = aIn * inputPrice + out * outputPrice
    let classifyCost = cIn * inputPrice + 5 * outputPrice     // 分類只輸出一個代號
    let faqCost = fIn * inputPrice + 5 * outputPrice
    let convCost = Double(turnIn.reduce(0, +)) * inputPrice + out * 5 * outputPrice   // 5 輪，帶對話紀錄

    print("【每則訊息】")
    print("  A 直接回答：　　　　　\(usd(answerCost))")
    print("  先分類再回答：　　　　\(usd(classifyCost + answerCost))")
    print("  B 挑 FAQ（回覆是固定答案，只付挑題的錢）：\(usd(faqCost))")
    print("\n【一段 5 輪的對話（帶對話紀錄）】 \(usd(convCost))")
    print("\n【1 萬則訊息】")
    print("  A 直接回答：　\(money(answerCost * 10_000))")
    print("  先分類再回答：\(money((classifyCost + answerCost) * 10_000))")
    print("  B 挑 FAQ：　　\(money(faqCost * 10_000))")
    print("\n【1 萬段 5 輪對話】")
    print("  A 直接回答（帶對話紀錄）：\(money(convCost * 10_000))")
    print("  先分類再回答：　　　　　　\(money((convCost + classifyCost * 5) * 10_000))")
    print("  B 挑 FAQ：　　　　　　　　\(money(faqCost * 5 * 10_000))")
} catch {
    print("❌ \(error.localizedDescription)")
    exit(1)
}
