package mylog

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "example.invalid/go-sql-driver/mysql" // driver
	"gopkg.in/mgo.v2"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// MyDBQueue ?
var MyDBQueue IDatabaseQueue

// ServerID ?
var ServerID int

// IDatabaseQueue ?
type IDatabaseQueue interface {
	addQueue(s *Statement)
	Start()
	writeDB(stmts []*Statement)
	getCodeStack() interface{}
	query(in interface{}) interface{}
}

// InitDatabaseQueue ?
func InitDatabaseQueue(connString string, connPool int, perWriteCount int, serverID int, log *mlog.DefaultLogger) {
	ServerID = serverID
	if connString != "" {
		MyDBQueue = NewDatabaseQueue(connString, connPool, perWriteCount, log)
	} else {
		MyDBQueue = NewEmptyDatabaseQueue()
	}
	MyDBQueue.Start()
}

// InitMongoDatabaseQueue ?
func InitMongoDatabaseQueue(connString string, connPool int, perWriteCount int, serverID int, log *mlog.DefaultLogger, database string, collection string) {
	ServerID = serverID
	if connString != "" && database != "" && collection != "" {
		MyDBQueue = NewMongoDatabaseQueue(connString, connPool, perWriteCount, log, database, collection)
	} else {
		MyDBQueue = NewEmptyDatabaseQueue()
	}
	MyDBQueue.Start()
}

type databaseQueue struct {
	db            *sql.DB
	queue         chan *Statement
	perWriteCount int
	log           *mlog.DefaultLogger
}

type mongoDatabaseQueue struct {
	db            *mgo.Database
	coll          *mgo.Collection
	queue         chan *Statement
	perWriteCount int
	log           *mlog.DefaultLogger
}

// NewDatabaseQueue ?
func NewDatabaseQueue(connString string, connPool int, perWriteCount int, log *mlog.DefaultLogger) IDatabaseQueue {
	db, _ := sql.Open("mysql", connString)
	db.SetMaxOpenConns(connPool)
	db.SetMaxIdleConns((connPool / 2) + 1)
	q := new(databaseQueue)
	q.db = db
	q.queue = make(chan *Statement, 10000)
	q.perWriteCount = perWriteCount
	q.log = log
	return q
}

// NewMongoDatabaseQueue ?
func NewMongoDatabaseQueue(connString string, connPool int, perWriteCount int, log *mlog.DefaultLogger, database string, collection string) IDatabaseQueue {
	uri := "mongodb://" + connString
	session, err := mgo.Dial(uri)
	if err != nil {
		mlog.Error("MongoDB Error : %v", err.Error())
	}

	db := session.DB(database)
	coll := db.C(collection)

	q := new(mongoDatabaseQueue)
	q.db = db
	q.coll = coll
	q.queue = make(chan *Statement, 10000)
	q.perWriteCount = perWriteCount
	q.log = log
	return q
}

func (q *databaseQueue) addQueue(s *Statement) {

	//code stack 程式func呼叫堆疊來源
	dataStr, _ := json.Marshal(q.log.GetCallStack())

	s.args = append(s.args, dataStr)
	q.queue <- s
	q.log.Reset()
}

func (q *databaseQueue) Start() {
	go func() {
		defer util.PanicErrStack()
		for {
			stmts := make([]*Statement, 0)
		FORLOOP:
			for i := 0; i < q.perWriteCount; i++ {

				select {
				case mystmt := <-q.queue:
					stmts = append(stmts, mystmt)
				default:
					break FORLOOP
				}
			}
			go func() {
				if len(stmts) > 0 {
					q.writeDB(stmts)
				}
			}()
			if len(stmts) >= q.perWriteCount {
				mlog.Error("bylog 每秒己經超過 %d 條, 還有 %d 條", q.perWriteCount, len(q.queue))
				continue
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func (q *databaseQueue) writeDB(stmts []*Statement) {
	tx, err := q.db.Begin()
	if err != nil {
		mlog.Error("mylog error %s", err)
		return
	}

	for _, stmt := range stmts {
		tx.Exec(stmt.sql, stmt.args...)
	}
	tx.Commit()
}

func (q *databaseQueue) getCodeStack() interface{} {
	return q.log.GetCallStack()
}
func (q *databaseQueue) query(in interface{}) interface{} {
	return nil
}

// empty databasequeue
type emptyDatabaseQueue struct{}

func (q *emptyDatabaseQueue) addQueue(s *Statement)      {}
func (q *emptyDatabaseQueue) Start()                     {}
func (q *emptyDatabaseQueue) writeDB(stmts []*Statement) {}
func (q *emptyDatabaseQueue) getCodeStack() interface{} {
	return nil
}
func (q *emptyDatabaseQueue) query(in interface{}) interface{} {
	return nil
}

// NewEmptyDatabaseQueue ?
func NewEmptyDatabaseQueue() IDatabaseQueue { return new(emptyDatabaseQueue) }

//MongoDB
func (q *mongoDatabaseQueue) addQueue(s *Statement) {
	//code stack 程式func呼叫堆疊來源
	dataStr := q.log.GetCallStack()
	s.args = append(s.args, dataStr)
	q.queue <- s
	q.log.Reset()
}
func (q *mongoDatabaseQueue) getCodeStack() interface{} {
	cloneLog := new(mlog.DefaultLogger)
	cloneLog.StdLog = q.log.StdLog
	cloneLog.CallDepLen = q.log.CallDepLen
	cloneLog.CallDepStart = q.log.CallDepStart
	return cloneLog.GetCallStack()
}

func (q *mongoDatabaseQueue) Start() {
	go func() {
		defer util.PanicErrStack()
		for {
			stmts := make([]*Statement, 0)
		FORLOOP:
			for i := 0; i < q.perWriteCount; i++ {

				select {
				case mystmt := <-q.queue:
					stmts = append(stmts, mystmt)
				default:
					break FORLOOP
				}
			}
			go func() {
				if len(stmts) > 0 && q.coll.Database.Session != nil {
					q.writeDB(stmts)
				}
			}()
			if len(stmts) >= q.perWriteCount {
				mlog.Error("bylog 每秒己經超過 %d 條, 還有 %d 條", q.perWriteCount, len(q.queue))
				continue
			}
			time.Sleep(1 * time.Second)
		}
	}()
}
func (q *mongoDatabaseQueue) writeDB(stmts []*Statement) {
	defer util.PanicErrStack()
	//mlog.System(ConnString)
	c := q.coll

	for _, stmt := range stmts {
		err := c.Insert(stmt.monglogargs)
		if err != nil {
			c.Database.Session.Refresh()
			panic(err)
		}
	}
}

func (q *mongoDatabaseQueue) query(query interface{}) interface{} {
	result := QueryMessage{}
	c := q.coll
	c.Find(query).Sort("-frameid").One(&result)
	return result
}

// QueryMessage ?
type QueryMessage struct {
	ServerID   int32 `bson:"serverid"`
	UserID     int32 `bson:"userid"`
	FrameID    int64 `bson:"frameid"`
	Msg        `bson:"message"`
	TableNo    int32     `bson:"tableno"`
	CreateTime time.Time `bson:"createtime"`
}

// Msg ?
type Msg struct {
	Seq int64 `bson:"seq"`
}
