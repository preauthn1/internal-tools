import Foundation

// MARK: - 分派：狀態 × 意圖 → 怎麼回（讀寫的 customer 是這次這位顧客，開的待處理事項記在 actions）

extension Turn {
    func handle(_ text: String) async -> (String, String) {
        original = text
        if customer.human { return ("silent", "") }

        // 預約問欄位時的短答（「Bill」「たなか」「김민수」）不切換語言，完整句子才跟著換
        if let lang = detectLanguage(text), customer.askingField == nil || !isShortAnswer(text) { customer.lang = lang }
        if !supportedLanguages.contains(customer.lang) {
            openTicket(.unsupported, "🌐 不支援的語言（\(customer.lang)）：\(text)")
            return ("unsupported", R("unsupported"))
        }

        // 只有表情、標點（沒有文字和數字）不用問 AI
        let symbolsOnly = text.allSatisfy { !$0.isLetter && !$0.isNumber }
        let (intent, name) = symbolsOnly ? ("small_talk", nil) : await classify(text)
        // 第一次知道名字就記住（之後不自動覆蓋，避免把朋友的名字當成顧客）
        var justIntroduced = false
        if customer.name == nil, let name { customer.name = name; justIntroduced = true }

        // 預約中：顧客多半是在回答剛問的欄位；明確是別的事（問問題、追問、聲稱）才跳出去處理
        // （身分類不跳出：正在問名字時只回「小康」，常被分成身分）
        if customer.askingField != nil {
            if intent == "cancel" {
                customer.booking = [:]; customer.askingField = nil
                return ("booking_cancel", R("booking_cancel"))
            }
            if !["store_question", "follow_up", "claim"].contains(intent),
               let r = await bookingStep(text, probe: intent != "booking", intent: intent) { return r }
        }

        if ["store_question", "booking"].contains(intent) { customer.offTopicStreak = 0 }

        switch intent {
        case "booking":
            // 一個欄位都沒抽到：可能只是想開始預約，也可能是在問預約規則（提前幾天、能不能包場）被分成預約。
            // 挑 FAQ 有正好回答的題目就回答，沒有才開始問日期（多的一次 AI 呼叫只發生在沒帶任何欄位時）
            if let r = await bookingStep(text, probe: true, intent: intent) { return r }
            let start = ("start_booking", "顧客要預約、訂位（不是在問預約或包場的規則）")
            if case let picked = await pickFAQ(text, extra: [start]), let faq = faqs.first(where: { $0.id == picked }) {
                return ("answer", faq.answers[customer.lang] ?? faq.answers["zh"]!)
            }
            return bookingNext()
        case "cancel":
            return ("small_talk", R("small_talk"))     // 沒有進行中的預約，當成一般回應
        case "follow_up":
            guard customer.openTicket else { return ("unclear", R("unclear")) }
            customer.followUps += 1
            if customer.followUps == 3 {
                // 引用還在等店員的那張（最近開的、不是緊急通知），店員回覆緊急通知時對應回那張
                let waitingFor = tickets.tickets.last { $0.customer == key && !$0.done && $0.kind != .urgent }
                let quote = waitingFor.map { "\n📌 原本的事項 #\($0.id)：\($0.original.isEmpty ? $0.text : $0.original)" } ?? ""
                openTicket(.urgent, "🔥 緊急：顧客已追問 \(customer.followUps) 次，請盡快回覆" + quote, waiting: false, relatedTo: waitingFor?.id)
            }
            return customer.followUps >= 3 ? ("nudge", R("nudge")) : ("waiting", R("waiting"))
        case "claim", "off_topic":
            customer.offTopicStreak += 1
            if intent == "claim" { openTicket(.claim, "❓ 顧客聲稱／更正：\(text)") }
            if customer.offTopicStreak >= 3 {
                customer.human = true
                openTicket(.handoff, "🙋 轉人工：連續 \(customer.offTopicStreak) 次無關或聲稱的訊息", waiting: false)
                return ("handoff", R("handoff"))
            }
            return (intent, R(intent))
        case "identity":
            if let name = customer.name { return ("identity", R("you_are", ["name": name]) + R("bot_identity")) }
            return ("identity", R("bot_identity") + R("ask_name_soft"))
        case "self_intro":
            guard let name = customer.name else { return ("greeting", R("greeting")) }
            return ("self_intro", R("self_intro", ["name": name]))
        case "greeting" where justIntroduced, "small_talk" where justIntroduced:   // where 要每個 pattern 各寫一次
            return ("self_intro", R("self_intro", ["name": customer.name!]))
        case "greeting", "thanks", "goodbye", "small_talk", "complaint", "unclear":
            return (intent, R(intent))
        default:
            // 這句話介紹了名字（「我是小康」常被分成店家問題，挑 FAQ 又會硬挑一題）：當成自我介紹
            if justIntroduced { return ("self_intro", R("self_intro", ["name": customer.name!])) }
            return await answerStoreQuestion(text)
        }
    }

    // 非文字訊息（圖片、語音、檔案）：看不懂，開待處理事項給店員（說明文字一起放進通知）；
    // 轉人工中照樣通知店員，但不回覆顧客
    func handleMedia(_ media: String, caption: String) -> (String, String) {
        original = caption.isEmpty ? "[\(media)]" : "[\(media)] \(caption)"
        openTicket(.question, "📎 顧客傳了\(media)" + (caption.isEmpty ? "" : "：\(caption)"))
        if customer.human { return ("silent", "") }
        return ("text_only", R("text_only"))
    }
}
