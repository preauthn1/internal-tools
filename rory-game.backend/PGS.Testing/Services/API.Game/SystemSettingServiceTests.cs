using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.Models.Models;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Linq;
using target.Models.CoreLibs;
using target.Models.ViewModels;

namespace target.Testing
{
    [TestClass]
    public class SystemSettingServiceTests : BaseTest
    {

        public ISystemSettingService SystemSettingService { get; set; }


        [TestMethod]
        public void GetDenyIpAddressList_Test()
        {
            var query = SystemSettingService.GetDenyIpAddressList();
            Assert_NotEmptyCollection(query);
            return;
        }

    }
}

