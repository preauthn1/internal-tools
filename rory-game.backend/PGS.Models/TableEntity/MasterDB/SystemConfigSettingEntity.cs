using target.Models.CoreLibs;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    [IsNotAutoKeyTable]
    public class SystemConfigSettingEntity
    {

        /// <summary>
        /// 請對應 SystemConfigSettingEnum
        /// </summary>
        [Key]
        public int ID { get; set; }
        public string SettingValue { get; set; }
        public string Memo { get; set; }


    }
}
