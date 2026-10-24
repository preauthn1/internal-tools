using Microsoft.VisualStudio.TestTools.UnitTesting;
using MSC.BL.Interface;
using MSC.DA.Repository;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.Testing
{
    [TestClass]
    public class CardInputServiceTests : BaseTest
    {

        public ICardInputService ServiceInstance { get; set; }

        public CardLinkRequestRecordRepository CardLinkRequestRecordRepository { get; set; }


        /// <summary>
        /// 這支API 有時候MSTest會壞掉?
        /// </summary>
        [TestMethod]
        public void FakeTestRandomInputCard_Test()
        {
            string testMerchantRequestID = "44ae26b6-2c28-4f02-87bc-d086ac85c762";

            CardLinkRequestRecordRepository.MSTest_ResetIsDisabledFlag(_TEST_DATA.Default.MerchantID , testMerchantRequestID);

            try
            {
                var query = ServiceInstance.FakeTestRandomInputCard(new CardInputFakeTestRequestModel
                {
                    MerchantID = _TEST_DATA.Default.MerchantID,
                    MerchantRequestID = testMerchantRequestID
                });

                Assert.IsNotNull(query);
                Assert.IsTrue(!string.IsNullOrWhiteSpace(query.OrderNo));
            }
            catch(Exception ex)
            {
                //這裡會發生錯誤 原因不明 需要debug下
                string msg = ex.ToString();
                Assert.Fail(msg);
            }
          
            return;
        }

    }
}

