package db

import (
	"database/sql"
	"fmt"
)

// Mssql2000 ?
type Mssql2000 struct {
	*sql.DB
	dataSource string
	database   string
	windows    bool
	sa         *SA
}

// SA ?
type SA struct {
	user   string
	passwd string
	port   int
}

// NewMssql2000 ?
func NewMssql2000(host string, port int32, user string, password string, database string) *Mssql2000 {
	mssql := new(Mssql2000)
	mssql = &Mssql2000{
		// 如果数据库是默认实例（MSSQLSERVER）则直接使用IP，命名实例需要指明。
		// dataSource: "0.0.0.0\\MSSQLSERVER",
		dataSource: host,
		database:   database,
		// windows: true 为windows身份验证，false 必须设置sa账号和密码
		windows: false,
		sa: &SA{
			user:   user,
			passwd: [REDACTED],
			port:   int(port),
		},
	}

	return mssql
}

// Open ?
func (m *Mssql2000) Open() error {
	config := fmt.Sprintf("Provider=SQLOLEDB;Initial Catalog=%s;Data Source=%s",
		m.database, m.dataSource)

	if m.windows {
		config = fmt.Sprintf("%s;Integrated Security=SSPI", config)
	} else {
		// sql 2000的端口寫法和sql 2005以上的有所不同，在Data Source 後以逗號隔開。
		config = fmt.Sprintf("%s,%d;user id=%s;password=%s",
			config, m.sa.port, m.sa.user, m.sa.passwd)
	}

	var err error
	m.DB, err = sql.Open("adodb", config)
	//fmt.Println(config)

	return err
}

// ExecSQL ?
func (m *Mssql2000) ExecSQL(sql string, args ...interface{}) (*RecordTable, error) {
	m.Open()
	defer m.Close()
	rows, err := m.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	if cols == nil {
		return nil, nil
	}

	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
	}

	rettable := &RecordTable{}
	rettable.Cols = cols
	for rows.Next() {
		err = rows.Scan(vals...)

		if err != nil {
			return nil, err
		}

		retrow := RecordRow{}

		for _, r := range vals {
			v := *r.(*interface{})

			retrow = append(retrow, v)
		}

		rettable.Rows = append(rettable.Rows, retrow)
	}

	return rettable, err
}
