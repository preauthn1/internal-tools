package by

import (
	"example.invalid/ServerTeam/Servers/internal/game/by/DataCfg"
	"math/rand"
	"strconv"
	"strings"
)

type fishinfo struct {
	fishid           int32   //鱼的种类
	name             string  //魚的名稱
	bosstype         int32   //鱼的BOSS类型，BOSS鱼，炸弹鱼等
	minhit           int32   //最小攻击次数
	money            int64   //掉钱
	oddType          int32   //倍率
	fishway          []int32 //鱼的线路
	maxcount         int32   //最多有多少条这类型的鱼
	returntime       int32   //鱼未死亡回游时间
	returntimeoffset int32   //回游时间偏移量
	maxreturn        int32   //最大回游次数
	rewardinfo       []*fishoddinfo
	bosstypeinfo     fishtypeinfo
	luckybonus       int32 //幸運倍數
}

type fishtypeinfo struct {
	id              int32
	fishtype        int32   //1：一般魚	2：Boss	3：特魚
	iseffect        int32   //是否被特殊效果影響	0否	1是
	spmaxfishmulti  int32   //倍數門檻	N倍以下(含)的普魚是可以被打死的	-1表示不使用
	spdeathrate     float32 //死亡概率	直接無視原始死亡概率	-1表示不使用
	rateincrease    float32 //沒打死增加概率 	0表示不使用
	increasemax     float32 //增加概率上限	0表示不使用
	spactiontime    int32   //演出時間	或鑽頭飛行時間	-1表示不使用
	spmaxfishcount  int32   //攻擊隻數上限	效果不能影響超過設定的隻數	e.g.鑽頭撞到魚會判定的隻數,電鰻魚可以判定的隻數上限…等	-1表示不管隻數
	spboomarea      float32 //爆炸半徑	-1表示不使用
	luckybonus      int32   //幸運倍數	幸運魚專屬欄位	輸入LuckyTable的群組編號	來設定幸運倍數	實際魚倍率=reward*LuckyBonus	-1表示不使用幸運倍數
	bulletnumber    string  //每吃多少子彈膨脹一次	維持空白表示不使用	有多組設定用逗號分隔
	swellrate       int32   //每次膨脹多少%	維持空白表示不使用
	swellmaxrate    int32   //每次膨脹多少%	維持空白表示不使用
	bodysep         string  //得分分配次數	維持空白表示不使用	有多組設定用逗號分隔
	winsep          string  //得分分配比例	維持空白表示不使用	有多組設定用逗號分隔	數字總和為分母	根據得分分配次數選擇分配比例的數字
	mrate           string  //倍增魚補正倍率參數	FishType是7才會使用到	維持空白表示不使用{範圍最小值,範圍最大值}
	issamekill      int32   //是否殺死同類	0否	1是	殺死同類魚須帶有特殊特效標示
	sammekillfishid int32   //根據IsSameKill判斷是否殺死指定的魚種	輸入要被殺死的FishId
}

type fishoddinfo struct {
	rewardgroupid int32
	reward        int32
	rewardrate    int32
}

type fishcount struct {
	fi    *fishinfo
	count int32
}

//鱼群配置
type fishshoal struct {
	fishshoaltype int32        //鱼群类型
	fishinfocount []*fishcount //各类鱼的个数
}

type fishconfig struct {
	fishinfoarr      []*fishinfo  //鱼信息
	fishshoalarr     []*fishshoal //鱼群信息
	createtime       int32        //生成时间
	createtimeoffset int32        //生成时间偏移值
}

//新建一个鱼配置实例
func newfishconfi() *fishconfig {
	fc := new(fishconfig)
	fc.fishinfoarr = make([]*fishinfo, 0)
	return fc
}

func getint32(s string) int32 {
	v, _ := strconv.ParseInt(s, 10, 32)
	return int32(v)
}

//获取创建时间
func (fc *fishconfig) getcreatetime() int32 {
	offset := rand.Int31n(fc.createtimeoffset) + 1
	b := rand.Int31n(2)
	if b == 0 {
		offset = -offset
	}

	return fc.createtime + offset
}

