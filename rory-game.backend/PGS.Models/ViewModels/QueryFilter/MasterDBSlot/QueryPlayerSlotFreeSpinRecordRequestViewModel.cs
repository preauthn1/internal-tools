using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryPlayerSlotFreeSpinRecordRequestViewModel : PaginationParameter
    {
        public int? PlayerID { get; set; }

        public int? GameID { get; set; }

        public int? FlagID { get; set; }

        public int? OriginID { get; set; }

        public int? FreeSpinGroup { get; set; }

        public int? ChooseRate { get; set; }

        public decimal? BetAmount { get; set; }

        public int? Times { get; set; }

        public decimal? WinAmount { get; set; }

        public int? WinTimes { get; set; }

        public string Data { get; set; }
    }
}
