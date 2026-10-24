import Foundation

// MARK: - 店家問題：挑 FAQ ／ 程式算價格 ／ 轉店員

public var lastPick: String?   // 最近一次挑 FAQ 的結果（BOT_TESTSET 統計用）

extension Turn {
    // extra：特定情境多給的選項（放在最前面），例如被分成預約時的 start_booking
    func pickFAQ(_ text: String, extra: [(String, String)] = []) async -> String {
        // 試過 FAQ 案例庫（examples.jsonl 加 faq 欄位，正反案例共 66 句，附最像的 6 句當參考）：完整題組每輪答錯 8.0 → 10.0
        // 正面案例會把無關的問題拉過去（「どこにありますか→address」讓「車で…どこに停めれば」變成回地址），已移除
        // 也試過程式比相似度把關：答錯的相似度 0.50～0.96、答對的低到 0.52，分不開，不採用
        // 附上答案內容：只看問法時 AI 會因主題相近而選（大杯、去冰→menu），看到答案才能判斷「答得出這件事嗎」
        let options = extra.map { "\($0.0)：\($0.1)" } + faqs.map { "\($0.id)：\($0.questions)（答案：\($0.answers["zh"]!)）" } + [
            "price_calc：計算買幾杯、多種飲品、自取或外送時總共多少錢",
            "unclear：看不懂顧客在問什麼（太短、沒頭沒尾、指代不明，例如「中」「那個呢」「還有嗎」）",
            "not_a_question：顧客不是在問問題（自我介紹、打招呼、回應、我想一下、要客服扮演別人）",
            "none：以上都無法回答（店家資料沒有這項資訊，例如吸菸區、禮盒、無障礙設施）",
        ]
        // 先逐字抄出顧客問的那件事再選：只看整句時 AI 會被主題帶走（帶環保杯→自取優惠、停車在哪→地址），
        // 抄出來的那件事（環保杯、停車）才是要在選項答案裡找的東西。欄位依序產生，所以 ask 要放在 faq 前面
        let fields: [Field] = [
            .text("ask", "顧客想知道的那件事，從訊息逐字照抄（不要改寫）"),
            .choice("faq", "選項的答案有正好講到上面那件事的代號", extra.map(\.0) + faqs.map(\.id) + ["price_calc", "unclear", "not_a_question", "none"]),
        ]
        // 診斷：AI 只要主題沾到邊就會選（停車→自取優惠、集點→付款、杯子大小→菜單），所以要講清楚「相關不算」
        let rule = "只有當選項的內容正好回答顧客問的那件事才選它；只是主題相關（同樣是飲料、折扣、到店、付款）不算，要選 none。"
        // 試過把分類時的 want（「詢問是否有大杯」）一起給它：反而更亂挑（快速題組每輪答錯 8.0 → 14.0），已移除
        let c = await ask("Pick", text, instructions: "從選項中挑出能回答顧客問題的一個。顧客可能使用任何語言。\(rule)不確定就選 none。\n選項：\n" + options.joined(separator: "\n"), fields: fields)
        let raw = c?.string("faq")
        debug("挑 FAQ：\(raw ?? "（AI 失敗，用預設值 none）")（問的事「\(c?.string("ask") ?? "")」）")
        lastPick = raw ?? "none"
        return raw ?? "none"
    }
}

// 讀訂單：AI 只抽杯數、取貨方式、付款方式和提到的金額；品項必須出現在原文、金額必須是原文的數字，總價由程式算
struct Order { let lines: [String]; let total: Int; let pickup: String; let payment: String; let amount: Int? }

extension Turn {
    func readOrder(_ text: String) async -> Order? {
        let fields: [Field] = store.menu.map { item in
            .integer("qty_\(item.id)", "顧客要買幾杯\(item.names["zh"]!)（\(item.aliases.joined(separator: "、"))）；沒提到就是 0")
        } + [
            // 這個欄位已不採用（取貨方式改由 readPickup 判斷），但拿掉會讓杯數和金額的判讀跟著變（「各一杯」「100塊」讀錯），所以保留
            .choice("pickup", "pickup=店內自取、delivery=外送、unknown=沒說", ["pickup", "delivery", "unknown"]),
            .choice("payment", "card=問能不能刷卡或用信用卡、other=問其他付款方式、unknown=沒問付款", ["card", "other", "unknown"]),
            .integer("amount", "顧客說的消費金額（例如「100塊」= 100）；沒說就是 0"),
        ]
        guard let c = await ask("Order", text, instructions: "從顧客訊息讀出訂單。「各一杯」是提到的每一樣都 1 杯。", fields: fields) else { return nil }
        debug("讀訂單杯數：" + store.menu.map { "\($0.id)=\(c.int("qty_\($0.id)").map(String.init) ?? "nil")" }.joined(separator: " "))

        let lower = text.lowercased()
        var lines: [String] = [], total = 0
        for item in store.menu {
            let qty = c.int("qty_\(item.id)") ?? 0
            // 防止 AI 自己加品項：訊息裡必須真的提到這個品項
            guard qty > 0, qty <= 50, item.aliases.contains(where: { lower.contains($0.lowercased()) }) else { continue }
            total += item.price * qty
            let name = item.names[customer.lang] ?? item.names["zh"]!
            lines.append(customer.lang == "en" ? "\(qty) × \(name)" : customer.lang == "ja" ? "\(name) \(qty)杯" : "\(name) \(qty) 杯")
        }
        // 金額也不能是 AI 編的：必須是訊息裡的數字
        let rawAmount = c.int("amount") ?? 0
        let amount: Int? = rawAmount > 0 && numbers(in: text).contains(rawAmount) ? rawAmount : nil
        let pickup = total > 0 ? await readPickup(text) : "unknown"
        return Order(lines: lines, total: total, pickup: pickup,
                     payment: c.string("payment") ?? "unknown", amount: amount)
    }
}

