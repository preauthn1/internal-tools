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
    public class LinkRequestController : BaseController
    {       
        public ILinkGenerateService LinkGenerateService { get; set; }


        [Route("card/link-request")]
        [Route("card/link")]
        [HttpPost]
        public LinkGenerateViewModel GenearateCardLink([FromBody]LinkGenerateRequestModel model)
        {
            var result = LinkGenerateService.GenerateCardLink(model);
            return result;
        }

    }
}
