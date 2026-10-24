using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class GameHealthSummaryRequestModel
    {
        public int MerchantID { get; set; }

        public double TimeZone { get; set; }
    }


    public class GameHealthSummaryViewModel
    {
        public List<GameHealthSummaryItemViewModel> Items { get; set; }


        public GameHealthSummaryItemViewModel RowTotal { get; set; }

    }


    public class GameHealthSummaryItemViewModel
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
