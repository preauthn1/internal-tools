using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using target.API.Vendor.Controllers;
using target.API.AdminSys;

namespace target.API.Vendor.Controllers
{
    [Route("v{version:int}")]
    public class WalletTransferController : BaseController
    {
        public IWalletCashService WalletCashService { get; set; }

        public ITransferRecordVendorService TransferRecordVendorService { get; set; }

        /// <summary>
        /// 2.4.1 Get player wallet (取得玩家钱包余额)
        /// </summary>
        [SwaggerApiDescriptionEN("Get player wallet balance. <BR> (the player's account should `exist` when access this API)")]
        [SwaggerApiDescriptionCN("取得玩家錢包金额。<BR> 呼叫此API前，`须先建立玩家账号` (至少调用一次游戏启动、或调用建立账号API)")]
        [Route("wallet/amount")]
        [ApiExplorerSettings(GroupName = "tfMode")]
        [HttpGet]
        public WalletCashPlayerAmountViewModel GetUserWalletCash([FromQuery]WalletCashPlayerAmountRequestModel model)
        {
            var query = WalletCashService.GetPlayerWalletCash(model);
            return query;
        }


        /// <summary>
        /// 2.4.2 Wallet transfer from platform to gameWallet (玩家钱包转账)
        /// </summary>
        [SwaggerApiDescriptionEN("(The player's account should `exist` when access this API) <BR> (Transfer is `not allowed` while the player is playing in game)")]
        [SwaggerApiDescriptionCN("呼叫此API前，须先`建立玩家账号` (至少调用一次游戏启动、或调用建立账号API) <BR> 玩家进行游戏过程中，`不允许进行钱包转帐`")]
        [Route("wallet/transfer-in")]
        [ApiExplorerSettings(GroupName = "tfMode")]
        [HttpPost]
        public WalletTransferResultViewModel TransferCashToGameWallet(WalletTransferRequestModel model)
        {
            var query = WalletCashService.TransferCashToGameWallet(model);
            return query;
        }


        /// <summary>
        /// 2.4.3 Wallet transfer from gameWallet to platform (玩家钱包转账)
        /// </summary>
        [SwaggerApiDescriptionEN("(The player's account should `exist` when access this API) <BR> (Transfer is `not allowed` while the player is playing in game)")]
        [SwaggerApiDescriptionCN("呼叫此API前，须先`建立玩家账号` (至少调用一次游戏启动、或调用建立账号API) <BR> 玩家进行游戏过程中，`不允许进行钱包转帐`")]
        [Route("wallet/transfer-out")]
        [ApiExplorerSettings(GroupName = "tfMode")]
        [HttpPost]
        public WalletTransferResultViewModel TransferCashToPlatformWallet(WalletTransferRequestModel model)
        {
            var query = WalletCashService.TransferCashToPlatformWallet(model);
            return query;
        }


        /// <summary>
        /// 2.4.4 Query transfer record (查询转账记录)
        /// </summary>
        [SwaggerApiDescriptionEN("Get transfer record by transactionID. (if data not exist. api will response `null`)")]
        [SwaggerApiDescriptionCN("根据转账单号取得轉帳紀錄。若查无资料返回data: null ")]
        [Route("transfer")]
        [ApiExplorerSettings(GroupName = "tfMode")]
        [HttpGet]
        public TransferRecordVendorViewModel GetTransferByMerchant([FromQuery] QueryTransferRecordRequestModel filter)
        {
            var query = TransferRecordVendorService.GetTransferByMerchant(filter);
            return query;
        }


    }
}
