using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.Models.Models;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Linq;
using target.Models.CoreLibs;
using target.Models.ViewModels;
using target.DA.Repository;
using target.Models.Enums;
using System.Collections.Generic;
using System;
using target.Utility;

namespace target.Testing
{
    [TestClass]
    public class WalletCashServiceTests : BaseTest
    {

        public IWalletCashService ServiceInstance { get; set; }


        [TestMethod]
        public void GetPlayerWalletCash_Test()
        {
            var query = ServiceInstance.GetPlayerWalletCash(new WalletCashPlayerAmountRequestModel
            {
                MerchantID = _TEST_DATA.Default.MerchantID,
                MerchantSecretKey = _TEST_DATA.Default.MerchantSecret,
                PlayerName = _TEST_DATA.Default.PlayerName
            });

            Assert.IsNotNull(query);
            Assert.IsTrue(query.Amount > 0);
            return;
        }



        [TestMethod]
        public void TransferCashToGameWallet_Test()
        {
            var query = ServiceInstance.TransferCashToGameWallet(new WalletTransferRequestModel
            {
                Amount = 10,
                MerchantOrderNo = RandomHelper.GetOrderNo("UNI"),
                MerchantID = _TEST_DATA.Default.MerchantID,
                MerchantSecretKey = _TEST_DATA.Default.MerchantSecret,
                PlayerName = _TEST_DATA.Default.PlayerName
            });

            Assert.IsNotNull(query);
            Assert.IsTrue(query.BeforeAmount > 0);
            Assert.IsTrue(!string.IsNullOrWhiteSpace(query.TransactionNo));
            return;
        }




        [TestMethod]
        public void TransferCashToPlatformWallet_Test()
        {
            var query = ServiceInstance.TransferCashToPlatformWallet(new WalletTransferRequestModel
            {
                Amount = 10,
                MerchantOrderNo = RandomHelper.GetOrderNo("UNI"),
                MerchantID = _TEST_DATA.Default.MerchantID,
                MerchantSecretKey = _TEST_DATA.Default.MerchantSecret,
                PlayerName = _TEST_DATA.Default.PlayerName
            });

            Assert.IsNotNull(query);
            Assert.IsTrue(query.BeforeAmount > 0);
            Assert.IsTrue(!string.IsNullOrWhiteSpace(query.TransactionNo));
            return;
        }



    }
}

