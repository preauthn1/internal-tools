package dialog

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// MARK: 預約（程式控制流程，AI 只抽欄位）。照搬 Swift 版的 Booking.swift；
// 不檢查時段、公休日和人數，一律由員工確認或拒絕（ADR 0003）

var bookingOrder = []string{"date", "time", "people", "name", "phone"}

// problemFields：問題代碼是哪個欄位的問題
var problemFields = map[string]string{
	"err_date": "date", "err_past": "date", "err_time": "time", "ask_time_midnight": "time", "err_people": "people", "err_phone": "phone",
}

var fieldNames = map[string]string{"date": "日期", "time": "時間", "people": "人數", "name": "預約人姓名", "phone": "電話"}

// 日期、時間也用整數：月、日、時、分由程式組成，並和原文出處的數字比對；相對日期（明天、下週二）由程式換算
var bookingFields = []model.Field{
	{Kind: model.KindText, Name: "date_quote", Description: "訊息中提到日期的原文，逐字照抄；未提及則為空字串"},
	{Kind: model.KindInteger, Name: "date_month", Description: "預約日期的月；有寫月日就照填；只說相對日期（明天、下週二…）或未提及則為 0"},
	{Kind: model.KindInteger, Name: "date_day", Description: "預約日期的日；有寫月日就照填；只說相對日期（明天、下週二…）或未提及則為 0"},
	{Kind: model.KindText, Name: "time_quote", Description: "訊息中提到時間的原文，逐字照抄；未提及則為空字串"},
	{Kind: model.KindInteger, Name: "time_hour", Description: "時，24 小時制（下午兩點=14）；未提及則為 0"},
	{Kind: model.KindInteger, Name: "time_minute", Description: "分（兩點半=30）；未提及則為 0"},
	// 人數用整數型別：用文字欄位時 AI 偶爾填成「{10}」「)」（guided generation 會強制只能輸出數字）
	{Kind: model.KindInteger, Name: "people", Description: "人數；未提及則為 0"},
	{Kind: model.KindText, Name: "people_quote", Description: "訊息中提到人數的原文，逐字照抄；未提及則為空字串"},
	{Kind: model.KindText, Name: "name", Description: "預約人姓名；未提及則為空字串"},
	{Kind: model.KindText, Name: "phone", Description: "電話號碼；未提及則為空字串"},
}

// bookingSpec 是依商家設定的預約欄位（#67）組出來的欄位清單；沒設定時就是上面的 bookingOrder、bookingFields、fieldNames
type bookingSpec struct {
	order  []string          // 要問的欄位，依序：日期、時間、（人數）、必填的自訂欄位、姓名、電話
	keys   []string          // 要抽的欄位：order 加上選填的自訂欄位（選填的不追問，顧客有講才記下來）
	fields []model.Field     // 給 AI 抽取的欄位
	names  map[string]string // 欄位名稱（提示詞、轉人工的原因、追問的預設句子）
	custom []string          // 自訂欄位的代號，照設定的順序
	people bool              // 有沒有問人數
}

func bookingSpecOf(cfg merchant.Config) bookingSpec {
	bf := merchant.EffectiveBookingFields(cfg)
	if bf.PeopleEnabled && len(bf.Custom) == 0 {
		return bookingSpec{order: bookingOrder, keys: bookingOrder, fields: bookingFields, names: fieldNames, people: true}
	}
	s := bookingSpec{order: []string{"date", "time"}, names: maps.Clone(fieldNames), people: bf.PeopleEnabled}
	if s.people {
		s.order = append(s.order, "people")
	}
	for _, f := range bf.Custom {
		if f.Required {
			s.order = append(s.order, f.Key)
		}
	}
	s.order = append(s.order, "name", "phone")
	s.keys = slices.Clone(s.order)
	for _, f := range bookingFields {
		if s.people || (f.Name != "people" && f.Name != "people_quote") {
			s.fields = append(s.fields, f)
		}
	}
	for _, f := range bf.Custom {
		if !f.Required {
			s.keys = append(s.keys, f.Key)
		}
		s.custom = append(s.custom, f.Key)
		s.names[f.Key] = f.Name
		s.fields = append(s.fields, model.Field{Kind: model.KindText, Name: f.Key,
			Description: f.Name + "：" + strings.TrimSpace(f.Description) + "；照原文抄寫，未提及則為空字串"})
	}
	return s
}

