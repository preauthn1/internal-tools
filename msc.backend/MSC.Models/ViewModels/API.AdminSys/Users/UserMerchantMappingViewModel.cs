using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class UserMerchantMappingViewModel
    {
        public List<UserMerchantMappingItemViewModel> Items { get; set; }

    }

    public class UserMerchantMappingItemViewModel
    {
        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

        public bool IsEnable { get; set; }

    }


}
