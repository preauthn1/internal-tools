using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.Models.Enums;
using target.Models.Models;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.Utility.Utility;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class ApiRequestHelperTests : BaseTest
    {

        [TestMethod]
        public void ApiRequestHelper_RequestApiResult_Test()
        {
            var apiResult = APIRequestHelper.RequestApiResult<string>(new RequestParameters
            {
                BaseUrl = "https://example.invalid/",
                ApiUrl = APIRequestHelper.No_Url_Keyword,
                Method = HttpMethodType.HttpGet
            });

            if(apiResult.Exception != null)
            {
                Assert.Fail("Http請求壞掉了" + apiResult.Exception.Message);
            }

            return;

        }
		
        
    }
}
