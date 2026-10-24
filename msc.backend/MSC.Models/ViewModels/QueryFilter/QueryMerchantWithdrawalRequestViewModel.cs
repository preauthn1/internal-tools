using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class QueryMerchantWithdrawalApiRequestParameterModel : QueryMerchantWithdrawalRequestViewModel
    {
     
    }



    public class QueryMerchantWithdrawalRequestViewModel : PaginationParameter
    {

        public List<int> IDs { get; set; }


        public List<int> MerchantIDs { get; set; }

        
    }


}
