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
    public class LinkGenerateServiceTests : BaseTest
    {

        public ILinkGenerateService ServiceInstance { get; set; }

        public ITokenCheckRefService TokenCheckRefService { get; set; }


        [TestMethod]
        public void LinkGenerateService_GenerateCardLink_Test()
        {
            var requestInfo = new LinkGenerateRequestInfoModel
            {
                Amount = 0,
                Lang = SystemConstant.DefaultLang,
                MemberCode = _TEST_DATA.Default.MemberCode,
                MerchantID = _TEST_DATA.Default.MerchantID,
                RandomID = Guid.NewGuid().ToString(),
                RequestID = RandomHelper.GetOrderNo("UNI"),
                NotifyUrl = "http://0.0.0.0",
                ReturnUrl = "http://0.0.0.0",
                TS = DateTimeHelper.GetCurrentTimeTimestamp(),
            };

            var singKeyResult = SignKeyRuleHelper.GenerateVendorApiIntegrationSignKey(requestInfo , _TEST_DATA.Default.MerchantSecret);

            var query = ServiceInstance.GenerateCardLink(new LinkGenerateRequestModel
            {
                RequestInfo = requestInfo,
                SignKey = singKeyResult.AfterMD5
            });

            Assert.IsTrue(query.Link.IsNotNullNotWhiteSpace());


            Uri myUri = new Uri(query.Link);
            var tokenValue = HttpUtility.ParseQueryString(myUri.Query).Get("token");


            Assert.IsTrue(tokenValue.IsNotNullNotWhiteSpace());


            var checkResulk = TokenCheckRefService.CheckIsValidToken(tokenValue);


            //若可以成功反解析 則為正常token
            Assert.IsNotNull(checkResulk);
            return;
        }



        /// <summary>
        /// 確保v1版本的欄位加密結果跟以往一樣
        /// </summary>
        [TestMethod]
        public void LinkGenerateService_CheckV1_ModelSignResultIsDefine_Test()
        {
            var requestInfo = new LinkGenerateRequestInfoModel
            {
                Amount = 150,
                Lang = SystemConstant.DefaultLang,
                MemberCode = _TEST_DATA.Default.MemberCode,
                MerchantID = _TEST_DATA.Default.MerchantID,
                RandomID = Guid.NewGuid().ToString(),
                RequestID = RandomHelper.GetOrderNo("UNI"),
                NotifyUrl = "http://0.0.0.0",
                ReturnUrl = "http://0.0.0.0",
                TS = DateTimeHelper.GetCurrentTimeTimestamp(),
            };


            Dictionary<string, object> defineColumns = new Dictionary<string, object>();

            //model 多定義一個欄位/少一個欄位/欄位改名  都可能導致原本已經對接v1版本的廠商壞掉
            //因此強制用MSTest確保最初版本的欄位
            defineColumns.Add("Amount", requestInfo.Amount.ToString("f2"));
            defineColumns.Add("Lang", requestInfo.Lang);
            defineColumns.Add("MemberCode", requestInfo.MemberCode);
            defineColumns.Add("MerchantID", requestInfo.MerchantID);
            defineColumns.Add("NotifyUrl", requestInfo.NotifyUrl);
            defineColumns.Add("RandomID", requestInfo.RandomID);
            defineColumns.Add("RequestID", requestInfo.RequestID);
            defineColumns.Add("ReturnUrl", requestInfo.ReturnUrl);
            defineColumns.Add("TS", requestInfo.TS);


            string secretKey = "886da1bb0f3bd6ed";

            var mySign = SignKeyRuleHelper.GenerateVendorApiIntegrationSignKey(requestInfo, secretKey);
            var expectSign = SignKeyRuleHelper.GenerateVendorApiIntegrationSignKey(defineColumns, secretKey);


            Assert.IsTrue(mySign.AfterMD5 == expectSign.AfterMD5);
            return;
        }


    }
}

