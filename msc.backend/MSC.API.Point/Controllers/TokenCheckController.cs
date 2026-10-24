using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.API.Game.Controllers;
using MSC.BL.Interface;
using MSC.Models.Models;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.RequestViewModel;

namespace MSC.API.Point.Controllers
{
    public class TokenCheckController : BaseController
    {
        public ITokenCheckRefService TokenCheckRefService { get; set; }


        [Route("token/valid-check")]
        [HttpPost]
        public LinkTokenMemberInfoModel CheckIsValidToken([FromBody]TokenCheckRequestModel model)
        {
            var result = TokenCheckRefService.CheckIsValidToken(model?.Token);
            return result;
        }
        

    }
}
