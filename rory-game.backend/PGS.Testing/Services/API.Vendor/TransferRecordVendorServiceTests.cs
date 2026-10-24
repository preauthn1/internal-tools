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
    public class TransferRecordVendorServiceTests : BaseTest
    {

        public ITransferRecordVendorService ServiceInstance { get; set; }


        public TransferRecordRepository TransferRecordRepository { get; set; }



        [TestMethod]
        public void GetTransferByMerchant_Test()
        {
            var queryByDB = TransferRecordRepository.GetPagingByFilter(new QueryTransferRecordRequestViewModel
            {
                MerchantIDs = _TEST_DATA.Default.MerchantID.ToListCollection(),
                PageSize = 5
            });
            
            var query = ServiceInstance.GetTransferByMerchant(new QueryTransferRecordRequestModel
            {
                MerchantID = _TEST_DATA.Default.MerchantID,
                MerchantSecretKey = _TEST_DATA.Default.MerchantSecret,
                MerchantTransactionNo = queryByDB.Items.FirstOrDefault().MerchantOrderNo
            });

            Assert.IsNotNull(query);
            return;
        }

    }
}

