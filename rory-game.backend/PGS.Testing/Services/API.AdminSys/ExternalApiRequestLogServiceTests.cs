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
    public class ExternalApiRequestLogServiceTests : BaseTest
    {

        public IExternalApiRequestLogService ServiceInstance { get; set; }


        [TestMethod]
        public void ExternalApiRequestLogService_GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(_TEST_DATA.Default.UserID , new QueryExternalApiRequestLogApiRequestParameterModel());
            Assert_NotEmptyCollection(query.Items);

            if(query.Items.IsNotEmptyNull())
            {
                ExternalApiRequestLogService_GetByRequestID_Test(query.Items.FirstOrDefault().RequestID);
            }

            return;
        }



        public void ExternalApiRequestLogService_GetByRequestID_Test(string requestID)
        {
            var query = ServiceInstance.GetByRequestID(requestID);
            Assert.IsNotNull(query);
            return;
        }



    }
}

