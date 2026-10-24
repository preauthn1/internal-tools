using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models;
using MSC.Models.Enums;
using MSC.Models.Paging;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.API.Point;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Web.Modules.Filter;

namespace MSC.API.AdminSys.Controllers
{
    [AdminSysTokenrAuthorizeFilter]
    public class MerchantController : BaseController
    {       
        public IMerchantService MerchantService { get; set; }

        public ILinkGenerateService LinkGenerateService { get; set; }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_Merchant_List)]
        [Route("merchants")]
        [HttpGet]
        public PagedList<MerchantViewModel> GetMerchantByPagingFilter([FromQuery]QueryMerchantApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = MerchantService.GetByPagingFilter(userID, filter);
            return query;
        }



        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_Merchant_List)]
        [Route("merchants/{merchantID}/settings")]
        [HttpGet]
        public SystemConfigSettingPageViewModel GetByPagingFilter(int merchantID)
        {
            int userID = GetCurrentUserID();
            var query = MerchantService.GetMerchantSettings(userID, merchantID);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_Merchant_List)]
        [Route("merchants/{merchantID}/settings")]
        [HttpPost]
        public void UpdateMerchantSettings(int merchantID , [FromBody]List<SystemConfigSettingUpdateRequestItemModel> models)
        {
            int userID = GetCurrentUserID();
            MerchantService.UpdateMerchantSettings(userID, merchantID , models);
        }



        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_Merchant_List)]
        [Route("merchants")]
        [HttpPost]
        public void Create([FromBody] MerchantCreateRequestModel filter)
        {
            int userID = GetCurrentUserID();
            MerchantService.Create(userID, filter);
            return;
        }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_Merchant_List)]
        [Route("merchants")]
        [HttpPut]
        public void Update([FromBody] MerchantUpdateRequestModel filter)
        {
            int userID = GetCurrentUserID();
            MerchantService.Update(userID, filter);
            return;
        }



        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_Merchant_List)]
        [Route("merchants/initialize/wallet")]
        [HttpPost]
        public void InitializeWalletData()
        {
            int userID = GetCurrentUserID();
            MerchantService.InitializeWalletData();
            return;
        }

        
        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.System_Merchant_List)]
        [Route("merchants/order-gen")]
        [HttpGet]
        public LinkGenerateViewModel GenerateTestOrderLinkByMerchantID([FromQuery]int merchantID)
        {
            int userID = GetCurrentUserID();
            var linkResult = LinkGenerateService.GenerateTestOrderLinkByMerchantID(merchantID);
            return linkResult;
        }


    }
}
