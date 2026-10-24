using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QueryVendorConfigMerchantMappingRequestViewModel : PaginationParameter
    {
        public List<int> VendorIDs { get; set; }

        public List<int> VendorConfigIDs { get; set; }

        public List<int> MerchantIDs { get; set; }

        public bool IsIncludeDeleteMerchant { get; set; }

    }


}
