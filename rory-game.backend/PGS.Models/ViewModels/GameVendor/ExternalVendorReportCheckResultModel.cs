using System;

namespace target.Models.ViewModels.GameVendor
{
    public class ExternalVendorReportCheckResultModel
    {

        public bool IsDataCorrect { get; set; }

        public decimal GameApiTotalBetAmount { get; set; }

        public decimal GameApiTotalWinLoss { get; set; }

        public decimal GameApiTotalPlayerCount { get; set; }

    }


}