func (s bookingSpec) isCustom(key string) bool { return slices.Contains(s.custom, key) }

var (
	digitsPattern = regexp.MustCompile(`\d+`)
	// 時間不看日期、電話的數字（「10/19 下午兩點」的 10、19 不是時間）
	dateOrPhoneDigits = regexp.MustCompile(`\d+\s*[/\-.月]\s*\d+\s*(日|號|号)?|\d{6,}`)
	// 下午、晚上的說法（時間要換成 24 小時制）
	afternoon = regexp.MustCompile(`下午|晚上|今晚|明晚|傍晚|午後|夕方|夜|\d\s*pm|\bp\.m\.|\bpm\b|afternoon|evening|tonight`)
	// 早上、上午的說法（AI 把「早上十點」填成 22 時換回上午）
	morning = regexp.MustCompile(`早上|上午|早晨|清晨|午前|\d\s*am|\ba\.m\.|\bam\b|morning`)
	// 中午的說法（#49：「中午一點」是 13 點，1～5 點才加 12；中午十一點、十二點照字面）
	noon = regexp.MustCompile(`中午|正午|お昼|昼|\bnoon\b`)
	// 半夜的一到五點算上午（#49：「半夜三點」有「夜」曾被加 12 存成 15:00；深夜十一點照樣是 23 點）
	lateNight = regexp.MustCompile(`凌晨|半夜|深夜|午夜`)
	// 晚上、半夜的十二點是跨日寫法（#49：「晚上十二點」曾存成中午 12:00），日期可能差一天，不收、再問一次。
	// 時段要緊接在「12點」前面（中文數字先換算），「今晚12個人」「中午12點，晚上不行」不算
	crossDay = regexp.MustCompile(`(?:晚|夜|凌晨)\D{0,3}12\s*[點点時时]|12\s*a\.?m\b|midnight`)
	// 明寫的 0 點（00:30）
	zeroClock = regexp.MustCompile(`(?:^|\D)0?0:\d{2}`)
	// 時間出處裡的「7點」「3點半」「10點15分」（中文數字先換算）
	clockTime = regexp.MustCompile(`(\d{1,2})\s*[點点時时]\s*(?:(半)|(\d{1,2})\s*分)?`)
	// 換算中文數字前把星期和後面的時間分開（「週五七點」不能變成 57 點）
	weekdayWord = regexp.MustCompile(`((?:週|周|星期|禮拜|礼拜)[一二三四五六日天])`)
	// 人數不看日期、時間、電話的數字（「晚上7點」的 7 曾被當成人數）
	notPeopleDigits = regexp.MustCompile(`\d+\s*[/:\-.月]\s*\d+|\d{6,}|\d+\s*(?:點|点|小?時|小?时|號|号|日|分|(?i:am|pm))`)
	// 大人小孩分開講（2大1小、2 adults 1 kid）
	adultsAndKids = regexp.MustCompile(`(\d+)\s*(?:位|個)?\s*(?:大人|大|成人|adults?)\D{0,4}?(\d+)\s*(?:位|個)?\s*(?:小孩|小朋友|小|兒童|kids?|child(?:ren)?)`)
	// 數字寫在後面（#53：大人兩位小孩一位、大人2小孩1）
	adultsAndKidsWordFirst = regexp.MustCompile(`(?:大人|成人)\s*(\d+)\s*(?:位|個)?\D{0,4}?(?:小孩|小朋友|兒童)\s*(\d+)`)
)

// peopleSource：拿來找人數的文字（中文數字換算，拿掉日期、時間、電話的數字）
func peopleSource(text string) string {
	return notPeopleDigits.ReplaceAllString(strings.ToLower(chineseNumerals(text, true)), " ")
}

// adultsAndKidsSums：大人小孩分開講的每一組加總（2大1小、大人兩位小孩一位）
func adultsAndKidsSums(q string) []int {
	var out []int
	for _, re := range []*regexp.Regexp{adultsAndKids, adultsAndKidsWordFirst} {
		for _, m := range re.FindAllStringSubmatch(q, -1) {
			adults, _ := strconv.Atoi(m[1])
			kids, _ := strconv.Atoi(m[2])
			out = append(out, adults+kids)
		}
	}
	return out
}

