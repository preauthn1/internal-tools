using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class SuperToolPoolRequestModel
    {

        public int GameID { get; set; }


        public int FlagID { get; set; }

    }


    public class SuperToolPoolUpdateModel
    {

        public int GameID { get; set; }


        public int FlagID { get; set; }


        public decimal UpdatePoolAmount { get; set; }


        public string GameData { get; set; }

    }



    public class SuperToolPoolViewModel
    {
        public int GameID { get; set; }
        public int FlagID { get; set; }
        public decimal WinLine { get; set; }
        public decimal PoolAmount { get; set; }
        public string GameData { get; set; }


    }


}
