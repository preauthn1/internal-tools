using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.Models.Models;
using target.Models.ViewModels.API.AdminSys;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Collections.Generic;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class ExtensionTests : BaseTest
    {


        [TestMethod]
        public void IsDictionary_Test()
        {

            Dictionary<string, object> dic1 = new Dictionary<string, object>();
            Dictionary<string, string> dic2 = new Dictionary<string, string>();
            Dictionary<string, int> dic3 = new Dictionary<string, int>();

            Assert.IsTrue(dic1.IsDictionary());
            Assert.IsTrue(dic2.IsDictionary());
            Assert.IsTrue(dic3.IsDictionary());


            Assert.IsFalse((new { qqq = "" }).IsDictionary());
            Assert.IsFalse((new DropdownValueViewModel()).IsDictionary());

            return;
        }


        [TestMethod]
        public void CutLastCharIfMatch_Test()
        {
            string test1 = "abcd_/";
            string test2 = "abcd_@";
            string test3 = null;

            var rrrr1 = test1.CutLastCharIfMatch('/');
            var rrrr2 = test2.CutLastCharIfMatch('/');
            var rrrr3 = test3.CutLastCharIfMatch('/');


            Assert.IsTrue((test1.Length - rrrr1.Length) == 1);
            Assert.IsTrue(test2 == rrrr2);
            Assert.IsTrue(test3 == rrrr3);

            return;
        }




        [TestMethod]
        public void GetInnerTextBetweenTwoWord_Test()
        {
            string text = "hello (this is my)";

            var rrrr1 = text.GetInnerTextBetweenTwoWord("(" , ")");
            var rrrr2 = text.GetInnerTextBetweenTwoWord("[", "]");

            return;
        }


        [TestMethod]
        public void DateTimeTest()
        {
            var date = DateTimeHelper.GetCurrentTime();

            var r1 = date.ToISO8601();
            return;
        }


    }
}
