using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryPlayerSlotAssignedModeRequestViewModel : PaginationParameter
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
        /// 特殊模式 (0: 預設, 1: 放分, 2: 殺分)
        /// </summary>
        public int? Mode { get; set; }

        /// <summary>
        /// 過期時間
        /// </summary>
        public DateTime ExpireTime { get; set; }

    }
}
