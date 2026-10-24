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
using target.Models.Enums;

namespace target.API.AdminSys.Controllers
{
    public class TransferRecordController : BaseController
    {
        public ITransferRecordService TransferRecordService { get; set; }


        /// <summary>
        /// 後台-轉帳紀錄
        /// </summary>
        [Route("transfer")]
        [HttpGet]
        [AdminTokenrAuthorizeFilter]
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.MembersManagementTransfer)]
        public PagedList<TransferRecordViewModel> GetTransferRecordPagingByFilter([FromQuery]QueryTransferRecordApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = TransferRecordService.GetByPagingFilter(userID , filter);
            return query;
        }



    }
}
