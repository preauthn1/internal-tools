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

namespace MSC.Testing
{
    [TestClass]
    public class UserOperateLogServiceTests : BaseTest
    {

        public IUserOperateLogService ServiceInstance { get; set; }


        [TestMethod]
        public void UserOperateLogService_GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID , new QueryUserOperateLogApiRequestParameterModel
            {
                PageSize = 10
            });
            Assert_NotEmptyCollection(query.Items);
            return;
        }

    }
}

