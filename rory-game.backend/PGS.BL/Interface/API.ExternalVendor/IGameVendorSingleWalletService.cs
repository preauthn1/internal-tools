using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IGameVendorSingleWalletService
    {
        SingleWalletGetWalletBalanceResultViewModel GetWalletBalance(GameVendorTypeEnum vendorID ,SingleWalletGetWalletBalanceRequestParameter model);

        SingleWalletVendorCustomizedResponseConvertModel PlaceBet(GameVendorTypeEnum vendorID, SingleWalletBetDebitAmountRequestParameter model);

        SingleWalletVendorCustomizedResponseConvertModel SettleBet(GameVendorTypeEnum vendorID, SingleWalletBetSettleRequestParameter model);

        SingleWalletVendorCustomizedResponseConvertModel CancelBet(GameVendorTypeEnum vendorID, SingleWalletCancelBetRequestParameter model);


    }
}
