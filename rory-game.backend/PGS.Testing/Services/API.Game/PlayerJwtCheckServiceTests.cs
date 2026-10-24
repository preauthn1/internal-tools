using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Models;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class PlayerJwtCheckServiceTests : BaseTest
    {

        public IPlayerJwtCheckService ServiceInstance { get; set; }


        public LobbyTokenLogRepository LobbyTokenLogRepository { get; set; }


        [TestMethod]
        public void CheckIsLobbyTokenExist_Test()
        {
            Guid guid = Guid.NewGuid();
            string guidTokenText = guid.ToString();

            LobbyTokenLogRepository.Create(new Models.Table.LobbyTokenLogEntity
            {
                PlayerID = _TEST_DATA.Default.PlayerID,
                MerchantID = _TEST_DATA.Default.MerchantID,
                CreateTime = DateTimeHelper.GetCurrentTime(),
                ExpiredTime = DateTimeHelper.GetCurrentTime().AddDays(7),
                IdentityHash = guidTokenText,
                IdentityGuid = guid
            });

            ServiceInstance.CheckIsLobbyTokenExist(guidTokenText);
            return;
        }

        [TestMethod]
        public void CheckIsLobbyTokenExist_Fail_Test()
        {
            try
            {
                ServiceInstance.CheckIsLobbyTokenExist("XXXXXXXXXXXXXX");
            }
            catch(Exception ex)
            {

                AssertExpectException(ex , new List<ApiErrorCode>
                {
                    ApiErrorCode.Authorize_AccessToken_NotValid ,
                    ApiErrorCode.Authorize_AccessToken_Expired
                });
            }

            return;
        }

    }
}

