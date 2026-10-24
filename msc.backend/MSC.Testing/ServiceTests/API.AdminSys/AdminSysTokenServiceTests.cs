using Microsoft.VisualStudio.TestTools.UnitTesting;
using MSC.BL.Interface;
using MSC.DA.Repository;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.Testing
{
    [TestClass]
    public class AdminSysTokenServiceTests : BaseTest
    {

        public IAdminSysTokenService ServiceInstance { get; set; }


        [TestMethod]
        public void AdminSysTokenService_CheckPermissionExist_Test()
        {
            ServiceInstance.CheckPermissionExist(_TEST_DATA.Default.UserID, new List<string>
            {
                AdminPermissionTokenEnum.Card_Management.ToString(),
                AdminPermissionTokenEnum.System_ExternalLink_List.ToString(),
            });

            return;
        }

    }
}

