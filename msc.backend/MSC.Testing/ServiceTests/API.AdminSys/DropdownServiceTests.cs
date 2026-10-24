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
    public class DropdownServiceTests : BaseTest
    {

        public IDropdownService ServiceInstance { get; set; }


        [TestMethod]
        public void DropdownService_GetMerchants_Test()
        {
            var query = ServiceInstance.GetMerchants(_TEST_DATA.Default.UserID);
            Assert_NotEmptyCollection(query);
            return;
        }


        [TestMethod]
        public void DropdownService_GetTransactionTypes_Test()
        {
            var query = ServiceInstance.GetTransactionTypes();
            Assert_NotEmptyCollection(query);
            return;
        }



        [TestMethod]
        public void DropdownService_GetCurrency_Test()
        {
            var query = ServiceInstance.GetCurrency();
            Assert_NotEmptyCollection(query);
            return;
        }


    }
}

