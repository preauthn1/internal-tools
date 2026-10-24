using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QuerySlot88MD5RequestModel
    {
        public string MD5 { get; set; }

    }


    public class Slot88MD5InfoResultModel
    {      
        public string RandomID { get; set; }

        public DateTime? ExpiredTime { get; set; }

        public Slot88MD5PlayerInfoModel PlayerInfo { get; set; }

        public Slot88MDServerInfoModel ServerInfo { get; set; }

    }



    public class Slot88MDServerInfoModel
    {
        public int ServerID { get; set; }

        public string ServerConnection { get; set; }

        public int GameID { get; set; }

        public int GameVendorType { get; set; }

        public string ConnType { get; set; }
    }



    public class Slot88MD5PlayerInfoModel
    {
        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public string Lang { get; set; }

        public string Access_Token { get; set; }

    }


}
