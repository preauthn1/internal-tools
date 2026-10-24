using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Services;
using target.Models.Enums;
using target.Models.Models;
using target.Utility.Extensions;
using target.Utility.Helper;
using System;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class RegexRuleHelperTests : BaseTest
    {

        [TestMethod]
        public void CheckIsAccountForamt()
        {
            //已確認要pass的案例
            Assert.IsTrue(IsAccountFormat("12345"));
            Assert.IsTrue(IsAccountFormat("abcd"));
            Assert.IsTrue(IsAccountFormat("AAAAAAAAABBBBBBBBB"));
            Assert.IsTrue(IsAccountFormat("user@example.invalid"));
            Assert.IsTrue(IsAccountFormat("user@example.invalid"));
            Assert.IsTrue(IsAccountFormat("user@example.invalid"));
            Assert.IsTrue(IsAccountFormat("user@example.invalid"));
            Assert.IsTrue(IsAccountFormat("user@example.invalid"));


            //已確認要fail的案例
            Assert.IsFalse(IsAccountFormat("1"));
            Assert.IsFalse(IsAccountFormat("にっぽんんごんごご"));
            Assert.IsFalse(IsAccountFormat("吃吃飯吃飯飯"));
            Assert.IsFalse(IsAccountFormat("~~~~~~!!!"));


            var match13 = IsAccountFormat("123abc?&*@@");
            var match63 = IsAccountFormat("hello-world");
            var match68 = IsAccountFormat("???????$#########");
            var match79 = IsAccountFormat("");
            return;
        }


        [TestMethod]
        public void CheckIsIntegerValue()
        {
            //已確認要pass的案例
            Assert.IsTrue((1.0M).IsIntegerValue());
            Assert.IsTrue((3500.0M).IsIntegerValue());
            Assert.IsTrue((173511.0M).IsIntegerValue());

            Assert.IsFalse((7.00008M).IsIntegerValue());
            Assert.IsFalse((3.60600M).IsIntegerValue());
            return;
        }



        [TestMethod]
        public void RuleFor_DecimalTwoPlaces()
        {
            RuleIsDecimalPointPlacesTwo(1);
            RuleIsDecimalPointPlacesTwo(5000);
            RuleIsDecimalPointPlacesTwo(5000.22M);
            RuleIsDecimalPointPlacesTwo(1.2M);
            RuleIsDecimalPointPlacesTwo(1.22M);

            AssertExceptionOccur(1.233M);
            AssertExceptionOccur(1.234567M);
            AssertExceptionOccur(5000.123M);
            return;

        }


        private void RuleIsDecimalPointPlacesTwo(decimal value)
        {
            var model = new
            {
                Amount = value
            };

            model.RuleFor(z => z.Amount).Rule_IsDecimalPointPlacesTwo();
        }

        private void AssertExceptionOccur(decimal value)
        {
            try
            {
                var model = new
                {
                    Amount = value
                };

                model.RuleFor(z => z.Amount).Rule_IsDecimalPointPlacesTwo();
            }
            catch(TerproException ex)
            {
                Assert.IsTrue(ex.ErrroCode == ApiErrorCode.InvalidParameter_NumberValueMustPlacesTwo);
                return;
            }

            Assert.Fail("unexpect !!");
        }


        private bool IsAccountFormat(string myValue)
        {
            return RegexRuleHelper.IsAccountFormat(myValue);
        }


    }
}
