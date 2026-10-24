using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using System;
using System.Collections.Generic;

namespace target.Models.ViewModels
{

    public class AGGLiveSingleWalletAPITransactionRequestModel
    {
        public string type { get; set; }
        public DateTime dtsent { get; set; }
        public AGGLiveBaseGame BaseGame { get; set; }
        public AGGLiveAuth Auth { get; set; }
        public AGGLiveFundTransferRequest FundTransferRequest { get; set; }
    }


    public class AGGLiveSingleWalletAPITransactionResponseViewModel
    {
        [JsonProperty(PropertyName = "fundtransferresponse")]
        public AGGLiveStatusOnlyResponseModel FundTransferResponse { get; set; }

        public decimal balance { get; set; }

        public string currencycode { get; set; }
    }



    public class AGGLiveFundInfo
    {
        public int gamestatemode { get; set; }
        public string transferid { get; set; }
        public string currencycode { get; set; }
        public decimal amount { get; set; }
        public bool jpwin { get; set; }
        public double jpcont { get; set; }
        public bool isbonus { get; set; }
        public DateTime dtevent { get; set; }
        public string initialdebittransferid { get; set; }
    }

    public class AGGLiveFunds
    {
        public bool debitandcredit { get; set; }
        public List<AGGLiveFundInfo> FundInfo { get; set; }
    }

    public class AGGLiveGameDetails
    {
        public string name { get; set; }
        public string keyname { get; set; }
        public int gametypeid { get; set; }
        public string gametypename { get; set; }
        public string brandgameid { get; set; }
        public string gamesessionid { get; set; }
        public string gameinstanceid { get; set; }
        public int friendlygameinstanceid { get; set; }
        public int channel { get; set; }
        public string device { get; set; }
        public string browser { get; set; }
    }

    public class AGGLiveFundTransferRequest
    {
        public string token { get; set; }
        public string accountid { get; set; }
        public int customplayertype { get; set; }
        public string gameinstanceid { get; set; }
        public int friendlygameinstanceid { get; set; }
        public bool isretry { get; set; }
        public int retrycount { get; set; }
        public bool isrefund { get; set; }
        public bool isrecredit { get; set; }
        public AGGLiveFunds Funds { get; set; }
        public AGGLiveGameDetails GameDetails { get; set; }
    }




}
