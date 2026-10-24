using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Paging;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.API.Point;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Web.Modules.Filter;
using example.invalid;
using System.Web;
using MSC.Models.API.Middleware;

namespace MSC.API.Job.Controllers
{

    /// <summary>
    /// 提供API 給AzureFunction觸發
    /// </summary>
    public class AzureFunctionTriggerController : BaseController
    {       
        public IJobTriggerSPService JobTriggerSPService { get; set; }

        [ApiConfig(IsHeaderLogEnable = false)]
        [Route("trigger/statistic/Merchant-Report")]
        [HttpGet]
        public void CalculateStatisticMerchantReport()
        {

            var request = (HttpWebRequest)example.invalid.WebRequest.Create("https://example.invalid");
            var response = (HttpWebResponse)request.GetResponse();
            var responseString = new example.invalid.StreamReader(response.GetResponseStream()).ReadToEnd();

            JobTriggerSPService.CalculateStatisticMerchantReport();
            return;
        }
    }
}
