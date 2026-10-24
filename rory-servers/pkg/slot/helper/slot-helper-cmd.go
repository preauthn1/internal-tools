package helper

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"example.invalid/spf13/cobra"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/gameserver"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/slot/baselogic"
)

// SlotHelperCmdP 查詢房間狀態指令
var SlotHelperCmdP *SlotHelperCmd

// SlotHelperCmd 查詢房間狀態指令
type SlotHelperCmd struct {
	gs            *gameserver.GameServer
	glm           *gamelogic.LogicManager
	BaseLogicFunc func(gl gamelogic.GameLogic) *baselogic.SlotLogicBase
	traceMap      map[int32]bool
	mu            sync.Mutex
}

// NewSlotHelperCmd Create SlotHelperCmd
func NewSlotHelperCmd(gs *gameserver.GameServer, glm *gamelogic.LogicManager) *SlotHelperCmd {
	cmd := new(SlotHelperCmd)
	cmd.gs = gs
	cmd.glm = glm
	cmd.traceMap = make(map[int32]bool)
	return cmd
}

// Name Command name
func (c *SlotHelperCmd) Name() string {
	return "slot"
}

// Help Command Description
func (c *SlotHelperCmd) Help() string {
	return fmt.Sprintf("Slot 遊戲相關工具\t輸入『slot help』可查詢功能說明")
}

// Exec Command Logic
func (c *SlotHelperCmd) Exec(args []string) string {
	cmd := c.CobraRoot()
	cmd.SetArgs(args)
	cmd.Execute()
	return ""
}

// CobraRoot 取得 cobra command
func (c *SlotHelperCmd) CobraRoot() *cobra.Command {
	cmd := c.rootCommand()
	cmd.AddCommand(c.infoCmd())
	cmd.AddCommand(c.protoCmd())
	cmd.AddCommand(c.traceCmd())
	cmd.AddCommand(c.ptnctrlCmd())
	return cmd
}

func (c *SlotHelperCmd) rootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "slot",
		Short: "Slot 遊戲相關工具",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}
}

func (c *SlotHelperCmd) infoCmd() *cobra.Command {
	infoCmd := &cobra.Command{
		Use:                   "info",
		Short:                 "查看遊戲設定",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			if c.BaseLogicFunc == nil {
				return fmt.Errorf("沒有設定取得房間的程式")
			}

			var tableno int32
			{
				uid, err := cmd.PersistentFlags().GetInt32("user")
				if err != nil {
					return fmt.Errorf("玩家編號輸入格式錯誤, %v", err)
				}
				if uid > 0 {
					p := c.gs.GameServerRunCtx.PlayerManager.GetPlayerByUserID(uid)
					if p == nil {
						return fmt.Errorf("查無玩家 %d", uid)
					}
					if u, ok := p.UserData.(*gamelogic.GameUser); ok {
						tableno = u.TableNo
					}
				}
			}

			tableno, err := cmd.PersistentFlags().GetInt32("table")
			if err != nil {
				return fmt.Errorf("桌號輸入格式錯誤, %v", err)
			}

			// 抓一間房間來印
			l := c.glm.GetGameLogic(tableno)
			if l == nil {
				return fmt.Errorf("沒有設置房間")
			}

			logic := c.BaseLogicFunc(l)

			cmd.Println("Table:", logic.TableNo())
			PrintLogicInfo(cmd.OutOrStdout(), logic)

			return nil
		},
	}
	infoCmd.PersistentFlags().Int32P("table", "t", 0, "指定桌號")
	infoCmd.PersistentFlags().Int32P("user", "u", 0, "指定玩家編號")

	return infoCmd
}

func (c *SlotHelperCmd) protoCmd() *cobra.Command {
	protoCmd := &cobra.Command{
		Use:          "proto",
		Short:        "查看遊戲機率設定",
		Hidden:       true,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			if c.BaseLogicFunc == nil {
				return fmt.Errorf("沒有設定取得房間的程式")
			}

			var tableno int32
			{
				uid, err := cmd.PersistentFlags().GetInt32("user")
				if err != nil {
					return fmt.Errorf("玩家編號輸入格式錯誤, %v", err)
				}
				if uid > 0 {
					p := c.gs.GameServerRunCtx.PlayerManager.GetPlayerByUserID(uid)
					if p == nil {
						return fmt.Errorf("查無玩家 %d", uid)
					}
					if u, ok := p.UserData.(*gamelogic.GameUser); ok {
						tableno = u.TableNo
					}
				}
			}

			tableno, err := cmd.PersistentFlags().GetInt32("table")
			if err != nil {
				return fmt.Errorf("桌號輸入格式錯誤, %v", err)
			}

			// 抓一間房間來印
			l := c.glm.GetGameLogic(tableno)
			if l == nil {
				return fmt.Errorf("沒有設置房間")
			}

			logic := c.BaseLogicFunc(l)
			pwd, _ := cmd.PersistentFlags().GetString("password")

			if len(pwd) > 0 && pwd != conf.Password {
				return fmt.Errorf("輸入密碼錯誤")
			}

			cmd.Println("Table:", logic.TableNo())
			PrintLogicProto(cmd.OutOrStdout(), logic.Proto, pwd == conf.Password)

			return nil
		},
	}
	protoCmd.PersistentFlags().StringP("password", "p", "", "看詳細資訊")
	protoCmd.PersistentFlags().Int32P("table", "t", 0, "指定桌號")
	protoCmd.PersistentFlags().Int32P("user", "u", 0, "指定玩家編號")

	return protoCmd
}

