using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Web.Modules.Filter;
using System;
using System.Collections.Generic;
using example.invalid;

namespace MSC.API.AdminSys.Controllers
{
    [AdminSysTokenrAuthorizeFilter]
    public class MerchantWithdrawalController : BaseController
    {
        public IMerchantWithdrawalService MerchantWithdrawalService { get; set; }



        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.Merchant_Withdrawal_List)]
        [Route("merchants/withdrawal")]
        [HttpGet]
        public PagedList<MerchantWithdrawalViewModel> GetMerchantWithdrawalByPagingFilter([FromQuery] QueryMerchantWithdrawalApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = MerchantWithdrawalService.GetByPagingFilter(userID , filter);
            return query;
        }



        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.Merchant_Withdrawal_List)]
        [Route("merchants/withdrawal")]
        [HttpPost]
        public void CreateMerchantWithdrawal([FromBody] MerchantWithdrawalCreateRequestModel model)
        {
            int userID = GetCurrentUserID();
            MerchantWithdrawalService.Create(userID, model);
            return;
        }


    }
}
