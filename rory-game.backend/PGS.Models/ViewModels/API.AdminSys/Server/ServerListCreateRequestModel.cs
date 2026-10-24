using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class ServerListCreateRequestModel
    {
        public int ServerID { get; set; }
        public int GameID { get; set; }
        public string ServerName { get; set; }
        public string Currency { get; set; }
        public int FlagLevel { get; set; }
        public int Status { get; set; }
        public string ClientAddr { get; set; }
        public string NotifyAddr { get; set; }
        public string HttpAddr { get; set; }    
        public decimal MinAvailableMoney { get; set; }
        public int RoomFlagLevel { get; set; }


    }


    public class ServerListUpdateRequestModel
    {
        public int ServerID { get; set; }
        public int GameID { get; set; }
        public string ServerName { get; set; }
        public string Currency { get; set; }
        public int FlagLevel { get; set; }
        public int Status { get; set; }
        public string ClientAddr { get; set; }
        public string NotifyAddr { get; set; }
        public string HttpAddr { get; set; }
        public decimal MinAvailableMoney { get; set; }
        public int RoomFlagLevel { get; set; }

    }



}
