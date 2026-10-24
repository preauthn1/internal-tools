using target.Models.Table;
using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class GameSlotHistoryLogRequestModel
    {
        public List<string> Data { get; set; }
    }
	
	/// <summary>
    /// GameSlotHistoryLog
    /// </summary>
    public class GameSlotHistoryLogViewModel
    {
        /// <summary>
        /// RowNum
        /// </summary>
        public int RowNum { get; set; }
        /// <summary>
        /// 玩家ID
        /// </summary>
        public int PlayerID { get; set; }

        /// <summary>
        /// LogID
        /// </summary>
        public long LogID { get; set; }

        /// <summary>
        /// BetAmount
        /// </summary>
        public decimal BetAmount { get; set; }

        /// <summary>
        /// FinalAmount
        /// </summary>
        public decimal FinalAmount { get; set; }

        /// <summary>
        /// 剩餘次數
        /// </summary>
        public DateTime GameTime { get; set; }
    }
}
