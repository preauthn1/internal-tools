using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class VendorConfigEntity
    {
        [Key]
        public int VendorConfigID { get; set; }
        public int VendorID { get; set; }
        public string Name { get; set; }
        public string Memo { get; set; }
        public bool IsDefault { get; set; }
        public DateTime CreateTime { get; set; }
        public string CreateBy { get; set; }
        public DateTime UpdateTime { get; set; }
        public string UpdateBy { get; set; }


    }
}
