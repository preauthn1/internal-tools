using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Models;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.Web.Modules.Filter;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class AdminSysAuthorizeHelperTests : BaseTest
    {

        [TestMethod]
        public void GetValidUserSessionByJwt_Test()
        {
            var jwtModel = AdminSysAuthorizeHelper.GetValidUserSessionByJwt();
            Assert.IsNotNull(jwtModel);

            return;
        }

    }
}

