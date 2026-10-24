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
    public class GameLaunchLobbyServiceTests : BaseTest
    {
        public LobbyTokenLogRepository LobbyTokenLogRepository { get; set; }

        public IInternalUtilityService InternalUtilityService { get; set; }


        [TestMethod]
        public void LaunchGameAndCheckPlayerExist_Test()
        {
            string uuid = Guid.NewGuid().ToString();

            var fakeJwt = MyJwtHelper.GeneratePlayerJwt(
                _TEST_DATA.Default.PlayerID, 
                _TEST_DATA.Default.MerchantID,
                _TEST_DATA.Default.PlayerName, 
                "zh-tw",
                uuid,
                uuid);

            MockHttpHeaderJwtToken(fakeJwt.JwtToken);

            Guid hashGuid = Guid.NewGuid();

            LobbyTokenLogRepository.Create(new Models.Table.LobbyTokenLogEntity
            {
                CreateTime = DateTimeHelper.GetCurrentTime(),
                ExpiredTime = DateTimeHelper.GetCurrentTime().AddHours(8),
                IdentityGuid = hashGuid,
                IdentityHash = fakeJwt.JwtToken.MD5Hash(),
                PlayerID = _TEST_DATA.Default.PlayerID,
                MerchantID = _TEST_DATA.Default.MerchantID,
                RawJwtToken = [REDACTED]
            });

            var result = InternalUtilityService.GetLobbyGames(new LobbyApiGetGamesRequestModel
            {
                SignKey = SystemConstant.PublicAPI_SignKey_MagicPassword,
                Ts = DateTimeHelper.GetCurrentTimeTimestamp(),
                PlayerName = _TEST_DATA.Default.PlayerName
            });

            Assert_NotEmptyCollection(result.Games);
            return;
        }


    }
}

