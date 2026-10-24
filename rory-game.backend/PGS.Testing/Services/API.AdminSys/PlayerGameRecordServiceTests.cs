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
    public class PlayerGameRecordServiceTests : BaseTest
    {

        public IPlayerGameRecordService ServiceInstance { get; set; }


        [TestMethod]
        public void PlayerGameRecord_GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID , new QueryPlayerGameRecordApiRequestParameterModel
            {
                PageSize = 20
            })
                ;
            Assert_NotEmptyCollection(query.Items);
            return;
        }


        [TestMethod]
        public void PlayerGameRecord_GetByPagingFilter_FilterColumns_Test()
        {
            var now = DateTimeHelper.GetCurrentTime();
            var previous_month = now.AddMonths(-2);

            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID, new QueryPlayerGameRecordApiRequestParameterModel
            {
                PageSize = 20,
                LoginTimeBegin = previous_month,
                LoginTimeEnd = now,
                GameID = GameListEnum.DFDC_多福多财.ToInteger()
            });

            Assert_NotEmptyCollection(query.Items);
            return;
        }



        [TestMethod]
        public void PlayerGameRecord_GetByPagingFilter_MustEmpty_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID, new QueryPlayerGameRecordApiRequestParameterModel
            {
                PageSize = 20,
                PlayerName = "ttttttttxxxxxxxxxxxxx",
            });

            Assert.IsTrue(query.Items.IsEmptyOrNull());
            return;
        }




    }
}

