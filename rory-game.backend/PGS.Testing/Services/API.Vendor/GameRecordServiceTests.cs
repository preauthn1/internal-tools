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
using target.Utility.IOC;

namespace target.Testing
{
    [TestClass]
    public class GameRecordServiceTests : BaseTest
    {

        public IGameRecordService ServiceInstance { get; set; }


        public BetRecordRepository BetRecordRepository { get; set; }


        [TestMethod]
        public void GetBetRecords_Test()
        {
            DateTime now = DateTimeHelper.GetCurrentTime();

            var query = ServiceInstance.GetBetRecords(new GameRecordRequestModel
            {
                PageIndex = 1,
                PageSize = 10,
                MerchantID = _TEST_DATA.Default.MerchantID,
                MerchantSecretKey = _TEST_DATA.Default.MerchantSecret,
                StartTime = now.AddDays(30 * -1),
                EndTime = now
            });

            Assert_NotEmptyCollection(query.Items);
            return;
        }


        [TestMethod]
        public void GetSingleWalletBetRecordByNo_NULL_Test()
        {           
            var query = ServiceInstance.GetBetRecordBySingleWalletTransactionNo(new SingleWalletGameRecordRequestModel
            {
                MerchantID = _TEST_DATA.Default.MerchantID,
                MerchantSecretKey = _TEST_DATA.Default.MerchantSecret,
                TransactionNo = "YYOOXX"
            });

            Assert.IsNull(query);
            return;
        }


        [TestMethod]
        public void GetSingleWalletBetRecordByNo_Test()
        {
            //清空其他測試注冊的mock
            IOCHelper.MockInstanceMapping = new Dictionary<string, object>();

            var query = ServiceInstance.GetBetRecordBySingleWalletTransactionNo(new SingleWalletGameRecordRequestModel
            {
                MerchantID = _TEST_DATA.Default.MerchantID,
                MerchantSecretKey = _TEST_DATA.Default.MerchantSecret,
                TransactionNo = _TEST_DATA.Default.TransactionNo_PgsBet
            });

            Assert.IsNotNull(query);
            return;
        }


        [TestMethod]
        public void GetReportSummary_Test()
        {
            DateTime now = DateTimeHelper.GetCurrentTime();

            var query = ServiceInstance.GetReportSummary(new VendorReportSummaryRequestModel
            {
                MerchantID = _TEST_DATA.Default.MerchantID,
                MerchantSecretKey = _TEST_DATA.Default.MerchantSecret,
                StartTime = now.AddDays(90 * -1),
                EndTime = now
            });

            Assert_NotEmptyCollection(query.GameSummary);
            return;
        }

    }
}

