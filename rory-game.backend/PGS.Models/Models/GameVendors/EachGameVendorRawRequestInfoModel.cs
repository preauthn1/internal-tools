using target.Models.Enums;
using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{



    /// <summary>
    /// 根據遊戲商不同 可能有不同的物件、原始欄位可操作
    /// 建議在 APIController將這些物件傳進class，以利後續interface操作
    /// </summary>
    public class EachGameVendorRawRequestInfoModel
    {

        public FunkyApiRawRequestInfoVendor Funky { get; set; }

        public AGGApiRawRequestInfoVendor AGG { get; set; }

        public AGGLiveApiRawRequestInfoVendor AGGLive { get; set; }

        public FishingApiRawRequestInfoVendor Fishing { get; set; }
    }


    public class FunkyApiRawRequestInfoVendor
    {
        public FunkySingleWalletAPIPlaceBetRequestModel PlaceBetInfo { get;set;}

        public FunkySingleWalletAPISettleBetRequestModel SettleBetInfo { get;set; }

        public FunkySingleWalletCancelBetRequestModel CancelBet { get; set; }

    }


    public class AGGApiRawRequestInfoVendor
    {
        public AGGSingleWalletAPITransactionRequestModel  PlaceBetInfo { get; set; }

        public AGGSingleWalletAPITransactionRequestModel SettleBetInfo { get; set; }

        public AGGSingleWalletAPITransactionRequestModel CancelBet { get; set; }

    }


    public class AGGLiveApiRawRequestInfoVendor
    {
        public AGGLiveSingleWalletAPITransactionRequestModel PlaceBetInfo { get; set; }

        public AGGLiveSingleWalletAPITransactionRequestModel SettleBetInfo { get; set; }

        public AGGLiveSingleWalletAPITransactionRequestModel CancelBet { get; set; }

    }

    public class FishingApiRawRequestInfoVendor
    {
        public FishingSingleWalletCheckAccountTokenRequestModel TokenAuth { get; set; }

        public FishingSingleWalletPlaceBetRequestModel PlaceBetInfo { get; set; }

        public FishingSingleWalletSettleBetRequestModel SettleBetInfo { get; set; }

        public AGGSingleWalletAPITransactionRequestModel CancelBet { get; set; }

    }


}
