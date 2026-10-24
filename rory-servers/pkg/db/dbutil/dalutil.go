package dbutil

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/httputil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// Settings
var (
	ConnString        string
	DealLogConnString string
	ClubConnString    string
	PoolSize          int = 100
	database          *db.Database
	dealLogDataBase   *db.Database
	clubDataBase      *db.Database
	wildCatDataBase   *db.Database
	SMSWebAPI         string
	DBUtil            db.DatabaseEngine
	DealLogDBUtil     db.DatabaseEngine
	ClubDBUtil        db.DatabaseEngine
	WildCatDBUtil     db.DatabaseEngine
	HDCodeMap         *util.SafeMap
	FastLogonStat     bool
)

// InitDatabase init
func InitDatabase() {
	//dao := new(db.Database)
	//dao.Connstring = ConnString
	//dao.Connpool = db.NewConnPool(PoolSize, ConnString)
	database = db.NewDatabase(ConnString, PoolSize)
	dealLogDataBase = db.NewDatabase(DealLogConnString, PoolSize)
	clubDataBase = db.NewDatabase(ClubConnString, PoolSize)
	DBUtil = db.NewDatabaseEngine(database)
	DealLogDBUtil = db.NewDatabaseEngine(dealLogDataBase)

	ClubDBUtil = db.NewDatabaseEngine(clubDataBase)
	WildCatDBUtil = db.NewDatabaseEngine(wildCatDataBase)

	HDCodeMap = util.NewSafeMap()
	FastLogonStat = false
	command.AddCommand(NewDBUtilCmd())
	command.AddCommand(httputil.NewAPIUtilCmd())
	//TickerStatus() //定時 output db 用量
}

//func TickerStatus() {
//	go func() {
//		for {
//			time.Sleep(time.Duration(10) * time.Second)
//			var buf bytes.Buffer
//			buf.WriteString("-----------------------\r\n")
//			buf.WriteString(database.GetStatus())
//			buf.WriteString(dealLogDataBase.GetStatus())
//			buf.WriteString(clubDataBase.GetStatus())
//			buf.WriteString("-----------------------\r\n")
//			fmt.Print(buf.String())
//		}
//	}()
//}

// GetErrMsgValue 獲取提示消息
func GetErrMsgValue(set db.RecordSet) string {
	rlen := len(set)

	return set[rlen-2].GetValueByColName(0, "chvErrMsg").(string)
}

// GetDBRetMsg 獲取數據庫執行結果消息
func GetDBRetMsg(retCode int32) *msg.Message {
	ret := new(netproto.DBServerRet)
	ret.Code = proto.Int32(retCode)

	msg := new(msg.Message)
	msg.BClassID = int32(netproto.MessageBClassID_DBServer)
	msg.SClassID = int32(netproto.DBServerClassID_DBRetID)
	msg.MsgData = ret

	return msg
}

