package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime/debug"
	"sort"
	"time"

	"example.invalid/gogo/protobuf/proto"
	"example.invalid/spf13/viper"

	"example.invalid/ServerTeam/Servers/cmd/slot88/auth/service"
	"example.invalid/ServerTeam/Servers/cmd/slot88/forward/internal/network"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// Logic 中繼服務器
type Logic struct {
	Name string

	sw   network.ServerWriter
	cw   network.ClientWriter
	user *service.PlayerInfo

	MoneyTransfer MoneyTransfer

	ServeState    ServeState
	GameState     slot.GameState
	LastGameState slot.GameState

	Round int

	UserData          *netproto.UserRoomLogonData
	Setting           *netproto.Slot_MSG_Setting
	GameResult        *netproto.Slot_MSG_GameResult
	StepResult        *netproto.Slot_MSG_StepGameResult
	LeftFreeGameData  *netproto.Slot_MSG_LeftFreeGameData
	PickGameRecord    *netproto.Slot_MSG_BonusGame
	JackpotGameRecord *netproto.Slot_MSG_BonusGame
	JackpotPool       map[slot.JackpotLevel]*netproto.Slot_JackpotInfo
	SelectGameOptions *netproto.Slot_MSG_ChooseGameOptions

	laststat     *netproto.UserStateChange
	lastspin     *netproto.Slot_MSG_Spin
	startedsend  bool
	lastWin      float64
	incfree      int
	freedatasend bool
	pickcnt      int
	leftshow     int
}

// New logic
func New(ctx context.Context, sw network.ServerWriter, cw network.ClientWriter, user *service.PlayerInfo) *Logic {
	logic := new(Logic)

	logic.Name = "Logic"

	logic.sw = sw
	logic.cw = cw
	logic.user = user

	switch viper.GetString("money.transfor") {
	default:
		logic.MoneyTransfer = new(ScoreTransfer)
	case "score":
		logic.MoneyTransfer = new(ScoreTransfer)
	case "dollar":
		logic.MoneyTransfer = new(DollarTransfer)
	}

	logic.ServeState = ServeStateUnlogin
	logic.GameState = slot.GameStateWait

	logic.JackpotPool = make(map[slot.JackpotLevel]*netproto.Slot_JackpotInfo)

	// LoginRoom
	{
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_LoginRoomID)

		login := new(netproto.LoginGameRoomInfo)
		login.HDType = proto.Int32(0)
		login.HDCode = proto.String("user.UserHDCode")
		login.UserID = proto.Int32(user.PlayerID)
		login.Cer = proto.String(user.AccessToken)
		m.MsgData = login

		logic.sw.WriteServer(m)
	}

	if viper.GetBool("debug") {
		go func() {
			for {
				select {
				default:
					logic.printStatus(os.Stdout)
					time.Sleep(5 * time.Second)
				case <-ctx.Done():
					break
				}
			}
		}()
	}

	return logic
}

