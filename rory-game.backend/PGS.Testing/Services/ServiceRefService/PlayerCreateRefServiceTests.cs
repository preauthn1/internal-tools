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
using target.Models;

namespace target.Testing
{
    [TestClass]
    public class PlayerCreateRefServiceTests : BaseTest
    {

        public IPlayerCreateRefService ServiceInstance { get; set; }


        [TestMethod]
        public void CreatePlayerOrReturnPlayerID_Test()
        {  
            //這隻測試每天第一次跑 是建立帳號、後面幾次是返回ID
            string dateText = DateTimeHelper.GetCurrentTime(8).ToString("yyyyMMdd");
            string new_player_name = $"Unitest_{dateText}";

            var new_player_id = ServiceInstance.CreatePlayerOrReturnPlayerID(new PlayerCreateParameterModel
            {
                Currency = SystemConstant.DefaultCurrency,
                MerchantID = _TEST_DATA.Default.MerchantID,
                PlayerName = new_player_name
            });

            Assert.IsTrue(new_player_id > 0);
            return;
        }

    }
}