// 取貨方式另外問一次（不放進讀訂單：讀訂單的欄位一多，杯數和金額的判讀就會跟著變）。
// AI 先逐字抄出講取貨方式的原文再選；原文要真的在訊息裡（忽略標點空白），沒有就當 unknown，不主動提折扣。
// 不用關鍵字清單：多語系會漏（「我自己過去拿」「swing by and grab」）
let pickupFields: [Field] = [
    .text("quote", "訊息中講到取貨方式（自己來拿或外送）的原文，逐字照抄；品項和杯數不算；沒講就是空字串"),
    .choice("pickup", "依上面的原文：pickup=店內自取、delivery=外送；沒有原文就是 unknown", ["pickup", "delivery", "unknown"]),
]

extension Turn {
    func readPickup(_ text: String) async -> String {
        let c = await ask("Pickup", text, instructions: "從顧客訊息找出取貨方式。", fields: pickupFields)
        let quote = (c?.string("quote") ?? "").trimmingCharacters(in: .whitespaces)
        let raw = c?.string("pickup") ?? "unknown"
        let pickup = !normalized(quote).isEmpty && normalized(text).contains(normalized(quote)) ? raw : "unknown"
        debug("取貨方式：\(raw)（出處「\(quote)」）→ \(pickup)")
        return pickup
    }


    // assumePickup：顧客問的就是自取優惠、沒說怎麼取貨時，直接回答自取有沒有打折
    func priceReply(_ o: Order, assumePickup: Bool = false) -> String {
        let d = store.pickupDiscount
        let discounted = String(Int((Double(o.total) * d.rate).rounded()))
        var parts = [R("calc_total", ["items": o.lines.joined(separator: customer.lang == "en" ? ", " : "、"), "total": String(o.total)])]
        switch o.pickup == "unknown" && assumePickup ? "pickup" : o.pickup {
        case "pickup": parts.append(o.total >= d.minTotal ? R("calc_pickup_discount", ["discounted": discounted]) : R("calc_pickup_none"))
        case "delivery": parts.append(R("calc_delivery"))
        default: break   // 店家決定：顧客沒說要自取，不主動提自取折扣或折後價
        }
        return parts.joined(separator: customer.lang == "en" ? " " : "")
    }

    func answerStoreQuestion(_ text: String) async -> (String, String) {
        let picked = await pickFAQ(text)
        let mentionsItem = store.menu.contains { $0.aliases.contains { text.lowercased().contains($0.lowercased()) } }
        // 挑到算價；或挑到單品價格 FAQ 但顧客要的不只一杯（「兩杯抹茶拿鐵多少錢」，單價答非所問）：由程式算總價（訂單只讀一次）
        if picked == "price_calc" || picked.hasPrefix("price_") && mentionsItem, let o = await readOrder(text), o.total > 0,
           picked == "price_calc" || !store.menu.contains(where: { $0.price == o.total }) {
            return ("answer", priceReply(o))
        }
        // 挑到有條件的 FAQ（自取優惠、付款），但顧客講了自己這一筆（品項或金額）：由程式判斷這一筆能不能打折、刷卡，
        // 不只回通則（「兩杯桂花烏龍自取有打折嗎」只回「滿 300 打 9 折」等於沒回答）
        // 外送 FAQ 只說「有外送」，顧客講了自己這一筆（外送幾杯某品項能不能打折）時同樣由程式算
        if picked == "pickup_discount" || picked == "delivery", mentionsItem, let o = await readOrder(text), o.total > 0 {
            return ("answer", priceReply(o, assumePickup: picked == "pickup_discount"))
        }
        if picked == "payment", mentionsItem || !numbers(in: text).isEmpty, let o = await readOrder(text), o.payment == "card",
           let amount = o.total > 0 ? o.total : o.amount {
            return ("answer", R(amount >= store.cardMinimum ? "calc_card_ok" : "calc_card_no", ["total": String(amount), "card": String(store.cardMinimum)]))
        }
        if picked == "unclear" { return ("unclear", R("unclear")) }
        if picked == "not_a_question" { return ("small_talk", R("small_talk")) }
        if let faq = faqs.first(where: { $0.id == picked }) { return ("answer", faq.answers[customer.lang] ?? faq.answers["zh"]!) }
        // 要轉店員之前，確認顧客真的在問問題（「好」「我想一下」不該變成待處理事項）
        // 只在 FAQ 挑不到時才檢查，所以不會擋掉答得出來的問題
        if !(await isQuestion(text)) { return ("small_talk", R("small_talk")) }
        openTicket(.question, "📝 待補問題（知識庫沒有）：\(text)")
        return ("fallback", R("fallback"))
    }
}

let messageKindFields: [Field] = [.choice("kind",
    "question=顧客在問問題或要求資訊（含沒有問號的問法，例如「有幾個老闆」「營業時間」）；pause=顧客在思考、請客服稍等（我想一下、等等、我再看看、let me think）；reply=只是回應（OK、好、收到）",
    ["question", "pause", "reply"])]

extension Turn {
    func isQuestion(_ text: String) async -> Bool {
        let c = await ask("Kind", text, instructions: "判斷顧客訊息屬於哪一種。顧客可能使用任何語言。", fields: messageKindFields)
        let raw = c?.string("kind")
        debug("是不是在發問：\(raw ?? "（AI 失敗，用預設值 question）")")
        return (raw ?? "question") == "question"
    }
}
