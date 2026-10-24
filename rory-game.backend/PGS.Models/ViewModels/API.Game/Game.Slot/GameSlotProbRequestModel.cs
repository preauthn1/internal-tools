using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class GameSlotProbRequestModel : GameSlotRequestModel
    {
        /// <summary>
        /// 機率表ID
        /// </summary>
        public int OriginID { get; set; }

        /// <summary>
        /// 剩餘次數
        /// </summary>
        public int Times { get; set; }

        /// <summary>
        /// 補充資料(自行做字串處理)
        /// </summary>
        public string Data { get; set; }

        /// <summary>
        /// 模式ID
        /// </summary>
        public int PatternID { get; set; }
    }

    /// <summary>
    /// 取得SlotProbRecord資料 (回應參數)
    /// </summary>
    public class GameSlotProbabilitySettingViewModel
    {
        public List<GameSlotProbabilitySettingItemViewModel> SlotProbRecord { get; set; }
    }

    /// <summary>
    /// SlotProbRecord細項
    /// </summary>
    public class GameSlotProbabilitySettingItemViewModel
    {
        /// <summary>
        /// 玩家ID
        /// </summary>
        public int PlayerID { get; set; }

        /// <summary>
        /// 遊戲ID
        /// </summary>
        public int GameID { get; set; }

        /// <summary>
        /// 房間類型ID
        /// </summary>
        public int FlagID { get; set; }

        /// <summary>
        /// 機率表ID
        /// </summary>
        public int OriginID { get; set; }

        /// <summary>
        /// 剩餘次數
        /// </summary>
        public int Times { get; set; }

        /// <summary>
        /// 補充資料(自行做字串處理)
        /// </summary>
        public string Data { get; set; }

        /// <summary>
        /// 此遊戲獲得的來源 Pattern
        /// </summary>
        public int PatternID { get; set; }
    }
}
