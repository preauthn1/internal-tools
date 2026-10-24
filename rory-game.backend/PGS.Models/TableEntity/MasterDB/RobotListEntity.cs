using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class RobotListEntity
    {
        [Key]
        public int PlayerID { get; set; }
        public int AssignServerID { get; set; }
        public int LevelID { get; set; }
        public string Group {get;set;}
        public int Sort { get; set; }

    }
}
