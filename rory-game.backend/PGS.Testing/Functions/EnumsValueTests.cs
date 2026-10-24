using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.Models.Models;
using target.Utility.Extensions;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class EnumsValueTests : BaseTest
    {

        /// <summary>
        /// ˬdErrroCodesO_
        /// </summary>
        [TestMethod]
        public void CheckDuplicateEnumValue()
        {
            var allEnums = ExtensionOfEnum.GetEnumAllValue<ApiErrorCode>();

            //ˬdO_O6Ʀrs
            foreach(ApiErrorCode item in allEnums)
            {
                //SwcodesS ˬd
                if(item == ApiErrorCode.Success || item == ApiErrorCode.SuccessOK)
                {
                    continue;
                }

                int myValue = (int)item;

                if(myValue < 400000)
                {
                    Assert.Fail($"{myValue} ƦrLpAswq֤F@ơH");
                }
                if (myValue > 900000)
                {
                    Assert.Fail($"{myValue} ƦrLjAswqhF@ơH");
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

            //ˬdsO_Ʃwq
            if(duplocateItems.Any(x=> x.Count > 1))
            {
                var matchKey = duplocateItems.FirstOrDefault().Key;
                Assert.Fail($"s {matchKey} ƤFC");
            }

        }
    }
}
