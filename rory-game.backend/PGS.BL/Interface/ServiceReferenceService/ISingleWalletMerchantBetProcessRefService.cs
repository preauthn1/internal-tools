using System;
using System.Collections.Generic;
using game.Entity.Game;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface ISingleWalletMerchantBetProcessRefService
    {

        SingleWalletModePlaceBetApiResponseModel PlaceBetProcessPassToMerchantApi(int merchantID, SingleWalletModePlaceBetParameter model);

        SingleWalletModeSettleApiResponseModel SettlemenetProcessPassToMerchantApi(SingleWalletModeSettleApiParameter model);

        SingleWalletModeSettleApiResponseModel BonusProcessPassToMerchantApi(SingleWalletModeSettleApiParameter model);

        SingleWalletModeSettleApiResponseModel CancelProcessPassToMerchantApi(SingleWalletModeCancelApiRequestModel model);
    }
}
