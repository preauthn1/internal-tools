package gamelogic

import "example.invalid/ServerTeam/Servers/pkg/netproto"

// Fish 共用捕魚資料
type Fish struct {
	SysRewardList     *netproto.SuperToolFishGetSystemRewardDetail  // 系統放分資料
	VIPUserRewardList *netproto.SuperToolFishGetVIPUserRewardDetail // 個人放分資料
}