// peopleNumbers：人數出處裡的數字（中文數字由程式換算，不看日期、時間、電話）；大人小孩分開講時只算加總
func peopleNumbers(quote string) map[int]bool {
	q := peopleSource(quote)
	out := map[int]bool{}
	for _, n := range adultsAndKidsSums(q) {
		out[n] = true
	}
	if len(out) > 0 {
		return out
	}
	return numbers(q)
}

// peopleFromText：整句寫了一組大人小孩（2大1小、兩大一小）時的加總；沒寫或寫了好幾組是 0
func peopleFromText(text string) int {
	if sums := adultsAndKidsSums(peopleSource(text)); len(sums) == 1 {
		return sums[0]
	}
	return 0
}

// hasTimeWording：有時間的寫法（阿拉伯數字、中文數字；7pm 也有數字），日期、電話的數字不算；
// 沒有數字的 am/pm 不算（「I am John」）
func hasTimeWording(text string) bool {
	return len(numbers(dateOrPhoneDigits.ReplaceAllString(chineseNumerals(text, false), " "))) > 0
}

func numbers(text string) map[int]bool {
	out := map[int]bool{}
	for _, s := range digitsPattern.FindAllString(text, -1) {
		if n, err := strconv.Atoi(s); err == nil {
			out[n] = true
		}
	}
	return out
}

func onlyDigits(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
}

func hasLetter(s string) bool { return strings.ContainsFunc(s, unicode.IsLetter) }

// numbersSupport：人數必須在訊息的數字裡找得到依據（日期、時間、電話的數字不能拿來當人數）
func numbersSupport(value, text string) bool {
	n, err := strconv.Atoi(onlyDigits(value))
	return err == nil && numbers(notPeopleDigits.ReplaceAllString(text, " "))[n]
}

const dayFormat = "2006-01-02"

// validate 檢查並整理欄位的值；有問題回傳問題代碼（固定回覆的 key）
func validate(key, value string, today time.Time) (string, string) {
	switch key {
	case "date":
		d, err := time.ParseInLocation(dayFormat, value, today.Location())
		if err != nil {
			return "err_date", value
		}
		if d.Before(startOfDay(today)) {
			return "err_past", value
		}
	case "time":
		parts := strings.Split(value, ":")
		h, err := strconv.Atoi(parts[0])
		if err != nil {
			return "err_time", value
		}
		m := 0
		if len(parts) > 1 {
			m, _ = strconv.Atoi(parts[1])
		}
		return "", fmt.Sprintf("%02d:%02d", h, m)
	case "people":
		n, err := strconv.Atoi(onlyDigits(value))
		if err != nil || n <= 0 {
			return "err_people", value
		}
		return "", strconv.Itoa(n)
	case "phone":
		if d := onlyDigits(value); len(d) >= 9 && len(d) <= 10 {
			return "", d
		}
		return "err_phone", value
	}
	return "", value
}

