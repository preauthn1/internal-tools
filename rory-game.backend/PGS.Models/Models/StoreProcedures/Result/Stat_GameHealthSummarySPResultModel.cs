using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Models.StoreProcedures
{
    public class Stat_GameHealthSummarySPResultModel
    {
        public int MerchantID { get; set; }
        public string MerchantName { get; set; }
        public int ServerID { get; set; }
        public string ServerName { get; set; }
        public int GameID { get; set; }
        public string GameName { get; set; }
        public int Status { get; set; }
        public int OnlinePlayer { get; set; }
        public decimal ProfitAmount { get; set; }
        public decimal FlowAmount { get; set; }
    }
}
