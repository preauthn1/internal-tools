// 互動式測試：扮演 LINE 顧客和本機 AI 客服對話（含意圖分流防護、多語系）
// 執行：swift eval/chat.swift            （預設讀 eval/shop.txt）
//       swift eval/chat.swift 你的檔案.txt
// 指令：/reset 換一位新顧客（清空記憶）  /reload 重新讀取店家設定
//       /me 名字 手動設定顧客名字  /pending 看待處理清單  /quit 離開
//
// 每則訊息的流程（模擬正式系統）：
//   1. 語言：Apple NaturalLanguage 辨識（比問 AI 準），太短判不出來就沿用上一句的語言；
//      只支援中英日，其他語言回固定的多語說明並交給店員
//   2. 同時請模型做「意圖細分類」與「名字抽取」（都是固定格式輸出，任何語言都適用）
//   3. 依意圖分流：店家問題 → AI 依資料回答；預約 → 程式流程收單；
//      其他（身分、刁難、閒聊…）→ 固定模板，依顧客語言回覆，不讓 AI 自由發揮
//   4. 檢查 AI 輸出，出現危險字眼就換成安全回覆
//   5. 連續 3 次無關/聲稱折扣 → 轉人工，AI 停止回覆
// 程式規則一律不依賴特定語言的關鍵字：只檢查數字、原文出處、符號等
import Foundation
import FoundationModels
import NaturalLanguage

