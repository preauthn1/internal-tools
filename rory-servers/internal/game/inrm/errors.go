package inrm

import "errors"

// ErrGameState 遊戲狀態錯誤
var ErrGameState = errors.New("Game State Error")

// ErrPlayerNotEnough 玩家人數不足
var ErrPlayerNotEnough = errors.New("Player Not Enough")

// ErrPlayerNotFound 查無玩家
var ErrPlayerNotFound = errors.New("Player Not Found")

// ErrPlayerState 玩家狀態錯誤
var ErrPlayerState = errors.New("Player State Error")

// ErrAllPlayerDead 玩家全死了
var ErrAllPlayerDead = errors.New("All Player Dead")

// ErrNoCardSelected 未選擇卡牌
var ErrNoCardSelected = errors.New("No Card Selected")

// ErrCardsLenghtCantPop 當前卡牌數量無法抽牌
var ErrCardsLenghtCantPop = errors.New("Cards Length Can't Take Card")

// ErrCardsLenghtCantThrow 當前卡牌數量無法棄牌
var ErrCardsLenghtCantThrow = errors.New("Cards Length Can't Throw Card")

// ErrHandsDiffDeck 手牌組合與玩家的牌不同
var ErrHandsDiffDeck = errors.New("Illegal Hands. Hands is diffrent with Deck cards")
