using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{



    public class QueryMerchantSummaryStatReportReportApiRequestParameterModel
    {
        public List<int> MerchantIDs { get; set; }

        public DateTime? Begin { get; set; }

        public DateTime? End { get; set; }

    }


    public class MerchantSummaryStatReportViewModel
    {
       public List<MerchantSummaryStatReportItemViewModel> Merchants { get; set; }

    }



    public class MerchantSummaryStatReportItemViewModel
    {
        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

        public MerchantSummaryStatReportItemBaseInfoViewModel BasicSummary { get; set; }

        public MerchantSummaryStatReportItemLoginInfoViewModel LoginSummary { get; set; }

        public MerchantSummaryStatReportItemGameInfoViewModel GameSummary { get; set; }

    }



    public class MerchantSummaryStatReportItemBaseInfoViewModel
    {
        public int TotalLoginPlayer { get; set; }
        public decimal TransferAmount { get; set; }
        public decimal TransferInAmount { get; set; }
        public decimal TransferOutAmount { get; set; }
        public decimal WalletAmount { get; set; }
    }



    public class MerchantSummaryStatReportItemLoginInfoViewModel
    {

        public int TotalLoginCount { get; set; }
        public int NewLoginPlayer { get; set; }
        public int MaxLoginPlayer { get; set; }
        public int AvgLoginPlayer { get; set; }


    }


    public class MerchantSummaryStatReportItemGameInfoViewModel
    {

        public int TotalRound { get; set; }
        public int TotalBetPlayer { get; set; }
        public decimal WinAmount { get; set; }
        public decimal LoseAmount { get; set; }
        public decimal MoneyChangeAmount { get; set; }
        public decimal BetAmount { get; set; }
        public decimal PayoutAmount { get; set; }
        public decimal ValidAmount { get; set; }
        public decimal WinLossAmount { get; set; }

    }

















}
