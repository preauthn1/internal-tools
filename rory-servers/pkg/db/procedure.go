package db

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// Procedure 存儲過程
type Procedure struct {
	Name       string // 存儲過程名稱
	Parameters *Parameters
}

// NewProcedure 創建存儲過程
func NewProcedure(name string, ps *Parameters) *Procedure {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	if ps == nil {
		mlog.Error("傳入參數資料為 nil")
		return nil
	}

	p := new(Procedure)
	p.Name = name
	p.Parameters = ps

	return p
}

// GetSQLString 獲取sql字符串
func (p *Procedure) GetSQLString() string {
	sqlstr := ""

	// 輸出參數申明
	for idx, pr := range p.Parameters.innerParameters {
		if pr.Direction == ParamDirectionOutput {
			sqlstr = sqlstr + fmt.Sprintf("declare @p%d %s;", idx, pr.GetParamType())

			if pr.Value != nil {
				sqlstr += fmt.Sprintf("set @p%d=%s;", idx, pr.GetParameterValue())
			}

		}
	}

	sqlstr += fmt.Sprintf("declare @intProcRet int; exec @intProcRet = %s ", p.Name)
	end := len(p.Parameters.innerParameters)
	for idx, pr := range p.Parameters.innerParameters {

		if pr.Direction == ParamDirectionInput {
			sqlstr += fmt.Sprintf("@%s=%s", pr.Name, pr.GetParameterValue())
		} else {
			sqlstr += fmt.Sprintf("@%s=@p%d output", pr.Name, idx)
		}

		if idx != end-1 {
			sqlstr += ","
		}
	}

	// 輸出參數申明
	optsize := p.Parameters.OutputParamCount()
	oidx := 0
	if optsize > 0 {
		sqlstr += ";select "

		for idx, pr := range p.Parameters.innerParameters {
			if pr.Direction == ParamDirectionOutput {
				sqlstr += fmt.Sprintf("@p%d %s", idx, pr.Name)
				oidx++

				if oidx != optsize {
					sqlstr += ","
				}
			}
		}
	}

	sqlstr += ";select @intProcRet;"
	return sqlstr
}
