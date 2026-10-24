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

namespace MSC.Testing
{
    [TestClass]
    public class TokenCheckRefServiceTests : BaseTest
    {

        public ITokenCheckRefService ServiceInstance { get; set; }


        [TestMethod]
        public void TokenCheckRefService_CheckIsValidToken_Test()
        {
            LanguageEnum lang = LanguageEnum.zh_cn;

            string systemRequestID = RandomHelper.GetOrderNo("LK");
            var tokenInfo = new LinkTokenMemberInfoModel
            {
                MSCRequestID = systemRequestID,
                MerchantRequestID = Guid.NewGuid().ToString(),
                MerchantID = _TEST_DATA.Default.MerchantID,
                SignKey = SignKeyRuleHelper.GeneratePointTokenSignKey(_TEST_DATA.Default.MerchantID, _TEST_DATA.Default.MemberCode),
                MemberCode = _TEST_DATA.Default.MemberCode,
                ExpiredTime = DateTimeHelper.GetCurrentTime().AddHours(1),
                Lang = lang,
                AllowAmount = 150,
                ReturnUrl = _TEST_DATA.Default.CallbackUrl
            };

            string originalJson = tokenInfo.ConvertToJsonString();
            var token = [REDACTED], SystemConstant.API_Vendor_AES_KEY).Base64Encode();


            var query = ServiceInstance.CheckIsTokenFormatCorrect(token);

            Assert.IsNotNull(query);

            //預期會發生錯誤Exception
            Assert_ExpectExceptionThrow(()=> {

                ServiceInstance.CheckIsDbExistToken(token);

            } , ErrorCode.ThirdPartyApiFail_TokenExpired);

            return;
        }

    }
}

