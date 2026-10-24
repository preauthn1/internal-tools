using Microsoft.VisualStudio.TestTools.UnitTesting;
using MSC.Utility;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using MSC.Utility.Utility;
using System.Linq;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;

namespace MSC.Testing
{
    [TestClass]
    public class OthersTests : BaseTest
    {
        private ConcurrentDictionary<long, object> dictionary = new ConcurrentDictionary<long, object>();

        private List<string> ListItems = new List<string>();


        [TestMethod]
        public void TestMethod1()
        {
            //uaq}o~|o TestMethod
            if (IsNotDeveloperComputerName())
                return;




            var list = new List<int>();

            var qqqqqq2 = list.Sum();
            var qqqqqq3 = list.DefaultIfEmpty().Average();




            return;
        }


        public static string GenerateVendorApiIntegrationSignKey(object model) 
        {
            List<string> keyValues = new List<string>();

            foreach (var item in model.GetType().GetProperties().OrderBy(x => x.Name))
            {
                var name = item.Name.ToLower();
                var value = item.GetValue(model);

                if (name == "pay_md5sign")
                    continue;

                keyValues.Add($"{name}={value}");
            }

            string beforeMD5 = string.Join("&", keyValues);

            beforeMD5 += $"&key=fi073jlzd524eqxt7wo5e2mp5shleul2";

            string MD5Result = beforeMD5.MD5Hash().ToUpper();
            return MD5Result;
        }

    }
}
