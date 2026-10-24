package network

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// MsgPacker 数据包处理接口定义
type MsgPacker interface {
	Packet(data []byte) []byte
	UnPack(data []byte) ([]byte, error)
}

// HeadMsgPacker 基于包头包体方式的分包
type HeadMsgPacker struct {
	maxMsgLen int32  //最大消息长度
	bigEndian bool   //是否大头字节序
	dataBuf   []byte //缓冲数据流
}

// NewHeadMsgPacker Create HeadMsgPacker
func NewHeadMsgPacker(maxMsgLen int32, bigEndian bool) *HeadMsgPacker {
	p := new(HeadMsgPacker)
	p.maxMsgLen = maxMsgLen
	p.bigEndian = bigEndian

	return p
}

// 包头长度
const (
	MSGHEADLEN int32 = 4
)

// Packet 打包
func (packer *HeadMsgPacker) Packet(data []byte) []byte {
	msgLen := int32(len(data))
	mlog.Trace("打包封包, 長度: %d, 前8碼: %v, %v:%v", msgLen, data[:8], util.BytesToInt(data[:4], packer.bigEndian), util.BytesToInt(data[4:8], packer.bigEndian))
	return append(util.IntToBytes(msgLen, packer.bigEndian), data...)
}

// UnPack 解包
func (packer *HeadMsgPacker) UnPack(data []byte) ([]byte, error) {
	var buf []byte
	if data != nil && len(data) > 0 {
		buf = append(packer.dataBuf, data...)
	} else {
		buf = packer.dataBuf
	}

	if int32(len(buf)) <= MSGHEADLEN {
		packer.dataBuf = buf
		return nil, nil
	}

	msglen := util.BytesToInt(buf[:MSGHEADLEN], packer.bigEndian)

	mlog.Trace("解開封包, 長度: %d, 前8碼: %v, %v:%v", msglen, buf[MSGHEADLEN:(MSGHEADLEN+8)], util.BytesToInt(buf[MSGHEADLEN:(MSGHEADLEN+4)], packer.bigEndian), util.BytesToInt(buf[MSGHEADLEN+4:(MSGHEADLEN+8)], packer.bigEndian))

	if msglen >= packer.maxMsgLen || msglen <= 0 {
		mlog.Error("消息長度超出範圍, %d[%d], 前8碼: %v, %v:%v", msglen, packer.maxMsgLen, buf[MSGHEADLEN:(MSGHEADLEN+8)], util.BytesToInt(buf[MSGHEADLEN:(MSGHEADLEN+4)], packer.bigEndian), util.BytesToInt(buf[MSGHEADLEN+4:(MSGHEADLEN+8)], packer.bigEndian))
		return nil, fmt.Errorf("消息長度超出範圍, 當前消息長度 %d[%d], 前8碼: %v, %v:%v", msglen, packer.maxMsgLen, buf[MSGHEADLEN:(MSGHEADLEN+8)], util.BytesToInt(buf[MSGHEADLEN:(MSGHEADLEN+4)], packer.bigEndian), util.BytesToInt(buf[MSGHEADLEN+4:(MSGHEADLEN+8)], packer.bigEndian))
	}

	if int32(len(buf)) < MSGHEADLEN+msglen {
		packer.dataBuf = buf
		return nil, nil
	}

	packer.dataBuf = buf[MSGHEADLEN+msglen:]
	return buf[MSGHEADLEN : MSGHEADLEN+msglen], nil
}

// EmptyMsgPacker 空的数据包处理器
type EmptyMsgPacker struct {
}

// Packet 打包
func (packer *EmptyMsgPacker) Packet(data []byte) []byte {
	return data
}

// UnPack 解包
func (packer *EmptyMsgPacker) UnPack(data []byte) ([]byte, error) {
	return nil, nil
}

// MsgPackerGroup 数据包处理器组
type MsgPackerGroup struct {
	innerPackers []MsgPacker
}

// Packet 打包
func (packer *MsgPackerGroup) Packet(data []byte) []byte {
	ndata := data
	for _, p := range packer.innerPackers {
		ndata = p.Packet(ndata)
	}

	return ndata
}

// UnPack 解包
func (packer *MsgPackerGroup) UnPack(data []byte) ([]byte, error) {
	dlen := len(packer.innerPackers)
	ndata := data
	var err error

	for i := dlen - 1; i >= 0; i-- {
		var nerr error
		p := packer.innerPackers[i]
		ndata, nerr = p.UnPack(ndata)

		if nerr != nil {
			err = nerr
		}
	}

	return ndata, err
}

// AppendPacker ..
func (packer *MsgPackerGroup) AppendPacker(p MsgPacker) *MsgPackerGroup {
	packer.innerPackers = append(packer.innerPackers, p)

	return packer
}
