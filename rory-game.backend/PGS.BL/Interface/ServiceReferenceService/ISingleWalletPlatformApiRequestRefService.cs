using System;
using System.Collections.Generic;
using game.Entity.Game;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Utility;

namespace target.BL.Interface
{

    public interface ISingleWalletPlatformApiRequestRefService
    {
        SingleWalletModeWalletApiResponseModel GetWalletBalance(SingleWalletModeWalletApiRequestParameter model);

        SingleWalletModePlaceBetApiResponseModel PlaceBet(SingleWalletModePlaceBetParameter param);

        SingleWalletModeSettleApiResponseModel Settlement(SingleWalletModeSettleApiParameter param);

        SingleWalletModeSettleApiResponseModel BonusSettlement(SingleWalletModeSettleApiParameter param);

        SingleWalletModeCancelApiResponseModel Cancel(SingleWalletModeCancelApiRequestModel param);

        /// <summary>
        /// 因為單元測試mocl需要 所以寫成public
        /// </summary>
        SingleWalletMerchantSessionTokenEntity CheckIsValidSessionToken(string sessionToken);


    }
}
