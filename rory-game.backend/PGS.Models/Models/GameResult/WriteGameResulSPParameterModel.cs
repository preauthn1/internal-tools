using target.Models.Enums;
using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class WriteGameResulSPParameterModel
    {
        public int PlayerID { get; set; }
        public int ServerID { get; set; }
        public int GameID { get; set; }
        public long TradeID { get; set; }
        public decimal CurrentAmount { get; set; }
        public decimal UpdateAmount { get; set; }
        public int WalletMode { get; set; }
        public int WalletType { get; set; }
        public string IPAddress { get; set; }


        /// <summary>
        /// output參數
        /// </summary>
        public string result_msg { get; set; }




    }


}
