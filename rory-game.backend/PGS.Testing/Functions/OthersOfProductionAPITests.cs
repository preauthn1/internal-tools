using Lck.Utility.Extensions.Cache;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.DA.Repository;
using target.Models;
using target.Models.Enums;
using target.Models.Models;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.Utility.IOC;
using target.Utility.Utility;
using System;
using System.Collections.Generic;
using System.Data.SqlClient;
using System.Diagnostics;
using System.Linq;
using System.Security.Cryptography;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading;
using System.Threading.Tasks;


namespace target.Testing
{
    [TestClass]
    public class OthersOfProductionAPITests : BaseTest
    {

        /// <summary>
        /// 各種雜項測試、驗證
        /// </summary>
        [TestMethod]
        public void SampleTest()
        {
            //非公司開發電腦 不執行這組test
            if (IsNotDeveloperComputerName())
                return;





            return;
        }

    }

}