// ParseUserLogonRet 解析用戶登錄結果
func ParseUserLogonRet(ret db.RecordSet) *netproto.UserLoginRet {

	logret := &netproto.UserLoginRet{}

	// 充值返利
	if ret.GetRetTableCount() >= 7 {
		grt := ret[5]
		if len(grt.Rows) < 1 {
			mlog.Error("充值返利信息為為空, userid=%d", logret.GetUserID())
			FlowAmount := new(netproto.FlowAmountInfo)
			FlowAmount.Status = proto.Int32(0)
			FlowAmount.UnionPlatform = proto.Int32(0)
			FlowAmount.FlowAmountNow = proto.Int32(0)
			FlowAmount.FlowAmountLimit = proto.Int32(0)
			FlowAmount.Achievement = proto.Int32(0)
			FlowAmount.FlowAmountRate = proto.Int32(0)
			FlowAmount.BonusRate = proto.Int32(0)
			FlowAmount.CustomerDomain = proto.String("")
			logret.FlowInfo = FlowAmount
		} else {
			grt.ForEachRow(func(r db.RecordRow) {
				FlowAmount := new(netproto.FlowAmountInfo)
				FlowAmount.Status = proto.Int32(int32(grt.GetValueByColName(0, "Status").(int64)))
				FlowAmount.UnionPlatform = proto.Int32(int32(grt.GetValueByColName(0, "UnionPlatform").(int64)))
				FlowAmount.FlowAmountNow = proto.Int32(int32(grt.GetValueByColName(0, "FlowAmountNow").(int64)))
				FlowAmount.FlowAmountLimit = proto.Int32(int32(grt.GetValueByColName(0, "FlowAmountLimit").(int64)))
				FlowAmount.Achievement = proto.Int32(int32(grt.GetValueByColName(0, "OpenStatus").(int64)))
				FlowAmount.FlowAmountRate = proto.Int32(int32(grt.GetValueByColName(0, "FlowAmountRate").(int64)))
				FlowAmount.BonusRate = proto.Int32(int32(grt.GetValueByColName(0, "BonusRate").(int64)))
				FlowAmount.CustomerDomain = proto.String(grt.GetSingleValueString("CustomerDomain"))
				logret.FlowInfo = FlowAmount
			})
		}
	}

	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().SetData(ret).Error("解析用戶登錄結果發生錯誤%v, %s", err, errorStr)
			panic(err)
		}
	}()

	code := 0
	if ret.GetReturnValue() == 1 {
		code = 1
	} else {
		// 登入失败
		errmsg := ret.GetOutputParamValue("chvErrMsg").(string)
		userid := ret.GetOutputParamValue("intUserId").(int64)
		mlog.Warn("parseUserLogon Error, code %d, userid %d, errmsg %s", ret.GetReturnValue(), userid, errmsg)
	}

	if ret.GetReturnValue() == -2 {
		mlog.Error("邀請碼錯誤")
		logret.Code = proto.Int32(-2)
		logret.Message = proto.String("遊戲重要數據不完整，請重新掃碼下載")
		return logret
	}

	if len(ret[0].Rows) < 1 { // 找不到用戶
		mlog.Error("缺少了用戶資料")
		logret.Code = proto.Int32(0)
		logret.Message = proto.String("服務器錯誤")
		return logret
	}

	logret.Code = proto.Int32(int32(code))
	logret.Message = proto.String(ret.GetOutputParamValue("chvErrMsg").(string))

	if code == 1 {
		logret.UserID = proto.Int32(int32(ret.GetOutputParamValue("intUserID").(int64)))
		logret.Cer = proto.String(ret.GetOutputParamValue("chvCertification").(string))

		userdata := new(netproto.UserHallLogonData)
		userdata.LoginName = proto.String(ret[0].GetSingleValue("LoginName").(string))
		userdata.NickName = proto.String(ret[0].GetSingleValue("NickName").(string))
		userdata.UserType = proto.Int32(ret[0].GetSingleValueInt32("UserType"))
		userdata.Sex = proto.Int32(ret[0].GetSingleValueInt32("Sex"))
		userdata.HeadID = proto.Int32(ret[0].GetSingleValueInt32("HeadID"))
		userdata.IsTopWindow = proto.Bool(ret[0].GetSingleValueInt32("IsTopWindow") == 1)
		userdata.CashAmount = proto.Int64(ret[0].GetSingleValueInt64("CashAmount"))
		userdata.BankAmount = proto.Int64(ret[0].GetSingleValueInt64("BankAmount"))
		userdata.IsGaming = proto.Bool(ret[0].GetSingleValue("IsGaming").(bool))
		userdata.ServerAddr = proto.String(ret[0].GetSingleValue("ServerAddr").(string))
		userdata.ServerName = proto.String(ret[0].GetSingleValue("ServerName").(string))
		userdata.IsKick = proto.Bool(ret[0].GetSingleValue("IsKick").(bool))
		userdata.IsBindGuest = proto.Bool(ret[0].GetSingleValueInt32("IsBindGuest") == 1)
		userdata.IsBindZhifubao = proto.Bool(ret[0].GetSingleValueInt32("IsBindZhifubao") == 1)
		userdata.Zhifubao = proto.String(ret[0].GetSingleValue("AlipayAccount").(string))
		userdata.RealName = proto.String(ret[0].GetSingleValue("RealName").(string))
		userdata.UserLevel = proto.Int32(ret[0].GetSingleValueInt32("UserLevel"))
		userdata.LevelKey = proto.String(ret[0].GetSingleValue("LevelKey").(string))
		userdata.AnnMsg = proto.String(ret[0].GetSingleValue("AnnMsg").(string))
		userdata.ConvertRateTipMsg = proto.String(ret[0].GetSingleValue("ConvertRateTipMsg").(string))
		userdata.BankPwdTipMsg = proto.String(ret[0].GetSingleValue("BankPwdTipMsg").(string))
		userdata.UIFlag = proto.String(ret[0].GetSingleValue("UIFlag").(string))
		userdata.PaySort = proto.String(ret[0].GetSingleValue("PaySort").(string))
		userdata.PayAmountConfig = proto.String(ret[0].GetSingleValue("PayAmountConfig").(string))
		userdata.PayNotifyMsg = proto.String(ret[0].GetSingleValue("PayNotifyMsg").(string))
		userdata.PayTips = proto.String(ret[0].GetSingleValue("PayTips").(string))
		userdata.NotifyFlag = proto.String(ret[0].GetSingleValue("NotifyFlag").(string))
		userdata.UpGradeMsg = proto.String(ret[0].GetSingleValue("UpGradeMsg").(string))
		userdata.LockGameID = proto.Int32(ret[0].GetSingleValueInt32("LockGameID"))
		userdata.IsBindBankCard = proto.Bool(ret[0].GetSingleValueInt32("IsBindUnionpayBank") == 1)
		userdata.BankCode = proto.String(ret[0].GetSingleValue("UnionpayBank").(string))
		userdata.BankCardNum = proto.String(ret[0].GetSingleValue("UnionpayAccount").(string))
		userdata.BankName = proto.String(ret[0].GetSingleValue("UnionpayBankName").(string))
		userdata.BankRealName = proto.String(ret[0].GetSingleValue("UnionpayRealName").(string))
		userdata.UISort = proto.String(ret[0].GetSingleValue("UISort").(string))
		userdata.SumCharge = proto.Int32(ret[0].GetSingleValueInt32("SumCharge"))
		userdata.ShowConvertLimit = proto.Int32(ret[0].GetSingleValueInt32("ShowConvertLimit"))
		userdata.ShowAllGameLimit = proto.Int32(ret[0].GetSingleValueInt32("ShowAllGameLimit"))
		userdata.QuestionContactTip = proto.String(ret[0].GetSingleValueString("QuestionContactTip"))
		userdata.AnnMsgShow = proto.Int32(ret[0].GetSingleValueInt32("AnnMsgShow"))
		userdata.ClubOpen = proto.Bool(ret[0].GetSingleValueBool("ClubOpen"))
		userdata.ConvertBankShow = proto.Int32(ret[0].GetSingleValueInt32("ConvertBankShow"))
		userdata.ConvertAlipayShow = proto.Int32(ret[0].GetSingleValueInt32("ConvertAlipayShow"))
		userdata.ShowVIPAgent = proto.Int64(ret[0].GetSingleValueInt64("ShowVIPAgent"))
		userdata.ShopNameID = proto.Int32(ret[0].GetSingleValueInt32("ShopNameID"))
		userdata.AgentNameID = proto.Int32(ret[0].GetSingleValueInt32("AgentNameID"))
		userdata.ShowCS = proto.Int64(ret[0].GetSingleValueInt64("ShowCS"))
		userdata.PlaytimeCS = proto.Int32(ret[0].GetSingleValueInt32("PlaytimeCS"))
		userdata.PlaytimeConvert = proto.Int32(ret[0].GetSingleValueInt32("PlaytimeConvert"))
		userdata.PlaytimeVIPAgent = proto.Int32(ret[0].GetSingleValueInt32("PlaytimeVIPAgent"))
		userdata.ShowConvertBindPhone = proto.Int32(ret[0].GetSingleValueInt32("ShowConvertBindPhone"))
		userdata.ShowAllGameBindPhone = proto.Int32(ret[0].GetSingleValueInt32("ShowAllGameBindPhone"))
		userdata.ShowVIPAgentBindPhone = proto.Int32(ret[0].GetSingleValueInt32("ShowVIPAgentBindPhone"))
		userdata.ShowCSBindPhone = proto.Int32(ret[0].GetSingleValueInt32("ShowCSBindPhone"))
		userdata.ShowConvertPopup = proto.Int32(ret[0].GetSingleValueInt32("ShowConvertPopup"))
		userdata.ConvertPopupText = proto.String(ret[0].GetSingleValueString("ConvertPopupText"))
		userdata.WinLossFlowAmount = proto.Int32(ret[0].GetSingleValueInt32("WinLossFlowAmount"))
		userdata.ExchangeFlowAmount = proto.Int32(ret[0].GetSingleValueInt32("ExchangeFlowAmount"))
		userdata.DayFlowAmount = proto.Int32(ret[0].GetSingleValueInt32("DayFlowAmount"))
		userdata.SevenDayFlowAmount = proto.Int32(ret[0].GetSingleValueInt32("SevenDayFlowAmount"))
		userdata.MonthFlowAmount = proto.Int32(ret[0].GetSingleValueInt32("MonthFlowAmount"))
		// 遊戲
		if ret.GetRetTableCount() >= 3 {
			grt := ret[1]

			grt.ForEachRow(func(r db.RecordRow) {
				gid := r[0].(int64)
				popular := r[1].(int64)
				iconName := r[2].(string)
				aliasName := r[3].(string)
				iconTitle := r[4].(string)
				skin := r[5].(string)
				isThirdParty := int64(0)
				if len(r) > 6 { // 兼容舊db
					isThirdParty = r[6].(int64)
				}

				userdata.GameList = append(userdata.GameList, int32(gid))
				// 新的GameList 協議
				gameList2 := new(netproto.GameList2)
				gameList2.GameID = proto.Int32(int32(gid))
				gameList2.Popular = proto.Int32(int32(popular))
				gameList2.IconName = proto.String(iconName)
				gameList2.AliasName = proto.String(aliasName)
				gameList2.IconTitle = proto.String(iconTitle)
				gameList2.Skin = proto.String(skin)
				gameList2.IsThirdParty = proto.Int32(int32(isThirdParty))
				userdata.GameList2 = append(userdata.GameList2, gameList2)
			})
		}

		// 支付方式最大值及最小值
		if ret.GetRetTableCount() >= 4 {
			payAmonutRange := ""
			grt := ret[2]
			grt.ForEachRow(func(r db.RecordRow) {
				payAmonutRange += fmt.Sprintf("%v:%v.%v", r[0], r[1], r[2]) + "|"
			})
			if len(payAmonutRange) > 0 {
				userdata.PayAmountRange = proto.String(payAmonutRange[:len(payAmonutRange)-1])
			}
		}

		// 輪播圖
		if ret.GetRetTableCount() >= 5 {
			grt := ret[3]
			if len(grt.Rows) < 1 {
				mlog.Error("輪播圖資訊為空, userid=%d", logret.GetUserID())
			}
			grt.ForEachRow(func(r db.RecordRow) {
				picture := r[0].(string)
				link := r[1].(string)
				behavior := r[2].(int64)
				banner := new(netproto.Banner)
				banner.Picture = proto.String(picture)
				banner.Link = proto.String(link)
				banner.Behavior = proto.Int32(int32(behavior))
				userdata.Banner = append(userdata.Banner, banner)
			})
		}

		// 支付信息
		if ret.GetRetTableCount() >= 6 {
			grt := ret[4]
			if len(grt.Rows) < 1 {
				mlog.Error("支付信息為空, userid=%d", logret.GetUserID())
			}
			grt.ForEachRow(func(r db.RecordRow) {
				gwID := r[0].(int64)
				typeID := r[1].(int64)
				typeName := r[2].(string)
				amountSort := r[3].(string)
				paymentInfo := new(netproto.PaymentInfo)
				paymentInfo.GWID = proto.Int32(int32(gwID))
				paymentInfo.TypeID = proto.Int32(int32(typeID))
				paymentInfo.TypeName = proto.String(typeName)
				paymentInfo.AmountSort = proto.String(amountSort)
				userdata.PaymentInfoList = append(userdata.PaymentInfoList, paymentInfo)
			})
		}

		logret.UserData = userdata
	}

	return logret
}

