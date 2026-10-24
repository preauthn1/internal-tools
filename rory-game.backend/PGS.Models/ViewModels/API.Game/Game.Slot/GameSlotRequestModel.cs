using target.Models.Table;
using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    /// <summary>
    /// Game Slot Request Model
    /// </summary>
    public class GameSlotRequestModel
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
    }

    public class GameSlotLogRequestModel : PaginationParameter
    {
        public int PlayerID { get; set; }
        public int GameID { get; set; }
        public DateTime? GameTimeBegin { get; set; }
        public DateTime? GameTimeEnd { get; set; }
    }
}
