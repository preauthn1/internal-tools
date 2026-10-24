using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryJackpotPoolRequestViewModel : PaginationParameter
    {
        public int? GameID { get; set; }
        public int? FlagID { get; set; }
        public int? LevelID { get; set; }
        public int? TableNo { get; set; }
        public bool? IsLink { get; set; }
    }
}
