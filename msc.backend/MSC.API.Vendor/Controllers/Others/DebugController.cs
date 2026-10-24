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
    public class DebugController : BaseController
    {       
        public ITokenCheckRefService TokenCheckRefService { get; set; }


        [Route("debug/is-valid-card")]
        [HttpPost]
        public object CheckIsValidToken(string token)
        {
            var result = TokenCheckRefService.CheckIsValidToken(token);
            return result;
        }

    }
}
