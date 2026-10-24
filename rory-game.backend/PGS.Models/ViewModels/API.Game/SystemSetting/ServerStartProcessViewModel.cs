using System.Collections.Generic;

namespace target.Models.ViewModels
{
    public class ServerStartProcessRequestModel
    {
        public int ServerID { get; set; }

        public int GameID { get; set; }

        public string ClientAddr { get; set; }

        public string HttpAddr { get; set; }

        public string NotifyAddr { get; set; }

        public decimal MinAvailableMoney { get; set; }

        public int FlagLevel { get; set; }

    }


    public class ServerCloseProcessRequestModel
    {
        public int ServerID { get; set; }

    }


    public class ServerStartProcessViewModel
    {
        public int GameID { get; set; }

        public int GameStatus { get; set; }

    }



}
