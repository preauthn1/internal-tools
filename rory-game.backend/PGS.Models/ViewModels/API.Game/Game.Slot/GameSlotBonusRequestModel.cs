using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    /// <summary>
    /// 取得SlotBonusRecord資料 (回應參數)
    /// </summary>
    public class PlayerSlotBonusRecordViewModel
    {
        public List<PlayerSlotBonusRecordItemViewModel> SlotBonusRecord { get; set; }
    }

    /// <summary>
    /// SlotBonusRecord細項
    /// </summary>
    public class PlayerSlotBonusRecordItemViewModel
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
        /// 第幾次選擇
        /// </summary>
        public int TurnID { get; set; }

        /// <summary>
        /// 選擇位置
        /// </summary>
        public int Location { get; set; }

        /// <summary>
        /// 下注金額
        /// </summary>
        public decimal BetAmount { get; set; }

        /// <summary>
        /// 下注倍率
        /// </summary>
        public int BetRate { get; set; }

        /// <summary>
        /// 補充資料(自行做字串處理)
        /// </summary>
        public string Data { get; set; }

        /// <summary>
        /// 此遊戲獲得的來源 Pattern
        /// </summary>
        public int PatternID { get; set; }
    }
	
    public class GameSlotBonusRequestModel : GameSlotRequestModel
    {
        public int OriginID { get; set; }
    }
	
	public class GameSlotBonusRecordRequestModel : GameSlotBonusRequestModel
    {
        public int TurnID { get; set; }

        public int Location { get; set; }

        public decimal BetAmount { get; set; }

        public int BetRate { get; set; }

        public string Data { get; set; }

        public int PatternID { get; set; }
    }
	
}
