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
using System.Threading;

namespace MSC.API.AdminSys.Controllers
{
    [AdminSysTokenrAuthorizeFilter]
    public class DropdownController : BaseController
    {
        public IDropdownService DropdownService { get; set; }


        [Route("dropdown/merchants")]
        [HttpGet]
        public List<DropdownSelectValueViewModel> GetMerchants()
        {
            int userID = GetCurrentUserID();
            var query = DropdownService.GetMerchants(userID);
            return query;
        }


        [Route("dropdown/merchantsWithCurrency")]
        [HttpGet]
        public List<DropdownSelectValueViewModel> GetMerchantsWithCurrency()
        {
            int userID = GetCurrentUserID();
            var query = DropdownService.GetMerchantsWithCurrency(userID);
            return query;
        }


        [Route("dropdown/status")]
        [HttpGet]
        public List<DropdownSelectValueViewModel> GetNormalStatus()
        {
            int userID = GetCurrentUserID();
            var query = DropdownService.GetNormalStatus();
            return query;
        }


        [Route("dropdown/transactionTypes")]
        [HttpGet]
        public List<DropdownSelectValueViewModel> GetTransactionTypes()
        {
            int userID = GetCurrentUserID();
            var query = DropdownService.GetTransactionTypes();
            return query;
        }


        [Route("dropdown/operateTypes")]
        [HttpGet]
        public List<DropdownSelectValueViewModel> GetUserOperateTypes()
        {
            int userID = GetCurrentUserID();
            var query = DropdownService.GetUserOperateTypes();
            return query;
        }


        [Route("dropdown/currency")]
        [HttpGet]
        public List<DropdownSelectValueViewModel> GetCurrency()
        {
            int userID = GetCurrentUserID();
            var query = DropdownService.GetCurrency();
            return query;
        }



        [Route("dropdown/roles")]
        [HttpGet]
        public List<DropdownSelectValueViewModel> GetRoles()
        {
            int userID = GetCurrentUserID();
            var query = DropdownService.GetRoles();
            return query;
        }

    }
}
