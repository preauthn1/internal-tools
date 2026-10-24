package by

type byLoadCfCmd struct {
	by *Bylogic
}

func NewByCmd(by *Bylogic) *byLoadCfCmd {
	cmd := new(byLoadCfCmd)
	cmd.by = by
	return cmd
}

func (byCmd *byLoadCfCmd) Name() string {
	return "byreload"
}

func (byCmd *byLoadCfCmd) Help() string {
	return "更新捕魚TableData. eg:bydata "
}

func (byCmd *byLoadCfCmd) Exec(args []string) string {
	if len(args) == 0 {
		return "error. seeref usage"
	}

	cmdname := args[0]

	switch cmdname {
	case "bydata":
		for _, table := range byCmd.by.TableManager.Room.Tables {
			glc := table.GameLogic.(*Bylogic)
			glc.loadTableData()
		}
		return "reload success"
	}

	return ""
}
