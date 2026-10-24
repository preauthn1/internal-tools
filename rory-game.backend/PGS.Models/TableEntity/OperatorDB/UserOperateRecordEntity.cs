using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class UserOperateRecordEntity
    {
        [Key]
        public long ID { get; set; }
        public int UserID { get; set; }
        public int MerchantID { get; set; }
        public int OperateType { get; set; }
        public long TargetID { get; set; }
        public string Memo { get; set; }
        public string ChangeLogInfo { get; set; }
        public DateTime CreateTime { get; set; }
        public string CreateIPAddress { get; set; }

    }
}
