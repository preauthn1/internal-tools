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
    public class ServerRobotConfigServiceTests : BaseTest
    {

        public IServerRobotConfigService ServiceInstance { get; set; }


        [TestMethod]
        public void ServerRobotConfigService_GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(new QueryServerRobotConfigApiRequestParameterModel
            {

            });

            Assert_NotEmptyCollection(query.Items);
            return;
        }

    }
}

