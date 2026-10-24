package poker

import "errors"

// ErrCardNotFound 查無卡牌
var ErrCardNotFound error = errors.New("Card Not Found")

// ErrCardNotEnough 卡牌不足
var ErrCardNotEnough error = errors.New("Card Not Enough")