// UpdateProtectNodeStatus ?
func UpdateProtectNodeStatus(userid int32, token string, status int32) {
	ps := db.NewSQLParameters()
	ps.AddVarcharInput("chvToken", token, 10)
	ps.AddIntInput("intLoginUserID", userid)
	ps.AddIntInput("intStatus", status)
	proc := db.NewProcedure("PrPs_UpdateProtectNodeStatus", ps)
	ret, err := DBUtil.ExecProc(proc)
	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().SetData(ret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
			panic(err)
		}
	}()
	//fmt.Println(proc.GetSQLString())
	if err != nil {
		panic(err)
	}
}

// GetRetMsg ?
func GetRetMsg(code int32, msgContent string) *netproto.RetMessage {
	retmsg := new(netproto.RetMessage)
	retmsg.Code = proto.Int32(code)
	retmsg.Message = proto.String(msgContent)

	return retmsg
}

// HandleDBParseErr ?
func HandleDBParseErr(proc *db.Procedure, ret db.RecordSet) {
	if err := recover(); err != nil {
		errorStr := util.GetStackInfo()
		mlog.NewDefaultLogger().SetData(ret).Error("解析存儲過程[%s] %s, 發生錯誤%s", proc.Name, proc.GetSQLString(), errorStr)
		panic(err)
	}
}