// bookingStep 從這句話抽預約欄位，回傳下一步怎麼回。
// probe：顧客在預約途中說了別的話，沒抽到新欄位就回傳 false，交回一般流程。intent：這句的分類，決定問姓名時能不能把短答整句當名字
func (t *turn) bookingStep(text string, probe bool, intent string) (Reply, bool) {
	c := &t.customer
	spec := bookingSpecOf(t.cfg)
	hint := ""
	if c.AskingField != "" {
		hint = "客服剛剛問了顧客的" + spec.names[c.AskingField] + "，這句話很可能就是在回答它。"
	}
	r, _ := t.bot.Model.Respond(t.ctx, model.Request{
		Name: "Booking", Text: text, Fields: spec.fields,
		Instructions: "從顧客訊息抽取預約資訊。今天是 " + t.now.Format(dayFormat) + "。" + hint +
			"只說相對日期（明天、下週二…）時照抄原文，月、日填 0，程式會換算；有寫月日就照填。顧客可能使用任何語言。",
	})
	str := func(k string) string {
		if r == nil {
			return ""
		}
		s, _ := r.String(k)
		return strings.TrimSpace(s)
	}
	num := func(k string) int {
		if r == nil {
			return 0
		}
		n, _ := r.Int(k)
		return n
	}
	quote := func(k string) string { return str(k + "_quote") }
	// 出處要逐字出現在訊息裡，而且不能和其他欄位的出處重疊（日期和時間常寫在一起「10/19 下午兩點」，彼此重疊不算）
	quoted := func(k string) bool {
		q := quote(k)
		others := []string{"people"}
		if k == "people" {
			others = []string{"date", "time"}
		}
		if q == "" || !strings.Contains(text, q) {
			return false
		}
		for _, o := range others {
			if oq := quote(o); oq != "" && (strings.Contains(oq, q) || strings.Contains(q, oq)) {
				return false
			}
		}
		return true
	}
	// 月日時分要和原文一致（每組至少一個數字出現）：有出處看出處，出處沒有數字不檢查；
	// 沒有出處就看整句的數字，剛問的欄位整句沒數字也採用（時間除外：要有時間的寫法），其他欄位一定要找到（防止 AI 自己編）
	supported := func(k string, needs [][]int) bool {
		source := text
		if quoted(k) {
			source = quote(k)
			// 時間出處的中文數字也要比對（#46：「晚上七點」AI 填 6，沒檢查就存成 18:00）
			if k == "time" {
				source = chineseNumerals(source, false)
			}
		}
		if k == "time" {
			source = dateOrPhoneDigits.ReplaceAllString(source, " ")
		}
		n := numbers(source)
		if len(n) == 0 {
			// 問時間時回的話沒有時間的寫法（#49：報名字），AI 自己填的時不收
			if k == "time" && k == c.AskingField && !hasTimeWording(text) {
				return false
			}
			return quoted(k) || k == c.AskingField
		}
		for _, group := range needs {
			if !slices.ContainsFunc(group, func(x int) bool { return n[x] }) {
				return false
			}
		}
		return true
	}
	date := func() string {
		// 出處不在原句裡（AI 改寫或沒抄）就看整句
		q := quote("date")
		inText := q != "" && strings.Contains(text, q)
		written := explicitMonthDays(text)
		// 原句是換算不了的說法（下下週三）又沒寫月日：AI 漏抄成「下週三」也不填
		if len(written) == 0 && inText && partOfUnresolved(q, text) {
			return ""
		}
		// 相對說法由程式依今天換算，AI 填的月日不用（AI 曾把週三說的「下週二」填成週六）；換算不了就不填，改問顧客
		source := text
		if inText {
			source = q
		}
		switch d, kind := readDate(source, t.now); kind {
		case dateResolved:
			return d.Format(dayFormat)
		case dateUnresolved:
			return ""
		}
		m, d := num("date_month"), num("date_day")
		// 寫了月日就由程式讀出來，AI 填的月日不用（AI 曾把「10/19」的月日填成 0，明確的日期被丟掉）；
		// 寫了好幾個日期時不自己猜，AI 填的要是其中之一
		candidates := written
		if inText {
			if fromQuote := explicitMonthDays(q); len(fromQuote) > 0 {
				candidates = fromQuote
			}
		}
		if len(candidates) == 1 {
			return absoluteDate(candidates[0][0], candidates[0][1], t.now)
		}
		if len(candidates) > 0 {
			if containsMD(candidates, [2]int{m, d}) {
				return absoluteDate(m, d, t.now)
			}
			return ""
		}
		// 出處沒寫日期（AI 抄錯段落），整句是相對說法就用整句換算
		if inText {
			switch day, kind := readDate(text, t.now); kind {
			case dateResolved:
				return day.Format(dayFormat)
			case dateUnresolved:
				return ""
			}
		}
		// 讀不到月日（「十月十九日」「19號」）才用 AI 填的整數，並和出處的數字比對
		if m < 1 || m > 12 || d < 1 || d > 31 || !supported("date", [][]int{{m}, {d}}) {
			return ""
		}
		return absoluteDate(m, d, t.now)
	}
	midnight, peopleFromWording := false, false
	timeOf := func() string {
		h, m := num("time_hour"), num("time_minute")
		source, period := text, text
		// 出處寫了一個「N點」（整句也只有一個）就由程式讀出來，AI 填的時分不用（#46：AI 把「晚上七點」只抄成「七點」、填 18，存成 18:00）；
		// AI 填的時和出處只差 12（沒寫下午的「兩點半」填 14）時用 AI 的
		withDigits := func(s string) string { return chineseNumerals(weekdayWord.ReplaceAllString(s, "$1 "), true) }
		var clock [][]string
		var restDigits map[int]bool
		if quoted("time") {
			source = quote("time")
			converted := withDigits(source)
			clock, restDigits = clockTime.FindAllStringSubmatch(converted, -1), numbers(clockTime.ReplaceAllString(converted, " "))
			if len(clockTime.FindAllString(withDigits(text), -1)) > 1 {
				clock = nil
			}
			// 下午、晚上連出處前面四個字（不算空白）一起看（AI 漏抄「晚上」，或寫成「晚上 大概七點」）
			before := []rune(strings.Join(strings.Fields(text[:strings.Index(text, source)]), ""))
			period = string(before[max(len(before)-4, 0):]) + source
		}
		if len(clock) == 1 {
			if quotedHour, _ := strconv.Atoi(clock[0][1]); h == 0 || quotedHour == 0 || quotedHour%12 != h%12 {
				h = quotedHour
			}
			// 「七點十五」沒寫分：AI 填的分要在出處「N點」以外的數字裡
			if clock[0][2] != "" {
				m = 30
			} else if clock[0][3] != "" {
				m, _ = strconv.Atoi(clock[0][3])
			} else if !restDigits[m] {
				m = 0
			}
		}
		// 晚上、半夜的十二點（AI 填 12 或 0）：存 00:00 日期可能差一天（明天晚上十二點是明天還是後天的 0 點），改問顧客
		lower := strings.ToLower(period)
		if h%12 == 0 && crossDay.MatchString(strings.ToLower(withDigits(period))) {
			midnight = true
			return ""
		}
		// 時填 0 通常是沒提到；原文明寫「0點」「00:30」才收 0 點
		zero := h == 0 && (len(clock) == 1 || zeroClock.MatchString(source))
		if (h < 1 && !zero) || h > 23 || m < 0 || m > 59 {
			return ""
		}
		// 讀不出來（7pm、7:30、寫了好幾個時間）才用 AI 填的時分，並和原文的數字比對；「2點半」的分沒有數字
		half := m == 30 && strings.Contains(source, "半")
		needs := [][]int{{h, h % 12}}
		if m != 0 && !half {
			needs = append(needs, []int{m})
		}
		if len(clock) != 1 && !supported("time", needs) {
			return ""
		}
		// AI 偶爾把「下午兩點」填成 2：原文寫了下午、晚上、pm 就換成 24 小時制。
		// （Swift 版用營業時間猜「兩點多半是下午」，ADR 0003 拿掉了營業時間，改看原文）
		pm, am := afternoon.MatchString(lower), morning.MatchString(lower)
		if lateNight.MatchString(period) && h%12 >= 1 && h%12 <= 5 {
			h %= 12
		} else if noon.MatchString(lower) && !pm && !am {
			if h >= 1 && h <= 5 {
				h += 12
			}
		} else if h < 12 && pm {
			h += 12
		} else if h > 12 && am && !pm {
			h -= 12
		}
		return fmt.Sprintf("%d:%02d", h, m)
	}
	field := func(k string) string {
		switch k {
		case "date":
			return date()
		case "time":
			return timeOf()
		case "people":
			// 大人小孩分開講由程式加總，AI 填的人數不用（#53：「2大1小」AI 只填了 2，或沒填）
			if n := peopleFromText(text); n > 0 {
				peopleFromWording = true
				return strconv.Itoa(n)
			}
			if n := num(k); n > 0 {
				return strconv.Itoa(n)
			}
			return ""
		case "name":
			raw := str(k)
			// 剛問姓名、AI 抽的名字不能用（空的或不在原文裡）時，短答整句就是名字
			// （AI 曾把「たなか」改寫掉，沒通過原文檢查，預約一直卡在問姓名）
			usable := raw != "" && containsFold(text, raw) && hasLetter(raw)
			// 只在分類成預約或自我介紹時才這樣做，而且整句都是漢字時不用（AI 會照抄中文名字）：
			// 「不知道」「好」曾被收成名字（Swift 版 #17 第二版驗收）
			if !usable && k == c.AskingField && (intent == "booking" || intent == "self_intro") {
				if name := nameFromShortAnswer(text); name != "" && !allHan(name) {
					return name
				}
			}
			return raw
		}
		if spec.isCustom(k) {
			// 自訂欄位要出現在原文裡（防止 AI 自己編）；剛問這一欄、AI 抽的不能用時，短答整句就是答案（「剪髮加護髮」）
			if v := str(k); v != "" && containsFold(text, v) {
				return v
			}
			if k == c.AskingField && isShortAnswer(text) && !askingBack(text) {
				return trimPunct(text)
			}
			return ""
		}
		return str(k)
	}

	var problems []string
	found := false
	for _, key := range spec.keys {
		raw := field(key)
		if key == "time" && midnight {
			found = true
			problems = append(problems, "ask_time_midnight")
			continue
		}
		if raw == "" || raw == "0" || c.Booking[key] == raw {
			continue
		}
		switch key {
		case "people":
			if peopleFromWording {
				break // 程式從原句讀出來的，不用再比對出處
			}
			q := quote(key)
			// 值要和原文出處一致：出處抄對「25 個人」、值卻填 10（取自日期 10/21）曾被收下。
			// 出處有數字時（「四個人」的中文數字由程式換算，「2大1小」加總），人數必須是其中之一；不一致就當作沒抓到，改成再問一次
			n, _ := strconv.Atoi(raw)
			qn := peopleNumbers(q)
			if strings.Contains(text, q) && len(qn) > 0 && !qn[n] {
				continue
			}
			// 出處讀不出人數（#42：問姓名那輪 AI 抄成「王小明」、人數填 1，蓋掉原本的 4 人）就當作沒有出處
			quoteHasPeople := quoted(key) && len(qn) > 0
			// 剛問的欄位直接採用；其他欄位要有原文出處或數字依據（防止 AI 自己編）。
			// 已經有人數時只看出處：整句碰巧有同一個數字（「晚1小時到」）不能蓋掉
			if key != c.AskingField && !quoteHasPeople && (c.Booking[key] != "" || !numbersSupport(raw, text)) {
				continue
			}
		case "phone":
			d := onlyDigits(raw)
			if d == "" || !strings.Contains(onlyDigits(text), d) {
				continue
			}
		case "name":
			// 姓名要出現在原文裡，而且要有文字（電話號碼曾被收成姓名）；角色稱呼（顧客、customer）不是姓名
			if !containsFold(text, raw) || !hasLetter(raw) || isRoleWord(raw) {
				continue
			}
		}
		found = true
		problem, cleaned := validate(key, raw, t.now)
		if problem != "" {
			problems = append(problems, problem)
			continue
		}
		if c.Booking == nil {
			c.Booking = map[string]string{}
		}
		c.Booking[key] = cleaned
	}
	if probe && !found {
		return Reply{}, false
	}
	return t.bookingNext(problems), true
}

