using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{
    public class VendorConfigMerchantMappingEntity
    {
        [Key]
        public int Id { get; set; }
        public int MerchantID { get; set; }
        public int VendorID { get; set; }
        public int VendorConfigID { get; set; }


    }
}
