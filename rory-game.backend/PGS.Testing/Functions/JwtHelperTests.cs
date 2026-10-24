using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.Models.Models;
using target.Utility;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class JwtHelperTests : BaseTest
    {
        [TestMethod]
        public void JwtTokenEncrypt()
        {
            var myObject = new
            {
                PlayerName = "yyyyMMdd"
            };

            string jwtText = BaseJWTTokenHelper.Generate_JWT_Token(myObject);

            var parseJwt = BaseJWTTokenHelper.GetObjectFromJWT<dynamic>(jwtText);
            
            Assert.IsNotNull(parseJwt);
            Assert.IsTrue(parseJwt.PlayerName == myObject.PlayerName);
            return;

        }
    }
}
