using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{
    public class CronJobSettingEntity
    {
        [Key]
        public int JobID { get; set; }
        public int JobType { get; set; }
        public string JobName { get; set; }
        public bool IsEnable { get; set; }
        public string CronExpression { get; set; }
        public string RequestUrl { get; set; }
        public DateTime StartTimeUTC { get; set; }
        public string Memo { get; set; }
        public DateTime CreateTime { get; set; }


    }
}
