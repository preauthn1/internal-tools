using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryPlayerSlotLightningRecordRequestViewModel : PaginationParameter
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
    }
}
