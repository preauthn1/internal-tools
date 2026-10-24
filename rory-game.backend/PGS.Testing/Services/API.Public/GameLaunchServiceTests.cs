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
using target.Models;
using target.Models.ViewModels.Base;
using NSubstitute;
using Microsoft.AspNetCore.Http;
using System.Collections.Generic;
using System;
using Microsoft.Extensions.Primitives;
using target.Utility.Helper.HttpContext;

namespace target.Testing
{
    [TestClass]
    public class GameLaunchServiceTests : BaseTest
    {

        public IGameLaunchService GameLaunchService { get; set; }


        [TestMethod]
        public void LaunchGameAndCheckPlayerExist_Test()
        {
            var result = GameLaunchService.LaunchGameAndCheckPlayerExist(new GameLaunchRequestModel
            {
                GameID = GameListEnum.DFDC_多福多财.ToInteger(),
                MerchantID = _TEST_DATA.Default.MerchantID,
                PlayerName = _TEST_DATA.Default.PlayerName,
            });

            Assert.IsNotNull(result);
            Assert.IsTrue(example.invalid.Access_Token.IsNotEmptyNull());
            Assert.IsTrue(example.invalid.PlayerID == _TEST_DATA.Default.PlayerID);
            return;
        }



        [TestMethod]
        public void GetGameServers_Test()
        {
            var myJwtToken = [REDACTED];

            MockHttpHeaderFakeAccessToken(myJwtToken.JwtToken);

            var result = GameLaunchService.GetGameServers(new GameServerRequestModel
            {
                GameID = _TEST_DATA.Default.GameID,
                Ts = DateTimeHelper.GetCurrentTimeTimestamp()
            });

            Assert.IsTrue(result.Servers.IsNotEmptyNull());
            Assert.IsTrue(result.PlayerID == _TEST_DATA.Default.PlayerID);
            return;
        }


        private void MockHttpHeaderFakeAccessToken(string myJwt)
        {
            var myHttpContext = Substitute.For<IHttpContextAccessor>();

            var mockHeaders = new HeaderDictionary(new Dictionary<String, StringValues>
                    {
                        { SystemConstant.HeaderKey_Lang ,LanguageEnum.zh_cn.ToString()},
                        { SystemConstant.HeaderKey_AccessToken ,myJwt}

                    }) as IHeaderDictionary;

            //設定假的實體的Add方法當傳入3,6 回傳 9
            myHttpContext.HttpContext.Request.Headers.Returns(mockHeaders);
            myHttpContext.HttpContext.Connection.RemoteIpAddress.Returns(new MockIPAddress(0));

            MyHttpContext.Configure(myHttpContext);
        }


        private PlayerJwtGenerateResultModel GenerateFakeJwt()
        {
            string rndID = Guid.NewGuid().ToString();

            var result = MyJwtHelper.GeneratePlayerJwt(
                _TEST_DATA.Default.PlayerID, 
                _TEST_DATA.Default.MerchantID, 
                _TEST_DATA.Default.PlayerName, 
                SystemConstant.DefaultLang , 
                rndID,
                rndID);

            return result;
        }


    }
}