// bookingNext 依目前收到的欄位決定下一步：回報問題、問下一個欄位，或收齊後送出預約申請
func (t *turn) bookingNext(problems []string) Reply {
	c := &t.customer
	if c.Booking == nil {
		c.Booking = map[string]string{}
	}
	if c.Booking["name"] == "" && c.Name != "" {
		c.Booking["name"] = c.Name
	}
	spec := bookingSpecOf(t.cfg)
	missing := ""
	for _, k := range spec.order {
		if c.Booking[k] == "" {
			missing = k
			break
		}
	}
	if len(problems) > 0 {
		// 回報問題等於再問一次那個欄位，一樣計次（#53）
		field := problemFields[problems[0]]
		if field == "" {
			field = missing
		}
		return t.askField(field, problems[0])
	}
	if missing != "" {
		return t.askField(missing, "ask_"+missing)
	}
	b := c.Booking
	// 摘要：日期 時間, 人數, 自訂欄位（有值的，照設定的順序）, 姓名, 電話；員工群組的通知、確認和拒絕的回覆都用它
	parts := []string{b["date"] + " " + b["time"]}
	people := ""
	if spec.people {
		people = b["people"]
		parts = append(parts, people+t.text("people_unit", nil))
	}
	var extra map[string]string
	for _, k := range spec.custom {
		if v := b[k]; v != "" {
			parts = append(parts, v)
			if extra == nil {
				extra = map[string]string{}
			}
			extra[k] = v
		}
	}
	summary := strings.Join(append(parts, b["name"], b["phone"]), ", ")
	// 預約申請留著等員工確認或拒絕（編號和待處理事項相同）
	id := t.openTicket("booking", "📅 預約申請："+summary, false, 0)
	if id != 0 {
		err := t.bot.DB.AddBooking(t.ctx, store.Booking{
			ID: id, Platform: t.msg.Platform, UserID: t.msg.UserID, ChatID: t.msg.ChatID,
			Date: b["date"], Time: b["time"], People: people, Name: b["name"], Phone: b["phone"], Extra: extra, Summary: summary,
		}, t.now)
		if err != nil && t.err == nil {
			t.err = err
		}
	}
	c.endBooking()
	return Reply{Key: "booking_done", Text: t.text("booking_done", map[string]string{"summary": summary})}
}

