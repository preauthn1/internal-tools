using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryServerMerchantRelationRequestViewModel : PaginationParameter
    {

        public List<int> ServerIDs { get; set; }

        public List<int> MerchantIDs { get; set; }

    }


}
