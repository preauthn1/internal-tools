using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{

    public class MerchantCustomizedSettingEntity
    {
        [Key]
        public int Id { get; set; }
        public int MerchantID { get; set; }
        public int CustomizedSettingKey { get; set; }
        public string SettingValue { get; set; }
        public string Memo { get; set; }
        public DateTime CreateTime { get; set; }


    }
}
