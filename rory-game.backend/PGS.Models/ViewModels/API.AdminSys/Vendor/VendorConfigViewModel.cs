using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class VendorConfigViewModel
    {
        public int VendorConfigID { get; set; }
        public int VendorID { get; set; }
        public string VendorName { get; set; }

        public string Name { get; set; }
        public string Memo { get; set; }
        public bool IsDefault { get; set; }
        public DateTime CreateTime { get; set; }
        public string CreateBy { get; set; }
        public DateTime UpdateTime { get; set; }
        public string UpdateBy { get; set; }

        /// <summary>
        /// 已分配使用的商戶數
        /// </summary>
        public int NumberOfMerchants { get; set; }


        /// <summary>
        /// 已分配使用的商戶
        /// </summary>
        public List<VendorConfigMappingMerchantViewModel> MappingMerchants { get; set; }

    }

    public class VendorConfigMappingMerchantViewModel
    {
        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

    }




    public class VendorConfigDeleteRequestModel
    {
        public int VendorConfigID { get; set; }

        public int MerchantID { get; set; }

    }


}
