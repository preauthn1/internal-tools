using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models;
using MSC.Models.Models;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.RequestViewModel;

namespace MSC.API.Vendor.Controllers
{
    public class ProviderCardGenerateController : BaseController
    {       

        public ICardGenerateService CardGenerateService { get; set; }


        /// <summary>
        /// 新版需求 希望可以在淘寶訂單
        /// </summary>
        [Route("provider/card/gen")]
        [HttpPost]
        public ProviderCardGenerateResultViewModel ProviderGenerateCard([FromBody] CardProviderCardGenerateRequestModel model)
        {
            var result = CardGenerateService.ProviderGenerateCard(model);
            return result;
        }

    }
}
