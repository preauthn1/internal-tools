using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{


    public class GameSlotPlayerAssignedRequestModel : GameSlotRequestModel
    {
        public int Mode { get; set; }
        public DateTime ExpireTime { get; set; }
    }

    /// <summary>
    /// 取得SlotUserAssignedMode資料 (回應參數)
    /// </summary>
    public class GameSlotPlayerAssignedModeViewModel
    {
        public List<GameSlotPlayerAssignedModeItemViewModel> SlotUserAssignedMode { get; set; }
    }

    /// <summary>
    /// SlotUserAssignedMode
    /// </summary>
    public class GameSlotPlayerAssignedModeItemViewModel
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
        /// 建立時間
        /// </summary>
        public DateTime CreateTime { get; set; }

        /// <summary>
        /// 特殊模式 (0: 預設, 1: 放分, 2: 殺分)
        /// </summary>
        public int Mode { get; set; }

        /// <summary>
        /// 過期時間
        /// </summary>
        public DateTime ExpireTime { get; set; }

    }
}
