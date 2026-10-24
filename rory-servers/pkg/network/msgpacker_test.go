package network

import (
	"fmt"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

func Test_pack(b *testing.T) {

	aa := util.IntToBytes(10, true)
	fmt.Println(aa)

	data := "this is just a test"
	data1 := "this is just a test2"
	packer := NewHeadMsgPacker(10240, true)
	ndata := packer.Packet([]byte(data))

	ndata1 := packer.Packet([]byte(data1))

	ndata = append(ndata, ndata1...)

	fmt.Println(len(ndata))

	nndata, _ := packer.UnPack(ndata[:3])
	nndata, _ = packer.UnPack(ndata[3:10])
	nndata, _ = packer.UnPack(ndata[10:])

	nndata2, _ := packer.UnPack(nil)

	fmt.Println(string(nndata))
	fmt.Println(string(nndata2))

}
