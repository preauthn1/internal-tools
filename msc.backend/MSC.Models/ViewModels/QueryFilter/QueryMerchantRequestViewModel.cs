using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class QueryMerchantApiRequestParameterModel : QueryMerchantRequestViewModel
    {

        public int? MerchantID { get; set; }

    }


    public class QueryMerchantRequestViewModel : PaginationParameter
    {

        public List<int> MerchantIDs { get; set; }

        public int? Status { get; set; }

    }


}
