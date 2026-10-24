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
using target.DA.Repository;
using target.Utility.Helper;

namespace target.BL.Services
{

    public class PlayerWalletInfoService : BaseGameService, IPlayerWalletInfoService
    {

        public PlayerWalletRepository PlayerWalletRepository { get; set; }

        public IPlayerWalletRefService PlayerWalletRefService { get; set; }


        public PlayerBriefInfoViewModel GetPlayerInfo(int playerID)
        {
            var queryPlayer = PlayerRepository.GetById(playerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            var viewModel = queryPlayer.ReturnInjecteObject<PlayerBriefInfoViewModel>();
            return viewModel;
        }


        public PlayerWalletViewModel GetWalletAmount(int playerID , string accessJwtToken)
        {
            var jwtParse = MyJwtHelper.GetPlayerModelFromJwtString(accessJwtToken);

            //refService裡面有做player的檢查防呆判斷了
            var query = PlayerWalletRefService.GetWalletAmountByMerchantWalletMode(playerID , jwtParse?.SessionTokenPGS);

            var queryPlayer = PlayerRepository.GetSingleByCache(playerID);

            decimal walletAmount = query.Amount;
            decimal walletGamePoint = ConvertMoneyToGamePoint(queryPlayer.MerchantID, walletAmount);

            return new PlayerWalletViewModel
            {
                Amount = walletGamePoint,
                MerchantWalletMode = (int)query.WalletMode
            };
        }


    }
}
