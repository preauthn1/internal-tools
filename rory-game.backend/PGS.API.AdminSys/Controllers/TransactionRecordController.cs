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
using target.Models;
using target.Utility.Helper;
using target.Utility.Extensions;
using target.Models.Models;

using System.Collections;
using System.Resources;
using target.Models.Enums;

namespace target.API.AdminSys.Controllers
{
    /// <summary>
    /// 帳變紀錄
    /// </summary>
    public class TransactionRecordController : BaseController
    {
        /// <summary>
        /// 帳變紀錄
        /// </summary>
        public ITransactionRecordService TransactionRecordService { get; set; }


        /// <summary>
        /// 後台-金幣紀錄(帳變紀錄)
        /// </summary>
        [Route("transaction")]
        [HttpGet]
        [AdminTokenrAuthorizeFilter]
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.MembersManagementMoneyRecord)]
        public PagedList<TransactionRecordViewModel> GetPagingByFilter([FromQuery]QueryTransactionRecordRequestViewModel model)
        {
            int userID = GetCurrentUserID();
            var query = TransactionRecordService.GetPagingByFilter(userID, model);
            return query;
        }
    }
}
