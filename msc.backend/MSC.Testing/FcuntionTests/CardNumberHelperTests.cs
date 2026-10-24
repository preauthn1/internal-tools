using Microsoft.VisualStudio.TestTools.UnitTesting;
using MSC.Utility;

namespace MSC.Testing
{
    [TestClass]
    public class CardNumberHelperTests
    {
        [TestMethod]
        public void GenerateCardNumber()
        {
            var card = CardNumberHelper.GenerateCardNumber();
            return;
        }
    }
}
