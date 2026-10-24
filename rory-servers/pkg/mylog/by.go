package mylog

import "time"

// Statement ?
type Statement struct {
	sql         string
	args        []interface{}
	monglogargs MongoLogData
}

// MongoLogData ?
type MongoLogData struct {
	ServerID   int         `bson:"serverid"`
	UserID     int32       `bson:"userid"`
	FrameID    int64       `bson:"frameid"`
	Message    interface{} `bson:"message"`
	TableNo    int32       `bson:"tableno"`
	CodeStack  interface{} `bson:"codestack"`
	CreateTime time.Time   `bson:"createtime"`
}

// MongoLog ?
func MongoLog(userid int32, frameid int64, tableno int32, message interface{}) {
	//sql := "INSERT INTO bylog (serverid, userid, frameid, message, tableno, codestack) VALUES (?, ?, ?, ?, ?, ?)"

	logData := MongoLogData{
		ServerID:   ServerID,
		UserID:     userid,
		FrameID:    frameid,
		Message:    message,
		TableNo:    tableno,
		CodeStack:  MyDBQueue.getCodeStack(),
		CreateTime: time.Now(),
	}

	query := &Statement{monglogargs: logData}

	MyDBQueue.addQueue(query)
}

// AddLog ?
func AddLog(userid int32, tableno int32, message interface{}) {
	logData := MongoLogData{
		ServerID:   ServerID,
		UserID:     userid,
		Message:    message,
		TableNo:    tableno,
		CodeStack:  MyDBQueue.getCodeStack(),
		CreateTime: time.Now(),
	}

	query := &Statement{monglogargs: logData}

	MyDBQueue.addQueue(query)
}

// BYQuery ?
func BYQuery(query interface{}) interface{} {
	return MyDBQueue.query(query)
}