var orgfunc func(userid int32, m *msg.Message)

func (c *SlotHelperCmd) traceCmd() *cobra.Command {

	examples := []string{
		"slot trace -t 0 \t(追蹤第0桌1分鐘)",
		"slot trace -t 10 -d 5m \t(追蹤第10桌5分鐘)",
	}

	traceCmd := &cobra.Command{
		Use:     "trace",
		Short:   "追蹤房間封包",
		Long:    "追蹤房間封包，顯示封包內容，追蹤期間暫時無法輸入指令",
		Example: strings.Join(examples, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {

			if c.BaseLogicFunc == nil {
				return fmt.Errorf("沒有設定取得房間的程式")
			}

			tableno, err := cmd.PersistentFlags().GetInt32("table")
			if err != nil {
				return fmt.Errorf("桌號輸入格式錯誤, %v", err)
			}

			if tableno < 0 {
				return fmt.Errorf("請指定桌號")
			}

			l := c.glm.GetGameLogic(tableno)
			if l == nil {
				return fmt.Errorf("查無房間(%d)", tableno)
			}

			d, err := cmd.PersistentFlags().GetDuration("duration")
			if err != nil {
				return fmt.Errorf("監控時間輸入格式錯誤, %v", err)
			}

			logic := c.BaseLogicFunc(l)

			if orgfunc == nil {
				orgfunc = logic.SendMessgeFunc
			}

			l.OnSendMessage(func(userid int32, m *msg.Message) {
				c.mu.Lock()
				if val, ok := c.traceMap[tableno]; ok && val {
					PrintMessage(os.Stdout, m)
				}
				c.mu.Unlock()
				orgfunc(userid, m)
			})

			ctx, cancel := context.WithTimeout(context.Background(), d)
			defer cancel()

			go func() {
				c.mu.Lock()
				defer c.mu.Unlock()
				c.traceMap[tableno] = true
				cmd.Printf("%s 開始追蹤遊戲封包, 桌號: %d, 時間: %v\n", time.Now().Format("2006-01-02 15:04:05"), tableno, d)
			}()

			select {
			case <-ctx.Done():
				c.mu.Lock()
				defer c.mu.Unlock()
				c.traceMap[tableno] = false
				cmd.Printf("%s 停止追蹤遊戲封包, 桌號: %d\n", time.Now().Format("2006-01-02 15:04:05"), tableno)
			}

			return nil
		},
	}
	traceCmd.PersistentFlags().Int32P("table", "t", 0, "桌號")
	traceCmd.PersistentFlags().DurationP("duration", "d", time.Minute, "追蹤時間")

	traceCmd.MarkPersistentFlagRequired("table")

	return traceCmd
}

func (c *SlotHelperCmd) ptnctrlCmd() *cobra.Command {

	defaultDuring := 5 * time.Minute

	ptnctrlCmd := &cobra.Command{
		Use:                   "ptnctrl <userid> <0/1/2>",
		Short:                 "控制玩家水線, 0: 基本(預設), 1: 放分, 2: 殺分",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			if c.BaseLogicFunc == nil {
				return fmt.Errorf("沒有設定取得房間的程式")
			}

			if len(args) != 2 {
				return fmt.Errorf("參數錯誤")
			}

			uid, err := strconv.ParseInt(args[0], 0, 32)
			if err != nil {
				return fmt.Errorf("玩家編號輸入格式錯誤, %v", err)
			}
			p := c.gs.GameServerRunCtx.PlayerManager.GetPlayerByUserID(int32(uid))
			if p == nil {
				return fmt.Errorf("查無玩家 %d", int32(uid))
			}

			l := c.glm.GetGameLogic(p.UserData.(*gamelogic.GameUser).TableNo)
			if l == nil {
				return fmt.Errorf("沒有設置房間")
			}

			logic := c.BaseLogicFunc(l)

			ptn, err := strconv.ParseInt(args[1], 0, 8)
			if err != nil {
				return fmt.Errorf("水線模式輸入格式錯誤, %v", err)
			}

			mode := baselogic.PatternMode(uint8(ptn))

			var during time.Duration
			d, err := cmd.PersistentFlags().GetDuration("during")
			if err != nil {
				during = defaultDuring
			} else {
				during = d
			}

			logic.ModeExpire = time.Now().Add(during)

			logic.SwitchProto(mode)
			logic.AssignUserMode(p.UserID, mode, logic.ModeExpire)

			return nil
		},
	}
	ptnctrlCmd.PersistentFlags().DurationP("during", "d", defaultDuring, "維持時間")

	return ptnctrlCmd
}

// GameLogic 取得遊戲邏輯
func (c *SlotHelperCmd) GameLogic(tid int32) gamelogic.GameLogic {
	return c.glm.GetGameLogic(tid)
}
