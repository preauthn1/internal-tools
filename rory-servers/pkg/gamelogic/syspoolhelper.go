package gamelogic

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// 水池模式
const (
	SysPoolModeNone   int = iota // 0 無水池
	SysPoolModeSingle            // 1 單機不寫檔模式, 重開服務水池就從贏錢線開始
	SysPoolModeDB                // 2 資料庫模式
)

// SysPoolHelper ?
type SysPoolHelper struct {
	glc      *LogicBase
	cfgs     *config.Group
	mode     int
	ctrlRate int   // 水池低於贏錢線開系統有利獎的機率
	pool     int64 // 水池
	winLine  int64 // 贏錢線
}

// NewSysPoolHelper ?
func NewSysPoolHelper(glc *LogicBase, cfgs *config.Group) *SysPoolHelper {
	h := new(SysPoolHelper)
	h.glc = glc
	h.cfgs = cfgs
	if cfgs != nil {
		gameServerConfig := cfgs.GetConfig("gameserver")
		h.mode = gameServerConfig.GetConfigInt("gameserver", "syspoolmode")
		h.ctrlRate = gameServerConfig.GetConfigInt("gameserver", "syspoolctrlrate")
		h.winLine = int64(gameServerConfig.GetConfigInt("gameserver", "winline"))
		h.pool = int64(gameServerConfig.GetConfigInt("gameserver", "winline"))
	} else {
		h.mode = SysPoolModeNone
		h.ctrlRate = 0
		h.winLine = 0
		h.pool = 0
	}

	return h
}

// GetMode ?
func (h *SysPoolHelper) GetMode() int {
	return h.mode
}

// GetPool 取得系統水池
func (h *SysPoolHelper) GetPool() int64 {
	if h.mode == SysPoolModeDB {
		return h.glc.GetSuperToolService().GetSuperToolHandler().GetSysPool()
	} else if h.mode == SysPoolModeSingle {
		return h.pool
	}
	return 0
}

// GetWinLine 取得贏錢線
func (h *SysPoolHelper) GetWinLine() int64 {
	if h.mode == SysPoolModeDB {
		return h.glc.GetSuperToolService().GetSuperToolHandler().GetWinLine()
	} else if h.mode == SysPoolModeSingle {
		return h.winLine
	}
	return 0
}

// UpdateSysPool 更新系統水池
func (h *SysPoolHelper) UpdateSysPool(updatePool int64) {
	if h.mode == SysPoolModeDB {
		h.glc.GetSuperToolService().GetSuperToolHandler().UpdateSysPool(updatePool)
	} else if h.mode == SysPoolModeSingle {
		h.pool += updatePool
	}
}

// IsCtrlSystemWin 判斷是否進行系統有利獎
func (h *SysPoolHelper) IsCtrlSystemWin() bool {
	mlog.System("pool mode %d, pool %d, winline %d, ctrlrate %d", h.mode, h.GetPool(), h.GetWinLine(), h.ctrlRate)
	if h.mode != SysPoolModeNone && h.GetPool() < h.GetWinLine() && util.IsChance(h.ctrlRate, 100) {
		return true
	}
	return false
}

// ToString 顯示水池狀態
func (h *SysPoolHelper) ToString() string {
	ctrlRate := h.ctrlRate
	if h.mode == SysPoolModeNone {
		ctrlRate = 0
	}
	return fmt.Sprintf("水池狀態 mode %d, ctrlRate %d, pool %d, winLine %d", h.mode, ctrlRate, h.GetPool(), h.GetWinLine())
}
