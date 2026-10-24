using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryVendorConfigSettingRequestViewModel : PaginationParameter
    {
        public List<int> VendorConfigIDs { get; set; }

        /// <summary>
        /// 對應 VendorConfigSettingEnum
        /// </summary>
        public List<int> SettingKeys { get; set; }

    }


}
