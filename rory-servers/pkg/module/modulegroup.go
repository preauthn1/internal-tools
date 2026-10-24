package module

import (
	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// Group 模塊組
type Group struct {
	modules map[string]Module
}

// NewGroup 建立模塊組
func NewGroup() *Group {
	mg := new(Group)
	mg.modules = make(map[string]Module)

	//添加模块命令
	cmd := NewCmd(mg)
	command.AddCommand(cmd)

	return mg
}

// Register 註冊模塊
func (mg *Group) Register(m Module) {
	if _, ok := mg.modules[m.Name()]; ok {
		mlog.Printf("模塊%s已經註冊", m.Name())
		return
	}

	mg.modules[m.Name()] = m

	return
}

// UnLoad 卸載模塊
func (mg *Group) UnLoad(mname string) Module {
	m := mg.GetModule(mname)
	if m == nil {
		return nil
	}

	m.Stop()

	return m
}

// GetModule 獲取模塊
func (mg *Group) GetModule(name string) Module {
	if m, ok := mg.modules[name]; ok {
		return m
	}

	return nil
}

// Range 遍歷所有模塊
func (mg *Group) Range(f func(m Module)) {
	for _, m := range mg.modules {
		f(m)
	}
}

// Run 運行
func (mg *Group) Run(ctx *MouduleCtx) error {
	for _, m := range mg.modules {
		err := m.Run(ctx)

		if err != nil {
			return err
		}
	}

	return nil
}

// Stop stop
func (mg *Group) Stop(ctx *MouduleCtx) {
	for _, m := range mg.modules {
		m.Stop()
	}
}

// Name 名稱
func (mg *Group) Name() string {
	return "modulegroup"
}
