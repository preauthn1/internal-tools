using MSC.Models.CoreLibs;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{

    [IsNotAutoKeyTable]
    public class SystemSettingConfigEntity
    {
        
        [Key]
        public int SettingKey { get; set; }
        public string SettingValue { get; set; }
        public string Memo { get; set; }


    }
}
