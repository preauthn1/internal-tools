package fixtures

import (
	"path"
	"runtime"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/mylog"
)

// MongoInitByTest ?
func MongoInitByTest() {
	_, filename, _, _ := runtime.Caller(1)
	mlog.System("[weihao] filename[%v]", filename)
	configDir := path.Dir(filename) + "/../../config"
	mlog.System("[weihao] configdir[%v]", configDir)
	cfgs := new(config.Group)
	cfgs.LoadConfig(configDir, 0)
	serverInfo := cfgs.GetConfig(gamelogic.ConfigFileNameServerInfo)
	ConnString := serverInfo.GetConfigValue("mongoconn", "connstring")
	database := serverInfo.GetConfigValue("mongoconn", "database")
	collection := serverInfo.GetConfigValue("mongoconn", "collection")
	InitMongoDatabaseQueue(ConnString, 100, 100, 0, mlog.NewDefaultLogger(), database, collection)
}

// InitMongoDatabaseQueue ?
func InitMongoDatabaseQueue(connString string, connPool int, perWriteCount int, serverID int, log *mlog.DefaultLogger, database string, collection string) {
	mylog.ServerID = serverID
	if connString != "" && database != "" && collection != "" {
		mylog.MyDBQueue = mylog.NewMongoDatabaseQueue(connString, connPool, perWriteCount, log, database, collection)
	} else {
		mylog.MyDBQueue = mylog.NewEmptyDatabaseQueue()
	}
	mylog.MyDBQueue.Start()
}
