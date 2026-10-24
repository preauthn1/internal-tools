using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryMerchantGameSettingRequestViewModel : PaginationParameter
    {

        public List<int> MerchantIDs { get; set; }


        public List<int> GameIDs { get; set; }


        public int? Status { get; set; }

        public bool EnableGameOnly { get; set; }

    }


}
