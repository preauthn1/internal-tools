using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.RequestViewModel;

namespace MSC.API.Vendor.Controllers
{
    public class RecordController : BaseController
    {       
        public IMerchantRecordService MerchantRecordService { get; set; }


        [Route("records/link-request")]
        [HttpPost]
        public MerchantRequestRecordViewModel GetMerchantRequestRecord([FromBody]VendorApiQueryMerchantRequestRecordRequestModel model)
        {
            var result = MerchantRecordService.GetMerchantRequestRecord(model);
            return result;
        }

    }
}
