using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryStatGameReportApiRequestParameterModel
    {
        public List<int> MerchantIDs { get; set; }

        public List<int> GameIDs { get; set; }

        public DateTime? Begin { get; set; }

        public DateTime? End { get; set; }
    }



    public class StatGameReportViewModel
    {

        public List<StatGameReportMerchantItemViewModel> Merchants { get; set; }

    }

    public class StatGameReportMerchantItemViewModel
    {
        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

        public List<StatGameReportItemViewModel> Items { get; set; }
    }



    public class StatGameReportItemViewModel
    {
        public int GameID { get; set; }
        public string GameName { get; set; }
        public int TotalRound { get; set; }
        public int OnlinePlayer { get; set; }
        public decimal UpdateAmount { get; set; }
        public decimal WinAmount { get; set; }
        public decimal LoseAmount { get; set; }
        public decimal TaxAmount { get; set; }
        public decimal FlowAmount { get; set; }
        public decimal BetAmount { get; set; }
        public decimal ValidAmount { get; set; }
        public decimal PayoutAmount { get; set; }
        public decimal RTP { get; set; }




        public decimal JackpotPoolAmount { get; set; }

    }


}