// OnMessage client message
func (logic *Logic) OnMessage(m *network.Message) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("OnMessage Error, %v\n%s", rec, debug.Stack())
		}
	}()

	if m.ID == 0 && m.ActionID == 1 {
		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Common)
		_m.SClassID = int32(netproto.PlatformCommonClassID_HeartBeatID)
		logic.sw.WriteServer(_m)
		return
	}

	mlog.Debug("OnMessage %d:%d %s", m.ID, m.ActionID, string(m.Data))

	if m.ID == 2 && m.ActionID == 6 {

		if logic.GameState != slot.GameStateBet && logic.GameState != slot.GameStateUserAction {
			mlog.Warn("User Bet Error, GameState not in Bet or UserAction")
			return
		}

		type SpinData struct {
			Param1 string `json:"param1"`
			Param2 string `json:"param2"`
			Param3 string `json:"param3"`

			LineCount  float64 `json:"-"`
			BetMoney   float64 `json:"-"`
			Multiplier float64 `json:"-"`
		}

		data := new(SpinData)
		if err := json.Unmarshal(m.Data, data); err != nil {
			mlog.Error("Trasfor %d:%d Error, %v", m.ID, m.ActionID, err)
			return
		}

		// 判斷各項操作
		// {"param1":"[0]","param2":"","param3":"","id":"2","actionId":"6"} 選 PickGame

		var t int = -1
		if len(data.Param1) > 0 {
			param1 := make([]int, 0, 1)
			if err := json.Unmarshal([]byte(data.Param1), &param1); err != nil {
				mlog.Error("Trasfor %d:%d Error, %v", m.ID, m.ActionID, err)
				return
			}
			if len(param1) > 0 {
				t = param1[0]
			}
		}

		mlog.Debug("Client Action Type %d, GameState, %v", t, logic.GameState)

		// 選擇的封包太靠北
		if logic.GameState == slot.GameStateUserAction {
			if logic.SelectGameOptions != nil {
				logic.srvSelectionChoose(t)
				logic.clientSelectionChoose(t)
				return
			}
		}

		switch t {
		case 0: // spin, pick

			// spin
			if logic.GameState == slot.GameStateBet && len(data.Param2) > 0 {
				param2 := make([]float64, 0, 3)
				if err := json.Unmarshal([]byte(data.Param2), &param2); err != nil {
					mlog.Error("Trasfor %d:%d Error, %v", m.ID, m.ActionID, err)
					return
				}
				if len(param2) != 3 {
					mlog.Error("Trasfor %d:%d Error, param2 is not allow", m.ID, m.ActionID)
					return
				}

				if logic.StepResult != nil && logic.leftshow > 0 {
					logic.clientGameResult()
					return
				}

				data.LineCount = param2[0]
				data.BetMoney = param2[1]
				data.Multiplier = param2[2]

				mlog.Debug("Trasfor %d:%d %+v", m.ID, m.ActionID, data)

				var betrate int32
				if logic.Setting != nil {
					betrate = logic.Setting.GetBetRates()[0]
				}
				if logic.LeftFreeGameData != nil {
					logic.srvSpin(logic.MoneyTransfer.S2C(float64(logic.LeftFreeGameData.GetBetMoney())), int32(logic.LeftFreeGameData.GetBetRate()))
				} else {
					logic.srvSpin(data.BetMoney, betrate)
				}

				return
			}

			// other
			if logic.GameState == slot.GameStateUserAction {
				if logic.PickGameRecord != nil {

					mlog.Debug("PickGameItem, Record %v", logic.PickGameRecord)
					// 檢查是否選過
					if len(logic.PickGameRecord.Records) < logic.pickcnt+1 {
						mlog.Error("PickGame has %d record, but this is %d times pick", len(logic.PickGameRecord.Records), logic.pickcnt+1)
						return
					}

					rec := logic.PickGameRecord.Records[logic.pickcnt]

					mlog.Debug("PickGameItem %v", rec)

					defer func() {
						logic.pickcnt++
					}()

					logic.clientPickGameItem(rec.GetID(), rec.GetRewareValue())

					// 沒表演過
					if rec.GetLocation() == 0 {
						rec.Location = proto.Int32(int32(logic.pickcnt + 1))
						logic.srvPickGameItem(rec.GetID(), rec.GetLocation(), rec.GetRewareValue())
						return
					}
				}

				if logic.JackpotGameRecord != nil {

					mlog.Debug("JackpotGameItem, Record %v", logic.JackpotGameRecord)

					// {"param1":"[0]","param2":"[6]","param3":"[7,3]","id":"2","actionId":"6"}
					// {"param1":"[0]","param2":"[7]","param3":"[7,0]","id":"2","actionId":"6"}

					// 檢查是否選過
					if len(logic.JackpotGameRecord.Records) < logic.pickcnt+1 {
						mlog.Error("JackpotGame has %d record, but this is %d times pick", len(logic.JackpotGameRecord.Records), logic.pickcnt+1)
						return
					}

					rec := logic.JackpotGameRecord.Records[logic.pickcnt]

					var loc int32
					if len(data.Param3) > 0 {
						param3 := make([]int32, 0, 2)
						if err := json.Unmarshal([]byte(data.Param3), &param3); err != nil {
							mlog.Error("Trasfor %d:%d Error, %v", m.ID, m.ActionID, err)
							return
						}
						if len(param3) > 0 {
							loc = param3[0]
						}
					}

					mlog.Debug("JackpotGameItem %v", rec)

					defer func() {
						logic.pickcnt++
					}()

					// logic.clientPickGameItem(rec.GetID(), rec.GetRewareValue())

					// 沒表演過
					if rec.GetLocation() == 0 {
						// rec.Location = proto.Int32(int32(logic.pickcnt + 1))
						rec.Location = proto.Int32(loc + 1)
						logic.srvJackpotGameItem(rec.GetID(), rec.GetLocation(), rec.GetRewareValue())

						if len(logic.JackpotGameRecord.Records) == logic.pickcnt+1 {
							lv := slot.NewJackpotLevel(int32(rec.GetRewareValue()))
							logic.srvJackpotPool(logic.JackpotPool[lv].GetID())
							logic.ServeState = ServeStateWaitResult
						}
						return
					}
				}
			}
		}

		return
	}

}

