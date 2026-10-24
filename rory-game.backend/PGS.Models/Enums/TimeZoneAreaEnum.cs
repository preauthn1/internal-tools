using target.Models.Attributes;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{

    public enum TimeZoneAreaEnum
    {

        [Description("(UTC-12:00) 國際換日線以西")]
        [ExtraAttributeValue("-12.0")]
        Dateline_minus_1200 = 99120001,

        [Description("(UTC-11:00) 國際標準時間-11")]
        [ExtraAttributeValue("-11.0")]
        UTC11_minus_1100 = 99110001,

        [Description("(UTC-10:00) 夏威夷")]
        [ExtraAttributeValue("-10.0")]
        Hawaiian_minus_1000 = 99100001,

        [Description("(UTC-10:00) 阿留申群島")]
        [ExtraAttributeValue("-10.0")]
        Aleutian_minus_1000 = 99100002,

        [Description("(UTC-09:30) 馬克沙斯群島")]
        [ExtraAttributeValue("-9.5")]
        Marquesas_minus_0930 = 99093001,

        [Description("(UTC-09:00) 國際標準時間-09")]
        [ExtraAttributeValue("-9.0")]
        UTC9_minus_0900 = 99090001,

        [Description("(UTC-09:00) 阿拉斯加")]
        [ExtraAttributeValue("-9.0")]
        Alaskan_minus_0900 = 99090002,

        [Description("(UTC-08:00) 下加利福尼亞")]
        [ExtraAttributeValue("-8.0")]
        PacificMexico_minus_0800 = 99080001,

        [Description("(UTC-08:00) 國際標準時間-08")]
        [ExtraAttributeValue("-8.0")]
        UTC8_minus_0800 = 99080002,

        [Description("(UTC-08:00) 太平洋時間 (美國和加拿大)")]
        [ExtraAttributeValue("-8.0")]
        Pacific_minus_0800 = 99080003,

        [Description("(UTC-07:00) 亞歷桑那")]
        [ExtraAttributeValue("-7.0")]
        USMountain_minus_0700 = 99070001,

        [Description("(UTC-07:00) 山區時間 (美國和加拿大)")]
        [ExtraAttributeValue("-7.0")]
        Mountain_minus_0700 = 99070002,

        [Description("(UTC-07:00) 齊驊華，拉帕茲，馬札特蘭")]
        [ExtraAttributeValue("-7.0")]
        MountainMexico_minus_0700 = 99070003,

        [Description("(UTC-06:00) 中美洲")]
        [ExtraAttributeValue("-6.0")]
        CentralAmerica_minus_0600 = 99060001,

        [Description("(UTC-06:00) 中部時間 (美國和加拿大)")]
        [ExtraAttributeValue("-6.0")]
        Central_minus_0600 = 99060002,

        [Description("(UTC-06:00) 伊斯特島")]
        [ExtraAttributeValue("-6.0")]
        EasterIsland_minus_0600 = 99060003,

        [Description("(UTC-06:00) 瓜達拉加若，墨西哥城，蒙特利")]
        [ExtraAttributeValue("-6.0")]
        CentralMexico_minus_0600 = 99060004,

        [Description("(UTC-06:00) 薩克其萬 (加拿大)")]
        [ExtraAttributeValue("-6.0")]
        CanadaCentral_minus_0600 = 99060005,

        [Description("(UTC-05:00) 切圖馬爾")]
        [ExtraAttributeValue("-5.0")]
        EasternMexico_minus_0500 = 99050001,

        [Description("(UTC-05:00) 印第安納 (東部)")]
        [ExtraAttributeValue("-5.0")]
        USEastern_minus_0500 = 99050002,

        [Description("(UTC-05:00) 哈瓦那")]
        [ExtraAttributeValue("-5.0")]
        Cuba_minus_0500 = 99050003,

        [Description("(UTC-05:00) 土克斯及開科斯")]
        [ExtraAttributeValue("-5.0")]
        TurksAndCaicos_minus_0500 = 99050004,

        [Description("(UTC-05:00) 東部時間 (美國和加拿大)")]
        [ExtraAttributeValue("-5.0")]
        Eastern_minus_0500 = 99050005,

        [Description("(UTC-05:00) 波哥大，利馬，基多，里約布蘭科")]
        [ExtraAttributeValue("-5.0")]
        SAPacific_minus_0500 = 99050006,

        [Description("(UTC-05:00) 海地")]
        [ExtraAttributeValue("-5.0")]
        Haiti_minus_0500 = 99050007,

        [Description("(UTC-04:00) 亞松森")]
        [ExtraAttributeValue("-4.0")]
        Paraguay_minus_0400 = 99040001,

        [Description("(UTC-04:00) 佐治敦，拉帕茲，瑪瑙斯，聖胡安")]
        [ExtraAttributeValue("-4.0")]
        SAWestern_minus_0400 = 99040002,

        [Description("(UTC-04:00) 卡拉卡斯")]
        [ExtraAttributeValue("-4.0")]
        Venezuela_minus_0400 = 99040003,

        [Description("(UTC-04:00) 古雅巴")]
        [ExtraAttributeValue("-4.0")]
        CentralBrazilian_minus_0400 = 99040004,

        [Description("(UTC-04:00) 大西洋時間 (加拿大)")]
        [ExtraAttributeValue("-4.0")]
        Atlantic_minus_0400 = 99040005,

        [Description("(UTC-04:00) 聖地牙哥")]
        [ExtraAttributeValue("-4.0")]
        PacificSA_minus_0400 = 99040006,

        [Description("(UTC-03:30) 紐芬蘭")]
        [ExtraAttributeValue("-3.5")]
        Newfoundland_minus_0330 = 99033001,

        [Description("(UTC-03:00) 巴西利亞")]
        [ExtraAttributeValue("-3.0")]
        ESouthAmerica_minus_0300 = 99030001,

        [Description("(UTC-03:00) 布宜諾斯艾利斯")]
        [ExtraAttributeValue("-3.0")]
        Argentina_minus_0300 = 99030002,

        [Description("(UTC-03:00) 格陵蘭")]
        [ExtraAttributeValue("-3.0")]
        Greenland_minus_0300 = 99030003,

        [Description("(UTC-03:00) 聖匹島")]
        [ExtraAttributeValue("-3.0")]
        SaintPierre_minus_0300 = 99030004,

        [Description("(UTC-03:00) 蒙特維多")]
        [ExtraAttributeValue("-3.0")]
        Montevideo_minus_0300 = 99030005,

        [Description("(UTC-03:00) 蓬塔阿雷納斯")]
        [ExtraAttributeValue("-3.0")]
        Magallanes_minus_0300 = 99030006,

        [Description("(UTC-03:00) 薩爾瓦多")]
        [ExtraAttributeValue("-3.0")]
        Bahia_minus_0300 = 99030007,

        [Description("(UTC-03:00) 開雲，福塔力沙")]
        [ExtraAttributeValue("-3.0")]
        SAEastern_minus_0300 = 99030008,

        [Description("(UTC-03:00) 阿拉瓜伊納")]
        [ExtraAttributeValue("-3.0")]
        Tocantins_minus_0300 = 99030009,

        [Description("(UTC-02:00) 國際標準時間-02")]
        [ExtraAttributeValue("-2.0")]
        UTC2_minus_0200 = 99020001,

        [Description("(UTC-01:00) 亞速爾群島")]
        [ExtraAttributeValue("-1.0")]
        Azores_minus_0100 = 99010001,

        [Description("(UTC-01:00) 維德角")]
        [ExtraAttributeValue("-1.0")]
        CapeVerde_minus_0100 = 99010002,

        [Description("(UTC) Coordinated Universal Time")]
        [ExtraAttributeValue("0.0")]
        UTC_0000 = 77000001,

        [Description("(UTC+00:00) 聖多美普林西比")]
        [ExtraAttributeValue("0.0")]
        SaoTome_0000 = 77000002,

        [Description("(UTC+00:00) 蒙羅維亞，雷克雅維克")]
        [ExtraAttributeValue("0.0")]
        Greenwich_0000 = 77000003,

        [Description("(UTC+00:00) 都柏林，愛丁堡，里斯本，倫敦")]
        [ExtraAttributeValue("0.0")]
        GMT_0000 = 77000004,

        [Description("(UTC+01:00) 卡薩布蘭卡")]
        [ExtraAttributeValue("0.0")]
        Morocco_0000 = 77000005,

        [Description("(UTC+01:00) 中西非")]
        [ExtraAttributeValue("1.0")]
        WCentralAfrica_0100 = 77010001,

        [Description("(UTC+01:00) 塞拉耶佛，斯高彼亞，華沙，札格雷布")]
        [ExtraAttributeValue("1.0")]
        CentralEuropean_0100 = 77010002,

        [Description("(UTC+01:00) 布魯塞爾，哥本哈根，馬德里，巴黎")]
        [ExtraAttributeValue("1.0")]
        Romance_0100 = 77010003,

        [Description("(UTC+01:00) 貝爾格勒，布拉提斯拉瓦，布達佩斯，盧布亞納，布拉格")]
        [ExtraAttributeValue("1.0")]
        CentralEurope_0100 = 77010004,

        [Description("(UTC+01:00) 阿姆斯特丹，柏林，伯恩，羅馬，斯德哥爾摩，維也納")]
        [ExtraAttributeValue("1.0")]
        WEurope_0100 = 77010005,

        [Description("(UTC+02:00) 加薩，希伯侖")]
        [ExtraAttributeValue("2.0")]
        WestBank_0200 = 77020001,

        [Description("(UTC+02:00) 卡里寧格勒")]
        [ExtraAttributeValue("2.0")]
        Kaliningrad_0200 = 77020002,

        [Description("(UTC+02:00) 哈拉雷，皮托里")]
        [ExtraAttributeValue("2.0")]
        SouthAfrica_0200 = 77020003,

        [Description("(UTC+02:00) 喀土穆")]
        [ExtraAttributeValue("2.0")]
        Sudan_0200 = 77020004,

        [Description("(UTC+02:00) 大馬士革")]
        [ExtraAttributeValue("2.0")]
        Syria_0200 = 77020005,

        [Description("(UTC+02:00) 奇西瑙")]
        [ExtraAttributeValue("2.0")]
        EEurope_0200 = 77020006,

        [Description("(UTC+02:00) 安曼")]
        [ExtraAttributeValue("2.0")]
        Jordan_0200 = 77020007,

        [Description("(UTC+02:00) 溫吐克")]
        [ExtraAttributeValue("2.0")]
        Namibia_0200 = 77020008,

        [Description("(UTC+02:00) 的黎波里")]
        [ExtraAttributeValue("2.0")]
        Libya_0200 = 77020009,
        
        [Description("(UTC+02:00) 耶路撒冷")]
        [ExtraAttributeValue("2.0")]
        Israel_0200 = 77020010,

        [Description("(UTC+02:00) 貝魯特")]
        [ExtraAttributeValue("2.0")]
        MiddleEast_0200 = 77020011,

        [Description("(UTC+02:00) 赫爾辛基，凱耶夫，里加，蘇非亞，塔林，維爾紐斯")]
        [ExtraAttributeValue("2.0")]
        FLE_0200 = 77020012,

        [Description("(UTC+02:00) 開羅")]
        [ExtraAttributeValue("2.0")]
        Egypt_0200 = 77020013,

        [Description("(UTC+02:00) 雅典，布加勒斯特")]
        [ExtraAttributeValue("2.0")]
        GTB_0200 = 77020014,

        [Description("(UTC+03:00) 伊斯坦堡")]
        [ExtraAttributeValue("3.0")]
        Turkey_0300 = 77030001,

        [Description("(UTC+03:00) 奈洛比")]
        [ExtraAttributeValue("3.0")]
        EAfrica_0300 = 77030002,

        [Description("(UTC+03:00) 巴格達")]
        [ExtraAttributeValue("3.0")]
        Arabic_0300 = 77030003,

        [Description("(UTC+03:00) 明斯克")]
        [ExtraAttributeValue("3.0")]
        Belarus_0300 = 77030004,

        [Description("(UTC+03:00) 科威特，利雅德")]
        [ExtraAttributeValue("3.0")]
        Arab_0300 = 77030005,

        [Description("(UTC+03:00) 莫斯科，聖彼得堡")]
        [ExtraAttributeValue("3.0")]
        Russian_0300 = 77030006,

        [Description("(UTC+03:30) 德黑蘭")]
        [ExtraAttributeValue("3.5")]
        Iran_0330 = 77033001,

        [Description("(UTC+04:00) 伊熱夫斯克，沙馬拉")]
        [ExtraAttributeValue("4.0")]
        RussiaTimeZone3_0400 = 77040001,

        [Description("(UTC+04:00) 伏爾加格勒")]
        [ExtraAttributeValue("4.0")]
        Volgograd_0400 = 77040002,

        [Description("(UTC+04:00) 巴庫")]
        [ExtraAttributeValue("4.0")]
        Azerbaijan_0400 = 77040003,

        [Description("(UTC+04:00) 第比利斯")]
        [ExtraAttributeValue("4.0")]
        Georgian_0400 = 77040004,

        [Description("(UTC+04:00) 葉里溫")]
        [ExtraAttributeValue("4.0")]
        Caucasus_0400 = 77040005,

        [Description("(UTC+04:00) 薩拉托夫")]
        [ExtraAttributeValue("4.0")]
        Saratov_0400 = 77040006,

        [Description("(UTC+04:00) 路易士港")]
        [ExtraAttributeValue("4.0")]
        Mauritius_0400 = 77040007,

        [Description("(UTC+04:00) 阿布達比，馬斯喀特")]
        [ExtraAttributeValue("4.0")]
        Arabian_0400 = 77040008,

        [Description("(UTC+04:00) 阿斯特拉罕，烏里揚諾夫斯克")]
        [ExtraAttributeValue("4.0")]
        Astrakhan_0400 = 77040009,

        [Description("(UTC+04:30) 喀布爾")]
        [ExtraAttributeValue("4.5")]
        Afghanistan_0430 = 77043001,

        [Description("(UTC+05:00) 伊卡特林堡")]
        [ExtraAttributeValue("5.0")]
        Ekaterinburg_0500 = 77050001,

        [Description("(UTC+05:00) 伊斯蘭馬巴德，克洛奇")]
        [ExtraAttributeValue("5.0")]
        Pakistan_0500 = 77050002,

        [Description("(UTC+05:00) 克孜勒奧爾達")]
        [ExtraAttributeValue("5.0")]
        Qyzylorda_0500 = 77050003,

        [Description("(UTC+05:00) 阿什哈巴德，塔什干")]
        [ExtraAttributeValue("5.0")]
        WestAsia_0500 = 77050004,

        [Description("(UTC+05:30) 斯里哈亞華登尼普拉")]
        [ExtraAttributeValue("5.5")]
        SriLanka_0530 = 77053001,

        [Description("(UTC+05:30) 辰內，加爾各答，孟拜，新德里")]
        [ExtraAttributeValue("5.5")]
        India_0530 = 77053002,

        [Description("(UTC+05:45) 加德滿都")]
        [ExtraAttributeValue("5.75")]
        Nepal_0545 = 77054501,

        [Description("(UTC+06:00) 達卡")]
        [ExtraAttributeValue("6.0")]
        Bangladesh_0600 = 77060001,

        [Description("(UTC+06:00) 鄂木斯克")]
        [ExtraAttributeValue("6.0")]
        Omsk_0600 = 77060002,

        [Description("(UTC+06:00) 阿斯坦納")]
        [ExtraAttributeValue("6.0")]
        CentralAsia_0600 = 77060003,

        [Description("(UTC+06:30) 仰光")]
        [ExtraAttributeValue("6.5")]
        Myanmar_0630 = 77063001,

        [Description("(UTC+07:00) 克拉斯諾亞爾斯克")]
        [ExtraAttributeValue("7.0")]
        NorthAsia_0700 = 77070001,

        [Description("(UTC+07:00) 哈霍德")]
        [ExtraAttributeValue("7.0")]
        WMongolia_0700 = 77070002,

        [Description("(UTC+07:00) 巴爾瑙爾，戈爾諾-阿爾泰斯克")]
        [ExtraAttributeValue("7.0")]
        Altai_0700 = 77070003,

        [Description("(UTC+07:00) 托木斯克")]
        [ExtraAttributeValue("7.0")]
        Tomsk_0700 = 77070004,

        [Description("(UTC+07:00) 新西伯利亞")]
        [ExtraAttributeValue("7.0")]
        NCentralAsia_0700 = 77070005,

        [Description("(UTC+07:00) 曼谷，河內，雅加達")]
        [ExtraAttributeValue("7.0")]
        SEAsia_0700 = 77070006,

        [Description("(UTC+08:00) 伊爾庫次克")]
        [ExtraAttributeValue("8.0")]
        NorthAsiaEast_0800 = 77080001,

        [Description("(UTC+08:00) 伯斯")]
        [ExtraAttributeValue("8.0")]
        WAustralia_0800 = 77080002,

        [Description("(UTC+08:00) 北京，重慶，香港特別行政區，烏魯木齊")]
        [ExtraAttributeValue("8.0")]
        China_0800 = 77080003,

        [Description("(UTC+08:00) 台北")]
        [ExtraAttributeValue("8.0")]
        Taipei_0800 = 77080004,

        [Description("(UTC+08:00) 吉隆坡，新加坡")]
        [ExtraAttributeValue("8.0")]
        Singapore_0800 = 77080005,

        [Description("(UTC+08:00) 庫倫")]
        [ExtraAttributeValue("8.0")]
        Ulaanbaatar_0800 = 77080006,

        [Description("(UTC+08:45) 尤克拉")]
        [ExtraAttributeValue("8.75")]
        AusCentralW_0845 = 77084501,

        [Description("(UTC+09:00) 亞庫茲克")]
        [ExtraAttributeValue("9.0")]
        Yakutsk_0900 = 77090001,

        [Description("(UTC+09:00) 大阪，北海道，東京")]
        [ExtraAttributeValue("9.0")]
        Tokyo_0900 = 77090002,

        [Description("(UTC+09:00) 平壤")]
        [ExtraAttributeValue("9.0")]
        NorthKorea_0900 = 77090003,

        [Description("(UTC+09:00) 赤塔")]
        [ExtraAttributeValue("9.0")]
        Transbaikal_0900 = 77090004,

        [Description("(UTC+09:00) 首爾")]
        [ExtraAttributeValue("9.0")]
        Korea_0900 = 77090005,

        [Description("(UTC+09:30) 愛德蘭")]
        [ExtraAttributeValue("9.0")]
        CenAustralia_0930 = 77093001,

        [Description("(UTC+09:30) 達爾文")]
        [ExtraAttributeValue("9.5")]
        AUSCentral_0930 = 77093002,

        [Description("(UTC+10:00) 坎培拉，墨爾本，雪梨")]
        [ExtraAttributeValue("10.0")]
        AUSEastern_1000 = 77100001,

        [Description("(UTC+10:00) 布里斯本")]
        [ExtraAttributeValue("10.0")]
        EAustralia_1000 = 77100002,

        [Description("(UTC+10:00) 海參威")]
        [ExtraAttributeValue("10.0")]
        Vladivostok_1000 = 77100003,

        [Description("(UTC+10:00) 關島，莫爾斯貝港")]
        [ExtraAttributeValue("10.0")]
        WestPacific_1000 = 77100004,

        [Description("(UTC+10:00) 霍巴特")]
        [ExtraAttributeValue("10.0")]
        Tasmania_1000 = 77100005,

        [Description("(UTC+10:30) 羅豪島")]
        [ExtraAttributeValue("10.5")]
        LordHowe_1030 = 77103001,

        [Description("(UTC+11:00) 喬庫爾達赫")]
        [ExtraAttributeValue("11.0")]
        RussiaTimeZone1_1100 = 77110001,

        [Description("(UTC+11:00) 布干維爾島")]
        [ExtraAttributeValue("11.0")]
        Bougainville_1100 = 77110002,

        [Description("(UTC+11:00) 庫頁島")]
        [ExtraAttributeValue("11.0")]
        Sakhalin_1100 = 77110003,

        [Description("(UTC+11:00) 索羅門群島，新喀里多尼亞群島")]
        [ExtraAttributeValue("11.0")]
        CentralPacific_1100 = 77110004,

        [Description("(UTC+11:00) 諾福克島")]
        [ExtraAttributeValue("11.0")]
        Norfolk_1100 = 77110005,

        [Description("(UTC+11:00) 馬加丹")]
        [ExtraAttributeValue("11.0")]
        Magadan_1100 = 77110006,

        [Description("(UTC+12:00) 國際標準時間+12")]
        [ExtraAttributeValue("12.0")]
        UTC12_1200 = 77120001,

        [Description("(UTC+12:00) 奧克蘭，威靈頓")]
        [ExtraAttributeValue("12.0")]
        NewZealand_1200 = 77120002,

        [Description("(UTC+12:00) 彼得保羅夫斯克-堪察加斯克 - 已過時")]
        [ExtraAttributeValue("12.0")]
        Kamchatka_1200 = 77120003,

        [Description("(UTC+12:00) 斐濟")]
        [ExtraAttributeValue("12.0")]
        Fiji_1200 = 77120004,

        [Description("(UTC+12:00) 阿納底，彼得保羅夫斯克-堪察加斯克")]
        [ExtraAttributeValue("12.0")]
        RussiaTimeZone11_1200 = 77120005,

        [Description("(UTC+12:45) 查塔姆群島")]
        [ExtraAttributeValue("12.75")]
        ChathamIslands_1245 = 77124501,

        [Description("(UTC+13:00) 努瓜婁發")]
        [ExtraAttributeValue("13.0")]
        Tonga_1300 = 77130001,

        [Description("(UTC+13:00) 國際標準時間+13")]
        [ExtraAttributeValue("13.0")]
        UTC13_1300 = 77130002,

        [Description("(UTC+13:00) 薩摩亞")]
        [ExtraAttributeValue("13.0")]
        Samoa_1300 = 77130003,

        [Description("(UTC+14:00) 刻里提瑪斯島")]
        [ExtraAttributeValue("14.0")]
        LineIslands_1400 = 77140001,


    }
}
