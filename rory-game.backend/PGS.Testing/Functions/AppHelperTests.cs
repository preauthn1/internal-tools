using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.Models.Enums;
using target.Models.Models;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class AppHelperTests : BaseTest
    {

        [TestMethod]
        public void GetMerchantSettingOrDefaultValue_Test()
        {
            var defaultValue1 = SettingKeyEnumLogicHelper.GetMerchantSettingOrDefaultValue(MerchantCustomizedSettingKeyEnum.DefaultWalletBalance, null);
            var defaultValue2 = SettingKeyEnumLogicHelper.GetMerchantSettingOrDefaultValue(MerchantCustomizedSettingKeyEnum.Style_LoadingStyleMode, null);

            Assert.IsNotNull(defaultValue1);
            Assert.IsNotNull(defaultValue2);

            return;

        }
        
    }
}
