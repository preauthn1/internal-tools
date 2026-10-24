package msg

import (
	"errors"
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

// Default
const (
	AllMsgBClassID = -1
	AllMsgSClassID = -1
)

// ParserGroup 消息解析服務器組
type ParserGroup struct {
	innerParsers map[string]Parser
}

// NewParserGroup 創建消息解析器組對象
func NewParserGroup() *ParserGroup {
	svr := new(ParserGroup)
	svr.innerParsers = make(map[string]Parser)

	return svr
}

func getParserKey(bClsID int32, sClsID int32) string {
	return fmt.Sprintf("%d_%d", bClsID, sClsID)
}

// RegisterParser 註冊解析器
func (server *ParserGroup) RegisterParser(bClsID int32, sClsID int32, parser Parser) error {
	key := getParserKey(bClsID, sClsID)
	if _, ok := server.innerParsers[key]; ok {
		return errors.New("消息解析器已經存在")
	}

	server.innerParsers[key] = parser
	return nil
}

// Marshal 序列化消息
func (server *ParserGroup) Marshal(msg *Message) ([]byte, error) {
	// 使用特定解析方式
	key := getParserKey(msg.BClassID, msg.SClassID)

	if parser, ok := server.innerParsers[key]; ok {
		return parser.Marshal(msg)
	}

	// 使用通用解析方式
	allkey := getParserKey(AllMsgBClassID, AllMsgSClassID)

	if allparser, ok := server.innerParsers[allkey]; ok {
		return allparser.Marshal(msg)
	}

	return nil, errors.New("消息解析器不存在")
}

// Unmarshal ?
func (server *ParserGroup) Unmarshal(data []byte) (*Message, error) {
	if len(data) < 8 {
		return nil, errors.New("消息不正确")
	}

	// 使用特定的解析方式
	// fmt.Println(data)
	bclsID := util.BytesToInt(data[:4], true)
	cclsID := util.BytesToInt(data[4:8], true)

	key := getParserKey(bclsID, cclsID)

	if parser, ok := server.innerParsers[key]; ok {
		return parser.Unmarshal(data)
	}

	// 使用通用解析方式
	allkey := getParserKey(AllMsgBClassID, AllMsgSClassID)

	if allparser, ok := server.innerParsers[allkey]; ok {
		return allparser.Unmarshal(data)
	}

	errStr := fmt.Sprintf("消息解析器不存在 bclsID[%v] cclsID[%v]", bclsID, cclsID)
	return nil, errors.New(errStr)
}
