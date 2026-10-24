using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
  
    public enum GameRoomFlagLevelEnum
    {
        //初級場
        //精英场
        //豪華場
        //富豪場

        [Description(i18nKeys.Enums_GameRoomFlagLevelEnum_Level_1_Basic)]
        Level_1_Basic = 1,

        [Description(i18nKeys.Enums_GameRoomFlagLevelEnum_Level_2_Elite)]
        Level_2_Elite = 2,

        [Description(i18nKeys.Enums_GameRoomFlagLevelEnum_Level_3_Luxury)]
        Level_3_Luxury = 3,

        [Description(i18nKeys.Enums_GameRoomFlagLevelEnum_Level_4_RichBoss)]
        Level_4_RichBoss = 4,


        //有到9的是鬥地主獨有
        //而且鬥地主是 菁英 高级 富豪
        //然后他有分 洗牌跟经典
        //两种玩法
    }
}
