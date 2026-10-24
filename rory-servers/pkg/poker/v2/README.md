# Poker

* [Cards 卡牌](#cards-卡牌)
* [Decks 牌堆](#decks-牌堆)
* [Shuffler 洗牌器](#shuffler-洗牌器)
* [Validator 驗證器](#validator-驗證器)
* [Protobuf](#protobuf)

## Cards 卡牌

以兩位的16進制表示，第一位為花色(Suit)，第二位為點數(Rank)

提供 `Unicode()` 顯示對應的Unicode編碼。
> Ref: https://example.invalid/charts/PDF/U1F0A0.pdf

* 花色(Suit)
  * 黑桃(Spades): `0x10`
  * 紅心(Hearts): `0x20`
  * 方塊(Diamonds): `0x30`
  * 梅花(Clubs): `0x40`
  * 鬼牌(Jokers): `0x50`
* 點數(Rank)
  * A: `0x01`
  * 2: `0x02`
  * 3: `0x03`
  * 4: `0x04`
  * 5: `0x05`
  * 6: `0x06`
  * 7: `0x07`
  * 8: `0x08`
  * 9: `0x09`
  * 10: `0x0A`
  * J: `0x0B`
  * Q: `0x0C`
  * K: `0x0D`

> 提供 `const` 直接使用 Card
> 
> Ex: 
> ``` go
> AceSpades   // 黑桃A
> NineClubs   // 梅花9
> QueenHearts // 方塊Q
> RedJoker    // 紅鬼牌
> ```

|    | 黑桃(Spades) | 紅心(Hearts) | 方塊(Diamonds) | 梅花(Clubs) | 鬼牌(Jokers) |
|:--:|:-----------:|:-----------:|:--------------:|:----------:|:-----------:|
| A  | `0x11`      | `0x21`      | `0x31`         | `0x41`     |             |
| 2  | `0x12`      | `0x22`      | `0x32`         | `0x42`     |             |
| 3  | `0x13`      | `0x23`      | `0x33`         | `0x43`     |             |
| 4  | `0x14`      | `0x24`      | `0x34`         | `0x44`     |             |
| 5  | `0x15`      | `0x25`      | `0x35`         | `0x45`     |             |
| 6  | `0x16`      | `0x26`      | `0x36`         | `0x46`     |             |
| 7  | `0x17`      | `0x27`      | `0x37`         | `0x47`     |             |
| 8  | `0x18`      | `0x28`      | `0x38`         | `0x48`     |             |
| 9  | `0x19`      | `0x29`      | `0x39`         | `0x49`     |             |
| 10 | `0x1A`      | `0x2A`      | `0x3A`         | `0x4A`     |             |
| J  | `0x1B`      | `0x2B`      | `0x3B`         | `0x4B`     |             |
| Q  | `0x1C`      | `0x2C`      | `0x3C`         | `0x4C`     |             |
| K  | `0x1D`      | `0x2D`      | `0x3D`         | `0x4D`     |             |
| 紅 |             |             |                |            | `0x50`      |
| 黑 |             |             |                |            | `0x51`      |
| 白 |             |             |                |            | `0x52`      |


## Decks 牌堆

牌堆就是一坨牌的概念

* 一副牌預設為 [黑桃 紅心 方塊 梅花] 四個花色
* 每個花色有 [A 2 3 4 5 6 7 8 9 10 J Q K]
* 有需要鬼牌的遊戲自行加入

Example:
``` go
// 建立一副牌(52張)
deck := NewDeck(1)

// 建立三副牌(156張)
deck := NewDeck(3)

// 建立空牌堆(沒有卡牌)
deck := NewEmptyDeck()                        // 空牌堆
deck := NewEmptyDeck(KingClubs, AceDiamonds)  // 帶入起始值

cards, err := ParseCards("♦2 ♠4 ♣5 ♠6 ♠7 ♦7 ♠8 ♦8 ♥9 ♦10 ♦Q ♥K ♣K")  // 傳入卡片文字轉換，空白作區隔，回傳 []Card
deck, err := ParseDeck("♦2 ♠4 ♣5 ♠6 ♠7 ♦7 ♠8 ♦8 ♥9 ♦10 ♦Q ♥K ♣K")    // 傳入卡片文字轉換，空白作區隔，回傳 *Deck

// 手動加入 2張紅鬼牌 2張黑鬼牌
deck.AddCard(RedJoker)
deck.AddCard(RedJoker)
deck.AddCards(BlackJoker, BlackJoker)

// 加入 指定 花色 及 點數 的卡牌
deck.Add(Diamonds, Queen)   // 加入方塊Q
deck.Add(Spades, King)      // 加入黑桃K
deck.Add(Diamonds, Ace)     // 加入方塊A
deck.Add(Clubs, King)       // 加入梅花K

// 抽出 指定 花色 及 點數 的卡牌 (抽出成功便從牌堆移除該卡牌)
card, err := deck.Pick(Diamonds, Queen) // 抽出方塊Q
card, err := deck.Pick(Spades, King)    // 抽出黑桃K

// 抽出 指定的卡牌 (抽出成功便從牌堆移除該卡牌)
card, err := deck.PickCard(KingClubs)            // 抽出梅花K
cards := deck.PickCards(KingClubs, AceDiamonds)  // 抽出梅花K與方塊A，沒有抽到也不會怎樣

// 抽出 指定的點數 (抽出成功便從牌堆移除該卡牌)
cards := deck.PickRank(Seven)            // 抽出7點，沒有抽到也不會怎樣

// 抽出 指定的花色 (抽出成功便從牌堆移除該卡牌)
cards := deck.PickSuit(Spades)           // 抽出黑桃，沒有抽到也不會怎樣

// 抽出第一張牌 (抽出成功便從牌堆移除該卡牌)
card, err := deck.Pop()

// 檢查卡牌是否存在
idx, ok := deck.Exist(Diamonds, Ace)
idx, ok := deck.ExistCard(RedJoker)

// 排序
deck.Sort(SortSameRank)         // 同點數排序，A-K
deck.Sort(SortSameRankAceLast)  // 同點數排序，2-A
deck.Sort(SortSameSuit)         // 同花色排序，A-K
deck.Sort(SortSameSuitAceLast)  // 同花色排序，2-A
```

## Shuffler 洗牌器

#### RankShuffler 亂數洗牌器

## Validator 驗證器

提供驗證牌組的工具，實作 `Validate([]Card) bool` 即可，多數牌型皆以驗證器組合即可達成

#### 基本驗證工具
* 固定長度 (`SizeValidator`)
* 最小長度 (`SizeLessValidator`)
* 順子(連續性) (`SequenceValidator`)
* 同花色 (`SameSuitValidator`)
* 同點數 (`SameRankValidator`)

## Protobuf

* 原生檔案: [Game_Poker.proto](http://example.invalid:1688/server-team/origin-protocol/blob/master/Game_Poker.proto)
* 轉換後檔案: [Game_Poker.pb.go](../../netproto/Game_Poker.pb.go)


提供轉換工具
* **TransePokerCard**: 將 Card 轉為 `netproto.Poker_Card`
* **TransePokerDeck**: 將 Deck 轉為 `netproto.Poker_Deck`
* **ReversePokerCard**: 將 `netproto.Poker_Card` 轉為 Card
* **ReversePokerDeck**: 將 `netproto.Poker_Deck` 轉為 Deck