//通过鱼类型获取鱼的信息
func (fc *fishconfig) getfishinfobytype(t int32) (*fishinfo, bool) {
	for _, f := range fc.fishinfoarr {
		if f.fishid == t {
			return f, true
		}
	}

	return nil, false
}

//通过鱼群类型获取鱼群的信息
func (fc *fishconfig) getfishshoalinfobytype(t int32) (*fishshoal, bool) {
	for _, f := range fc.fishshoalarr {
		if f.fishshoaltype == t {
			return f, true
		}
	}

	return nil, false
}

//通过鱼的类型设置鱼的类型
func (fc *fishconfig) setfishinfobytype(ft *DataCfg.FishTable, flag int32, fishTypeArr *DataCfg.FishTypeTable_ARRAY) {
	fi, ok := fc.getfishinfobytype(int32(ft.GetFishId()))

	if !ok || fi == nil {
		fi = new(fishinfo)
		fc.fishinfoarr = append(fc.fishinfoarr, fi)
	}

	fi.fishid = int32(ft.GetFishId())
	fi.name = string(ft.GetName()[:])

	fi.oddType = int32(ft.GetReward())
	fi.minhit = int32(ft.GetMinHit())

	for _, v := range fishTypeArr.Items {
		if ft.GetFishType() == v.GetID() {
			fi.bosstype = int32(v.GetFishType())
			fi.bosstypeinfo.id = int32(v.GetID())
			fi.bosstypeinfo.fishtype = int32(v.GetFishType())
			fi.bosstypeinfo.iseffect = int32(v.GetIsEffect())
			fi.bosstypeinfo.spmaxfishmulti = int32(v.GetSPMaxFishMulti())
			fi.bosstypeinfo.spdeathrate = float32(v.GetSPDeathRate()) / 100
			fi.bosstypeinfo.rateincrease = float32(v.GetRateIncrease()) / 100
			fi.bosstypeinfo.increasemax = float32(v.GetIncreaseMax()) / 100
			fi.bosstypeinfo.spactiontime = int32(v.GetSPActionTime())
			fi.bosstypeinfo.spmaxfishcount = int32(v.GetSPMaxFishCount())
			fi.bosstypeinfo.spboomarea = v.GetSPBoomArea()
			fi.bosstypeinfo.luckybonus = int32(v.GetLuckyBonus())
			fi.bosstypeinfo.issamekill = int32(v.GetIsSameKill())
			fi.bosstypeinfo.sammekillfishid = int32(v.GetSameKillFishId())
			fi.bosstypeinfo.bulletnumber = v.GetBulletNumber()
			fi.bosstypeinfo.swellrate = v.GetSwellRate()
			fi.bosstypeinfo.swellmaxrate = v.GetSwellMaxRate()
			fi.bosstypeinfo.bodysep = v.GetBodySep()
			fi.bosstypeinfo.winsep = v.GetWinSep()
			fi.bosstypeinfo.mrate = v.GetMRate()
		}
	}

	fi.maxcount = int32(ft.GetNumberLimit())

	//mlog.Debug("魚死亡率 fishid[%v] bosstype[%v] bosstypeinfo[%+v]",
	//	fi.fishid, fi.bosstype, fi.bosstypeinfo)
}

//鱼群配置
func (fc *fishconfig) setfishshoalbytype(ft *DataCfg.FishShoalTable) {
	fs := new(fishshoal)
	fc.fishshoalarr = append(fc.fishshoalarr, fs)

	fs.fishshoaltype = int32(ft.GetShoalId())

	strarr := strings.Split(ft.GetFishInfo(), ",")
	for _, arr := range strarr {
		infoarr := strings.Split(arr, ":")
		fishid := getint32(infoarr[0])
		fi, find := fc.getfishinfobytype(fishid)
		if find {
			ficount := new(fishcount)
			ficount.fi = fi
			ficount.count = getint32(infoarr[1])
			fs.fishinfocount = append(fs.fishinfocount, ficount)
			//mlog.Debug("FishShoalTable鱼群配置 鱼群序号[%v] fishSeq[%v] count[%v] 最小攻击次数[%v] 掉钱[%v] fish倍率[%v] 最多有多少条这类型的鱼[%v] ",
			//	fs.fishshoaltype, fishid, ficount.count, fi.minhit, fi.money, fi.rewardinfo[0].reward, fi.maxcount)
		}
	}
}

