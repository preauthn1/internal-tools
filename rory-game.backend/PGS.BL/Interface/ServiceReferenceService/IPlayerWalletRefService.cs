using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IPlayerWalletRefService
    {
        PlayerWalletAmountBriefModel GetWalletAmountByMerchantWalletMode(int playerID, WalletTypeEnum walletType);

        PlayerWalletAmountBriefModel GetWalletAmountByMerchantWalletMode(int playerID, string sessionTokenPGS);

        PlayerWalletAmountBriefModel GetWalletAmountByLatestToken(int playerID);
    }
}
