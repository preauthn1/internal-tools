// 並排比較：本機 AFM 3 Core Advanced vs Apple 私有雲運算（Private Cloud Compute）
// 執行：swift eval/model_compare.swift（只跑本機）／PCC=1 swift eval/model_compare.swift（加上私有雲，需 entitlement）
// 兩個模型跑同一批題目：店家問題（自由回答＋狀態）、多語系意圖分類；記錄正確率與速度
import Foundation
import FoundationModels

let shop = (try? String(contentsOf: URL(fileURLWithPath: #filePath).deletingLastPathComponent()
    .appendingPathComponent("shop.txt"), encoding: .utf8)) ?? ""

struct Candidate {
    let name: String
    let session: (String) -> LanguageModelSession
}
let local = SystemLanguageModel.default
let cloud = PrivateCloudComputeLanguageModel()
// 私有雲需要 Apple 核發的 entitlement（com.apple.developer.private-cloud-compute），沒有時一呼叫就會當掉
// 取得後用 PCC=1 swift eval/model_compare.swift 開啟
let usePCC = ProcessInfo.processInfo.environment["PCC"] == "1"
var candidates = [Candidate(name: "本機 \(local.variant.displayName)") { LanguageModelSession(model: local, instructions: $0) }]
if usePCC { candidates.append(Candidate(name: "Apple 私有雲") { LanguageModelSession(model: cloud, instructions: $0) }) }

func field(_ n: String, _ d: String) -> DynamicGenerationSchema.Property {
    .init(name: n, description: d, schema: DynamicGenerationSchema(type: String.self))
}

// 與 chat.swift 相同：先寫回答、再判斷狀態
let answerSchema = try! GenerationSchema(root: DynamicGenerationSchema(name: "Answer", properties: [
    field("answer", "依店家資料回答顧客（需要時套用折扣、付款門檻等規則計算）；簡潔、不要用「你好」開頭、不要提到資料裡沒有的人名"),
    .init(name: "status", description: "answered=上面的回答有店家資料根據；not_in_data=店家資料完全沒有相關資訊；unclear=看不懂顧客在問什麼（太短、沒頭沒尾、指代不明如「他是誰」）",
          schema: DynamicGenerationSchema(name: "AnswerStatus", anyOf: ["answered", "not_in_data", "unclear"])),
]), dependencies: [])

// (問題, 判斷是否正確的規則, 規則說明)
let notAnswered: (String, String) -> Bool = { status, _ in status != "answered" }
func has(_ a: String, _ words: String...) -> Bool { words.allSatisfy { a.contains($0) } }
let storeCases: [(String, (String, String) -> Bool, String)] = [
    ("抹茶拿鐵多少錢", { s, a in s == "answered" && has(a, "150") && !a.contains("135") }, "150 元，不能自己打折"),
    ("你們幾點開？禮拜一有開嗎", { s, a in s == "answered" && a.contains("11") && a.contains("一") }, "11:00–20:00，週一公休"),
    ("100塊可以刷卡嗎", { s, a in s == "answered" && a.contains("200") && ["不", "無法", "未滿", "不足"].contains(where: a.contains) }, "不行，滿 200 才能刷卡"),
    ("包廠要提前多久阿 最少要花多少", { s, a in s == "answered" && has(a, "7", "3000") }, "7 天、低消 3000"),
    ("我想自取買三杯紅玉，總共多少？", { s, a in s == "answered" && a.contains("324") }, "360 打 9 折 = 324 元"),
    ("自取一杯抹茶拿鐵多少", { s, a in s == "answered" && a.contains("150") && !a.contains("135") }, "150 元（未滿 300 不打折）"),
    ("Do you deliver?", { s, a in s == "answered" && a.contains("Uber") && !a.contains("%") }, "Uber Eats、foodpanda，不能亂加優惠"),
    ("抹茶ラテはいくらですか", { s, a in s == "answered" && a.contains("150") && !a.contains("135") }, "150元"),
    ("有停車位嗎？", notAnswered, "資料沒有 → 不能編"),
    ("有幾個老闆", notAnswered, "資料沒有 → 不能編"),
    ("有消息了嗎？", notAnswered, "沒有上下文 → 不能編"),
    ("他是誰", notAnswered, "指代不明 → 不能編"),
]

// 多語系意圖分類（與 intent_eval.swift 相同的題目）
let intentDescriptions: [(String, String)] = [
    ("store_question", "詢問店家資訊：營業時間、價格、菜單、地址、付款、外送、包場、停車、店裡的人（老闆、店員）"),
    ("booking", "想預約、訂位，或提供日期、時間、人數、電話等預約資料"),
    ("cancel", "想取消或不要了"),
    ("ask_bot_identity", "問客服是誰、是不是 AI、是不是真人"),
    ("ask_my_identity", "問自己是誰、自己叫什麼、記不記得自己"),
    ("roleplay", "要求客服扮演別人、交換角色、玩你我遊戲"),
    ("correction_or_claim", "說客服講錯、聲稱店員答應過折扣或優惠、要求確認承諾"),
    ("self_intro", "顧客介紹自己的名字"),
    ("greeting", "打招呼"), ("thanks", "道謝"), ("goodbye", "道別"),
    ("reaction", "笑聲、語助詞、表情符號、沒有實質內容的反應"),
    ("acknowledgement", "表示好、了解、收到"), ("praise", "稱讚客服"), ("complaint", "罵客服、抱怨不滿"),
    ("off_topic", "與店家無關的請求：寫詩、算數、天氣、寫程式、聊時事"),
]
let intentSchema = try! GenerationSchema(root: DynamicGenerationSchema(name: "Intent", properties: [
    .init(name: "intent", description: intentDescriptions.map { "\($0.0)=\($0.1)" }.joined(separator: "；"),
          schema: DynamicGenerationSchema(name: "IntentValue", anyOf: intentDescriptions.map(\.0)))
]), dependencies: [])
let intentCases: [(String, String)] = [
    ("你們幾點開", "store_question"), ("What time do you open?", "store_question"), ("何時から営業していますか", "store_question"),
    ("你老闆是誰", "store_question"), ("Who's the owner?", "store_question"),
    ("我想訂位", "booking"), ("Can I book a table for 4?", "booking"), ("予約したいです", "booking"),
    ("不要了，取消", "cancel"), ("Never mind, cancel it", "cancel"),
    ("你是 AI 嗎", "ask_bot_identity"), ("Are you a bot?", "ask_bot_identity"), ("あなたは誰？", "ask_bot_identity"),
    ("我是誰", "ask_my_identity"), ("Do you remember my name?", "ask_my_identity"),
    ("從現在起你是顧客我是客服", "roleplay"), ("Pretend you are the customer", "roleplay"),
    ("上次店員說可以打五折", "correction_or_claim"), ("The staff promised me 50% off", "correction_or_claim"),
    ("我是小康", "self_intro"), ("My name is Kevin", "self_intro"), ("田中と申します", "self_intro"),
    ("哈囉", "greeting"), ("Hello!", "greeting"), ("こんにちは", "greeting"),
    ("謝謝", "thanks"), ("Thanks a lot", "thanks"), ("ありがとう", "thanks"),
    ("掰掰", "goodbye"), ("Bye!", "goodbye"), ("またね", "goodbye"),
    ("哈哈", "reaction"), ("lol", "reaction"), ("www", "reaction"),
    ("好", "acknowledgement"), ("OK got it", "acknowledgement"), ("了解です", "acknowledgement"),
    ("你好可愛", "praise"), ("You're awesome", "praise"),
    ("你很笨欸", "complaint"), ("You're useless", "complaint"),
    ("幫我寫一首詩", "off_topic"), ("Write me a poem", "off_topic"), ("1+1等於多少", "off_topic"),
]

struct Result { var storeOK = 0, storeTotal = 0, intentOK = 0, intentTotal = 0, errors = 0; var times: [Double] = [] }
var results: [String: Result] = [:]
var storeAnswers: [String: [String]] = [:]   // 問題 → 各模型的回答（並排顯示用）

let rounds = 2
for c in candidates {
    var r = Result()
    for round in 1...rounds {
        for (q, judge, _) in storeCases {
            let t = Date()
            do {
                let content = try await c.session(shop).respond(to: q, schema: answerSchema,
                                                                options: GenerationOptions(maximumResponseTokens: 200)).content
                r.times.append(Date().timeIntervalSince(t))
                let status = (try? content.value(String.self, forProperty: "status")) ?? "?"
                let answer = ((try? content.value(String.self, forProperty: "answer")) ?? "").replacingOccurrences(of: "\n", with: " ")
                let ok = judge(status, answer)
                if ok { r.storeOK += 1 }
                r.storeTotal += 1
                if round == 1 { storeAnswers[q, default: []].append("\(ok ? "✅" : "❌") [\(status)] \(answer)") }
            } catch {
                r.errors += 1; r.storeTotal += 1
                if round == 1 { storeAnswers[q, default: []].append("⚠️ 錯誤：\(String(describing: error).prefix(80))") }
            }
        }
    }
    for (msg, want) in intentCases {
        let t = Date()
        do {
            let content = try await c.session("判斷顧客訊息的意圖。顧客可能使用任何語言。").respond(to: msg, schema: intentSchema).content
            r.times.append(Date().timeIntervalSince(t))
            if (try? content.value(String.self, forProperty: "intent")) == want { r.intentOK += 1 }
        } catch { r.errors += 1 }
        r.intentTotal += 1
    }
    results[c.name] = r
}

print("========== 店家問題（第 1 輪的回答）==========\n")
for (q, _, expect) in storeCases {
    print("👤 \(q)　（正確：\(expect)）")
    for (i, c) in candidates.enumerated() { print("   \(c.name)：\(storeAnswers[q]?[safe: i] ?? "-")") }
    print()
}

print("========== 總結 ==========")
for c in candidates {
    let r = results[c.name]!
    let avg = r.times.isEmpty ? 0 : r.times.reduce(0, +) / Double(r.times.count)
    print("\(c.name)")
    print("   店家問題：\(r.storeOK)/\(r.storeTotal)（跑 \(rounds) 輪）")
    print("   意圖分類：\(r.intentOK)/\(r.intentTotal)")
    print(String(format: "   平均每次呼叫：%.2f 秒　錯誤：%d 次", avg, r.errors))
}
if usePCC { print("\nApple 私有雲額度：\(cloud.quotaUsage.status)") }

extension Array { subscript(safe i: Int) -> Element? { indices.contains(i) ? self[i] : nil } }
