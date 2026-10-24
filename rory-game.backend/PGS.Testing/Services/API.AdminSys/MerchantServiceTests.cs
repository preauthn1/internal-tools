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
using target.Utility;

namespace target.Testing
{
    [TestClass]
    public class MerchantServiceTests : BaseTest
    {

        public IMerchantService MerchantService { get; set; }

        public IMerchantCustomizedSettingService MerchantCustomizedSettingService { get; set; }


        [TestMethod]
        public void MerchantGetByPagingByFilter_Test()
        {
            var query = MerchantService.GetByPagingFilter(new QueryMerchantApiRequestParameterModel { PageSize = 10 });
            Assert_NotEmptyCollection(query.Items);
            return;
        }


        [TestMethod]
        public void MerchantGetByPagingByFilte_FilterOrgID_Test()
        {
            var filter = new QueryMerchantApiRequestParameterModel 
            {
                PageSize = 10, 
                OrganizationIDs = _TEST_DATA.Default.OrganizationID.ToListCollection()
            };

            var query = MerchantService.GetByPagingFilter(filter);
            Assert_NotEmptyCollection(query.Items);
            return;
        }



        [TestMethod]
        public void Update_Test()
        {
            var queryExist = MerchantService.GetByID(_TEST_DATA.Not_Super_User.MerchantID);

            Assert.IsNotNull(queryExist);

            var updateModel = queryExist.ReturnInjecteObject<MerchantUpdateRequestModel>();
            MerchantService.Update(updateModel);
            return;
        }




        [TestMethod]
        public void GetMerchantAdvanceSettings_Test()
        {
            var query = MerchantCustomizedSettingService.GetMerchantAdvanceSettings(_TEST_DATA.Default.UserID , _TEST_DATA.Not_Super_User.MerchantID);

            Assert.IsNotNull(query);
            Assert_NotEmptyCollection(query.Groups);
            return;
        }




    }
}

