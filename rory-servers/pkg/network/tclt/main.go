package main

import (
	"flag"
	"fmt"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

func main() {

	var port = flag.Int("port", 8011, "input server port")

	flag.Parse()

	packer := network.NewHeadMsgPacker(1024, true)
	evtChan := network.NewClientEventChan(100)
	c := network.NewTCPClient(packer, evtChan)
	clt := network.NewReConnectTCPClient(c, 10)

	go func() {
		for {
			select {
			case data := <-evtChan.DataSig:
				fmt.Printf("recv data ..., len:%d\r\n", len(data))
			case data := <-evtChan.MsgSig:
				fmt.Printf("recv msg ..., %s\r\n", string(data))
			case <-evtChan.ConnectedSig:
				fmt.Println("connected ...")
			case <-evtChan.DisconnectedSig:
				fmt.Println("disconnected ...")
			}
		}
	}()

	err := clt.Connect(fmt.Sprintf(":%d", *port))
	fmt.Println(err)

	go func() {
		i := 0
		for {
			if !clt.IsConnected() {
				return
			}

			time.Sleep(time.Second * 10)
			str := fmt.Sprintf("this is just a test . %d", i)
			fmt.Printf("发送数据%s\r\n", str)
			b1 := util.IntToBytes(0, true)
			b2 := util.IntToBytes(0, true)
			b3 := append(b1, b2...)
			b4 := append(b3, ([]byte(str))...)

			clt.WriteMsg(b4)
			i++
		}
	}()

	time.Sleep(time.Minute * 5)
}
