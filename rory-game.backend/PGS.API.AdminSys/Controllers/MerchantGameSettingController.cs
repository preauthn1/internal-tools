using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using target.Web.Modules.Filter;
using target.Models.Paging;

namespace target.API.AdminSys.Controllers
{
    [AdminTokenrAuthorizeFilter]
    public class MerchantGameSettingController : BaseController
    {
        public IViewPermissionCheckService ViewPermissionCheckService { get; set; }


        public IMerchantGameSettingService MerchantGameSettingService { get; set; }


        [Route("merchant/{merchantID}/game-setting")]
        [HttpGet]
        public MerchantGameSettingViewModel GetMerchantGameSettingPagingByFilter(int merchantID)
        {
            int userID = GetCurrentUserID();
            ViewPermissionCheckService.CheckViewMerchantPermission(userID , merchantID);

            var query = MerchantGameSettingService.GetGameSettingByMerchantID(merchantID);
            return query;
        }


        [Route("merchant/{merchantID}/game-setting")]
        [HttpPut]
        public void Update(int merchantID , [FromBody]MerchantGameSettingUpdateRequestModel model)
        {
            int userID = GetCurrentUserID();
            ViewPermissionCheckService.CheckViewMerchantPermission(userID, merchantID);

            MerchantGameSettingService.Update(merchantID , model);
        }


    }
}
