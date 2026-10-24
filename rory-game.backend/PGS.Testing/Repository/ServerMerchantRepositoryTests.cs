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
    public class ServerMerchantRepositoryTests : BaseTest
    {

        public ServerMerchantMappingRepository ServerMerchantMappingRepository { get; set; }


        [TestMethod]
        public void IsServerMerchantEnable_Test()
        {
            var queryMine = ServerMerchantMappingRepository.GetByFilter(new QueryServerMerchantRelationRequestViewModel
            {
                MerchantIDs = _TEST_DATA.Default.MerchantID.ToListCollection()
            });

            if(queryMine.IsNotEmptyNull())
            {
                var mustMatch = ServerMerchantMappingRepository.IsServerMerchantEnable(queryMine.FirstOrDefault().ServerID, _TEST_DATA.Default.MerchantID);
                Assert.IsTrue(mustMatch);
            }

            var mustNotFound = ServerMerchantMappingRepository.IsServerMerchantEnable(666666, _TEST_DATA.Default.MerchantID);
            Assert.IsFalse(mustNotFound);
            return;
        }


    }
}

