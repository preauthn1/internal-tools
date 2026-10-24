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

namespace target.Testing
{
    [TestClass]
    public class MerchantGameSettingRepositoryTests : BaseTest
    {

        public GameListRepository GameListRepository { get; set; }

        public MerchantGameSettingRepository MerchantGameSettingRepository { get; set; }


        [TestMethod]
        public void IsMerchantGameEnableInSystem_MustFalse_Test()
        {
            var mustDisable = MerchantGameSettingRepository.IsMerchantGameEnableInSystem(_TEST_DATA.Default.MerchantID , 6553621);
            Assert.IsFalse(mustDisable);
            
            var queryDisabledGames = GameListRepository.GetByFilter(new QueryGameListRequestViewModel
            {
                Status = NormalActiveStatus.Disable.ToInteger()
            });
            
            if(queryDisabledGames.IsNotEmptyNull())
            {
                var mustDisable2 = MerchantGameSettingRepository.IsMerchantGameEnableInSystem(_TEST_DATA.Default.MerchantID, queryDisabledGames.FirstOrDefault().GameID);
                Assert.IsFalse(mustDisable2);
            }

            return;
        }



        [TestMethod]
        public void IsMerchantGameEnableInSystem_MustExist_Test()
        {
            var queryDisabledGames = GameListRepository.GetByFilter(new QueryGameListRequestViewModel
            {
                Status = NormalActiveStatus.Enable.ToInteger()
            });

            var queryGameSetting = MerchantGameSettingRepository.GetByFilter(new QueryMerchantGameSettingRequestViewModel
            {
                Status = NormalActiveStatus.Enable.ToInteger(),
                GameIDs = queryDisabledGames.Select(x=> x.GameID).ToList(),
            });

            
            if (queryGameSetting.IsNotEmptyNull())
            {
                var matchFirst = queryGameSetting.FirstOrDefault();
                var mustEnable = MerchantGameSettingRepository.IsMerchantGameEnableInSystem(matchFirst.MerchantID, matchFirst.GameID);
                Assert.IsTrue(mustEnable);
            }

            return;
        }




    }
}

