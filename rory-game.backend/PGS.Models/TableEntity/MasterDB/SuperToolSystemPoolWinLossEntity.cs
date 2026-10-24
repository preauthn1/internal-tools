using target.Models.CoreLibs;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    [IsNotAutoKeyTable]
    public class SuperToolSystemPoolWinLossEntity
    {        
        public int GameID { get; set; }
        public int FlagID { get; set; }
        public decimal WinLine { get; set; }
        public decimal PoolAmount { get; set; }
        public string GameData { get; set; }


    }
}
