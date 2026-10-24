using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class QueryMerchantWalletLogApiRequestParameterModel : QueryMerchantWalletLogRequestViewModel
    {

        public int? MerchantID { get; set; }
    }

    public class QueryMerchantWalletLogRequestViewModel : PaginationParameter
    {
        public List<int> IDs { get; set; }

        public List<int> MerchantIDs { get; set; }

        public DateTime? CreateTimeBegin { get; set; }

        public DateTime? CreateTimeEnd { get; set; }


    }


}
