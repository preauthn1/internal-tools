using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{

    public class UserOperateLogEntity
    {
        [Key]
        public long Id { get; set; }
        public int UserID { get; set; }
        public string UserName { get; set; }
        public int OperationType { get; set; }
        public long? ReferenceID { get; set; }
        public string DebugInfo { get; set; }
        public string Memo { get; set; }

        public string IPAddress { get; set; }
        public string DeviceID { get; set; }
        public DateTime CreateTime { get; set; }


    }
}
