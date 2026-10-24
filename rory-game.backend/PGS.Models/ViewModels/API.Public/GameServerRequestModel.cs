using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{



    public class GameServerRequestModel 
    {
        public int GameID { get; set; }

        public long Ts { get; set; }

    }

    public class GameServerViewModel
    {   
        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public int GameID { get; set; }

        public string GameName { get; set; }

        public decimal BalanceAmount { get; set; }

        public List<GameServerItemViewModel> Servers { get; set; }
    }


    public class GameServerItemViewModel
    {
        public int RoomLevel { get; set; }

        public string RoomLevelName { get; set; }

        public int FlagID { get; set; }

        public decimal MinAvailableMoney { get; set; }

        public string Url { get; set; }
    }




}
