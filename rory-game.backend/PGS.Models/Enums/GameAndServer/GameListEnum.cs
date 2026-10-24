using target.Models.Attributes;
using target.Models.CultureText;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;

namespace target.Models.Enums
{
    /// <summary>
    /// 遊戲名稱有用table管理，但是管理上比較像列舉
    /// </summary>
    public enum GameListEnum
    {
        /// <summary>
        /// 多語系請參考類似這樣的鍵值，總之就是 
        /// 多語系務必follow這個格式， Enums(固定)_GameListEnum(固定)_{程式碼列舉這邊的命名}
        /// 只標記這一個範例，剩下就不一一列舉了
        /// </summary>
        
        [ReferenceText(i18nKeys.Enums_GameListEnum_WLZB_五龙争霸)]
        
        WLZB_五龙争霸 = 1001,        
        
        XGPD_星光派对 = 1002,
        
        XMBX_熊猫报喜 = 1003,
        
        JJBX_金吉报喜_无尽宝藏 = 1004,
        
        DHZW_东海之王 = 1005,
        
        WLZBJP_五龍爭霸_彩金版 = 1006,
        
        CSDJP_財神到_彩金版 = 1007,
        
        WLHC_五龍皇朝 = 1008,
        
        WLJR_百萬金輪五龍爭霸 = 1009,
        
        JGYF_擊鼓迎福 = 1010,
        
        JXRY_吉祥如意 = 1011,
        
        CFLL_金吉報喜_財富連連 = 1012,
        
        S88F_88Fortune = 1013,

        JKJL_JokersJewels = 1014,
        
        LKB_LuckyBats = 1015,
        
        FAF_FrogsAndFlies = 1016,
        
        QOW_QueenOfWands = 1017,
        
        BB_BuffaloBonus = 1018,
        
        S888_888 = 1019,

        S5LL_5LuckyLions = 1020,      
        
        HHF_HotHotFruit = 1021,
        
        DGNO_DragonOrb = 1022,       
        
        DFDC_多福多财 = 1023,
        
        TGW_糖果屋 = 1024,
        
        ALD_阿拉丁 = 1025,
        
        ZCJB_招财进宝 = 1026,
        
        WLYC_五龍運財 = 1027,

        HYXM_鴻運熊貓 = 1028,
        
        JLBPK_加勒比扑克 = 2001,

        DCZZ_赌场战争 = 2003,
        
        HDC_海盗船 = 2004,

        DCHJK_賭場黑傑克 = 2005,

        HBSL_红包扫雷 = 2006,

        FGSB_富贵三宝 = 2007,
        
        CJL_免傭百家樂 = 2008,

        hljbp_欢乐聚宝盆 = 2009,
        
        DDZ_鬥地主 = 2010,
        
        QZNN_搶庄牛牛 = 2020,
        
        ZJH_炸金花 = 2030,

        dzpk_德州撲克 = 2040,
        
        wzp_梭哈 = 2050,

        QZ_搶庄21 = 2060,
        
        tbnn_通比牛牛 = 2070,
        
        sango_三公 = 2080,

        SGJ_股市大亨 = 2090,
        
        sss_十三水 = 2100,

        LHD_龙虎斗 = 2110,
        
        ttz_推筒子 = 2120,

        BJL_百家乐 = 2130,
        
        scmj_四川麻将 = 2150,
        
        srnn_雙人牛牛 = 2160,

        ermj_雙人麻將 = 2170,

        FKEB_二八杠 = 2180,

        BRNN_百人牛牛 = 2190,

        HHDZ_红黑大战 = 2200,
        
        HHNN_红黑牛牛 = 2210,
        
        INRM_Rummy = 2220,
        
        SGBJL_性感百家乐 = 2230,       
        
        JD21D_经典21点 = 2400,

        OKMA_海王3 = 3001,

        BY_捕鱼 = 3014,

        INFC_FruitChop = 4001,

        INLD_Ludo = 4002,

        INCRM_Carrom = 4003,

    }
}
