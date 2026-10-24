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

namespace MSC.API.AdminSys.Controllers
{
    [AdminSysTokenrAuthorizeFilter]
    public class UserMerchantMappingController : BaseController
    {       
        public IUserMerchantMappingService UserMerchantMappingService { get; set; }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_User_List)]
        [Route("users/merchant-mapping/{userID}")]
        [HttpGet]
        public UserMerchantMappingViewModel GetUserMerchants(int userID)
        {
            var query = UserMerchantMappingService.GetUserMerchantsEnableFlag(userID);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_User_List)]
        [Route("users/merchant-mapping")]
        [HttpPost]
        public void UpdateUserMerchants([FromBody] UserMerchantMappingRequestModel model)
        {
            int userID = GetCurrentUserID();
            UserMerchantMappingService.UpdateUserMerchants(userID, model);
            return;
        }


    }
}
