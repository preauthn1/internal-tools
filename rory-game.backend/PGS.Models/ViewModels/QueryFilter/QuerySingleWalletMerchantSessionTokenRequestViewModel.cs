using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QuerySingleWalletMerchantSessionTokenRequestViewModel : PaginationParameter
    {

        public List<int> IDs { get; set; }

        public List<int> MerchantIDs { get; set; }

        public List<int> PlayerIDs { get; set; }

        public string SessionTokenPGS { get; set; }

        public string SessionTokenRaw { get; set; }

    }


}
