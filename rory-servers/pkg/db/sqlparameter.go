package db

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// ParamDirection
const (
	ParamDirectionInput = iota
	ParamDirectionOutput
)

// SQL DB Type
const (
	SQLDBTypeVarchar = iota
	SQLDBTypeNVarchar
	SQLDBTypeInt
	SQLDBTypeTinyint
	SQLDBTypeBigint
	SQLDBTypeSmallint
	SQLDBTypeChar
	SQLDBTypeNChar
	SQLDBTypeDecimal
	SQLDBTypeDateTime
)

// SQLParameter sql參數
type SQLParameter struct {
	Name      string      // 參數名
	Value     interface{} // 參數值
	DBType    int32       // 參數類型
	Direction byte        // 方向
	Size      int32
}

// GetParamType 獲取參數類型描述字符串
func (p *SQLParameter) GetParamType() string {
	ptype := ""

	if p.DBType == SQLDBTypeBigint {
		ptype = "bigint"
	} else if p.DBType == SQLDBTypeChar {
		ptype = fmt.Sprintf("char(%d)", p.Size)
	} else if p.DBType == SQLDBTypeDateTime {
		ptype = "datetime"
	} else if p.DBType == SQLDBTypeDecimal {
		ptype = "decimal(18,2)"
	} else if p.DBType == SQLDBTypeInt {
		ptype = "int"
	} else if p.DBType == SQLDBTypeNChar {
		ptype = fmt.Sprintf("nchar(%d)", p.Size)
	} else if p.DBType == SQLDBTypeNVarchar {
		ptype = fmt.Sprintf("nvarchar(%d)", p.Size)
	} else if p.DBType == SQLDBTypeSmallint {
		ptype = "smallint"
	} else if p.DBType == SQLDBTypeTinyint {
		ptype = "tinyint"
	} else if p.DBType == SQLDBTypeVarchar {
		ptype = fmt.Sprintf("varchar(%d)", p.Size)
	}

	return ptype
}

// GetParameterValue 獲取參數值描述字符串
func (p *SQLParameter) GetParameterValue() string {
	pvalue := ""

	if p.DBType == SQLDBTypeBigint || p.DBType == SQLDBTypeDecimal || p.DBType == SQLDBTypeInt || p.DBType == SQLDBTypeSmallint || p.DBType == SQLDBTypeTinyint {
		pvalue = fmt.Sprintf("%v", p.Value)
	} else {
		pvalue = fmt.Sprintf("'%v'", p.Value)
	}

	return pvalue
}

// NewSQLParameter 創建sql參數
func NewSQLParameter(name string, value interface{}, dbtype int32, direction byte, size int32) *SQLParameter {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := new(SQLParameter)
	p.Name = name
	p.Value = value
	p.DBType = dbtype
	p.Direction = direction
	p.Size = size

	return p
}

// NewSQLInputParameter 創建sql輸入參數
func NewSQLInputParameter(name string, value interface{}, dbtype int32, size int32) *SQLParameter {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := new(SQLParameter)
	p.Name = name
	p.Value = value
	p.DBType = dbtype
	p.Direction = ParamDirectionInput
	p.Size = size

	return p
}

// NewSQLOutputParameter 創建sql輸出參數
func NewSQLOutputParameter(name string, value interface{}, dbtype int32, size int32) *SQLParameter {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := new(SQLParameter)
	p.Name = name
	p.Value = value
	p.DBType = dbtype
	p.Direction = ParamDirectionOutput
	p.Size = size

	return p
}

// NewSQLVarcharParameter 創建sql varchar參數
func NewSQLVarcharParameter(name string, value string, size int32) *SQLParameter {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := new(SQLParameter)
	p.Name = name
	p.Value = value
	p.Direction = ParamDirectionInput
	p.DBType = SQLDBTypeVarchar
	p.Size = size

	return p
}

// NewSQLIntParameter 創建sql int類型參數
func NewSQLIntParameter(name string, value int32, size int32) *SQLParameter {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := new(SQLParameter)
	p.Name = name
	p.Value = value
	p.Direction = ParamDirectionInput
	p.DBType = SQLDBTypeInt
	p.Size = size

	return p
}

// NewSQLDateTimeParameter 創建sql datetime類型參數
func NewSQLDateTimeParameter(name string, value string, size int32) *SQLParameter {
	if name == "" {
		mlog.Error("傳入的欄位名稱為空字串")
		return nil
	}

	p := new(SQLParameter)
	p.Name = name
	p.Value = value
	p.Direction = ParamDirectionInput
	p.DBType = SQLDBTypeDateTime
	p.Size = size

	return p
}