// 同一個欄位最多問幾次（#53）：問了這麼多次還拿不到就轉人工，不再重問
const maxFieldAsks = 2

// askField 問預約的 field（回固定回覆 key）：換了欄位就重新計次；同一個欄位已經問了 maxFieldAsks 次就轉人工
func (t *turn) askField(field, key string) Reply {
	c := &t.customer
	if field != c.AskingField {
		c.AskingField, c.AskCount = field, 1
		return t.reply(key)
	}
	// 這一輪有收到新的預約資料（例如問人數時報了姓名、電話）不算一次：顧客在配合，只是沒先回這一題（aichatbot!71）
	if filledFields(c.Booking) > t.filled {
		return t.reply(key)
	}
	if c.AskCount >= maxFieldAsks {
		return t.bookingHandoff(field)
	}
	c.AskCount++
	return t.reply(key)
}

// filledFields：已經收到幾個預約欄位
func filledFields(b map[string]string) int {
	n := 0
	for _, v := range b {
		if v != "" {
			n++
		}
	}
	return n
}

// bookingHandoff：預約的 field 一直拿不到，轉人工；原因附上已收到的預約資料讓員工接手，預約狀態清掉
func (t *turn) bookingHandoff(field string) Reply {
	c := &t.customer
	spec := bookingSpecOf(t.cfg)
	var got []string
	for _, k := range spec.keys {
		if v := c.Booking[k]; v != "" {
			got = append(got, spec.names[k]+" "+v)
		}
	}
	reason := "預約資料收不齊（" + spec.names[field] + "）"
	if len(got) > 0 {
		reason += "：" + strings.Join(got, "、")
	}
	c.endBooking()
	return t.handoff(reason)
}

