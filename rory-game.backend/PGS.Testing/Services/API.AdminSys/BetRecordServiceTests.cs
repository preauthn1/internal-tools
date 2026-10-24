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
    public class BetRecordServiceTests : BaseTest
    {

        public IBetRecordService ServiceInstance { get; set; }


        [TestMethod]
        public void BetRecordService_GetUnionBetRecordsByPagingFilter_Test()
        {
            var query = ServiceInstance.GetUnionBetRecordsByPagingFilter(_TEST_DATA.Default.UserID, new QueryUnionBetRecordRequestViewModel
            {
                PageSize = 1
            });

            Assert.IsNotNull(query);
            Assert_NotEmptyCollection(query.Items);
            return;
        }


        [TestMethod]
        public void GetByTransactionNo_Test()
        {
            var query = ServiceInstance.GetByTransactionNo(_TEST_DATA.Default.UserID , new BetRecordQueryRequestModel
            {
                TransactionNo = _TEST_DATA.Default.TransactionNo_PgsBet
            });

            Assert.IsNotNull(query);
            return;
        }



    }
}

