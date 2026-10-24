using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class VendorConfigSettingEntity
    {
        [Key]
        public int Id { get; set; }
        public int VendorConfigID { get; set; }

        /// <summary>
        /// 對應 VendorConfigSettingEnum
        /// </summary>
        public int SettingKey { get; set; }
        public string Value { get; set; }
        public string Memo { get; set; }
        public DateTime UpdateTime { get; set; }
        public string UpdateBy { get; set; }


    }
}
