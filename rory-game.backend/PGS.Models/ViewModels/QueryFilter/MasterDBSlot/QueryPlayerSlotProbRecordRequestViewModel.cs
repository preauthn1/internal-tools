using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryGameSlotProbabilitySettingRequestViewModel : PaginationParameter
    {
        public int? PlayerID { get; set; }

        public int? GameID { get; set; }

        public int? FlagID { get; set; }

        public int? OriginID { get; set; }
    }
}
