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
    public class PerGameDetailServiceTests : BaseTest
    {

        public IPerGameDetailService ServiceInstance { get; set; }


        [TestMethod]
        public void PerGameDetail_GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID , new QueryPerGameDetailLogApiRequestParameterModel { PageSize = 1});
            Assert_NotEmptyCollection(query.Items);


            var matchFirst = query.Items.FirstOrDefault();

            var queryFull = ServiceInstance.GetFullInformationInPerGame(_TEST_DATA.Default.UserID, matchFirst.PerGameLogID);

            Assert.IsNotNull(queryFull);
            Assert_NotEmptyCollection(queryFull.PlayerDetails);
            return;
        }



        [TestMethod]
        public void GetSlotHistoryFromPerGameLog_Test()
        {
            var query = ServiceInstance.GetSlotHistoryFromPerGameLog(_TEST_DATA.Default.UserID, 3734835);

            Assert.IsNotNull(query);
            return;
        }

    }
}

