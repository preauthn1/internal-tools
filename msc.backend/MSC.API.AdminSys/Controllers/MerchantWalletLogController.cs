using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using MSC.BL.Interface;
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
    public class MerchantWalletLogController : BaseController
    {
        public IMerchantWalletLogService MerchantWalletLogService { get; set; }


        [Route("records/merchant/wallet-logs")]
        [HttpGet]
        public PagedList<MerchantWalletLogViewModel> GetMerchantWalletLogPagingByFilter([FromQuery] QueryMerchantWalletLogApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = MerchantWalletLogService.GetByPagingFilter(userID , filter);
            return query;
        }

    }
}
