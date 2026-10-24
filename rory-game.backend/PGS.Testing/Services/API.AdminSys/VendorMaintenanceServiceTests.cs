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
    public class VendorMaintenanceServiceTests : BaseTest
    {

        public IVendorMaintenanceService ServiceInstance { get; set; }


        [TestMethod]
        public void GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID , new QueryVendorMaintenanceApiRequestParameterModel
            {
            });

            Assert_NotEmptyCollection(query.Items);
            return;
        }



        [TestMethod]
        public void Create_Test()
        {
            var query = ServiceInstance.Create(_TEST_DATA.Default.UserID, new VendorMaintenanceCreateViewModel
            {
                VendorID = GameVendorTypeEnum.Funky.ToInteger(),
                StartTime = new DateTime(2020 , 1 ,1),
                EndTime = new DateTime(2020, 1, 2),
            });

            return;
        }


    }
}

