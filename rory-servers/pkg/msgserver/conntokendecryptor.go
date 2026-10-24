package msgserver

import (
	"encoding/json"
	"errors"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

// Token ?
type Token struct {
	Value string `json:"value"`
	Time  int    `json:"time"`
}

// TFDataReq ?
type TFDataReq struct {
	UserID     int
	VIP        int
	GUI        string
	Channel    string
	Simulation int
	Action     int // 1正常走節點 2強迫取備用節點
	Time       int `json:"time"`
	AreaID     int
	Ver        int
}

// TFDataRet ?
type TFDataRet struct {
	Result  int    `json:"result"`
	Message string `json:"message"`
	Token   string `json:"token"`
	Data    struct {
		Addr         string `json:"Addr"`
		TransferPort int    `json:"TransferPort"`
		GiveUserIP   int    `json:"GiveUserIP"`
		From         string `json:"From"`
		UserIP       string `json:"UserIP"`
	} `json:"data"`
}

// NodeDefDecryptor ?
func NodeDefDecryptor(data string) (*TFDataRet, error) {
	if len(data) == 0 {
		return nil, errors.New("Decrypt data is empty")
	}
	dataKey := "V1M5WE7X8WE2TWEK"
	iv := "DNW123EXGMTR83K7"
	s, _ := util.AesDecryptString(data, dataKey, iv)
	//mlog.System(s)
	ret := new(TFDataRet)
	if err := json.Unmarshal([]byte(s), ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// ConnTokenDecryptor ?
func ConnTokenDecryptor(data string) (*Token, error) {
	tokenKey := "X2LGWS7XDFE2YWRU"
	iv := "DNW123EXGMTR83K7"
	ret, err := NodeDefDecryptor(data)
	if err != nil {
		return nil, err
	}
	if ret.Token != "" {
		t, _ := util.AesDecryptString(ret.Token, tokenKey, iv)
		var tokenResult map[string]interface{}
		if err := json.Unmarshal([]byte(t), &tokenResult); err != nil {
			return nil, err
		}
		token := new(Token)
		if tokenResult["value"] != nil {
			token.Value = tokenResult["value"].(string)
		}
		if tokenResult["time"] != nil {
			foo := int(tokenResult["time"].(float64))
			token.Time = foo

		}
		return token, nil
	}

	return nil, errors.New("Could not fetch token")
}
