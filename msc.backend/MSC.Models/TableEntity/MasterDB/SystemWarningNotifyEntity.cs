using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{
    public class SystemWarningNotifyEntity
    {
        [Key]
        public int Id { get; set; }
        public int NotifyType { get; set; }
        public string NotifyChannel { get; set; }
        public string Message { get; set; }
        public bool IsSend { get; set; }
        public DateTime CreateTime { get; set; }
        public DateTime? SendTime { get; set; }

    }
}
