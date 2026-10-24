using Microsoft.VisualStudio.TestTools.UnitTesting;
using MSC.BL.Interface;
using MSC.DA.Repository;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.Testing
{
    [TestClass]
    public class CardGenerateServiceTests : BaseTest
    {

        public ICardGenerateService ServiceInstance { get; set; }


        [TestMethod]
        public void CardGenerateService_GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID, new QueryCardInfoApiRequestParameterModel
            {
                
            });

            Assert_NotEmptyCollection(query.Items);
            Assert.IsTrue(query.Items.All(z=> z.Currency?.Length > 0));
            return;
        }


        [TestMethod]
        public void CardGenerateService_BatchGenerateCardNumber_Test()
        {
            var query = ServiceInstance.BatchGenerateCardNumber(_TEST_DATA.Default.UserID, new CardGenerateRequestModel
            {
                Currency = CurrencyEnum.USD.ToString(),
                ExpectCount = 1,
                Point = 50,
                Memo = "MSTest"
            });

            return;
        }


    }
}

