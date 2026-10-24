using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class PlayerGameRecordViewModel
    {

        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

        public int WalletMode { get; set; }

        public int GameID { get; set; }

        public string GameName { get; set; }

        public int ServerID { get; set; }

        public string ServerName { get; set; }

        public long TradeID { get; set; }

        public decimal LoginAmount { get; set; }

        public decimal UpdateAmount { get; set; }

        public decimal TaxAmount { get; set; }


        public int WinCount { get; set; }

        public int LoseCount { get; set; }

        public int DrawCount { get; set; }

        public string IPAddress { get; set; }

        public long GameTimeSecond { get; set; }

        public DateTime LoginTime { get; set; }

        public DateTime? LogoutTime { get; set; }

    }



}
