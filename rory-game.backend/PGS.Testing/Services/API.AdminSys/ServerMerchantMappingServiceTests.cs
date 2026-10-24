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
    public class ServerMerchantMappingServiceTests : BaseTest
    {

        public IServerMerchantMappingService ServiceInstance { get; set; }


        [TestMethod]
        public void GetServerMerchantMappings_Test()
        {
            var query = ServiceInstance.GetServerMerchantMappings(_TEST_DATA.Default.ServerID);
            Assert_NotEmptyCollection(query.AvailableSettings);
            return;
        }


        [TestMethod]
        public void GetAllServerMerchantMappings_Test()
        {
            var query = ServiceInstance.GetAllServerMerchantMappings();

            Assert_NotEmptyCollection(query.CurrencyGroups);
            Assert_NotEmptyCollection(query.CurrencyGroups.SelectMany(x => x.Merchants));
            Assert_NotEmptyCollection(query.CurrencyGroups.SelectMany(x=> x.Servers));
            return;
        }


        [TestMethod]
        public void ReversMerchantMappings_Test()
        {
            ServiceInstance.ReversMerchantMappings(new ReversServerMerchantRequestViewModel
            {
                ServerID = _TEST_DATA.Default.ServerID,
                MerchantID = _TEST_DATA.Default.MerchantID
            });

            return;
        }



    }
}