// OnSrvMessage server message
func (logic *Logic) OnSrvMessage(m *msg.Message) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("OnSrvMessage Error, %v, %s", rec, debug.Stack())
		}
	}()

	if m.BClassID == int32(netproto.MessageBClassID_Common) && m.SClassID == int32(netproto.PlatformCommonClassID_HeartBeatReturnID) {
		_m := new(network.Message)
		_m.ID = 0
		_m.ActionID = 2
		logic.cw.WriteClient(_m)
		return
	}

	mlog.Debug("OnSrvMessage %d:%d", m.BClassID, m.SClassID)

	if m.BClassID == int32(netproto.MessageBClassID_Common) && m.SClassID == int32(netproto.PlatformCommonClassID_TipMessageID) {
		data := &netproto.TipMessage{}
		if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
			mlog.Error("Transfor TipMessage Error. %v", err)
			logic.cw.WriteError(&network.Netty{Code: 3000, Data: "System Error ..."})
			return
		}

		logic.cw.WriteError(&network.Netty{Code: 3000, Data: data.GetMessage()})
		return
	}

	if m.BClassID == int32(netproto.MessageBClassID_GameRoom) && m.SClassID == int32(netproto.GameRoomClassID_LoginRoomRetID) {

		data := &netproto.LoginGameRoomRet{}
		if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
			mlog.Error("Transfor LoginGameRoomRet Error. %v", err)
			return
		}

		// 登入失敗
		if data.GetCode() != 1 {
			logic.cw.WriteError(&network.Netty{Code: 3000, Data: data.GetMessage()})
			mlog.Error("User Login Error, Code: %d, Message: %s", data.GetCode(), data.GetMessage())
			return
		}

		mlog.Debug("OnSrvMessage %d:%d LoginRoomRet %v", m.BClassID, m.SClassID, data)

		logic.UserData = data.GetUserData()
		logic.ServeState = ServeStateLogined

		logic.clientLogin()

		return
	}

	if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GameMessageClassID_UserStateChangeID) {
		data := &netproto.UserStateChange{}
		if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
			mlog.Error("Transfor UserStateChange Error. %v", err)
			return
		}

		logic.laststat = data
		// {"actionId":"3","eventData":[0,1000130],"id":"0"}
		// _m := new(network.Message)
		// _m.ID = 0
		// _m.ActionID = 3
		// _m.EventData = []interface{}{0, data.User.GetScore()}

		// logic.cw.WriteClient(_m)

		return
	}

	if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SettingID) {
		data := &netproto.Slot_MSG_Setting{}
		if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
			mlog.Error("Transfor Slot_MSG_Setting Error. %v", err)
			return
		}
		logic.Setting = data
		return
	}

	if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_StatusID) {
		data := &netproto.Slot_MSG_Status{}
		if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
			mlog.Error("Transfor Slot_MSG_Status Error. %v", err)
			return
		}
		mlog.Debug("OnSrvMessage %d:%d %v", m.BClassID, m.SClassID, data)

		state := slot.GameState(data.GetState())

		if logic.Round == 0 {

			switch state {
			default:
				return
			case slot.GameStateBet, slot.GameStateUserAction:

				logic.clientSetting()

				logic.LastGameState = logic.GameState
				logic.GameState = state

				result, callback := logic.gameresult()
				// defer callback()

				var data interface{}

				if state == slot.GameStateUserAction {
					if logic.PickGameRecord != nil {

						logic.GameResult = logic.genFackResult("b1")
						logic.GameResult.BetMoney = proto.Int64(logic.PickGameRecord.GetBetMoney() * logic.PickGameRecord.GetBetRate())

						for _, pos := range logic.PickGameRecord.Positions {
							for _, res := range logic.GameResult.ReelResult {
								if reflect.DeepEqual(pos, res.Position) {
									res.Symbol.SymbolName = proto.String("b1")
								}
							}
						}

						result, _ = logic.gameresult()

						// {"actionId":"1","eventData":[[null,null,0,0],["kgs800087",1114470],[[4,-1],[-1,1,[-1]],[[1,0,1,-1,10,1,0,0,1,0,[]],[[["s3","s4","s2","b1","s3"],["s1","b1","b1","s1","s1"],["s2","s3","s5","s4","s2"]],[0,0,0,0,0],[-1,0]],[],[1]]],-1],"id":"2"}
						data = []interface{}{
							[]interface{}{4, -1},
							[]interface{}{-1, 1, []interface{}{-1}},
							result,
						}
					}

					if logic.JackpotGameRecord != nil {
						callback()
						data = -1
					}

					if logic.SelectGameOptions != nil && logic.SelectGameOptions.GetTimes() > 0 {
						defer callback()

						logic.GameResult = logic.genFackResult()
						logic.GameResult.ReelResult = logic.SelectGameOptions.ReelResult
						logic.GameResult.BetMoney = proto.Int64(logic.SelectGameOptions.GetBetMoney() * logic.SelectGameOptions.GetBetRate())
						result, _ = logic.gameresult()

						data = []interface{}{
							[]interface{}{2, -1},
							[]interface{}{},
							result,
						}
					}
				}

				if state == slot.GameStateBet {
					data = result
					callback()
				}

				if !logic.startedsend {
					logic.clientStart(data)

					logic.startedsend = true
				}

				logic.ServeState = ServeStateWaitBet
			}

			return
		}

		switch state {
		default:
			return
		case slot.GameStateBet, slot.GameStateUserAction:
			logic.LastGameState = logic.GameState
			logic.GameState = state

			if state == slot.GameStateBet {
				logic.pickcnt = 0
				logic.PickGameRecord = nil
				logic.JackpotGameRecord = nil
				logic.SelectGameOptions = nil

				if logic.LeftFreeGameData != nil && logic.LeftFreeGameData.GetCount() > 0 && logic.LastGameState == slot.GameStateUserAction {
					logic.clientEventFreeSpin(int(logic.LeftFreeGameData.GetCount()))
				}
			}
			if logic.ServeState == ServeStateWaitResult {
				logic.clientGameResult()
				logic.ServeState = ServeStateWaitBet
			}
			if logic.laststat != nil && logic.StepResult != nil {
				logic.clientUpdateMoney(logic.laststat.User.GetScore(), true)
			}
		}

		return
	}

	if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_BetStatusID) {
		data := &netproto.Slot_MSG_BetStatus{}
		if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
			mlog.Error("Transfor Slot_MSG_BetStatus Error. %v", err)
			return
		}
		if data.GetRet() == 0 {
			if logic.laststat != nil {
				logic.clientUpdateMoney(logic.laststat.User.GetScore()-data.GetMoney(), false)
			}
			logic.ServeState = ServeStateWaitResult
		}
		return
	}

	if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_JackpotPoolRetID) {
		data := &netproto.Slot_MSG_JackpotPool{}
		if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
			mlog.Error("Transfor Slot_MSG_JackpotPool Error. %v", err)
			return
		}

		var lvs []slot.JackpotLevel
		for _, row := range data.Jackpots {
			lv := slot.JackpotLevel(row.GetName())
			lvs = append(lvs, lv)
			logic.JackpotPool[lv] = row
		}

		logic.clientJackpotPool(lvs...)

		return
	}

	if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LeftFreeGameID) {
		if logic.GameState == slot.GameStateAward {
			return
		}
		data := &netproto.Slot_MSG_LeftFreeGameData{}
		if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
			mlog.Error("Transfor Slot_MSG_LeftFreeGameData Error. %v", err)
			return
		}
		if logic.LeftFreeGameData == nil {
			logic.incfree = int(data.GetCount())
		} else if logic.LeftFreeGameData.GetCount() == data.GetCount() {
			return
		} else {
			logic.incfree = int(data.GetCount() - logic.LeftFreeGameData.GetCount() + 1)
		}
		logic.LeftFreeGameData = data
		return
	}

	if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_BonusGameID) {
		data := &netproto.Slot_MSG_BonusGame{}
		if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
			mlog.Error("Transfor Slot_MSG_BonusGame Error. %v", err)
			return
		}

		switch slot.GameMode(data.GetType()) {
		case slot.GameModePickGame:
			logic.PickGameRecord = data
		case slot.GameModeJackpotGame:
			logic.JackpotGameRecord = data
		}
	}

	if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SelectGameOptionsID) {
		data := &netproto.Slot_MSG_ChooseGameOptions{}
		if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
			mlog.Error("Transfor Slot_MSG_ChooseGameOptions Error. %v", err)
			return
		}

		logic.SelectGameOptions = data

		return
	}

	if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_GameResultID) {

		// 如果是StepResult
		{
			data := new(netproto.Slot_MSG_StepGameResult)
			if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
				mlog.Warn("Transfor Slot_MSG_StepGameResult Error. %v", err)
			} else {
				// fmt.Println("Slot_MSG_StepGameResult", data)
				logic.StepResult = data
				logic.leftshow = len(logic.StepResult.SlotResult)
				if logic.leftshow > 0 {
					return
				}
			}
		}

		data := &netproto.Slot_MSG_GameResult{}
		if err := proto.Unmarshal(m.MsgData.([]byte), data); err != nil {
			mlog.Error("Transfor Slot_MSG_GameResult Error. %v", err)
			return
		}

		// fmt.Println("Slot_MSG_GameResult", data)

		logic.GameResult = data
	}

}

