using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.Models.Models;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class PlayerWalletInfoServiceTests : BaseTest
    {

        public IPlayerWalletInfoService PlayerWalletInfoService { get; set; }



        [TestMethod]
        public void GetWalletCash_Test()
        {
            var query = PlayerWalletInfoService.GetWalletAmount(_TEST_DATA.Default.PlayerID , "");

            Assert.IsNotNull(query);
            Assert.IsNotNull(query.Amount > 10);
            return;
        }



        [TestMethod]
        public void GetPlayerInfo_Test()
        {
            var query = PlayerWalletInfoService.GetPlayerInfo(_TEST_DATA.Default.PlayerID);

            Assert.IsNotNull(query);
            return;
        }



    }
}
