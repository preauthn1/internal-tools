using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using MSC.BL.Interface;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Web.Modules.Filter;
using MSC.Models.Paging;
using MSC.Models.Enums;

namespace MSC.API.AdminSys.Controllers
{
    [AdminSysTokenrAuthorizeFilter]
    public class NotifyTriggerController : BaseController
    {
        public INotifyCardDepositService NotifyCardDepositService { get; set; }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.Record_CardConsumeMerchantNotifyQueue_List)]
        [Route("trigger/success-notify")]
        [HttpPost]
        public object TriggerNotifyAgainByAdminSys([FromBody]TriggerNotifyAgainRequestModel model)
        {
            int userID = GetCurrentUserID();
            NotifyCardDepositService.TriggerNotifyAgainByAdminSys(userID , model);
            return "OK";
        }


    }
}
