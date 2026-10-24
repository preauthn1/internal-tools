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
    public class MerchantSettingController : BaseController
    {       

        public IMerchantSettingService MerchantSettingService { get; set; }

        [Route("merchant/card/available-amounts")]
        [HttpPost]
        public MerchantAvailableAmountViewModel GetAvialbleAmountSetting([FromBody]VendorApiMerchantAmountSettingRequestModel model)
        {
            var result = MerchantSettingService.GetAvialbleAmountSetting(model);
            return result;
        }

    }
}
