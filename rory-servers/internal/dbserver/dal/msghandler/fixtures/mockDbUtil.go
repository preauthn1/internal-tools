package fixtures

import (
	"testing"

	"example.invalid/golang/mock/gomock"
	"example.invalid/ServerTeam/Servers/pkg/db"
)

// MockDBUtil ?
func MockDBUtil(t *testing.T) *db.MockDatabaseEngine {

	mockCtrl := gomock.NewController(t)
	//defer mockCtrl.Finish()

	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)

	return mockDbUtil
}

// GetExpectLogonRobotRet ?
func GetExpectLogonRobotRet() db.RecordSet {

	mockretset := db.RecordSet{}

	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Cols = []string{"UserId", "UserType", "NickName", "FaceId", "Sex", "CashAmount", "WinCount", "LoseCount", "DrawCount", "LevelID"}
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{int64(9000000), int64(2), "武汉市", int64(3), int64(0), int64(10000000), int64(0),
		int64(0), int64(0), int64(3)})
	mockretset = append(mockretset, &dataSetResult1)

	dataSetResult2 := db.RecordTable{}
	dataSetResult2.Cols = []string{"intLogoutRobotCount"}
	dataSetResult2.Rows = append(dataSetResult2.Rows, db.RecordRow{int64(1)})
	mockretset = append(mockretset, &dataSetResult2)

	dataSetResult3 := db.RecordTable{}
	dataSetResult3.Rows = append(dataSetResult3.Rows, db.RecordRow{int64(1)})
	mockretset = append(mockretset, &dataSetResult3)

	return mockretset
}

// GetExpectAddPerGameLogRet ?
func GetExpectAddPerGameLogRet() db.RecordSet {

	mockretset := db.RecordSet{}

	dataSetResult3 := db.RecordTable{}
	dataSetResult3.Rows = append(dataSetResult3.Rows, db.RecordRow{int64(1)})
	mockretset = append(mockretset, &dataSetResult3)

	return mockretset
}

// GetExpectApplyAgentRet ?
func GetExpectApplyAgentRet() db.RecordSet {

	mockretset := db.RecordSet{}
	dataSetResult2 := db.RecordTable{}
	dataSetResult2.Cols = []string{"chvErrMsg"}
	dataSetResult2.Rows = append(dataSetResult2.Rows, db.RecordRow{"申請成功"})
	mockretset = append(mockretset, &dataSetResult2)
	dataSetResult3 := db.RecordTable{}
	dataSetResult3.Rows = append(dataSetResult3.Rows, db.RecordRow{int64(1)})
	mockretset = append(mockretset, &dataSetResult3)

	return mockretset
}

// GetExpectAgentListRet ?
func GetExpectAgentListRet() db.RecordSet {

	mockretset := db.RecordSet{}
	dataSetResult2 := db.RecordTable{}
	dataSetResult2.Cols = []string{"QQ", "WXNo", "Name", "Alipay", "Jiuliao"}
	dataSetResult2.Rows = append(dataSetResult2.Rows, db.RecordRow{"1", "1", "4444", "", ""})
	dataSetResult2.Rows = append(dataSetResult2.Rows, db.RecordRow{"0", "111111", "123456", "test_alipay", ""})
	mockretset = append(mockretset, &dataSetResult2)
	dataSetResult3 := db.RecordTable{}
	dataSetResult3.Rows = append(dataSetResult3.Rows, db.RecordRow{int64(0)})
	mockretset = append(mockretset, &dataSetResult3)

	return mockretset
}
