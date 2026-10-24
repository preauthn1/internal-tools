package db

import "example.invalid/ServerTeam/Servers/pkg/mlog"

// Parameters 參數數組
type Parameters struct {
	innerParameters []*SQLParameter
}

// NewSQLParameters 創建Sql參數集合
func NewSQLParameters() *Parameters {
	ps := new(Parameters)

	return ps
}

// AddParameter 添加参数
func (ps *Parameters) AddParameter(p *SQLParameter) *Parameters {
	if p == nil {
		mlog.Error("傳入的參數為 nil")
		return nil
	}

	ps.innerParameters = append(ps.innerParameters, p)

	return ps
}

// AddIntInput 添加int參數
func (ps *Parameters) AddIntInput(name string, value int32) *Parameters {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := NewSQLIntParameter(name, value, 4)

	ps.AddParameter(p)

	return ps
}

// AddIntOutput 添加int輸出參數
func (ps *Parameters) AddIntOutput(name string, value int32) *Parameters {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := NewSQLOutputParameter(name, value, SQLDBTypeInt, 4)

	ps.AddParameter(p)

	return ps
}

// AddTinyIntInput 添加tinyint參數
func (ps *Parameters) AddTinyIntInput(name string, value byte) *Parameters {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := NewSQLParameter(name, value, SQLDBTypeTinyint, ParamDirectionInput, 1)

	ps.AddParameter(p)

	return ps
}

// AddTinyIntOutput 添加tinyint輸出參數
func (ps *Parameters) AddTinyIntOutput(name string, value byte) *Parameters {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := NewSQLOutputParameter(name, value, SQLDBTypeTinyint, 1)

	ps.AddParameter(p)

	return ps
}

// AddSmallIntInput 添加smallint參數
func (ps *Parameters) AddSmallIntInput(name string, value int16) *Parameters {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := NewSQLParameter(name, value, SQLDBTypeSmallint, ParamDirectionInput, 2)

	ps.AddParameter(p)

	return ps
}

// AddSmallIntOutput 添加smallint輸出參數
func (ps *Parameters) AddSmallIntOutput(name string, value int16) *Parameters {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := NewSQLOutputParameter(name, value, SQLDBTypeSmallint, 2)

	ps.AddParameter(p)

	return ps
}

// AddBigIntInput 添加bigint參數
func (ps *Parameters) AddBigIntInput(name string, value int64) *Parameters {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := NewSQLParameter(name, value, SQLDBTypeBigint, ParamDirectionInput, 8)

	ps.AddParameter(p)

	return ps
}

// AddBigIntOutput 添加bigint輸出參數
func (ps *Parameters) AddBigIntOutput(name string, value int64) *Parameters {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := NewSQLOutputParameter(name, value, SQLDBTypeBigint, 8)

	ps.AddParameter(p)

	return ps
}

// AddVarcharInput 添加varchar參數
func (ps *Parameters) AddVarcharInput(name string, value string, size int32) *Parameters {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := NewSQLParameter(name, value, SQLDBTypeVarchar, ParamDirectionInput, size)

	ps.AddParameter(p)

	return ps
}

// AddVarcharOutput 添加varchar輸出參數
func (ps *Parameters) AddVarcharOutput(name string, value string) *Parameters {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := NewSQLOutputParameter(name, value, SQLDBTypeVarchar, 256)

	ps.AddParameter(p)

	return ps
}

// AddVarcharSizeOutput 添加varchar輸出參數
func (ps *Parameters) AddVarcharSizeOutput(name string, value string, size int32) *Parameters {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := NewSQLOutputParameter(name, value, SQLDBTypeVarchar, size)

	ps.AddParameter(p)

	return ps
}

// OutputParamCount 輸出參數數量
func (ps *Parameters) OutputParamCount() int {
	cn := 0

	for _, p := range ps.innerParameters {
		if p.Direction == ParamDirectionOutput {
			cn++
		}
	}

	return cn
}