// isShortAnswer：很短的回答（名字、電話、人數）。預約問欄位時不依它切換語言（Swift 版 #17：中文預約問到姓名回「Bill」，之後都變成英文）。
// 用空白分詞的文字（英文…）不超過 2 個詞；中日韓（不用空白）去掉標點後不超過 6 個字；只有數字、符號也算短
func isShortAnswer(text string) bool {
	if strings.ContainsFunc(text, func(r rune) bool {
		return unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana, unicode.Hangul)
	}) {
		n := 0
		for _, r := range text {
			if unicode.IsLetter(r) || unicode.IsNumber(r) {
				n++
			}
		}
		return n <= 6
	}
	words := 0
	for _, w := range strings.FieldsFunc(text, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r) }) {
		if hasLetter(w) {
			words++
		}
	}
	return words <= 2
}

// nameFromShortAnswer：問姓名時的短答整句就是名字（去掉頭尾標點和「我是」「です」這類前後綴）；有數字或是完整句子就不算
func nameFromShortAnswer(text string) string {
	name := trimPunct(text)
	for _, prefix := range []string{"我是", "我叫", "叫我", "My name is ", "I'm ", "I am "} {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			name = name[len(prefix):]
			break
		}
	}
	for _, suffix := range []string{"と申します", "と言います", "です"} {
		if s, ok := strings.CutSuffix(name, suffix); ok {
			name = s
			break
		}
	}
	name = strings.TrimSpace(name)
	if name == "" || !isShortAnswer(text) || !hasLetter(name) || strings.ContainsFunc(name, unicode.IsNumber) {
		return ""
	}
	return name
}

// trimPunct 去掉頭尾的空白、標點、符號
func trimPunct(text string) string {
	return strings.TrimFunc(text, func(r rune) bool { return unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r) })
}

// askingBackPattern：反問的說法（「多少錢？」「有哪些」），問自訂欄位時這種短句不是答案
var askingBackPattern = regexp.MustCompile(`[?？]|嗎|多少|幾|哪|什麼|甚麼|(?i:\b(?:how|what|which)\b)`)

func askingBack(text string) bool { return askingBackPattern.MatchString(text) }

func allHan(s string) bool {
	for _, r := range s {
		if !unicode.Is(unicode.Han, r) {
			return false
		}
	}
	return true
}
