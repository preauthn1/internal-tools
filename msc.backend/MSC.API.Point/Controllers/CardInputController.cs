using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.API.Game.Controllers;
using MSC.BL.Interface;
using MSC.Models;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Utility.Helper;
using MSC.Web.Modules.Filter;

namespace MSC.API.Point.Controllers
{
    public class CardInputController : BaseController
    {       
        public ICardInputService CardInputService { get; set; }


        [Route("card/input")]
        [HttpPost]
        public CardInputViewModel InputCardNo([FromBody]CardInputRequestModel model)
        {
            var result = CardInputService.InputCard(model);
            return result;
        }


    }
}