let path = CommandLine.arguments.count > 1
    ? CommandLine.arguments[1]
    : URL(fileURLWithPath: #filePath).deletingLastPathComponent().appendingPathComponent("shop.txt").path

func loadInstructions() -> String {
    (try? String(contentsOfFile: path, encoding: .utf8)) ?? "你是親切的客服。"
}

// 對話紀錄存在 eval/logs/（可用環境變數 CHAT_LOG_DIR 改位置，自動測試用，避免動到使用者的紀錄）
let logDir = ProcessInfo.processInfo.environment["CHAT_LOG_DIR"].map { URL(fileURLWithPath: $0) }
    ?? URL(fileURLWithPath: path).deletingLastPathComponent().appendingPathComponent("logs")
let _ = try? FileManager.default.createDirectory(at: logDir, withIntermediateDirectories: true)
let logURL = logDir.appendingPathComponent("chat_log_\(Int(Date().timeIntervalSince1970)).md")
var log = "# 對話紀錄（\(Date())）\n設定檔：\(path)\n\n"
func save() { try? log.write(to: logURL, atomically: true, encoding: .utf8) }

guard case .available = SystemLanguageModel.default.availability else {
    print("❌ 本機模型無法使用：\(SystemLanguageModel.default.availability)"); exit(1)
}

// MARK: - 狀態（正式系統會以 LINE userId 為 key 存進資料庫）

// 只有「店家問題」的回合會進入這個對話紀錄，刁難的回合不留，避免越聊越歪
var session = LanguageModelSession(instructions: loadInstructions())
session.prewarm()
var customerName: String?
var customerLang = "zh"               // 預設店家語言；之後跟著顧客走
var offTopicStreak = 0
var humanMode = false
var booking: [String: String] = [:]   // 進行中的預約
var bookingActive = false
var lastAsked: String?                // 預約流程剛問了哪個欄位（顧客的下一句多半就是答案）
var awaitingStaff = false             // 已經跟顧客說過「轉給店員」，還沒有店員回覆
var lastBotReply = ""
var faqs = loadFAQ()
let staffReplyKeys = ["fallback", "claim", "waiting", "nudge"]  // 這些回覆代表「已轉給店員」
var followUpStreak = 0                // 連續追問次數
var turnFollowUpReply: String?        // 這一輪的追問回覆（A、B 共用，只算一次）
var pending: [String] = []            // 待店員處理：聲稱/更正、轉人工、預約申請

// MARK: - 語言

// 本機模型只有中英日的回答品質可靠；其他語言一律回固定的多語說明並轉店員
let supportedLanguages: Set<String> = ["zh", "en", "ja"]
let unsupportedNotice = "抱歉，目前只能用中文、English、日本語為您服務，已幫您通知店員～\nSorry, we currently support Chinese, English and Japanese only. Our staff will get back to you.\n申し訳ありません、現在は中国語・英語・日本語のみ対応しております。"

func detectLanguage(_ text: String) -> String? {
    // 先看文字系統：有假名就是日文、有韓文字母就是韓文、有漢字就是中文
    // （「可以用LINE Pay嗎」這種夾英文品牌名的句子，語言辨識常誤判成英文）
    let scalars = text.unicodeScalars
    if scalars.contains(where: { (0x3040...0x30FF).contains($0.value) }) { return "ja" }
    if scalars.contains(where: { (0xAC00...0xD7AF).contains($0.value) || (0x1100...0x11FF).contains($0.value) }) { return "ko" }
    if scalars.contains(where: { $0.properties.isIdeographic }) { return "zh" }
    guard text.filter(\.isLetter).count >= 2 else { return nil }  // 表情、標點、單一字：沿用上一句
    let r = NLLanguageRecognizer()
    r.languageConstraints = [.traditionalChinese, .simplifiedChinese, .english, .japanese, .korean, .vietnamese, .thai,
                             .spanish, .french, .german, .indonesian, .malay, .portuguese, .italian, .russian]
    r.processString(text)
    let hypotheses = r.languageHypotheses(withMaximum: 3).sorted { $0.value > $1.value }
    guard let best = hypotheses.first else { return nil }
    func code(_ l: NLLanguage) -> String { l.rawValue.hasPrefix("zh") ? "zh" : String(l.rawValue.prefix(2)) }
    // 短句很容易誤判（「Hi, I am Kevin」被判成 80% 印尼文）。誤擋真正的中英日顧客代價比較大，
    // 所以短句只要候選裡有支援的語言（≥10%），就優先用支援的語言
    if !supportedLanguages.contains(code(best.key)), text.count < 40,
       let alt = hypotheses.first(where: { supportedLanguages.contains(code($0.key)) && $0.value >= 0.1 }) {
        return code(alt.key)
    }
    if best.value < 0.5 {
        // 「Hello!」「Bye」「lol」這類太短、信心低：全是拉丁字母就當英文，否則沿用上一句
        return text.unicodeScalars.filter({ $0.properties.isAlphabetic }).allSatisfy({ $0.isASCII }) ? "en" : nil
    }
    return best.key.rawValue.hasPrefix("zh") ? "zh" : String(best.key.rawValue.prefix(2))
}

func englishName(_ code: String) -> String {
    code == "zh" ? "Traditional Chinese" : (Locale(identifier: "en").localizedString(forLanguageCode: code) ?? code)
}

// 固定文案：中英日手寫
// 一個 key 可放多句，用 | 分隔，隨機挑一句比較不像機器人
let strings: [String: [String: String]] = [
    "ask_name_soft": ["zh": "還不知道怎麼稱呼您呢？", "en": " May I have your name?", "ja": "お名前を伺ってもよろしいですか？"],
    "nudge": ["zh": "我再幫您催一下店員，造成不便很抱歉 🙏", "en": "I'll nudge our staff again — sorry for the wait 🙏",
              "ja": "スタッフにもう一度催促しますね。お待たせして申し訳ありません 🙏"],
    "not_question": ["zh": "好的～有需要再跟我說喔！", "en": "Sure~ Just let me know if you need anything!",
                     "ja": "かしこまりました～何かあればお気軽にどうぞ！"],
    "waiting": ["zh": "已經通知店員了，店員看到後會盡快回覆您，請稍等喔～", "en": "I've already notified our staff — they'll get back to you as soon as they can. Thanks for waiting~",
                "ja": "スタッフにはすでに連絡済みです。確認でき次第ご返信しますので、少々お待ちください～"],
    "unclear": ["zh": "不太確定您的意思，可以再說詳細一點嗎？", "en": "Sorry, I'm not quite sure what you mean. Could you tell me a bit more?",
                "ja": "すみません、もう少し詳しく教えていただけますか？"],
    "fallback": ["zh": "這個問題我幫您轉給店員確認喔！", "en": "Let me check with our staff and get back to you!",
                 "ja": "こちらはスタッフに確認してご連絡しますね！"],
    "off_topic": ["zh": "這個我不太清楚耶～", "en": "Sorry, I'm not sure about that~",
                  "ja": "すみません、それはちょっと分かりません～"],
    "bot_identity": ["zh": "我是森日茶屋的 AI 客服小茶～可以幫您查營業時間、價格，或協助預約喔！",
                     "en": "I'm Xiao Cha, the AI assistant of Senri Tea House~ I can help with opening hours, prices, or reservations!",
                     "ja": "森日茶屋の AI アシスタント、小茶です～営業時間や価格のご案内、ご予約のお手伝いができます！"],
    "you_are": ["zh": "您是{name}～", "en": "You're {name}~ ", "ja": "{name}様ですね～"],
    "claim": ["zh": "這部分我沒辦法確認，已幫您記錄下來，會請店員回覆您喔！",
              "en": "I can't confirm that, but I've noted it down and our staff will get back to you!",
              "ja": "こちらでは確認できないため、記録してスタッフからご連絡しますね！"],
    "handoff": ["zh": "我先幫您轉給店員，請稍候喔～", "en": "Let me connect you with our staff, please wait a moment~",
                "ja": "スタッフにおつなぎしますので、少々お待ちください～"],
    "greeting": ["zh": "您好～我是森日茶屋的小茶，可以幫您查營業時間、價格或預約喔！|您好！今天想喝點什麼呢？有問題都可以問我～",
                 "en": "Hi there~ I'm Xiao Cha from Senri Tea House. I can help with hours, prices or reservations!|Hello! What would you like to drink today? Feel free to ask me anything~",
                 "ja": "こんにちは～森日茶屋の小茶です。営業時間や価格、ご予約のご案内ができます！|こんにちは！今日は何を飲みたいですか？何でも聞いてくださいね～"],
    "self_intro": ["zh": "{name}您好～很高興認識您！有什麼可以幫您的嗎？", "en": "Nice to meet you, {name}! How can I help you today?",
                   "ja": "{name}様、はじめまして！何かお手伝いできることはありますか？"],
    "thanks": ["zh": "不客氣～還有需要再跟我說喔！|不會～很高興能幫上忙 😊", "en": "You're welcome~ Let me know if you need anything else!|No problem, happy to help 😊",
               "ja": "どういたしまして～他に何かあればお気軽にどうぞ！|お役に立ててうれしいです 😊"],
    "goodbye": ["zh": "謝謝您，掰掰～歡迎再來森日茶屋！|好的～祝您有美好的一天，掰掰！", "en": "Thank you, bye~ Hope to see you at Senri Tea House again!|Have a great day, bye!",
                "ja": "ありがとうございました、またお越しください～！|素敵な一日を、またね！"],
    "reaction": ["zh": "😄|😊 還有什麼需要幫忙的嗎？", "en": "😄|😊 Anything else I can help with?",
                 "ja": "😄|😊 他に何かお手伝いできますか？"],
    "acknowledgement": ["zh": "好的～還有其他需要再跟我說喔！|👌 有問題隨時找我～", "en": "Got it~ Let me know if you need anything else!|👌 Ask me anytime~",
                        "ja": "承知しました～他にもあればどうぞ！|👌 いつでも聞いてくださいね～"],
    "praise": ["zh": "謝謝誇獎～ 😊 有什麼可以幫您的嗎？", "en": "Aww, thank you~ 😊 How can I help?",
               "ja": "ありがとうございます～ 😊 何かお手伝いできますか？"],
    "complaint": ["zh": "抱歉讓您不滿意 🙏 需要的話我可以幫您轉給店員喔！", "en": "Sorry for the trouble 🙏 I can connect you with our staff if you'd like!",
                  "ja": "ご不便をおかけしてすみません 🙏 必要ならスタッフにおつなぎします！"],
    // 預約流程
    "ask_date": ["zh": "想預約哪一天呢？", "en": "Which date would you like to book?", "ja": "ご希望の日にちはいつですか？"],
    "ask_time": ["zh": "想約幾點呢？（營業時間 11:00–20:00）", "en": "What time would you like? (We're open 11:00–20:00)",
                 "ja": "何時がご希望ですか？（営業時間 11:00–20:00）"],
    "ask_people": ["zh": "請問幾位呢？", "en": "How many people?", "ja": "何名様ですか？"],
    "ask_name": ["zh": "請問預約人大名？", "en": "May I have your name for the reservation?", "ja": "ご予約のお名前をお願いします。"],
    "ask_phone": ["zh": "請留一支聯絡電話～", "en": "Could you leave a contact phone number?", "ja": "ご連絡先の電話番号をお願いします～"],
    "err_date": ["zh": "日期看不太懂，可以再說一次嗎？", "en": "Sorry, I didn't catch the date. Could you say it again?",
                 "ja": "日付がよく分かりませんでした。もう一度お願いします。"],
    "err_past": ["zh": "這天已經過了喔，請選今天以後的日期～", "en": "That date has passed. Please pick a date from today onward~",
                 "ja": "その日は過ぎています。今日以降の日付をお選びください～"],
    "err_monday": ["zh": "週一公休喔，請選其他天～", "en": "We're closed on Mondays, please pick another day~",
                   "ja": "月曜日は定休日です。別の日をお選びください～"],
    "err_time": ["zh": "時間看不太懂，可以再說一次嗎？", "en": "Sorry, I didn't catch the time. Could you say it again?",
                 "ja": "時間がよく分かりませんでした。もう一度お願いします。"],
    "err_hours": ["zh": "營業時間是 11:00–20:00，請選這段時間內喔～", "en": "We're open 11:00–20:00, please pick a time within those hours~",
                  "ja": "営業時間は 11:00–20:00 です。その時間内でお選びください～"],
    "err_people": ["zh": "請問幾位呢？", "en": "How many people?", "ja": "何名様ですか？"],
    "err_too_many": ["zh": "店內最多 20 位，人數較多可以詢問包場喔～", "en": "We seat up to 20 people. For larger groups, ask us about private booking~",
                     "ja": "店内は最大 20 名様までです。それ以上は貸切をご相談ください～"],
    "err_phone": ["zh": "電話號碼好像不完整，可以再確認一下嗎？", "en": "The phone number looks incomplete. Could you double-check it?",
                  "ja": "電話番号が不完全なようです。もう一度ご確認いただけますか？"],
    "people_unit": ["zh": " 位", "en": " people", "ja": "名"],
    "booking_done": ["zh": "收到您的預約申請：{summary}。店員確認後會再通知您喔！", "en": "Got your reservation request: {summary}. We'll let you know once our staff confirms it!",
                     "ja": "ご予約リクエストを承りました：{summary}。スタッフが確認後ご連絡します！"],
    "booking_cancel": ["zh": "好的，已取消這次預約～", "en": "OK, I've cancelled this reservation request~",
                       "ja": "かしこまりました、今回のご予約は取り消しました～"],
]

@MainActor func L(_ key: String, _ vars: [String: String] = [:]) async -> String {
    guard let table = strings[key] else { return key }
    var s = (table[customerLang] ?? table["zh"]!).components(separatedBy: "|").randomElement()!
    for (k, v) in vars { s = s.replacingOccurrences(of: "{\(k)}", with: v) }
    return s
}

// MARK: - 結構化輸出（模型最準的能力）

func stringField(_ name: String, _ desc: String) -> DynamicGenerationSchema.Property {
    .init(name: name, description: desc, schema: DynamicGenerationSchema(type: String.self))
}

// 本機模型偶爾會卡住好幾分鐘：超過時間就放棄，改用安全回覆
@Sendable func withTimeout<T: Sendable>(_ seconds: Double, _ op: @escaping @Sendable () async -> T?) async -> T? {
    await withTaskGroup(of: T?.self) { group in
        group.addTask { await op() }
        group.addTask { try? await Task.sleep(for: .seconds(seconds)); return nil }
        let first = await group.next() ?? nil
        group.cancelAll()
        return first
    }
}

@Sendable func extract(_ text: String, instructions: String, schema: GenerationSchema) async -> GeneratedContent? {
    await withTimeout(10) {
        try? await LanguageModelSession(instructions: instructions)
            .respond(to: text, schema: schema, options: GenerationOptions(maximumResponseTokens: 120)).content
    }
}

// 細分類讓模型處理（任何語言都適用）；容易互相混淆的「問你是誰/問我是誰/角色扮演」合併成 identity，回同一個安全模板
let intentDescriptions: [(String, String)] = [
    ("store_question", "詢問店家資訊：營業時間、價格、菜單、地址、付款、外送、包場、停車、店裡的人（老闆、店員）"),
    ("booking", "想預約、訂位，或提供日期、時間、人數、電話等預約資料"),
    ("cancel", "想取消或不要了"),
    ("identity", "問客服是誰、是不是 AI、問客服的個人資訊（幾歲、住哪、男生女生、喜歡什麼）、問自己是誰、記不記得自己、要求客服扮演別人或交換角色"),
    ("follow_up", "追問先前問題的進度、問店員回覆了沒、催促（例如：有回應了嗎、好了沒、還要等多久、any update、still waiting、まだですか）"),
    ("correction_or_claim", "說客服講錯、聲稱店員答應過折扣或優惠、要求確認承諾"),
    ("self_intro", "顧客介紹自己的名字"),
    ("greeting", "打招呼"),
    ("thanks", "道謝"),
    ("goodbye", "道別"),
    ("reaction", "笑聲、語助詞、表情符號、沒有實質內容的反應"),
    ("acknowledgement", "表示好、了解、收到"),
    ("praise", "稱讚客服"),
    ("complaint", "罵客服、抱怨不滿"),
    ("off_topic", "與店家無關的請求：寫詩、算數、天氣、寫程式、聊時事"),
]
let intentSchema = try! GenerationSchema(root: DynamicGenerationSchema(name: "Intent", properties: [
    .init(name: "intent", description: intentDescriptions.map { "\($0.0)=\($0.1)" }.joined(separator: "；"),
          schema: DynamicGenerationSchema(name: "IntentValue", anyOf: intentDescriptions.map(\.0)))
]), dependencies: [])

let nameSchema = try! GenerationSchema(root: DynamicGenerationSchema(name: "Intro", properties: [
    stringField("name", "顧客自我介紹自己的名字或希望被稱呼的名字（含姓氏稱呼如王小姐）；若訊息不是在介紹顧客自己（例如「我是想問…」、提到朋友的名字、問客服的名字）則為空字串")
]), dependencies: [])

// *_quote 要求模型照抄訊息原文出處，程式再檢查原文真的在訊息裡 → 不靠任何語言的關鍵字也能擋住模型自己編的值
let bookingSchema = try! GenerationSchema(root: DynamicGenerationSchema(name: "Booking", properties: [
    stringField("date", "預約日期，格式 YYYY-MM-DD；未提及則為空字串"),
    stringField("date_quote", "訊息中提到日期的原文片段，逐字照抄；未提及則為空字串"),
    stringField("time", "時段，24 小時制 HH:mm（下午兩點=14:00）；未提及則為空字串"),
    stringField("time_quote", "訊息中提到時間的原文片段，逐字照抄；未提及則為空字串"),
    stringField("people", "人數，只填阿拉伯數字；未提及則為空字串"),
    stringField("people_quote", "訊息中提到人數的原文片段，逐字照抄；未提及則為空字串"),
    stringField("name", "預約人姓名；未提及則為空字串"),
    stringField("phone", "電話號碼；未提及則為空字串"),
]), dependencies: [])

let followUpSchema = try! GenerationSchema(root: DynamicGenerationSchema(name: "FollowUp", properties: [
    .init(name: "is_follow_up", description: "yes=顧客在追問那件事的進度、問有沒有消息、催促；no=顧客在講別的事（新問題、道謝、打招呼等）",
          schema: DynamicGenerationSchema(name: "YesNo", anyOf: ["yes", "no"]))
]), dependencies: [])

// 只有在「已經說過轉給店員」時才問：帶著上下文的是非題，比 14 選 1 的大分類準，而且不靠任何語言的關鍵字
@Sendable func isFollowUp(_ text: String, previousReply: String, enabled: Bool) async -> Bool {
    guard enabled else { return false }
    let c = await extract(text, instructions: "客服上一句對顧客說：「\(previousReply)」。判斷顧客接下來這句話是不是在追問那件事。顧客可能使用任何語言。", schema: followUpSchema)
    return (try? c?.value(String.self, forProperty: "is_follow_up")) == "yes"
}

// MARK: - B 版本：FAQ 挑選（AI 只挑題，不寫字）

struct FAQ: Sendable { let id: String; let questions: String; let answers: [String: String] }

func loadFAQ() -> [FAQ] {
    let faqPath = URL(fileURLWithPath: path).deletingLastPathComponent().appendingPathComponent("faq.txt").path
    guard let text = try? String(contentsOfFile: faqPath, encoding: .utf8) else { return [] }
    return text.split(separator: "\n").compactMap { raw in
        let line = raw.trimmingCharacters(in: .whitespaces)
        if line.isEmpty || line.hasPrefix("#") { return nil }
        let cols = line.components(separatedBy: " | ").map { $0.trimmingCharacters(in: .whitespaces) }
        guard cols.count >= 3 else { return nil }
        var answers = ["zh": cols[2]]
        if cols.count > 3 { answers["en"] = cols[3] }
        if cols.count > 4 { answers["ja"] = cols[4] }
        return FAQ(id: cols[0], questions: cols[1], answers: answers)
    }
}

@Sendable func pickFAQ(_ text: String, _ faqs: [FAQ], previousReply: String) async -> String {
    guard !faqs.isEmpty else { return "none" }
    let ids = faqs.map(\.id) + ["none", "not_a_question"]
    let list = faqs.map { "\($0.id)：\($0.questions)" }.joined(separator: "\n")
    let schema = try! GenerationSchema(root: DynamicGenerationSchema(name: "Pick", properties: [
        .init(name: "faq", description: "能直接回答顧客問題的 FAQ 代號；沒有任何一題能回答就選 none；顧客不是在發問（例如：我想一下、等等、OK）就選 not_a_question",
              schema: DynamicGenerationSchema(name: "FAQId", anyOf: ids))
    ]), dependencies: [])
    // 帶上客服上一句：單看「什麼時候呀」像在問營業時間，接在「已轉給店員」後面其實是在追問
    let context = previousReply.isEmpty ? "" : "客服上一句對顧客說：「\(previousReply)」。請依這個上下文理解顧客的話。\n"
    let c = await extract(text, instructions: "\(context)從 FAQ 清單挑出能回答顧客問題的一題。顧客可能使用任何語言。不確定就選 none。\nFAQ 清單：\n\(list)", schema: schema)
    return (try? c?.value(String.self, forProperty: "faq")) ?? "none"
}

let messageKindSchema = try! GenerationSchema(root: DynamicGenerationSchema(name: "Kind", properties: [
    .init(name: "kind", description: "question=顧客在問問題或要求資訊（含沒有問號的問法，例如「有幾個老闆」「營業時間」）；pause=顧客在思考、請客服稍等（我想一下、等等、我再看看、let me think）；reply=只是回應（OK、好、收到）",
          schema: DynamicGenerationSchema(name: "KindValue", anyOf: ["question", "pause", "reply"]))
]), dependencies: [])

// 「我想一下」這類不是問題的話。實測 40/45，仍會把少數真問題判成 pause，
// 所以只在 A、B 都答不出來時才用，不會擋掉答得出來的問題
@Sendable func isQuestion(_ text: String) async -> Bool {
    let c = await extract(text, instructions: "判斷顧客訊息屬於哪一種。顧客可能使用任何語言。", schema: messageKindSchema)
    return ((try? c?.value(String.self, forProperty: "kind")) ?? "question") == "question"
}

@Sendable func classify(_ text: String) async -> String {
    let c = await extract(text, instructions: "判斷顧客訊息的意圖。顧客可能使用任何語言。", schema: intentSchema)
    return (try? c?.value(String.self, forProperty: "intent")) ?? "off_topic"
}

@Sendable func extractName(_ text: String) async -> String? {
    let c = await extract(text, instructions: "從顧客訊息抽取顧客本人的名字。顧客可能使用任何語言。", schema: nameSchema)
    let name = ((try? c?.value(String.self, forProperty: "name")) ?? "").trimmingCharacters(in: .whitespaces)
    // 模型有時會憑空編名字（例如「他是誰」→ 張先生）：名字本體必須真的出現在訊息裡
    var core = name
    for title in ["先生", "小姐", "女士", "太太", "様", "さん", "님", "Mr. ", "Ms. ", "Mrs. "] {
        core = core.replacingOccurrences(of: title, with: "")
    }
    core = core.trimmingCharacters(in: .whitespaces)
    if core.isEmpty || !text.localizedCaseInsensitiveContains(core) { return nil }
    // 只說了姓氏的中文名（「我姓林」）：加上稱謂
    return name.count == 1 && name.unicodeScalars.allSatisfy({ $0.properties.isIdeographic }) ? "\(name)先生／小姐" : name
}

// MARK: - 預約流程（由程式控制，AI 只負責抽欄位）

let bookingOrder = ["date", "time", "people", "name", "phone"]

let dayFormatter: DateFormatter = {
    let f = DateFormatter(); f.dateFormat = "yyyy-MM-dd"; f.locale = Locale(identifier: "en_US_POSIX"); return f
}()

// 回傳 (錯誤文案 key, 清理後的值)；規則判斷交給程式，不交給 AI
func validate(_ key: String, _ value: String) -> (String?, String) {
    switch key {
    case "date":
        guard let d = dayFormatter.date(from: value) else { return ("err_date", value) }
        if d < Calendar.current.startOfDay(for: Date()) { return ("err_past", value) }
        if Calendar.current.component(.weekday, from: d) == 2 { return ("err_monday", value) }
        return (nil, value)
    case "time":
        let parts = value.split(separator: ":").compactMap { Int($0) }
        guard var h = parts.first else { return ("err_time", value) }
        let m = parts.count > 1 ? parts[1] : 0
        if h < 11, h + 12 <= 20 { h += 12 }  // 「兩點」多半是下午兩點
        if h < 11 || h > 19 { return ("err_hours", value) }
        return (nil, String(format: "%02d:%02d", h, m))
    case "people":
        guard let n = Int(value.filter(\.isNumber)), n > 0 else { return ("err_people", value) }
        if n > 20 { return ("err_too_many", value) }
        return (nil, String(n))
    case "phone":
        let digits = value.filter(\.isNumber)
        if digits.count < 9 || digits.count > 10 { return ("err_phone", value) }
        return (nil, digits)
    default:
        return (nil, value)
    }
}

// 訊息裡獨立出現的數字（「0912345678」只算一個數字，不會讓人數 1 混進來）
func numbers(in text: String) -> Set<Int> {
    let regex = try! NSRegularExpression(pattern: #"\d+"#)
    return Set(regex.matches(in: text, range: NSRange(text.startIndex..., in: text)).compactMap {
        Range($0.range, in: text).flatMap { Int(text[$0]) }
    })
}

// 模型給的值能不能在訊息的數字裡找到依據
func numbersSupport(_ key: String, _ value: String, in text: String) -> Bool {
    let nums = numbers(in: text)
    switch key {
    case "people":
        // 先拿掉日期、時間、電話這類數字（「10/20」的 10 不是人數）
        let stripped = text.replacingOccurrences(of: #"\d+\s*[/:\-.月]\s*\d+|\d{6,}"#, with: " ", options: .regularExpression)
        return Int(value.filter(\.isNumber)).map(numbers(in: stripped).contains) ?? false
    case "date":
        let parts = value.split(separator: "-").compactMap { Int($0) }
        return parts.count == 3 && nums.contains(parts[1]) && nums.contains(parts[2])
    case "time":
        guard let h = value.split(separator: ":").first.flatMap({ Int($0) }) else { return false }
        return nums.contains(h) || nums.contains(h % 12)
    default: return false
    }
}

// probe=true：預約進行中收到其他類訊息時試探用，沒抽到任何新欄位就回傳 nil，交回一般流程
@MainActor func handleBooking(_ text: String, probe: Bool = false) async -> String? {
    let today = dayFormatter.string(from: Date())
    let fieldNames = ["date": "日期", "time": "時間", "people": "人數", "name": "預約人姓名", "phone": "電話"]
    let context = lastAsked.map { "客服剛剛問了顧客的\(fieldNames[$0]!)，顧客這句話很可能就是在回答\(fieldNames[$0]!)。" } ?? ""
    let c = await extract(text, instructions: "從顧客訊息抽取預約資訊。今天是 \(today)。\(context)顧客可能使用任何語言。", schema: bookingSchema)
    func field(_ k: String) -> String { ((try? c?.value(String.self, forProperty: k)) ?? "").trimmingCharacters(in: .whitespaces) }

    var problems: [String] = []
    var found = false
    for key in bookingOrder {
        let raw = field(key)
        if raw.isEmpty || raw == "0" || booking[key] == raw { continue }
        // 防止模型自己編：日期/時間/人數要有原文出處；電話數字、名字必須真的出現在訊息裡
        switch key {
        case "date", "time", "people":
            if key == "people", !raw.contains(where: \.isNumber) { continue }
            let quote = field(key + "_quote")
            let otherQuotes = ["date", "time", "people"].filter { $0 != key }.map { field($0 + "_quote") }.filter { !$0.isEmpty }
            let overlaps = otherQuotes.contains { $0.contains(quote) || quote.contains($0) }
            let quoted = !quote.isEmpty && text.contains(quote) && !overlaps
            // 剛問的欄位直接採用（顧客在回答問題，例如「四個人」）；其他欄位要有依據才採用
            if key != lastAsked && !quoted && !numbersSupport(key, raw, in: text) { continue }
        case "phone":
            let digits = raw.filter(\.isNumber)
            if digits.isEmpty || !text.filter(\.isNumber).contains(digits) { continue }
        case "name":
            if !text.localizedCaseInsensitiveContains(raw) { continue }
        default: break
        }
        found = true
        let (problem, cleaned) = validate(key, raw)
        if let problem { problems.append(problem) } else { booking[key] = cleaned }
    }
    if probe && !found { return nil }
    if booking["name"] == nil, let name = customerName { booking["name"] = name }
    bookingActive = true
    if let p = problems.first { return await L(p) }
    if let next = bookingOrder.first(where: { booking[$0] == nil }) { lastAsked = next; return await L("ask_" + next) }

    let summary = "\(booking["date"]!) \(booking["time"]!), \(booking["people"]!)\(await L("people_unit")), \(booking["name"]!), \(booking["phone"]!)"
    pending.append("📅 預約申請：\(summary)")
    booking = [:]; bookingActive = false; lastAsked = nil
    return await L("booking_done", ["summary": summary])
}

// 追問進度的回覆：連續第 3 次起改成「再催一次」，並在待處理清單標成緊急
@MainActor func followUpReply() async -> String {
    if let r = turnFollowUpReply { return r }
    followUpStreak += 1
    if followUpStreak == 3 { pending.append("🔥 緊急：顧客已連續追問 \(followUpStreak) 次，請盡快回覆") }
    let r = followUpStreak >= 3 ? await L("nudge") : await L("waiting")
    turnFollowUpReply = r
    return r
}

// MARK: - 固定模板與輸出檢查

@MainActor func template(_ intent: String) async -> String? {
    switch intent {
    case "identity":
        if let name = customerName { return (await L("you_are", ["name": name])) + (await L("bot_identity")) }
        return (await L("bot_identity")) + (await L("ask_name_soft"))
    case "correction_or_claim": return await L("claim")
    case "follow_up": return awaitingStaff ? await followUpReply() : await L("unclear")
    case "self_intro": return customerName == nil ? await L("greeting") : await L("self_intro", ["name": customerName!])
    case "cancel": return await L("acknowledgement")  // 沒有進行中的預約
    case "greeting", "thanks", "goodbye", "reaction", "acknowledgement", "praise", "complaint", "off_topic":
        return await L(intent)
    default: return nil
    }
}

// AI 回答裡不該出現的說法（支援的中英日各列一份）
let bannedPhrases = [
    "不是ai", "不是 ai", "已確認", "已完成", "已幫您預約", "預約成功", "五折", "免費", "我是顧客",
    "not an ai", "not a bot", "i'm human", "i am human", "have confirmed", "is confirmed", "been confirmed",
    "booked for you", "50% off", "half price", "for free",
    // A 會自己承諾回覆時間（「我們會盡快回覆」），但它根本沒有轉給任何人
    "盡快回覆", "很快回覆", "會回覆你", "會回覆您", "get back to you", "返信します",
    "aiではありません", "確認しました", "予約しました", "半額", "無料",
]

func safetyCheck(_ reply: String) -> String? {
    let lower = reply.lowercased()
    return bannedPhrases.first { lower.contains($0) }
}

func buildPrompt(_ text: String) -> String {
    var p = "【回覆語言】請用 \(englishName(customerLang)) 回答。\n"
    if let name = customerName { p += "【顧客資料】正在跟你對話的顧客名字是「\(name)」。\n" }
    return p + "\n顧客訊息：\(text)"
}

enum AIOutcome: Sendable { case ok(status: String, answer: String), failed(String) }

// 店家問題的回答：先選狀態，「不知道」「看不懂」由程式回固定句子，AI 不會自己編奇怪的說法
let answerSchema = try! GenerationSchema(root: DynamicGenerationSchema(name: "Answer", properties: [
    // 先寫回答再判斷狀態：讓模型先查資料、算價格，再決定是不是真的答不出來
    stringField("answer", "依店家資料回答顧客（需要時套用折扣、付款門檻等規則計算）；簡潔、不要用「你好」開頭、不要提到資料裡沒有的人名"),
    .init(name: "status", description: "answered=上面的回答有店家資料根據；not_in_data=店家資料完全沒有相關資訊；unclear=看不懂顧客在問什麼（太短、沒頭沒尾、指代不明如「他是誰」）；not_a_question=顧客不是在發問（例如：我想一下、等等、OK、我再看看）",
          schema: DynamicGenerationSchema(name: "AnswerStatus", anyOf: ["answered", "not_in_data", "unclear", "not_a_question"])),
]), dependencies: [])

// MARK: - 主迴圈

print("""
🍵 本機 AI 客服測試（設定檔：\(path)）
   直接輸入訊息扮演顧客（支援中文、English、日本語），試著玩壞它吧！
   /reset 換新顧客  /reload 重讀設定  /me 名字  /pending 待處理清單  /quit 離開
   📋=固定模板  🤖=AI 回答  📅=預約流程  🌐=不支援的語言  [灰字]=分類·語言
   只有「店家問題」（價格、營業時間、地址、付款、外送、包場…）會分別顯示 A（自由回答）與 B（FAQ 挑選）
   其他（打招呼、身分、預約、閒聊…）兩個版本一樣，只顯示一句
   紀錄會存到 \(logURL.path)

""")

// MARK: - 流程測試追蹤（FLOW_TRACE=檔案路徑 時啟用，平常使用不受影響）
// 每次回覆寫一行 JSON：回覆屬於哪一種（模板的 key、answer=AI 自由回答、silent=轉人工不回覆）

let traceURL = ProcessInfo.processInfo.environment["FLOW_TRACE"].map { URL(fileURLWithPath: $0) }

func trace(_ obj: [String: Any]) {
    guard let traceURL, let data = try? JSONSerialization.data(withJSONObject: obj),
          let line = String(data: data, encoding: .utf8) else { return }
    if let h = try? FileHandle(forWritingTo: traceURL) {
        h.seekToEndOfFile(); h.write((line + "\n").data(using: .utf8)!); try? h.close()
    } else {
        try? (line + "\n").write(to: traceURL, atomically: true, encoding: .utf8)
    }
}

// 由回覆文字反推是哪一個模板（各語言、多句變體、{name} 之類的變數都比對）
func replyKey(_ text: String) -> String {
    func matches(_ key: String) -> Bool {
        strings[key]!.values.flatMap { $0.components(separatedBy: "|") }.contains { variant in
            // 先把 {name} 這類變數換成佔位字，跳脫特殊字元後再換回「任意文字」
            let marked = variant.replacingOccurrences(of: #"\{[a-z]+\}"#, with: "\u{1}", options: .regularExpression)
            let pattern = "^" + NSRegularExpression.escapedPattern(for: marked).replacingOccurrences(of: "\u{1}", with: ".+") + "$"
            return text.range(of: pattern, options: .regularExpression) != nil
        }
    }
    if text == unsupportedNotice { return "unsupported" }
    // 身分回覆是兩段組合起來的
    if strings["bot_identity"]!.values.contains(where: { text.contains($0) }) { return "identity" }
    for key in strings.keys.sorted() where !["bot_identity", "you_are", "ask_name_soft", "people_unit"].contains(key) {
        if matches(key) { return key }
    }
    return "answer"
}

func say(_ icon: String, _ text: String, _ intent: String, _ t: Date, who: String = "小茶（A、B 相同）") {
    trace(["key": replyKey(text), "reply": text, "tag": intent, "lang": customerLang])
    let secs = String(format: "%.2f", Date().timeIntervalSince(t))
    print("\(icon) \(who)：\(text)  \u{1B}[90m[\(intent)·\(customerLang)] \(secs)s\u{1B}[0m\n")
    log += "**\(who)**（\(icon) \(intent)）：\(text)\n\n"
    lastBotReply = text
    save()
}

while true {
    print("👤 你：", terminator: ""); fflush(stdout)
    guard let line = readLine()?.trimmingCharacters(in: .whitespaces) else { break }
    if line.isEmpty { continue }
    switch line {
    case "/quit", "/exit": save(); print("👋 紀錄已存：\(logURL.path)"); exit(0)
    case "/reset":
        trace(["key": "reset"])
        session = LanguageModelSession(instructions: loadInstructions())
        customerName = nil; customerLang = "zh"; offTopicStreak = 0; humanMode = false; booking = [:]; bookingActive = false; lastAsked = nil
        awaitingStaff = false; lastBotReply = ""
        log += "---\n*（換一位新顧客）*\n\n"; print("🔄 已換一位新顧客（對話、名字、語言、狀態都清空）\n"); continue
    case "/reload":
        session = LanguageModelSession(instructions: loadInstructions())
        faqs = loadFAQ()
        log += "---\n*（重新讀取設定）*\n\n"; print("📄 已重新讀取設定並重開對話\n"); continue
    case "/pending":
        trace(["key": "pending", "pending": pending])
        print(pending.isEmpty ? "📭 沒有待處理項目\n" : "📬 待店員處理：\n" + pending.map { "   \($0)" }.joined(separator: "\n") + "\n"); continue
    case _ where line.hasPrefix("/me "):
        customerName = String(line.dropFirst(4)).trimmingCharacters(in: .whitespaces)
        log += "*（顧客名字設為 \(customerName!)）*\n\n"; print("🪪 已記住顧客名字：\(customerName!)\n"); continue
    default: break
    }

    log += "**顧客**：\(line)\n\n"
    if humanMode {
        trace(["key": "silent", "reply": "", "tag": "human_mode", "lang": customerLang])
        print("🙋 （轉人工中，AI 不回覆；/reset 結束）\n")
        log += "*（轉人工中，AI 不回覆）*\n\n"; save(); continue
    }

    let t = Date()
    if turnFollowUpReply == nil { followUpStreak = 0 }  // 上一輪不是追問：連續追問次數歸零
    turnFollowUpReply = nil
    if let lang = detectLanguage(line) { customerLang = lang }
    if !supportedLanguages.contains(customerLang) {
        if pending.last.map({ !$0.contains("（\(customerLang)）") && !$0.hasPrefix("   └") }) ?? true
            || !pending.contains(where: { $0.contains("（\(customerLang)）") }) {
            pending.append("🌐 不支援的語言（\(customerLang)）需店員處理：\(customerName ?? "（未知顧客）")")
        }
        pending.append("   └ \(line)")
        awaitingStaff = true
        say("🌐", unsupportedNotice, "unsupported", t); continue
    }
    let currentLine = line
    async let intentResult = classify(currentLine)
    async let nameResult = extractName(currentLine)
    // 只有「上一句」在講轉店員時才檢查追問；否則模型很容易把新問題也判成追問
    let replyBeforeThisTurn = lastBotReply
    let previousWasStaff = staffReplyKeys.contains { key in strings[key]!.values.contains(lastBotReply) } || lastBotReply == unsupportedNotice
    async let followUpResult = isFollowUp(currentLine, previousReply: lastBotReply, enabled: previousWasStaff)
    var intent = await intentResult
    // 只有表情、標點、符號（沒有任何文字或數字）：不分語言都是單純反應
    if line.allSatisfy({ !$0.isLetter && !$0.isNumber }) { intent = "reaction" }
    // 店家問題先讓 A/B 試著回答，兩個都答不出來才看是不是追問（是非題偏向回答 yes）
    let looksLikeFollowUp = await followUpResult
    if looksLikeFollowUp, !["booking", "cancel", "store_question"].contains(intent) { intent = "follow_up" }
    // 「我想一下」「等等」常被分到 off_topic：確認不是在發問就回「好的」（off_topic 本來就不會回答，誤判也沒損失）
    if intent == "off_topic", !(await isQuestion(line)) {
        say("📋", await L("not_question"), "off_topic→not_a_question", t); continue
    }

    var justIntroduced = false
    if let name = await nameResult, customerName == nil {  // 已記住名字就不自動覆蓋（避免把朋友的名字當成顧客），要改用 /me
        customerName = name
        justIntroduced = true
        print("🪪 （程式已記住顧客名字：\(name)）")
        log += "*（程式記住顧客名字：\(name)）*\n\n"
    }

    if bookingActive, intent == "cancel" {
        booking = [:]; bookingActive = false; lastAsked = nil
        say("📅", await L("booking_cancel"), intent, t); continue
    }
    // 預約進行中：顧客補資料（例如只打「4個人」）常被分到其他類，先試著當預約資料處理
    if bookingActive, !["booking", "correction_or_claim"].contains(intent),
       let reply = await handleBooking(line, probe: true) {
        offTopicStreak = 0
        say("📅", reply, "booking（補資料）", t); continue
    }

    if ["correction_or_claim", "off_topic"].contains(intent) {
        offTopicStreak += 1
    } else if ["store_question", "booking", "cancel"].contains(intent) {
        offTopicStreak = 0
    }
    if intent == "correction_or_claim" { pending.append("❓ 顧客聲稱/更正：\(line)"); awaitingStaff = true }

    if offTopicStreak >= 3 {
        humanMode = true
        pending.append("🙋 轉人工：連續 \(offTopicStreak) 次無關/聲稱訊息")
        say("📋", await L("handoff"), intent, t); continue
    }

    if intent == "booking" {
        var reply = await handleBooking(line)
        if reply == nil { reply = await L("ask_date") }
        say("📅", reply!, intent, t); continue
    }

    if justIntroduced, ["identity", "self_intro", "greeting", "reaction", "acknowledgement"].contains(intent) {
        say("📋", await L("self_intro", ["name": customerName!]), intent, t); continue
    }

    if let reply = await template(intent) {
        say("📋", reply, intent, t); continue
    }

    // store_question：A（AI 自由回答）與 B（FAQ 挑選）依序跑（同時跑 A 會出錯），兩種回覆都顯示出來比較
    let whoA = "小茶 A（自由回答）"
    // 先挑 FAQ：像追問、而且 FAQ 沒有一題對得上 → A、B 都回追問（「Do you deliver?」對得上外送 FAQ，不會被誤判）
    let picked = await pickFAQ(currentLine, faqs, previousReply: replyBeforeThisTurn)
    if looksLikeFollowUp, ["none", "not_a_question"].contains(picked) {
        say("📋", await followUpReply(), "\(intent)→follow_up", t); continue
    }

    // 每題用新的 session：保留對話歷史時，小模型會一直抄自己上一句（例如之後每句都回營業時間）
    let prompt = buildPrompt(line), current = LanguageModelSession(instructions: loadInstructions())
    let outcome: AIOutcome? = await withTimeout(15) {
        do {
            let c = try await current.respond(to: prompt, schema: answerSchema, options: GenerationOptions(maximumResponseTokens: 200)).content
            return .ok(status: (try? c.value(String.self, forProperty: "status")) ?? "unclear",
                       answer: ((try? c.value(String.self, forProperty: "answer")) ?? "").trimmingCharacters(in: .whitespacesAndNewlines))
        }
        catch { return .failed(String(describing: error)) }
    }
    switch outcome {
    case .ok(let status, let answer)? where (status != "answered" || answer.isEmpty) && justIntroduced:
        say("📋", await L("self_intro", ["name": customerName!]), "\(intent)→self_intro", t, who: whoA)
    case .ok(let status, let answer)? where (status != "answered" || answer.isEmpty) && looksLikeFollowUp:
        say("📋", await followUpReply(), "\(intent)→follow_up", t, who: whoA)
    case .ok(let status, let answer)? where status != "answered" || answer.isEmpty:
        // A 答不出來，而且顧客其實不是在發問（我想一下、等等…）：不轉店員、不記待補問題，B 也一樣
        if !(await isQuestion(line)) {
            say("📋", await L("not_question"), "\(intent)→not_a_question", t); continue
        }
        if status == "not_in_data" { pending.append("📝 待補問題（知識庫沒有）：\(line)"); awaitingStaff = true }
        say("📋", await L(status == "not_in_data" ? "fallback" : "unclear"), "\(intent)→\(status)", t, who: whoA)
    case .ok(_, let answer)?:
        var reply = answer
        if let bad = safetyCheck(reply) {
            print("⚠️ （AI 原本回答「\(reply)」含有「\(bad)」，已換成安全回覆）")
            log += "*（攔截 AI 原回答：\(reply)）*\n\n"
            reply = await L("fallback")
        }
        say("🤖", reply, intent, t, who: whoA)
    case .failed(let error)?:
        print("❌ \(error)")
        log += "**錯誤**：\(error)\n\n"
        if error.contains("exceededContextWindowSize") {
            session = LanguageModelSession(instructions: loadInstructions())
            print("⚠️ 對話太長超過模型上限，已自動重開對話")
        }
        say("📋", await L("fallback"), intent, t, who: whoA)
    case nil:
        print("⏱️ （AI 超過 15 秒沒回應，已改用安全回覆）")
        log += "*（AI 逾時）*\n\n"
        say("📋", await L("fallback"), intent, t, who: whoA)
    }

    // 這一輪已判定是追問：B 用同一句追問回覆（追問屬於對話流程，不是 FAQ 要比較的部分）
    if let r = turnFollowUpReply {
        print("📚 小茶 B（FAQ 挑選）：\(r)  \u{1B}[90m[follow_up，與 A 相同]\u{1B}[0m\n")
        log += "**小茶 B（FAQ 挑選）**（follow_up）：\(r)\n\n"; save(); continue
    }
    if let faq = faqs.first(where: { $0.id == picked }) {
        print("📚 小茶 B（FAQ 挑選）：\(faq.answers[customerLang] ?? faq.answers["zh"]!)  \u{1B}[90m[faq:\(faq.id)]\u{1B}[0m\n")
        log += "**小茶 B（FAQ 挑選）**（faq:\(faq.id)）：\(faq.answers[customerLang] ?? faq.answers["zh"]!)\n\n"
    } else {
        let reply = justIntroduced ? await L("self_intro", ["name": customerName!])
            : looksLikeFollowUp ? await followUpReply()
            : picked == "not_a_question" ? await L("not_question") : await L("fallback")
        print("📚 小茶 B（FAQ 挑選）：\(reply)  \u{1B}[90m[faq:\(picked)]\u{1B}[0m\n")
        log += "**小茶 B（FAQ 挑選）**（faq:\(picked)）：\(reply)\n\n"
    }
    save()
}
save()
