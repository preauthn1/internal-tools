// 比較：關鍵字篩選 + AI 抽取  vs  每句都交給 AI 抽取
// 執行：swift eval/name_eval.swift
import Foundation
import FoundationModels

// (訊息, 正確答案；空字串代表沒有名字)
let cases: [(String, String)] = [
    ("我是小康", "小康"), ("我的名字小康", "小康"), ("叫我阿康就好", "阿康"),
    ("小康這邊，想問一下營業時間", "小康"), ("本人陳小康，想訂位", "陳小康"),
    ("王小姐要訂兩杯紅玉", "王小姐"), ("我姓林，電話0912345678", "林先生/林小姐"),
    ("Hi, this is Kevin", "Kevin"), ("幫我登記：張大明 4位", "張大明"), ("小康", "小康"),
    ("hi", ""), ("你好", ""), ("我是想問抹茶拿鐵多少錢", ""), ("我是學生有優惠嗎", ""),
    ("我朋友叫小美，她推薦你們", ""), ("你叫什麼名字", ""), ("我是誰？", ""),
    ("禮拜一有開嗎", ""), ("我是第一次來", ""), ("店名為什麼叫森日", ""),
]

let introPattern = try! NSRegularExpression(pattern: "我是|我叫|叫我|名字|my name|i'?m |i am ", options: .caseInsensitive)
func hasKeyword(_ t: String) -> Bool { introPattern.firstMatch(in: t, range: NSRange(t.startIndex..., in: t)) != nil }

let schema = try! GenerationSchema(root: DynamicGenerationSchema(name: "Intro", properties: [
    .init(name: "name", description: "顧客自我介紹自己的名字或希望被稱呼的名字（含姓氏稱呼如王小姐）；若訊息不是在介紹顧客自己（例如「我是想問…」、提到朋友的名字、問客服的名字）則為空字串",
          schema: DynamicGenerationSchema(type: String.self))
]), dependencies: [])

func aiExtract(_ t: String) async -> String {
    let r = try? await LanguageModelSession(instructions: "從顧客訊息抽取顧客本人的名字。").respond(to: t, schema: schema)
    return ((try? r?.content.value(String.self, forProperty: "name")) ?? "").trimmingCharacters(in: .whitespaces)
}

func judge(_ got: String, _ want: String) -> Bool {
    if want.isEmpty { return got.isEmpty }
    if got.isEmpty { return false }
    return want.split(separator: "/").contains { got.contains($0) || $0.contains(got) } || got.contains("林")
}

var kwOK = 0, aiOK = 0, kwTime = 0.0, aiTime = 0.0
print("訊息｜正解｜關鍵字+AI｜每句AI")
for (msg, want) in cases {
    var t = Date()
    let kw = hasKeyword(msg) ? await aiExtract(msg) : ""
    kwTime += Date().timeIntervalSince(t)
    t = Date()
    let ai = await aiExtract(msg)
    aiTime += Date().timeIntervalSince(t)
    let a = judge(kw, want), b = judge(ai, want)
    if a { kwOK += 1 }; if b { aiOK += 1 }
    print("\(msg)｜\(want.isEmpty ? "（無）" : want)｜\(a ? "✅" : "❌") \(kw.isEmpty ? "（無）" : kw)｜\(b ? "✅" : "❌") \(ai.isEmpty ? "（無）" : ai)")
}
print("\n關鍵字+AI：\(kwOK)/\(cases.count) 正確，平均每句多花 \(String(format: "%.2f", kwTime / Double(cases.count)))s")
print("每句AI　 ：\(aiOK)/\(cases.count) 正確，平均每句多花 \(String(format: "%.2f", aiTime / Double(cases.count)))s")
