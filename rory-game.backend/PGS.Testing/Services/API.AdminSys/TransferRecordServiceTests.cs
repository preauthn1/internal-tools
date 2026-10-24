using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.Models.Models;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Linq;
using target.Models.CoreLibs;
using target.Models.ViewModels;
using target.DA.Repository;
using target.Models.Enums;
using System.Collections.Generic;
using System;

namespace target.Testing
{
    [TestClass]
    public class TransferRecordServiceTests : BaseTest
    {

        public ITransferRecordService ServiceInstance { get; set; }


        [TestMethod]
        public void TransferRecordService_GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID, new QueryTransferRecordApiRequestParameterModel
            {
                PageSize = 20
            });

            Assert_NotEmptyCollection(query.Items);
            return;
        }



        [TestMethod]
        public void TransferRecordService_GetByPagingFilter_FilterColumns_Test()
        {
            DateTime now = DateTimeHelper.GetCurrentTime();

            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID, new QueryTransferRecordApiRequestParameterModel
            {
                PageSize = 20,
                CreateTimeBegin = now.AddMonths(-2),
                CreateTimeEnd = now,
            });

            Assert_NotEmptyCollection(query.Items);
            return;
        }



    }
}

