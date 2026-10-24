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
    public class FakeDataGenerateServiceTests : BaseTest
    {

        public IFakeDataGenerateService ServiceInstance { get; set; }


        [TestMethod]
        public void GenerateJwtTokenForQATesting_Test()
        {
            var result = ServiceInstance.GenerateJwtTokenForQATesting(new FakeDataGenerateRequestModel
            {
                PlayerName = _TEST_DATA.Default.PlayerName
            });

            Assert.IsNotNull(result);
            Assert.IsTrue(!string.IsNullOrWhiteSpace(example.invalid.Access_Token));
            return;
        }



        [TestMethod]
        public void AddPlayerWallet_Test()
        {
            var result = ServiceInstance.AddPlayerWallet(_TEST_DATA.Default.PlayerID);

            Assert.IsNotNull(result);
            Assert.IsTrue(result.BeforeAmount > 0);
            Assert.IsTrue(result.AfterAmount > 0);
            return;
        }


    }
}

