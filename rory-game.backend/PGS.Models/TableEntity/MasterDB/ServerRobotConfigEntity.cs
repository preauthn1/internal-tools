using target.Models.CoreLibs;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{
    [IsNotAutoKeyTable]
    public class ServerRobotConfigEntity
    {
        [Key]
        public int ServerID { get; set; }
        public int MinRobotCount { get; set; }
        public int MaxRobotCount { get; set; }
        public int MinLevel { get; set; }
        public int MaxLevel { get; set; }

        public bool IsRobotEnable { get; set; }
        public string BeginTime { get; set; }
        public string EndTime { get; set; }


        public DateTime CreateTime { get; set; }

        public string CreateBy { get; set; }

        public DateTime UpdateTime { get; set; }

        public string UpdateBy { get; set; }



    }
}
