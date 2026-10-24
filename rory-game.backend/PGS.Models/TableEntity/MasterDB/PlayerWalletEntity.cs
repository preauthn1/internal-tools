using target.Models.CoreLibs;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    [IsNotAutoKeyTable]
    public class PlayerWalletEntity
    {
        [Key]
        public int PlayerID { get; set; }

        public decimal Amount { get; set; }


        public int WalletType { get; set; }



    }
}
