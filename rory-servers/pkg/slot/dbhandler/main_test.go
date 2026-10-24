package dbhandler

import (
	"testing"

	apiutil "example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	dbutil "example.invalid/ServerTeam/Servers/pkg/db/dbutil"
)

func TestMain(m *testing.M) {

	dbutil.ConnString = "server=0.0.0.0;database=CenterDB;user id=sa;password=[REDACTED];port=1433;encrypt=disable;"
	dbutil.DealLogConnString = "server=0.0.0.0;database=PerGameLog;user id=sa;password=[REDACTED];port=1433;encrypt=disable;" //牌局记录库数据库连接字符串
	dbutil.ClubConnString = "server=0.0.0.0;database=Club;user id=sa;password=[REDACTED];port=1433;encrypt=disable;"          //Club数据库连接字符串

	dbutil.SMSWebAPI = "0.0.0.0:80/SMS/Send"
	dbutil.PoolSize = 200

	dbutil.InitDatabase()

	apiutil.InitTestService()

	m.Run()
}
