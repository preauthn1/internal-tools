using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class ServerRobotConfigCreateRequestModel : ServerRobotConfigUpdateRequestModel
    {        
      
    }


    public class ServerRobotConfigUpdateRequestModel
    {
        public int ServerID { get; set; }
        public int MinRobotCount { get; set; }
        public int MaxRobotCount { get; set; }
        public int MinLevel { get; set; }
        public int MaxLevel { get; set; }
        public bool IsRobotEnable { get; set; }
        public string BeginTime { get; set; }
        public string EndTime { get; set; }
    }



}
