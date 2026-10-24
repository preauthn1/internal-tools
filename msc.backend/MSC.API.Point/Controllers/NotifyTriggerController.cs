using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.API.Game.Controllers;
using MSC.BL.Interface;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Web.Modules.Filter;

namespace MSC.API.Point.Controllers
{
    public class NotifyTriggerController : BaseController
    {       
        public INotifyCardDepositService NotifyCardDepositService { get; set; }

        
        [Route("card/trigger-notify-again")]
        [HttpPost]
        public string InputCardNo([FromBody]TriggerNotifyAgainRequestModel model)
        {
            NotifyCardDepositService.TriggerNotifyAgain(model);
            return "trigger end";
        }        

    }
}
