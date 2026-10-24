using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryPlayerSlotBonusRecordRequestViewModel : PaginationParameter
    {
        /// <summary>
        /// 玩家ID
        /// </summary>
        public int? PlayerID { get; set; }

        /// <summary>
        /// GameID
        /// </summary>
        public int? GameID { get; set; }

        /// <summary>
        /// FlagID
        /// </summary>
        public int? FlagID { get; set; }

        /// <summary>
        /// 規則ID
        /// </summary>
        public int? OriginID { get; set; }

        /// <summary>
        /// 第幾次選擇
        /// </summary>
        public int? TurnID { get; set; }

        /// <summary>
        /// 選擇位置
        /// </summary>
        public int? Location { get; set; }

        /// <summary>
        /// 下注金額
        /// </summary>
        public decimal? BetAmount { get; set; }

        /// <summary>
        /// 下注倍率
        /// </summary>
        public int? BetRate { get; set; }



        public bool IsLockQuery { get; set; }



    }
}
