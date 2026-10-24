using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Models;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.Web.Modules.Filter;
using System.Collections;
using System.Collections.Generic;
using System.Runtime.Serialization;

namespace target.Testing
{
    [TestClass]
    public class ExceptionHandlerTests : BaseTest
    {

        [TestMethod]
        public void ErrorHandler_IsSkipSlackNotifyExceptionType_Test()
        {
            var sqlEx1 = new SqlExceptionBuilder().WithErrorNumber(2601).WithErrorMessage("duplicate PK").Build();
            var sqlEx2 = new SqlExceptionBuilder().WithErrorNumber(-2).WithErrorMessage("timeout").Build();
            var sqlEx3 = new SqlExceptionBuilder().WithErrorNumber(777).WithErrorMessage("test").Build();


            //預期要跳過的exception
            bool isSkip_1 = ExceptionHandler.IsSkipSlackNotifyExceptionType(sqlEx1, "/game/slot/bonus");
            bool isSkip_2 = ExceptionHandler.IsSkipSlackNotifyExceptionType(sqlEx2, "game/slot/bonus");
            bool isSkip_3 = ExceptionHandler.IsSkipSlackNotifyExceptionType(sqlEx2, "game/slot/bonusXXXXX");

            Assert.IsTrue(isSkip_1);
            Assert.IsTrue(isSkip_2);
            Assert.IsTrue(isSkip_3);


            //預期 [不應該] 跳過的exception情境
            bool isNotSkip_1 = ExceptionHandler.IsSkipSlackNotifyExceptionType(sqlEx3, "game/slot/bonus");
            bool isNotSkip_2 = ExceptionHandler.IsSkipSlackNotifyExceptionType(sqlEx1, "xxxxxxxx");
            bool isNotSkip_3 = ExceptionHandler.IsSkipSlackNotifyExceptionType(sqlEx2, "xxxxxxxx");

            Assert.IsFalse(isNotSkip_1);
            Assert.IsFalse(isNotSkip_2);
            Assert.IsFalse(isNotSkip_3);

            return;
        }


        [TestMethod]
        public void ErrorHandler_GetNotifySlackChannel_Test()
        {
            var channel = ExceptionHandler.GetNotifySlackChannel();
            return;
        }



    }



}

