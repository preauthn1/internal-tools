using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{


    public class GameSlotLightningRequestModel : GameSlotRequestModel
    {
        public string Data { get; set; }
        public int PatternID { get; set; }
    }

    /// <summary>
    /// 取得SlotLightingRecord資料 (回應參數)
    /// </summary>
    public class PlayerSlotLightingRecordViewModel
    {
        public List<PlayerSlotLightingRecordItemViewModel> SlotLightingRecord { get; set; }
    }

    /// <summary>
    /// SlotLightingRecord細項
    /// </summary>
    /// PlayerSlotLightingRecordViewModel
    public class PlayerSlotLightingRecordItemViewModel
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
        /// 遊戲資料(自行做字串處理)
        /// </summary>
        public string Data { get; set; }

        /// <summary>
        /// 此遊戲獲得的來源 Pattern
        /// </summary>
        public int PatternID { get; set; }
    }
}
