using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class QueryMerchantWalletRequestViewModel : PaginationParameter
    {

        public List<int> MerchantIDs { get; set; }

        public List<int> WalletTypes { get; set; }


    }


}
