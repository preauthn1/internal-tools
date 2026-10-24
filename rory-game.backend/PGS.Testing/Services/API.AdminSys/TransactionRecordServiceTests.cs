using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Models;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class TransactionRecordServiceTests : BaseTest
    {

        public ITransactionRecordService ServiceInstance { get; set; }


        [TestMethod]
        public void TransactionRecordService_ServiceInstance_Test()
        {
            var query = ServiceInstance.GetPagingByFilter(_TEST_DATA.Default.UserID , new QueryTransactionRecordRequestViewModel());
            Assert_NotEmptyCollection(query.Items);
            Assert.IsNotNull(new CRUDResultModel());
            return;
        }

    }
}

