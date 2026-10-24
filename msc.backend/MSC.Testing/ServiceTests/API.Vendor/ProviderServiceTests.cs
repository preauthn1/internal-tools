using Microsoft.VisualStudio.TestTools.UnitTesting;
using MSC.BL.Interface;
using MSC.Models.Models;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System.Linq;
using MSC.Models.CoreLibs;
using MSC.Models.ViewModels;
using MSC.DA.Repository;
using MSC.Models.Enums;
using System.Collections.Generic;
using System;
using MSC.Utility;
using MSC.Models;
using System.Web;

namespace MSC.Testing
{
    [TestClass]
    public class ProviderServiceTests : BaseTest
    {

        public ICardGenerateService CardGenerateService { get; set; }


        /// <summary>
        /// 確保v1版本的欄位加密結果跟以往一樣
        /// </summary>
        [TestMethod]
        public void CardGenerateService_GenerateProviderApiIntegrationSignKey_Test()
        {
            var requestInfo = new CardProviderCardGenerateRequestDetailModel
            {
                Memo = "Test",
                ProviderID = 552001,
                TransactionNo = Guid.NewGuid().ToString(),
                Point = 150,
                TS = DateTimeHelper.GetCurrentTimeTimestamp(),
            };


            Dictionary<string, object> defineColumns = new Dictionary<string, object>();

            //model 多定義一個欄位/少一個欄位/欄位改名  都可能導致原本已經對接v1版本的廠商壞掉
            //因此強制用MSTest確保最初版本的欄位
            defineColumns.Add("ProviderID", requestInfo.ProviderID);
            defineColumns.Add("TransactionNo", requestInfo.TransactionNo);
            defineColumns.Add("Point", requestInfo.Point.ToString("f2"));
            defineColumns.Add("Memo", requestInfo.Memo);
            defineColumns.Add("TS", requestInfo.TS);


            string secretKey = "a2af45661a184a1da7dc4273f0217a88";

            var mySign = SignKeyRuleHelper.GenerateProviderApiIntegrationSignKey(requestInfo, secretKey);
            var expectSign = SignKeyRuleHelper.GenerateProviderApiIntegrationSignKey(defineColumns, secretKey);


            Assert.IsTrue(mySign.AfterMD5 == expectSign.AfterMD5);
            return;
        }


    }
}