type fishwavefish struct {
	fishType int32  //鱼种类
	count    int32  //鱼潮有多少条鱼
	pathName string //路線
	interval int32  //間隔時間
}

type fishwaveinfo struct {
	fisharr      []*fishwavefish //鱼潮的数据
	fishwavetype int32           //鱼潮类型
	weight       int32           //鱼潮权重
	entertime    int32           //进入时间
	lasttime     int32           //持续时间

}

type fishwaveconfig struct {
	fishwaveinfoarr []*fishwaveinfo //鱼潮的数组
}

//新建鱼潮配置文件
func newFishWaveConfig() *fishwaveconfig {
	f := new(fishwaveconfig)
	f.fishwaveinfoarr = make([]*fishwaveinfo, 0)
	return f
}

//通过类型获取鱼潮
func (fw *fishwaveconfig) getfishwavebytype(t int32) (*fishwaveinfo, bool) {
	for _, f := range fw.fishwaveinfoarr {
		if f.fishwavetype == t {
			return f, true
		}
	}

	return nil, false
}

//通过类型设置鱼潮内容
func (fw *fishwaveconfig) setfishwavebytype(ft *DataCfg.FishWaveTable, pathTable *DataCfg.PathTable_ARRAY) {
	fi, _ := fw.getfishwavebytype(int32(ft.GetWaveID()))

	if fi == nil {
		fi = new(fishwaveinfo)
		fi.fisharr = make([]*fishwavefish, 0)
		fw.fishwaveinfoarr = append(fw.fishwaveinfoarr, fi)
	}

	fi.fisharr = fw.parsefishstr(ft.GetPathTableId(), pathTable)
	fi.lasttime = int32(ft.GetDuration())
	//fi.pathinterval = ft.GetPathInterval()

	//for _, v := range fi.fisharr {
	//	mlog.Debug("FishWaveTable 鱼潮内容 鱼潮类型[%v] 鱼种类[%v] 鱼潮有多少条鱼[%v] 进入时间[%v] 持续时间[%v]",
	//		fi.fishwavetype, v.fishType, v.count, fi.entertime, fi.lasttime)
	//}

}

//解析鱼的字符串
func (fw *fishwaveconfig) parsefishstr(str string, pathTable *DataCfg.PathTable_ARRAY) []*fishwavefish {
	arr := make([]*fishwavefish, 0)
	strarr := strings.Split(str, ",")
	for _, v := range strarr {
		pathNo, _ := strconv.Atoi(v)
		for _, path := range pathTable.Items {
			if int32(pathNo) == path.GetPathNo() {
				fishinfo := new(fishwavefish)
				fishinfo.fishType = path.GetFishType()
				fishinfo.count = path.GetFishCount()
				fishinfo.pathName = path.GetPathName()
				fishinfo.interval = path.GetInterval()
				arr = append(arr, fishinfo)
				//mlog.Debug("魚潮 fishinfo 種類[%v] count[%v]", fishinfo.fishType, fishinfo.count)
			}
		}
	}

	return arr
}

//生成新的鱼潮
func (fw *fishwaveconfig) newFishWave() *fishwaveinfo {
	weight := int32(0)
	for _, f := range fw.fishwaveinfoarr {
		weight += f.weight
	}

	if weight <= 0 {
		return fw.fishwaveinfoarr[0]
	}

	w := rand.Int31n(weight)
	weight = 0
	for _, f := range fw.fishwaveinfoarr {
		weight += f.weight
		if w < weight {
			return f
		}
	}

	return fw.fishwaveinfoarr[0]
}
