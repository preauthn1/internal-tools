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
    public class UserMerchantMappingServiceTests : BaseTest
    {
        public MerchantRepository MerchantRepository { get; set; }

        public IUserMerchantMappingService ServiceInstance { get; set; }


        [TestMethod]
        public void UserMerchantMappingService_GetUserMerchantsEnableFlag_Test()
        {
            var query = ServiceInstance.GetUserMerchantsEnableFlag(_TEST_DATA.Default.UserID);
            Assert_NotEmptyCollection(query.Items);
            return;
        }


        [TestMethod]
        public void UserMerchantMappingService_UpdateUserMerchants_Test()
        {
            var queryAllMerchants = MerchantRepository.GetAll();

            ServiceInstance.UpdateUserMerchants(_TEST_DATA.Default.UserID , new UserMerchantMappingRequestModel
            {
                UserID = _TEST_DATA.Default.UserID,
                MerchantIDs = queryAllMerchants.Select(z=> z.MerchantID).ToList()
            });

            return;
        }

    }
}

