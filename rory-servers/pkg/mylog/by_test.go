package mylog

//func TestWorker(t *testing.T) {
//	ConnString := "root:123456@tcp(0.0.0.0:3306)/mylog_dev?charset=utf8"
//	//ConnPool = 20
//	InitDatabaseQueue(ConnString, 0, 0, 0, mlog.NewDefaultLogger())
//	mlog.System("test worker start")
//	n := 50
//	JobChans := make(chan int64, n*2)
//	DoneChans := make(chan bool, n*2)
//	exitChan := make(chan bool, 1)
//	userid := int32(1)
//	for w := 1; w <= n; w++ {
//		go func(id int, jChans <-chan int64, dChans chan<- bool) {
//			for frameid := range jChans {
//				_ = userid
//				_ = frameid
//				//mlog.System("Run Worker [%v] job %d", id, frameid)
//				MongoLog(userid, frameid, 0, "Run Worker")
//				dChans <- true
//			}
//		}(w, JobChans, DoneChans)
//	}
//	jobsNum := 1000000
//	go func() {
//		for {
//			<-DoneChans
//		}
//	}()
//
//	//limiter := time.Tick(time.Duration(10) * time.Millisecond)
//	for i := 0; i < jobsNum; i++ {
//		//<-limiter
//		JobChans <- int64(i)
//		//mlog.System("============job %d", i)
//	}
//	mlog.System("test worker finish")
//	close(JobChans)
//	select {
//	case <-exitChan:
//	}
//}
//
//func TestMgo(t *testing.T) {
//	ConnString := "chiou.w:terpro123!$!@0.0.0.0:27017/admin"
//	database := "FishLog_Dev"
//	collection := "FishLogFlag1"
//
//	InitMongoDatabaseQueue(ConnString, 100, 100, 0, mlog.NewDefaultLogger(), database, collection)
//
//	userid := int32(1)
//	//frameid := int64(0)
//	//MongoLog(userid, frameid, 0, "Run Worker")
//
//	mlog.System("test worker start")
//	n := 50
//	JobChans := make(chan int64, n*2)
//	DoneChans := make(chan bool, n*2)
//	exitChan := make(chan bool, 1)
//
//	for w := 1; w <= n; w++ {
//		go func(id int, jChans <-chan int64, dChans chan<- bool) {
//			for frameid := range jChans {
//				_ = userid
//				_ = frameid
//				//mlog.System("Run Worker [%v] job %d", id, frameid)
//				AddLog(userid, 0, "Run Worker")
//				dChans <- true
//			}
//		}(w, JobChans, DoneChans)
//	}
//	jobsNum := 1 //1000000
//	go func() {
//		for {
//			<-DoneChans
//		}
//	}()
//
//	//limiter := time.Tick(time.Duration(10) * time.Millisecond)
//	for i := 0; i < jobsNum; i++ {
//		//<-limiter
//		JobChans <- int64(i)
//		//mlog.System("============job %d", i)
//	}
//	mlog.System("test worker finish")
//	close(JobChans)
//	select {
//	case <-exitChan:
//	}
//}
