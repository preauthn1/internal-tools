using Microsoft.VisualStudio.TestTools.UnitTesting;
using MSC.BL.Interface;
using MSC.DA.Repository;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.Testing
{
    [TestClass]
    public class MerchantWithdrawalServiceTests : BaseTest
    {

        public IMerchantWithdrawalService ServiceInstance { get; set; }


        [TestMethod]
        public void MerchantWithdrawalService_GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID, new QueryMerchantWithdrawalApiRequestParameterModel
            {
                
            });

            Assert_NotEmptyCollection(query.Items);
            return;
        }

    }
}

