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
    public class PlayerLoginLogServiceTests : BaseTest
    {

        public IPlayerLoginLogService ServiceInstance { get; set; }


        [TestMethod]
        public void PlayerLoginLogService_GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID , new QueryPlayerLoginLogApiRequestParameterModel
            {
                PageSize = 20,
            });

            Assert_NotEmptyCollection(query.Items);
            return;
        }



        [TestMethod]
        public void PlayerLoginLogService_GetByPagingFilter_FilterDateTime_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID, new QueryPlayerLoginLogApiRequestParameterModel
            {
                CreateTimeBegin = DateTimeHelper.GetCurrentTime().AddDays(-90),
                CreateTimeEnd = DateTimeHelper.GetCurrentTime(),
                PageSize = 20,
            });

            Assert_NotEmptyCollection(query.Items);
            return;
        }


    }
}

