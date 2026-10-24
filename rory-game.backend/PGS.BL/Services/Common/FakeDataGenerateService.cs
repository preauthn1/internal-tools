using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.BL.Interface;
using target.Models.Enums;
using target.DA.Repository;
using target.Models.Models;
using System.Linq;
using target.Utility.Helper;
using target.Models;
using target.Utility;

namespace target.BL.Services
{

    public class FakeDataGenerateService : BaseService, IFakeDataGenerateService
    {
        public PlayerRepository PlayerRepository { get; set; }

        public PlayerWalletRepository PlayerWalletRepository { get; set; }

        public PlayerCurrentServerRepository PlayerCurrentServerRepository { get; set; }


        /// <summary>
        /// 提供測試用 快速產生jwt
        /// </summary>
        public GameLaunchViewModel GenerateJwtTokenForQATesting(FakeDataGenerateRequestModel model)
        {
            model.RuleFor(x => x.PlayerName).Rule_NotEmptyText();

            var queryPlayers = PlayerRepository.GetByPlayerName(model.PlayerName , model.MerchantID);

            if (queryPlayers.IsEmptyOrNull())
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            if(queryPlayers.Count() >= 2)
                throw new TerproException(ApiErrorCode.BusinessLogicDeny, "此用户ID有两笔以上资料，请指定商户号");

            var matchUser = queryPlayers.FirstOrDefault();

            string rndSession = Guid.NewGuid().ToString();

            //產生jwtToken
            var jwtResult = MyJwtHelper.GeneratePlayerJwt(matchUser.PlayerID,  matchUser.MerchantID, matchUser.PlayerName, SystemConstant.DefaultLang , rndSession , rndSession);

            string jwtHash = jwtResult.JwtToken.MD5Hash();

            //更新token (讓前次的token自動過期)
            PlayerRepository.UpdateTokenAndLastTime(matchUser.PlayerID, jwtHash , IPHelper.GetIP());

            return new GameLaunchViewModel
            {
                Info = new GameLaunchRequiredInfoViewModel
                {
                    PlayerName = matchUser.PlayerName,
                    PlayerID = matchUser.PlayerID,
                    Access_Token = [REDACTED]
                }             
            };
        }



        /// <summary>
        /// 提供測試用 快速幫用戶錢包加錢
        /// </summary>
        public PlayerWalletAddResult AddPlayerWallet(int playerID)
        {
            if (SettingHelper.IsProduction)
                throw new TerproException(ApiErrorCode.Forbidden_DoNotAllowWhenProduction);

            var queryPlayer = PlayerRepository.GetById(playerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            decimal updateAmount = 10000;
            var insertModel = new PlayerWalletTransferModel
            {
                Amount = updateAmount,
                OrderNo_PGS = RandomHelper.GetOrderNo(""),
                OrderNo_Merchant = RandomHelper.GetOrderNo("FakeData"),
                PlayerID = queryPlayer.PlayerID,
                TranferType = Models.Enums.WalletTransferTypeEnum.Transfer_In,
            };

            //遊戲房中 禁止轉帳
            CheckPlayerIsPlayingGame(queryPlayer.PlayerID);

            //呼叫SP 進行轉帳
            var result = PlayerWalletRepository.WalletTransfer(insertModel);

            if (!result.IsSuccess)
                throw new TerproException(ApiErrorCode.WalletTranferFail_WalletAmountNotEnough);

            return new PlayerWalletAddResult
            {
                UpdateAmount = updateAmount,
                BeforeAmount = result.BeforeAmount,
                AfterAmount = result.AfterAmount,
            };
        }


        /// <summary>
        /// 檢查用戶是否正在遊戲中
        /// </summary>
        private void CheckPlayerIsPlayingGame(int playerID)
        {
            bool isExist = PlayerCurrentServerRepository.IsExist(playerID);

            if (isExist)
            {
                throw new TerproException(ApiErrorCode.Forbidden_WhenPlayingGame_NotAllowTransfer);
            }
        }




    }
}
