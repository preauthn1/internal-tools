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
    public class VendorConfigServiceTests : BaseTest
    {

        public IVendorConfigService ServiceInstance { get; set; }


        [TestMethod]
        public void GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID , new QueryVendorConfigApiRequestParameterModel
            {
                PageSize = 10
            });


            if(query.Items.IsNotEmptyNull())
            {
                var first = query.Items.FirstOrDefault();

                //順手測試 另個案例2
                GetVendorConfigSettinsAndAssert(first.VendorConfigID);

                //順手測試 另個案例3
                ServiceInstance.SyncDefaultVendorConfigMapping(first.VendorConfigID);
            }


            Assert_NotEmptyCollection(query.Items);
            Assert.IsNotNull(new CRUDResultModel());
            return;
        }


        public void GetVendorConfigSettinsAndAssert(int vendorConfigID)
        {
            var query = ServiceInstance.GetVendorConfigSettins(vendorConfigID);
            Assert_NotEmptyCollection(query.Groups);
            return;
        }


    }
}

