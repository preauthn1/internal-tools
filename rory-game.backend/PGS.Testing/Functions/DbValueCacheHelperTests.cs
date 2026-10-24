using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Services;
using target.Models.Enums;
using target.Models.Models;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class DbValueCacheHelperTests : BaseTest
    {

        [TestMethod]
        public void CacheHelper_GetGameNameFromLanguageText()
        {
            var q1 = DbValueCacheHelper.GetGameNameFromLanguageText(_TEST_DATA.Default.GameID);
            var q2 = DbValueCacheHelper.GetGameNameFromLanguageText(_TEST_DATA.Default.ThirdPartyGameID);
            return;

        }


        [TestMethod]
        public void CacheHelper_GetGame_Test()
        {
            var q1 = DbValueCacheHelper.GetGame(_TEST_DATA.Default.GameID);
            var q2 = DbValueCacheHelper.GetGame(_TEST_DATA.Default.ThirdPartyGameID);

            Assert.IsNotNull(q1);
            Assert.IsNotNull(q2);
            return;

        }


    }
}
