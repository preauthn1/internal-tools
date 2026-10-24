using Microsoft.VisualStudio.TestTools.UnitTesting;
using MSC.BL.Interface;
using MSC.Models.Models;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System.Linq;
using MSC.Models.CoreLibs;
using MSC.Models.ViewModels;
using MSC.DA.Repository;
using MSC.Models.Enums;
using System.Collections.Generic;
using System;

namespace MSC.Testing
{
    [TestClass]
    public class CardLinkRequestRecordServiceTests : BaseTest
    {

        public ICardLinkRequestRecordService ServiceInstance { get; set; }


        [TestMethod]
        public void CardLinkRequestRecordService_GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID, new QueryCardLinkRequestRecordApiRequestParameterModel
            {
                PageSize = 10
            });

            Assert_NotEmptyCollection(query.Items);


            if(query.Items.IsNotEmptyNull())
            {
                var matchFirst = query.Items.FirstOrDefault();
                var queryDetail = ServiceInstance.GetProcessDetail(_TEST_DATA.Default.UserID, new MerchantRequestProcessDetailInfoRequestModel
                {
                    MerchantID = matchFirst.MerchantID,
                    MerchantRequestID = matchFirst.MerchantRequestID
                });

                Assert.IsNotNull(queryDetail);
            }

            return;
        }

    }
}

