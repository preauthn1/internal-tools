package by

import (
	"fmt"
	"example.invalid/ServerTeam/Servers/internal/game/by/DataCfg"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/mylog"
	"example.invalid/ServerTeam/Servers/pkg/util"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type byfishdata struct {
	fishSeq     int32 //鱼的流水號id
	fishtype    int32 //鱼的种类
	fishhit     int32 //攻击次数
	returncount int32 //回游次数
	fishnumber  int32 //魚的id
	fishcount   int32 //魚的數量
	fishodd     int32 //倍率
}

type fishweight struct {
	fishid int32 //鱼种类 fishSeq or shalId魚群
	weight int32 //鱼权重
}

type byfishmanager struct {
	fishdata             map[int32]*byfishdata //所有鱼数据
	fishweightarr        []*fishweight         //生成鱼权重数组
	config               *fishconfig           //配置文件
	fishwavecf           *fishwaveconfig       //鱼潮配置
	fishSeq              int32                 //新生成鱼的id，注ID不能为0，客户端不能识别ID是空还是0
	fishwavetime         int32                 //生成鱼潮间隔
	fishwave             bool                  //现在生成的是鱼潮
	ice                  bool                  //现在在冰冻状态
	fishwavetype         int32                 //鱼潮的类型
	fishwavecount        int32                 //魚潮數量
	enterfishwavetime    int32                 //进入鱼潮的海浪时间
	fishwavelasttime     int32                 //鱼潮持续时间
	bossfishweight       []*fishweight         //boss鱼权重
	bossentertime        int32                 //boss鱼入场提示
	bossfishtime         int32                 //boss鱼出现的时间
	teshufishweight      []*fishweight         //特殊鱼权重
	teshufishtime        int32                 //特殊鱼出现的时间
	fishshoaltime        int32                 //魚群出現時間
	goldfishtime         int32
	totalfishweight      int64         //普通鱼的总权重
	totalbossfishweight  int64         //BOSS鱼的总权重
	totalteshufishweight int64         //特殊鱼的总权重
	fishshoalweight      []*fishweight //魚群權重
	totalfishshoalweight int64         //魚群權重
	goldfishweight       []*fishweight
	totalgoldfishweight  int64
}

//所有炮彈管理器(BoomTable)
type byboommanager struct {
	boomNumber    int32   //砲彈編號
	boomName      string  //砲彈名稱
	boomUIType    int32   //砲彈UI類型
	boomOdds      int32   //倍率
	boomMoney     int32   //砲彈金幣
	boomSpeed     float32 //子彈速度(30ms的移動間隔)
	shootInterval float32 //子彈發射間隔(ms)
	boomType      int32   //boomType
	boomMaxLimit  int32   //子彈上限
	rebound       int32   //反彈次數
}

//當魚出現指定次數增加魚死亡率
type fishdynamichitratemanager struct {
}

type byicemanager struct {
	iceNumber    int32  //冰凍編號
	iceName      string //冰凍名稱
	iceUIType    int32  //冰凍UI類型
	iceTime      int32  //冰凍時間(s)
	iceResetTime int32  //冰冻冷却时间(s)
	iceMoney     int32  //使用金幣
}

//新建鱼管理器实例
func NewFishManager() *byfishmanager {
	m := new(byfishmanager)
	m.fishdata = make(map[int32]*byfishdata)
	m.fishweightarr = make([]*fishweight, 0)
	m.bossfishweight = make([]*fishweight, 0)
	m.teshufishweight = make([]*fishweight, 0)
	m.config = newfishconfi()
	m.fishwavecf = newFishWaveConfig()
	m.fishSeq = 1
	m.fishwave = false
	m.ice = false
	return m
}

//生成一条新的鱼
func (m *byfishmanager) newfish(by *Bylogic) (*byfishdata, int32) {
	var fish *byfishdata

	t := m.getCreateFishType()
	fi, find := m.config.getfishinfobytype(t)

	if find {
		mlog.Debug("需要生成一条没达到上限的鱼 fi.maxcount[%v]", fi.maxcount)
		//需要生成一条没达到上限的鱼，最多循环10次
		for i := 0; i < 10; i++ {
			if find && fi.maxcount > 0 &&
				m.getFishCountByType(fi.fishid) >= fi.maxcount {
				t = m.getCreateFishType()
				fi, find = m.config.getfishinfobytype(t)
				//if find {
				//	mlog.Debug(" 新增魚1 fishSeq[%v] maxcount[%v]", fi.fishid, fi.maxcount)
				//}

				continue
			}

			if !find {
				t = m.getCreateFishType()
				fi, find = m.config.getfishinfobytype(t)
				continue
			}

			break
		}
		//10次还没找到一条新的鱼就选择第一个鱼作为鱼的模板
		if fi == nil {
			fi = m.config.fishinfoarr[0]
		}

		if m.getFishCountByType(fi.fishid) < fi.maxcount && !by.isTotalFishMaxLimit() {
			//mlog.Debug("m.getFishCountByType(fi.fishid)[%v]  fi.maxcount[%v] !by.isTotalFishMaxLimit()[%v] ",
			//	m.getFishCountByType(fi.fishid), fi.maxcount, !by.isTotalFishMaxLimit())
			fish = m.createFish(fi, by)
		}
	} else {
		//生成鱼群
		fish = m.createFishShoal(t, by)
		mlog.Debug("生成鱼群 fishSeq[%v] fishGUID[%v] fishnumber[%v] fishcount[%v]", fish.fishSeq, fish.fishtype, fish.fishnumber, fish.fishcount)
	}

	return fish, m.config.getcreatetime()
}

//创建一条BOSS鱼
func (m *byfishmanager) newBossFish(by *Bylogic) *byfishdata {
	t := m.getCreateBossFishType()
	fi, ok := m.config.getfishinfobytype(t)
	if !ok || fi == nil {
		fi = m.config.fishinfoarr[0]
	}

	if by.fishmanager.getFishCountByType(fi.fishid) >= fi.maxcount {
		return nil
	}

	//m.setNewFishLuckyBonus(by, fi)

	mlog.Debug("创建一条BOSS鱼 fishSeq[%v] luckybonus[%v]", fi.fishid, fi.luckybonus)
	return m.createFish(fi, by)
}

//設定新產生魚的LuckyNumber
func (m *byfishmanager) setNewFishLuckyBonus(by *Bylogic, fi *fishinfo) {
	for _, v := range by.fishtypetable.Items {
		if fi.fishid == int32(v.GetID()) {
			sum := int32(0)
			for _, lucky := range by.luckyTable.Items {
				if v.GetLuckyBonus() != -1 && v.GetLuckyBonus() == int32(lucky.GetLuckyGroupID()) {
					sum += int32(lucky.GetLuckyBonusRate())
				}
			}

			if sum > 0 {
				rnd := util.RndInt(0, int(sum))

				layerrewardrate := int32(0)
				for _, lucky := range by.luckyTable.Items {
					if v.GetLuckyBonus() == int32(lucky.GetLuckyGroupID()) {
						layerrewardrate += int32(lucky.GetLuckyBonusRate())
						if int32(rnd) <= layerrewardrate {
							logmsg := fmt.Sprintf("設定新產生魚Lucky fishid[%v][%v] Lucky[%v] ", fi.fishid, fi.name, lucky.GetLuckyBonus())
							mlog.Debug(logmsg)
							mylog.MongoLog(0, *by.TableManager.Room.Seq	, by.TableNo(), logmsg)
							fi.luckybonus = int32(lucky.GetLuckyBonus())
							break
						}
					}
				}

			}
			if fi.luckybonus != 0 {
				break
			}
		}
	}
}

//创建一条特殊鱼
func (m *byfishmanager) newTeshuFish(by *Bylogic) *byfishdata {
	t := m.getCreateTeshuFishType()
	fi, ok := m.config.getfishinfobytype(t)

	if !ok {
		return nil
	}

	if by.fishmanager.getFishCountByType(fi.fishid) >= fi.maxcount {
		return nil
	}

	if fi == nil {
		f := m.createFishShoal(t, by)
		if f != nil {
			return f
		}

		fi = m.config.fishinfoarr[0]
	}
	//mlog.System("新增魚 fishid[%v]", fi.fishid)
	return m.createFish(fi, by)
}

func (m *byfishmanager) newGoldFish(by *Bylogic) *byfishdata {
	t := m.getCreateGoldFishType()
	fi, ok := m.config.getfishinfobytype(t)
	if !ok || fi == nil {
		fi = m.config.fishinfoarr[0]
	}

	if m.getFishCountByType(fi.fishid) >= fi.maxcount {
		return nil
	}

	mlog.Debug("创建一条Gold鱼 fishSeq[%v]", fi.fishid)
	return m.createFish(fi, by)
}

//创建鱼群
func (m *byfishmanager) createFishShoal(t int32, by *Bylogic) *byfishdata {
	//生成鱼群
	fs, _ := m.config.getfishshoalinfobytype(t)
	var f *byfishdata
	mlog.Debug("创建鱼群 t[%v] ", t)
	for i := 0; i < len(fs.fishinfocount); i++ {
		count := fs.fishinfocount[i].count
		mlog.Debug("创建鱼群 fish count[%v]", count)
		//这个信息是发送给客户端，表示这是一个什么类型的鱼群
		if f == nil {
			//需要将数据拷贝出来做改变，否则会导致鱼群的第一条鱼找不到相应的类型
			mlog.Debug("创建鱼群 fishSeq[%v] type[%v] bosstype[%v] fishinfocount[%v]",
				fs.fishinfocount[i].fi.fishid, fs.fishshoaltype, fs.fishinfocount[i].fi.bosstype, len(fs.fishinfocount))

			f = new(byfishdata)
			mlog.Debug("新增魚 fishSeq[%v]", fs.fishinfocount[i].fi)
			tmpf := m.createFish(fs.fishinfocount[i].fi, by)
			f.fishSeq = tmpf.fishSeq
			f.fishtype = fs.fishshoaltype

			f.fishcount = fs.fishinfocount[i].count
			f.fishnumber = fs.fishinfocount[i].fi.fishid
			count--
		}
		//逐个生成鱼
		for n := 0; n < int(count); n++ {
			mlog.Debug("创建鱼群 fishSeq[%v] type[%v] count[%v]", fs.fishinfocount[i].fi.fishid, fs.fishshoaltype, count)
			m.createFish(fs.fishinfocount[i].fi, by)
		}
	}

	return f
}

//创建一条鱼
func (m *byfishmanager) createFish(fi *fishinfo, by *Bylogic) *byfishdata {
	f := new(byfishdata)

	f.fishtype = fi.fishid
	f.fishSeq = m.fishSeq

	if len(fi.rewardinfo) > 1 {

		//TODO 優化 抽離rewardinfo不要放在fi裡面
		sum := int32(0)
		for _, r := range fi.rewardinfo {
			sum += r.rewardrate
		}
		rnd := util.RndInt(0, int(sum))
		rndRewardIndex := 0
		layerrewardrate := int32(0)
		for k, v := range fi.rewardinfo {
			layerrewardrate += v.rewardrate
			if int32(rnd) <= layerrewardrate {
				rndRewardIndex = k
				//mlog.System("fishid[%v] rewardrate[%v] rnd[%v] k[%v] sum[%v]", f.fishtype, v.rewardrate, rnd, rndRewardIndex, sum)
				break
			}
		}
		layerrewardrate = 0

		f.fishodd = fi.rewardinfo[rndRewardIndex].reward
		//mlog.Debug("创建特殊魚 fishid[%v] rnd[%v] fishodd[%v]", f.fishtype, rnd, f.fishodd)
		if f.fishodd == 0 {
			mlog.Error("fishodd is 0 ,  创建特殊魚 QueryMessage[%v] fishid[%v]", f.fishSeq, f.fishtype)
		}
	} else {
		f.fishodd = fi.rewardinfo[0].reward
		//mlog.Debug("创建一条鱼 QueryMessage[%v] fishid[%v] fishodd[%v]", f.fishSeq, f.fishtype, f.fishodd)
		if f.fishodd == 0 {
			mlog.Debug("fishodd is 0 ,  创建一条鱼 QueryMessage[%v] fishid[%v]", f.fishSeq, f.fishtype)
		}
	}

	m.setNewFishLuckyBonus(by, fi)

	//将新生成的鱼放入鱼数据管理中
	m.fishdata[f.fishSeq] = f
	mlog.Debug("創建一條魚 fishSeq[%v] fishOdd[%v] fishtype[%v]", f.fishSeq, f.fishodd, f.fishtype)

	m.fishSeq++

	if m.fishSeq >= 0xfffffff {
		m.fishSeq = 1
	}

	return f
}

//生成新的鱼潮
func (m *byfishmanager) newFishWave() {
	m.fishwave = true

	//m.delAllFish()
	fw := m.fishwavecf.newFishWave()

	m.enterfishwavetime = fw.entertime
	m.fishwavelasttime = fw.lasttime
	m.fishwavetype = fw.fishwavetype

	exsitFishType := make(map[int32]int32, 0)
	for _, fa := range fw.fisharr {
		exsitNum := int32(0)
		if result, ok := exsitFishType[fa.fishType]; ok {
			exsitNum += result
			//mlog.System("已經產生魚朝 fishType[%v] result[%v]",fa.fishType, result)
		} else {
			exsitFishType[fa.fishType] = 0
		}
		for i := int32(1); i <= fa.count; i++ {
			f := new(byfishdata)
			f.fishtype = fa.fishType
			//mlog.Debug("鱼的类型：%d", fa.fishType)
			seqId := 1*100000 + fa.fishType*1000 + (exsitNum + i)
			f.fishSeq = seqId
			m.fishdata[f.fishSeq] = f
			exsitFishType[fa.fishType] = i
			//mlog.Debug("產生魚朝 fishSeq[%v] fishid[%v] fishwavetype[%v] count[%v] i[%v]", f.fishSeq, f.fishtype, fw.fishwavetype, fa.count, i)
		}
	}
}

//鱼潮结束
func (m *byfishmanager) endFishWave() {
	m.fishwave = false
	mlog.Debug("鱼潮结束")
	//m.delAllFish()
}

//设置鱼的配置文件
func (m *byfishmanager) setFishConfig(rt *DataCfg.RoomTable) {
	defer util.PanicErrStack()
	m.praseFishWeight(rt.GetFishWeight())
	m.praseFishShoalWeight(rt.GetFishShoalWeight())
	m.parseGoldFishWeight(rt.GetGoldFishWeight())
	m.praseBossFishWeight(rt.GetBossWeight())
	m.bossentertime = int32(rt.GetEnterTime())
	m.bossfishtime = int32(rt.GetBossInterval())
	m.praseTeShuFishWeight(rt.GetTeshuWeight())
	m.teshufishtime = int32(rt.GetTeshuInterval())
	m.fishshoaltime = int32(rt.GetGenFishShoalInerval())
	m.goldfishtime = int32(rt.GetGoldFishInterval())
	m.config.createtime = int32(rt.GetGenFishInterval())
	m.config.createtimeoffset = int32(rt.GetGenFishIntervalRndRange())
}

//解析鱼生成权重
func (m *byfishmanager) praseFishWeight(str string) {
	m.fishweightarr = make([]*fishweight, 0)
	strarr := strings.Split(str, ",")
	for _, v := range strarr {
		varr := strings.Split(v, ":")
		fw := new(fishweight)
		fw.fishid = getint32(varr[0])
		fw.weight = getint32(varr[1])
		m.totalfishweight += int64(fw.weight)
		m.fishweightarr = append(m.fishweightarr, fw)
		//mlog.Debug("RoomTable.fishWeight解析鱼生成权重 鱼种类[%v] 鱼权重[%v]", fw.fishid, fw.weight)
	}
}

//解析魚群生成權重
func (m *byfishmanager) praseFishShoalWeight(str string) {
	m.fishshoalweight = make([]*fishweight, 0)
	strarr := strings.Split(str, ",")
	for _, v := range strarr {
		varr := strings.Split(v, ":")
		fw := new(fishweight)
		fw.fishid = getint32(varr[0])
		fw.weight = getint32(varr[1])
		m.totalfishshoalweight += int64(fw.weight)
		m.fishshoalweight = append(m.fishshoalweight, fw)
		//mlog.Debug("RoomTable.bossWeight解析魚群生成权重 鱼种类[%v] 鱼权重[%v]", fw.fishid, fw.weight)
	}
}

func (m *byfishmanager) parseGoldFishWeight(str string) {
	m.goldfishweight = make([]*fishweight, 0)
	strarr := strings.Split(str, ",")

	for _, v := range strarr {
		varr := strings.Split(v, ":")
		fw := new(fishweight)
		fw.fishid = getint32(varr[0])
		fw.weight = getint32(varr[1])
		m.totalgoldfishweight += int64(fw.weight)
		m.goldfishweight = append(m.goldfishweight, fw)
		//mlog.Debug("RoomTable.GoldFishWeight解析魚群生成权重 鱼种类[%v] 鱼权重[%v]", fw.fishid, fw.weight)
	}
}

//解析BOSS鱼生成权重
func (m *byfishmanager) praseBossFishWeight(str string) {
	defer util.PanicErrStack()
	m.bossfishweight = make([]*fishweight, 0)
	strarr := strings.Split(str, ",")
	for _, v := range strarr {
		varr := strings.Split(v, ":")
		fw := new(fishweight)
		fw.fishid = getint32(varr[0])
		fw.weight = getint32(varr[1])
		m.totalbossfishweight += int64(fw.weight)
		m.bossfishweight = append(m.bossfishweight, fw)
		//mlog.Debug("RoomTable.bossWeight解析BOSS鱼生成权重 鱼种类[%v] 鱼权重[%v]", fw.fishid, fw.weight)
	}
}

//解析特殊鱼生成权重
func (m *byfishmanager) praseTeShuFishWeight(str string) {
	m.teshufishweight = make([]*fishweight, 0)
	strarr := strings.Split(str, ",")
	for _, v := range strarr {
		varr := strings.Split(v, ":")
		fw := new(fishweight)
		fw.fishid = getint32(varr[0])
		fw.weight = getint32(varr[1])
		m.totalteshufishweight += int64(fw.weight)
		m.teshufishweight = append(m.teshufishweight, fw)
		//mlog.Debug("RoomTable.teshuWeight解析特殊鱼生成权重 鱼种类[%v] 鱼权重[%v]", fw.fishid, fw.weight)
	}
}

//解析鱼潮字符串
func (m *byfishmanager) parsefishwavestring(fstr string) {
	fstrarr := strings.Split(fstr, ",")
	for i := 0; i < len(fstrarr); i++ {
		fishstr := strings.Split(fstrarr[i], ":")
		t, _ := strconv.ParseInt(fishstr[0], 10, 32)
		w, _ := strconv.ParseInt(fishstr[1], 10, 32)

		fa := new(fishwaveinfo)
		fa.fisharr = make([]*fishwavefish, 0)
		fa.fishwavetype = int32(t)
		fa.weight = int32(w)

		m.fishwavecf.fishwaveinfoarr = append(m.fishwavecf.fishwaveinfoarr, fa)
	}
}

//判断是否能補到鱼
func (m *byfishmanager) canfish(frameid int64, fishid int32, rate float32, fished bool, ud *byuserdata, basescore int32, allfishid []int32) (bool, int32) {
	mlog.Debug("判断是否能補到鱼 fishId[%v] rate[%v]", m.fishdata[fishid].fishtype, rate)
	if f, ok := m.fishdata[fishid]; ok {
		v, ok, fishOdd := m.getFishLuckyOdds(f)

		if ok {
			f.fishhit++

			//必中的情况,炸彈蟹
			if fished {
				if ud.diedbosstype == BoomFish {
					n, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", rand.Float64()), 64)
					if rate > float32(n) {
						logMsg := fmt.Sprintf("userid[%v]必中的情况,炸彈蟹(or核彈頭) allfish 抓鱼fishid[%v][%v]的概率 rand[%v] 要小餘 (子彈威力[%v]",
							ud.user.UserID, fishid, v.name, n, rate)
						mlog.Debug(logMsg)
						mylog.MongoLog(ud.user.UserID, frameid, ud.user.TableNo, logMsg)
						mylog.MongoLog(ud.user.UserID, 0, ud.user.TableNo, fmt.Sprintf("userid[%v]必中的情况,炸彈蟹(or核彈頭) allfish 抓鱼fishid[%v][%v]的概率 rand[%v] 要小餘 (子彈威力[%v]",
							ud.user.UserID, fishid, v.name, n, rate))
						delete(m.fishdata, fishid)
						return true, fishOdd
					} else {
						mlog.Debug("allfish不能捕到魚")
						return false, 0
					}
				}

				if ud.diedbosstype == LeiLong || ud.diedbosstype == Zuantoulongxia {
					logMsg := fmt.Sprintf("userid[%v] 特殊魚[%v] allfish 抓鱼閃電鰻魚or钻头龙虾 fishid[%v]的概率 (子彈威力[%v]",
						ud.user.UserID, ud.diedbosstype, fishid, rate)
					mlog.Debug(logMsg)
					mylog.MongoLog(ud.user.UserID, frameid, ud.user.TableNo, logMsg)
					mylog.MongoLog(ud.user.UserID, 0, ud.user.TableNo, fmt.Sprintf("userid[%v] 特殊魚[%v] allfish 抓鱼閃電鰻魚or钻头龙虾 fishid[%v]的概率 (子彈威力[%v]",
						ud.user.UserID, ud.diedbosstype, fishid, rate))
					if fishOdd <= 0 {
						mlog.Error("allfish 抓鱼閃電鰻魚or钻头龙虾 fishid[%v]的概率 (子彈威力[%v]",
							fishid, rate)
					}
					delete(m.fishdata, fishid)
					return true, fishOdd
				}
			}

			//对鱼的攻击次数需要大于最小攻击次数
			if f.fishhit >= v.minhit {
				mlog.Debug("对鱼的攻击次数需要大于最小攻击次数[%v] fishSeq[%v] fishSeq[%v]", f.fishhit, f.fishSeq, f.fishtype)
				//捕获鱼的概率为0~0.999
				rand := m.getRand(rate)
				diedRate := m.getDiedRate(rate, fishOdd, allfishid, v, basescore, ud)

				logMsg := fmt.Sprintf("userid[%v] if [%v]抓鱼前fishid[%v][%v] lucky[%v]的概率 rand[%v] 要小餘 (子彈威力[%v]/魚倍率[%v]) = %v",
					ud.user.UserID, (diedRate > rand), v.fishid, v.name, v.luckybonus, rand, rate, fishOdd, diedRate)
				mlog.Debug(logMsg)
				mylog.MongoLog(ud.user.UserID, frameid, ud.user.TableNo, logMsg)

				if diedRate > rand {
					logMsg = fmt.Sprintf("userid[%v] 抓到魚fishid[%v][%v] lucky[%v]的概率 rand[%v] 要小餘 (子彈威力[%v]/魚倍率[%v]) = %v",
						ud.user.UserID, fishid, v.name, v.luckybonus, rand, rate, fishOdd, diedRate)
					mlog.Debug(logMsg)
					mylog.MongoLog(ud.user.UserID, frameid, ud.user.TableNo, logMsg)

					if f.fishtype == 42 {
						ud.lastshoottime = int64(time.Now().UnixNano()) + 30*int64(time.Second)
					}

					isTeshufishDied := m.hasTeshfishDied(f, v, ud, fishOdd, basescore)
					if isTeshufishDied {
						delete(m.fishdata, fishid)
						return true, 0
					}

					delete(m.fishdata, fishid)
					return true, fishOdd
				} else if fishOdd >= (50*basescore) {
					logmsg := fmt.Sprintf("userid[%v]未打死 fishid[%v][%v] 倍率[%v]",
						ud.user.UserID, fishid, v.name, fishOdd)
					mylog.MongoLog(ud.user.UserID, frameid, ud.user.TableNo, logmsg)
				}
			}
		}
	} else {
		mlog.Debug("不能捕到魚")
	}

	return false, 0
}

//取得魚lucky後的倍率
func (m *byfishmanager) getFishLuckyOdds(f *byfishdata) (*fishinfo, bool, int32) {
	v, ok := m.config.getfishinfobytype(f.fishtype)
	fishOdd := v.rewardinfo[0].reward
	if f.fishodd != 0 {
		fishOdd = f.fishodd
	}
	if v.luckybonus > 1 {
		fishOdd = fishOdd * v.luckybonus
	}
	return v, ok, fishOdd
}

//取得魚的死亡率,
//核彈魚傳入allfishid, teshuFish
func (m *byfishmanager) getDiedRate(rate float32, fishOdd int32, allfishid []int32, teshuFish *fishinfo, basescore int32, ud *byuserdata) float32 {
	sumFishOdd := float32(0)

	if teshuFish.bosstypeinfo.fishtype == Dujiaojing {
		newRate := rate / float32(1800*basescore)
		n, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", newRate), 64)
		rand := float32(n)
		mlog.Debug("昆魚died rand[%v] = rate[%v] sumOdd[%v]", rand, rate, sumFishOdd)
		logmsg := fmt.Sprintf("昆魚 死亡率[%v] = rate[%v] / %v", rand, rate, 1800*basescore )
		mylog.MongoLog(ud.user.UserID, 0, ud.user.TableNo, logmsg)
		return rand
	}

	if allfishid != nil && len(allfishid) != 0 {
		for _, fishid := range allfishid {
			if f, ok := m.fishdata[fishid]; ok {
				v, _ := m.config.getfishinfobytype(f.fishtype)
				fishinfo, _, fishOdd := m.getFishLuckyOdds(f)
				spmaxfishmulti := teshuFish.bosstypeinfo.spmaxfishmulti / 1000 * basescore
				if spmaxfishmulti != 0 && spmaxfishmulti != -1 && fishOdd > spmaxfishmulti {
					logmsg := fmt.Sprintf("核彈頭 fishid[%v] 特殊於不能打中%v倍以上的魚[%v]倍率[%v]lucky[%v] ", teshuFish.fishid, spmaxfishmulti, fishinfo.name, fishOdd, fishinfo.luckybonus)
					mlog.Debug(logmsg)
					mylog.MongoLog(ud.user.UserID, 0, ud.user.TableNo, logmsg)
					continue
				}

				//不会炸死特殊鱼
				if v.bosstypeinfo.fishtype >= 3 {
					mlog.Debug("核彈頭不会炸死特殊鱼 fishtype[%v]", f.fishtype)
					continue
				}

				if f.fishodd != 0 {
					fishinfo, _, fishOdd := m.getFishLuckyOdds(f)
					sumFishOdd += float32(fishOdd)
					logmsg := fmt.Sprintf("核彈頭 打中[%v][%v] 倍率[%v] lucky[%v]", fishinfo.fishid, fishinfo.name, fishOdd, fishinfo.luckybonus)
					mylog.MongoLog(ud.user.UserID, 0, ud.user.TableNo, logmsg)
				}
			}
		}

		newRand := float32(0)
		if sumFishOdd > 0 {
			newRand = rate / sumFishOdd
		}
		n, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", newRand), 64)
		rand := float32(n)

		logmsg := fmt.Sprintf("核彈頭 死亡率[%v] = 攻擊率[%v] / 所有倍數總和[%v]", rand, rate, sumFishOdd)
		mylog.MongoLog(ud.user.UserID, 0, ud.user.TableNo, logmsg)

		return rand
	} else {
		diedRate := util.Decimal(float64(rate/float32(fishOdd)), 8)
		return diedRate
	}
}

//取得亂數
func (m *byfishmanager) getRand(rate float32) float32 {
	n, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", rand.Float64()), 64)
	rand := float32(n)
	return rand
}

func (m *byfishmanager) hasTeshfishDied(f *byfishdata, v *fishinfo, ud *byuserdata, fishOdd int32, basescore int32) bool {
	result := false

	fishinfotype, ok := m.config.getfishinfobytype(f.fishtype)

	if ok && fishinfotype.bosstypeinfo.fishtype >= 3 {
		ud.bossfishid = f.fishtype
		ud.diedbosstype = fishinfotype.bosstypeinfo.fishtype
		ud.teshufishfishodd = fishOdd
		ud.spmaxfishmulti = v.bosstypeinfo.spmaxfishmulti / 1000 * basescore
		ud.candiedrate = float32(v.bosstypeinfo.spdeathrate)

		logMsg := fmt.Sprintf("userid[%v] 抓到特殊魚[%v] 倍率[%v] 倍數門檻[%v] ",
			ud.user.UserID, f.fishtype, fishOdd, ud.spmaxfishmulti)
		mlog.Debug(logMsg)
		mylog.MongoLog(ud.user.UserID, 0, 0, logMsg)
		result = true
	}

	return result
}

func (m *byfishmanager) findFish(fishid int32) bool {
	_, ok := m.fishdata[fishid]
	return ok
}

//设置鱼的数据
func (m *byfishmanager) setfishtabledata(data *DataCfg.FishTable_ARRAY, flag int32, fishTypeArr *DataCfg.FishTypeTable_ARRAY) {
	for _, d := range data.GetItems() {
		m.config.setfishinfobytype(d, flag, fishTypeArr)
	}
}

func (m *byfishmanager) setfishtypetabledata(data *DataCfg.RewardTable_ARRAY, basescore int32) {
	for _, v := range m.config.fishinfoarr {
		v.rewardinfo = make([]*fishoddinfo, 0)
		for _, d := range data.Items {
			if v.oddType == int32(d.GetRewardGroupId()) {
				rewardInt, _ := strconv.ParseFloat(d.GetReward(), 64)
				v.rewardinfo = append(v.rewardinfo, &fishoddinfo{
					rewardgroupid: int32(d.GetRewardGroupId()),
					reward:        int32(rewardInt / 1000 * float64(basescore)),
					rewardrate:    int32(d.GetRewardRate())})
			}
		}
	}

	//for _, v := range m.config.fishinfoarr {
	//	for _, r := range v.rewardinfo {
	//		mlog.Debug("oddType[%v] rewardgroupid[%v] reward[%v] rewardrate[%v]", v.oddType, r.rewardgroupid, r.reward, r.rewardrate)
	//	}
	//}
}

//设置鱼群的数据
func (m *byfishmanager) setfishshoaltabledata(data *DataCfg.FishShoalTable_ARRAY) {
	m.config.fishshoalarr = make([]*fishshoal, 0)
	for _, d := range data.GetItems() {
		m.config.setfishshoalbytype(d)
	}
}

//检测鱼潮是否结束
func (m *byfishmanager) checkFishWaveEnd() bool {
	if !m.fishwave {
		return true
	}

	if len(m.fishdata) <= 0 {
		return true
	}
	return false
}

//设置鱼潮的数据
func (m *byfishmanager) setFishWaveTalbeData(data *DataCfg.FishWaveTable_ARRAY, pathTable *DataCfg.PathTable_ARRAY) {

	for _, d := range data.GetItems() {
		m.fishwavecf.setfishwavebytype(d, pathTable)
	}
}

//获取此类型的鱼的个数，在鱼多了以后可能会卡顿。需要优化
func (m *byfishmanager) getFishCountByType(t int32) int32 {
	count := int32(0)
	for _, f := range m.fishdata {
		if f.fishtype == t {
			count++
		}
	}

	return count
}

//删掉所有的鱼
func (m *byfishmanager) delAllFish() {
	m.fishdata = make(map[int32]*byfishdata, 0)
}

//通过鱼ID获取鱼的信息
func (m *byfishmanager) getFishInfoByFishID(fishid int32) *byfishdata {
	for k, v := range m.fishdata {
		mlog.Debug("通过鱼ID获取鱼的信息 fishid[%v] key[%v]  fishSeq[%v] fishtype[%v] fishnumber[%v] ", fishid, k, v.fishSeq, v.fishtype, v.fishnumber)
	}
	f, ok := m.fishdata[fishid]
	if ok {
		return f
	}

	return nil
}

//获取生成鱼的类型
func (m *byfishmanager) getCreateFishType() int32 {
	w := rand.Int63n(m.totalfishweight)

	weight := int64(0)

	for _, f := range m.fishweightarr {
		weight += int64(f.weight)
		if w < weight {
			return f.fishid
		}
	}

	return m.fishweightarr[0].fishid
}

//獲取生成魚群的類型
func (m *byfishmanager) getCreateFishShoalType() int32 {
	w := rand.Int63n(m.totalfishshoalweight)
	weight := int64(0)

	for _, f := range m.fishshoalweight {
		weight += int64(f.weight)
		if w < weight {
			return f.fishid
		}
	}
	return m.fishshoalweight[0].fishid
}

func (m *byfishmanager) getCreateGoldFishType() int32 {
	w := rand.Int63n(m.totalgoldfishweight)
	weight := int64(0)

	for _, f := range m.goldfishweight {
		weight += int64(f.weight)
		if w < weight {
			mlog.Debug("獲取生成Gold魚的類型 fishSeq[%v]", f.fishid)
			return f.fishid
		}
	}
	return m.goldfishweight[0].fishid
}

//获取生成BOSS鱼的类型
func (m *byfishmanager) getCreateBossFishType() int32 {
	w := rand.Int63n(m.totalbossfishweight)

	weight := int64(0)
	for _, f := range m.bossfishweight {
		weight += int64(f.weight)
		if w < weight {
			return f.fishid
		}
	}

	return m.bossfishweight[0].fishid
}

//获取生成特殊鱼的类型
func (m *byfishmanager) getCreateTeshuFishType() int32 {
	w := rand.Int63n(m.totalteshufishweight)

	weight := int64(0)
	for _, f := range m.teshufishweight {
		weight += int64(f.weight)
		if w < weight {
			return f.fishid
		}
	}

	return m.teshufishweight[0].fishid
}
