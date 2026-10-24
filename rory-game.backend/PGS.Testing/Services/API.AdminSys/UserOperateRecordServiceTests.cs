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
    public class UserOperateRecordServiceTests : BaseTest
    {

        public IUserOperateRecordService ServiceInstance { get; set; }


        [TestMethod]
        public void UserOperateRecordService_GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID , new QueryUserOperateRecordApiRequestParameterModel
            {
                
            });

            Assert_NotEmptyCollection(query.Items);
            return;
        }
        

        [TestMethod]
        public void UserOperateRecordService_GetByPagingFilter_FilterMerchantID_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID, new QueryUserOperateRecordApiRequestParameterModel
            {
                MerchantID = _TEST_DATA.Default.MerchantID
            });

            Assert_NotEmptyCollection(query.Items);
            return;
        }


    }
}

