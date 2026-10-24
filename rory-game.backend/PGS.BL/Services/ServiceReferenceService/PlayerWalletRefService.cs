using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.BL.Interface;
using target.Models.Enums;
using target.Models.Models;
using System.Linq;
using target.Utility.Helper;
using target.Models;
using target.DA.Repository;
using target.Utility;
using game.Entity.Game;

namespace target.BL.Services
{

    public class PlayerWalletRefService : BaseGameService, IPlayerWalletRefService
    {
        public SingleWalletMerchantSessionTokenRepository SingleWalletMerchantSessionTokenRepository { get; set; }

        public PlayerWalletRepository PlayerWalletRepository { get; set; }

        public ISingleWalletPlatformApiRequestRefService SingleWalletCallApiRefService { get; set; }



        public PlayerWalletAmountBriefModel GetWalletAmountByMerchantWalletMode(int playerID , WalletTypeEnum walletType)
        {
            var queryPlayer = GetValidPlayerAndMerchant(playerID);

            MerchantWalletModeEnum walletMode = queryPlayer.WalletMode;

            switch(walletMode)
            {
                case MerchantWalletModeEnum.TransferMode:
                    return GetTransferModePlayerWallet(playerID);

                case MerchantWalletModeEnum.SingleWalletMode:
                    throw new TerproException(ApiErrorCode.NotImplement);

                default:
                    throw new TerproException(ApiErrorCode.NotImplement);
            }
        }



        public PlayerWalletAmountBriefModel GetWalletAmountByMerchantWalletMode(int playerID, string sessionToken)
        {
            var queryPlayer = GetValidPlayerAndMerchant(playerID);

            MerchantWalletModeEnum walletMode = queryPlayer.WalletMode;

            switch (walletMode)
            {
                case MerchantWalletModeEnum.TransferMode:
                    return GetTransferModePlayerWallet(playerID);

                case MerchantWalletModeEnum.SingleWalletMode:

                    var apiResult = SingleWalletCallApiRefService.GetWalletBalance(new SingleWalletModeWalletApiRequestParameter
                    {
                        PlayerID = playerID,
                        SessionTokenPGS = sessionToken
                    });

                    if (apiResult == null)
                        throw new TerproException(ApiErrorCode.ThirdPartyApiFail_ParsingFail);

                    return new PlayerWalletAmountBriefModel
                    {
                        Amount = apiResult.Balance,
                        WalletMode = walletMode
                    };

                default:
                    throw new TerproException(ApiErrorCode.NotImplement);
            }
        }


        public PlayerWalletAmountBriefModel GetWalletAmountByLatestToken(int playerID)
        {
            var queryPlayer = GetValidPlayerAndMerchant(playerID);

            MerchantWalletModeEnum walletMode = queryPlayer.WalletMode;

            switch (walletMode)
            {
                case MerchantWalletModeEnum.TransferMode:
                    return GetTransferModePlayerWallet(playerID);

                case MerchantWalletModeEnum.SingleWalletMode:

                    var queryLatest = SingleWalletMerchantSessionTokenRepository.GetPagingByFilter(new QuerySingleWalletMerchantSessionTokenRequestViewModel
                        {
                            PlayerIDs = playerID.ToListCollection(),
                            PageSize = 1
                        })
                        .Items.FirstOrDefault();

                    if (queryLatest == null)
                        throw new TerproException(ApiErrorCode.NotFoundSessionToken);

                    var apiResult = SingleWalletCallApiRefService.GetWalletBalance(new SingleWalletModeWalletApiRequestParameter
                    {
                        SessionTokenPGS = queryLatest.SessionTokenPGS,
                        PlayerID = playerID
                    });
                    return new PlayerWalletAmountBriefModel
                    {
                        Amount = apiResult.Balance,
                        WalletMode = walletMode
                    };

                default:
                    throw new TerproException(ApiErrorCode.NotImplement);
            }
        }


        private PlayerWalletAmountBriefModel GetTransferModePlayerWallet(int playerID)
        {
            var queryWallet = PlayerWalletRepository.GetPlayerWalletByType(playerID, WalletTypeEnum.GameWallet);

            if (queryWallet == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayerWallet);

            return new PlayerWalletAmountBriefModel
            {
                Amount = queryWallet.Amount,
                WalletMode = MerchantWalletModeEnum.TransferMode
            };
        }


    }
}
