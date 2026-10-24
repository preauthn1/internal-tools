using Lck.Utility.Extensions.Cache;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.Models.Enums;
using target.Models.Models;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class CacheHelperTests : BaseTest
    {

        [TestMethod]
        public void RedisCache_Test()
        {
            //非公司開發電腦 不執行這組test
            if (IsNotDeveloperComputerName())
                return;

            var redis = new RedisCache();

            string key = "MS_Test";
            redis.Set(key, "MS_Test_66666", new System.TimeSpan(24, 0, 0));

            var qqq = redis.Get<string>(key);
            return;

        }
        
    }
}
