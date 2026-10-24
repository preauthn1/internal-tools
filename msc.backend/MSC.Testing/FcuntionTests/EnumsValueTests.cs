using Microsoft.VisualStudio.TestTools.UnitTesting;
using MSC.Models.Models;
using MSC.Utility.Extensions;
using System.Linq;

namespace MSC.Testing
{
    [TestClass]
    public class EnumsValueTests 
    {

        /// <summary>
        /// 檢查ErrroCode編號是否重複
        /// </summary>
        [TestMethod]
        public void CheckDuplicateEnumValue_Test()
        {
            var allEnums = ExtensionOfEnum.GetEnumAllValue<ErrorCode>();

            //檢查是否是6位數字編號
            foreach(ErrorCode item in allEnums)
            {
                int myValue = (int)item;

                if(myValue < 400000 && myValue > 20000)
                {
                    Assert.Fail($"{myValue} 數字過小，編號定義少了一位數？");
                }
                if (myValue > 900000)
                {
                    Assert.Fail($"{myValue} 數字過大，編號定義多了一位數？");
                }
            }

            var duplocateItems = allEnums
                    .GroupBy(x => (int)x)
                    .Where(x => x.Count() > 1)
                    .Select(x => new
                    {
                        Key = x.Key,
                        Count = x.Count()
                    });

            //檢查編號是否有重複定義
            if(duplocateItems.Any(x=> x.Count > 1))
            {
                var matchKey = duplocateItems.FirstOrDefault().Key;
                Assert.Fail($"編號 {matchKey} 重複了。");
            }

        }
    }
}