func (logic *Logic) clean() {
	// logic.GameResult = nil
}

func (logic *Logic) printStatus(w io.Writer) {

	fmt.Fprintln(w, "================== Info ==================")
	fmt.Fprintf(w, "User: %d\n", logic.user.PlayerID)

	fmt.Fprintln(w, "================= Status =================")
	fmt.Fprintf(w, "ServeState: %v\n", logic.ServeState)
	fmt.Fprintf(w, "GameState: %v\n", logic.GameState)
	fmt.Fprintf(w, "Round: %d\n", logic.Round)

	fmt.Fprintf(w, "Last Win: %f\n", logic.lastWin)

	if logic.JackpotPool != nil {
		fmt.Fprintln(w, "================ Jackpot =================")
		var jackpots []*netproto.Slot_JackpotInfo

		for _, row := range logic.JackpotPool {
			jackpots = append(jackpots, row)
		}

		sort.Slice(jackpots, func(i, j int) bool {
			return jackpots[i].GetID() < jackpots[j].GetID()
		})

		for _, row := range jackpots {
			fmt.Fprintf(w, "Level: %s, Money: %d\n", row.GetName(), row.GetMoney())
		}
	}

	fmt.Fprintln(w, "================ FreeSpin ================")
	if logic.LeftFreeGameData == nil {
		fmt.Fprintln(w, "Count: 0")
	} else {
		fmt.Fprintf(w, "Count: %d\n", logic.LeftFreeGameData.GetCount())
	}
	fmt.Fprintf(w, "Inc: %d\n", logic.incfree)
	fmt.Fprintf(w, "Send: %t\n", logic.freedatasend)

	fmt.Fprintln(w, "==========================================")
}

func (logic *Logic) onStateChange(before, after slot.GameState) {

	defer func() {
		logic.GameState = after
	}()

}
