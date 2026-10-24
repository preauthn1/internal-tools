using System;
using System.Collections.Generic;
using target.BL.Interface;
using target.BL.Services;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Models.Models;
using System.Linq;
using target.DA.Repository;
using target.Utility;

namespace target.BL.Services
{

    public class WalletCashService : BaseVendorService, IWalletCashService
    {
        public PlayerWalletRepository PlayerWalletRepository { get; set; }

        public PlayerRepository PlayerRepository { get; set; }

        public PlayerCurrentServerRepository PlayerCurrentServerRepository { get; set; }

        public SystemConfigSettingRepository SystemConfigSettingRepository { get; set; }

        public VendorSessionTokenRepository VendorSessionTokenRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }



        private const decimal Min_Transfer_Amount = 0.01M;


        public WalletCashPlayerAmountViewModel GetPlayerWalletCash(WalletCashPlayerAmountRequestModel model)
        {
            model.RuleFor(x => x.PlayerName).Rule_NotEmptyText(); 
            
            //檢查廠商token合法
            this.GetMerchantAndCheckIsValid(model.MerchantID, model.MerchantSecretKey);
            
            var queryPlayer = PlayerRepository.GetByPlayerNameAndMerchantID(model.PlayerName , model.MerchantID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            var queryWallet = PlayerWalletRepository.GetPlayerWalletByType(queryPlayer.PlayerID , Models.Enums.WalletTypeEnum.GameWallet);

            //這邊不顯示gamePoint 避免商戶誤以為有多100倍的金幣可以提現

            return new WalletCashPlayerAmountViewModel
            {
                Amount = queryWallet.Amount,
                Currency = queryPlayer.Currency
            };
        }


        public WalletTransferResultViewModel TransferCashToGameWallet(WalletTransferRequestModel model)
        {
            model.RuleFor(x => x.Amount).Rule_NotMinus();
            model.RuleFor(x => x.Amount).Rule_NotTooSmall(minValue: Min_Transfer_Amount);
            model.RuleFor(x => x.Amount).Rule_IsDecimalPointPlacesTwo();
            model.RuleFor(x => x.PlayerName).Rule_NotEmptyText();
            model.RuleFor(x => x.MerchantOrderNo).Rule_NotEmptyText();

            this.GetMerchantAndCheckIsValid(model.MerchantID, model.MerchantSecretKey);

            var queryPlayer = PlayerRepository.GetByPlayerNameAndMerchantID(model.PlayerName, model.MerchantID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            //檢查系統維護
            SystemConfigSettingRepository.CheckSystemMaintainAndThrowException();

            var insertModel = new PlayerWalletTransferModel
            {                
                Amount = model.Amount,
                OrderNo_PGS = RandomHelper.GetOrderNo(""),
                OrderNo_Merchant = model.MerchantOrderNo,
                PlayerID = queryPlayer.PlayerID,
                TranferType = Models.Enums.WalletTransferTypeEnum.Transfer_In,
            };

            //遊戲房中 禁止轉帳
            CheckPlayerIsPlayingGame(queryPlayer.PlayerID);

            //呼叫SP 進行轉帳
            var result = PlayerWalletRepository.WalletTransfer(insertModel);

            if (!result.IsSuccess)
                throw new TerproException(ApiErrorCode.WalletTranferFail_WalletAmountNotEnough);

            return new WalletTransferResultViewModel
            {
                Amount = model.Amount,
                BeforeAmount = result.BeforeAmount,
                AfterAmount = result.AfterAmount,
                TransactionNo = insertModel.OrderNo_PGS
            };
        }


        public WalletTransferResultViewModel TransferCashToPlatformWallet(WalletTransferRequestModel model)
        {
            model.RuleFor(x => x.Amount).Rule_NotZeroNotMinus().Rule_NotTooSmall(minValue: Min_Transfer_Amount);
            model.RuleFor(x => x.PlayerName).Rule_NotEmptyText();

            this.GetMerchantAndCheckIsValid(model.MerchantID, model.MerchantSecretKey);

            var queryPlayer = PlayerRepository.GetByPlayerNameAndMerchantID(model.PlayerName, model.MerchantID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            //檢查系統維護
            SystemConfigSettingRepository.CheckSystemMaintainAndThrowException();

            var insertModel = new PlayerWalletTransferModel
            {
                Amount = model.Amount,
                OrderNo_PGS = RandomHelper.GetOrderNo(""),
                OrderNo_Merchant = model.MerchantOrderNo,
                PlayerID = queryPlayer.PlayerID,
                TranferType = Models.Enums.WalletTransferTypeEnum.Transfer_Out,
            };


            //為了避免潛在風險 將錢轉出後 立刻登出三方遊戲
            VendorSessionTokenRepository.ResetAllPlayerSessionToken(queryPlayer.PlayerID);

            //遊戲房中 禁止轉帳
            CheckPlayerIsPlayingGame(queryPlayer.PlayerID);

            //呼叫SP 進行轉帳
            var result = PlayerWalletRepository.WalletTransfer(insertModel);
            
            if (!result.IsSuccess)
                throw new TerproException(ApiErrorCode.WalletTranferFail_WalletAmountNotEnough);


            return new WalletTransferResultViewModel
            {
                Amount = model.Amount,
                BeforeAmount = result.BeforeAmount,
                AfterAmount = result.AfterAmount,
                TransactionNo = insertModel.OrderNo_PGS
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
