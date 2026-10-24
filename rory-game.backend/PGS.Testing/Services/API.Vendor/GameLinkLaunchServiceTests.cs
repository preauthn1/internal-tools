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
using target.Utility.IOC;
using target.Models.ViewModels.GameVendor;
using target.BL.Services;

namespace target.Testing
{
    [TestClass]
    public class GameLinkLaunchServiceTests : BaseTest
    {

        public IGameLinkService GameLinkService { get; set; }

        public IPlayerJwtCheckService IPlayerJwtCheckService { get; set; }


        [TestMethod]
        public void GenerateGameUrl_LobbyPage_Test()
        {
            var result = GameLinkService.GenerateTransferModeLaunchLink(new GameLaunchVendorApiTransferModeRequestModel
            {
                GameID = 0,
                MerchantID = _TEST_DATA.Default.MerchantID,
                PlayerName = _TEST_DATA.Default.PlayerName,
                Lang = SystemConstant.DefaultLang,
                MerchantSecretKey = _TEST_DATA.Default.MerchantSecret,
            });

            Assert.IsTrue(result.Link.IsNotNullEmptyText());

            CheckTokenIsValidAndParsingSuccess(result);
            return;
        }

        private void CheckTokenIsValidAndParsingSuccess(GameLinkGenerateResultViewModel result)
        {
            var tokenValue = result.Link.GetQueryStringValueFromUrlFormat("token");

            //這是跟前端約定好的模式
            tokenValue = tokenValue.Base64Decode();

            var tokenParsing = IPlayerJwtCheckService.ParsingLobbyToken(tokenValue);
        }


        [TestMethod]
        public void GenerateGameUrlForVendorApi_Test()
        {
            string fakeUrl = "https://example.invalid";

            //使用mock
            MockGameVendorService(fakeUrl);

            var result = GameLinkService.GenerateTransferModeLaunchLink(new GameLaunchVendorApiTransferModeRequestModel
            {
                GameID = _TEST_DATA.ExternalGameVendor.GameID,
                MerchantID = _TEST_DATA.Default.MerchantID,
                PlayerName = _TEST_DATA.Default.PlayerName,
                Lang = SystemConstant.DefaultLang,
                MerchantSecretKey = _TEST_DATA.Default.MerchantSecret,
            });

            Assert.IsTrue(fakeUrl == result.Link);
            return;
        }


        [TestMethod]
        public void GenerateGameUrlForVendorApi_SingleWalletMerchant_Test()
        {
            string fakeUrl = "https://example.invalid";

            //使用mock
            MockGameVendorService(fakeUrl);

            var result = GameLinkService.GenerateSingleWalletLaunchLink(new GameLaunchVendorApiSingleWalletModeRequestModel
            {
                GameID = _TEST_DATA.ExternalGameVendor.GameID,
                MerchantID = _TEST_DATA.SingleWallet.MerchantID,
                MerchantSecretKey = _TEST_DATA.SingleWallet.MerchantSecret,
                PlayerUniqueID = _TEST_DATA.SingleWallet.PlayerName,
                SessionToken = [REDACTED],
                Lang = SystemConstant.DefaultLang,
            });

            Assert.IsTrue(fakeUrl == result.Link);
            return;
        }


        private void MockGameVendorService(string fakeUrl)
        {
            var substitute = Substitute.For<IGameVendorIntegrationService>();
            substitute.GetGameLaunchUrl(null).ReturnsForAnyArgs(new ExternalVendorLaunchUrlViewModel
            {
                Url = fakeUrl,
            });

            IOCHelper.AddMockInstance<IGameVendorIntegrationService>(substitute);
        }

    }
}

