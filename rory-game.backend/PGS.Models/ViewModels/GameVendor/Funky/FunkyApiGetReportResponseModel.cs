using System.Collections.Generic;

namespace target.Models.ViewModels.Base
{      

    public class FunkyApiGetReportResponseModel
    {
        public int errorCode { get; set; }
        public string errorMessage { get; set; }
        public int totalBetCount { get; set; }
        public int totalActivePlayer { get; set; }
        public List<FunkyEachCurrencyInfoModel> data { get; set; }
    }

    public class FunkyEachCurrencyInfoModel
    {
        public string currency { get; set; }
        public int betCount { get; set; }
        public decimal betAmount { get; set; }
        public decimal winloss { get; set; }
    }




}
