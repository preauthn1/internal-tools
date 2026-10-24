package defender

import (
	"fmt"
	"os"
	"strconv"

	"example.invalid/spf13/cobra"
)

// JackpotDefenderCmd 查詢JackpotDefender指令
type JackpotDefenderCmd struct{}

// Name Command name
func (c *JackpotDefenderCmd) Name() string {
	return "jp-def"
}

// Help Command Description
func (c *JackpotDefenderCmd) Help() string {
	return fmt.Sprintf("Jackpot 防護工具")
}

// Exec Command Logic
func (c *JackpotDefenderCmd) Exec(args []string) string {
	cmd := c.CobraRoot()
	cmd.SetArgs(args)
	cmd.Execute()
	return ""
}

// CobraRoot 取得 cobra command
func (c *JackpotDefenderCmd) CobraRoot() *cobra.Command {
	cmd := c.rootCommand()
	cmd.AddCommand(c.infoCommand())
	cmd.AddCommand(c.incCommand())

	cmd.SetOutput(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.SetIn(os.Stdin)

	return cmd
}

func (c *JackpotDefenderCmd) rootCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "jp-def",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if JPDefender == nil {
				return fmt.Errorf("此遊戲尚未設定 Jackpot Defender")
			}
			cmd.Usage()
			return nil
		},
	}
}

func (c *JackpotDefenderCmd) infoCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "info",
		Short:                 "取得 Jackpot Defender 資訊",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if JPDefender == nil {
				return fmt.Errorf("此遊戲尚未設定 Jackpot Defender")
			}
			example.invalid(cmd.OutOrStdout())
			return nil
		},
	}
}

func (c *JackpotDefenderCmd) incCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "inc <amount>",
		Short:                 "增加 Pool",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if JPDefender == nil {
				return fmt.Errorf("此遊戲尚未設定 Jackpot Defender")
			}

			i, err := strconv.ParseInt(args[0], 0, 64)
			if err != nil {
				return fmt.Errorf("Contribution 格式錯誤, %v", err)
			}

			JPDefender.contribution += i * 1e4

			example.invalid(cmd.OutOrStdout())
			return nil
		},
	}
}
