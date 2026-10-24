using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
	
	public class GameSlotFreeSpinRequestModel : GameSlotRequestModel
    {
        /// <summary>
        /// 機率表ID
        /// </summary>
        public int OriginID { get; set; }
        /// <summary>
        /// 玩家下注金額
        /// </summary>
        public decimal BetAmount { get; set; }

        /// <summary>
        /// FreeSpin次數
        /// </summary>
        public int Times { get; set; }

        /// <summary>
        /// FreeSpin群組
        /// </summary>
        public int FreeSpinGroup { get; set; }

        /// <summary>
        /// 玩家選擇賠率
        /// </summary>
        public int ChooseRate { get; set; }

        /// <summary>
        /// 玩家輸贏金額
        /// </summary>
        public decimal WinAmount { get; set; }

        /// <summary>
        /// 累計贏的次數
        /// </summary>
        public int WinTimes { get; set; }

        /// <summary>
        /// 額外資訊
        /// </summary>
        public string Data { get; set; }

        /// <summary>
        /// 模式ID
        /// </summary>
        public int PatternID { get; set; }
    }
	
    /// <summary>
    /// 取得SlotFreeSpinRecord資料 (回應參數)
    /// </summary>
    public class PlayerSlotFreeSpinRecordViewModel
    {
        public List<PlayerSlotFreeSpinRecordItemViewModel> SlotFreeSpinRecord { get; set; }
    }

    /// <summary>
    /// SlotFreeSpinRecord細項
    /// </summary>
    public class PlayerSlotFreeSpinRecordItemViewModel
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
        /// 
        /// </summary>
        public int FreeSpinGroup { get; set; }

        /// <summary>
        /// 
        /// </summary>
        public int ChooseRate { get; set; }

        /// <summary>
        /// 玩家投注
        /// </summary>
        public decimal BetAmount { get; set; }

        /// <summary>
        /// 
        /// </summary>
        public int Times { get; set; }

        /// <summary>
        /// 累積贏的金額
        /// </summary>
        public decimal WinAmount { get; set; }

        /// <summary>
        /// 
        /// </summary>
        public int WinTimes { get; set; }

        /// <summary>
        /// 額外資訊
        /// </summary>
        public string Data { get; set; }

        /// <summary>
        /// 此遊戲獲得的來源 Pattern
        /// </summary>
        public int PatternID { get; set; }
    }
}
