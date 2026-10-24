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
using target.Models.Enums;

namespace target.API.AdminSys.Controllers
{
    [AdminTokenrAuthorizeFilter]
    public class RecordsController : BaseController
    {
        public IPlayerLoginLogService PlayerLoginLogService { get; set; }

        public IUserOperateRecordService UserOperateRecordService { get; set; }

        public IPlayerGameRecordService PlayerGameRecordService { get; set; }

        public IExternalApiRequestLogService ExternalApiRequestLogService { get; set; }

        public IBetRecordService BetRecordService { get; set; }


        /// <summary>
        /// 後台-進出房紀錄
        /// </summary>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.MembersManagementBetRecord)]
        [Route("records/game-play")]
        [HttpGet]
        public PagedList<PlayerGameRecordViewModel> GetPlayerGameRecordPagingByFilter([FromQuery]QueryPlayerGameRecordApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = PlayerGameRecordService.GetByPagingFilter(userID, filter);
            return query;
        }


        /// <summary>
        /// 後台-玩家遊戲啟動紀錄、三方遊戲啟動紀錄
        /// </summary>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.MembersManagementDetail)]
        [Route("records/player-login")]
        [HttpGet]
        public PagedList<PlayerLoginLogViewModel> GetPlayerLoginLogPagingByFilter([FromQuery]QueryPlayerLoginLogApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = PlayerLoginLogService.GetByPagingFilter(userID , filter);
            return query;
        }


        /// <summary>
        /// 後台-玩家封號紀錄
        /// </summary>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.MembersManagementLock)]
        [Route("records/modify-player-status")]
        [HttpGet]
        public PagedList<UserOperateRecordViewModel> GetPlayerStatusChangeRecords([FromQuery]QueryUserOperateRecordApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = UserOperateRecordService.GetByPagingFilter(userID, filter);
            return query;
        }


        [Route("records/bet/find")]
        [HttpGet]
        public BetRecordGeneralInfoViewModel GetBetRecordByTransactionNo([FromQuery] BetRecordQueryRequestModel filter)
        {
            int userID = GetCurrentUserID();
            var query = BetRecordService.GetByTransactionNo(userID, filter);
            return query;
        }


        [Route("records/ex-api-logs")]
        [HttpGet]
        public PagedList<ExternalApiRequestLogViewModel> GetExternalLogsByRequestID([FromQuery] QueryExternalApiRequestLogApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = ExternalApiRequestLogService.GetByPagingFilter(userID , filter);
            return query;
        }


        [Route("records/ex-api-logs/id/{id}")]
        [HttpGet]
        public ExternalApiRequestLogViewModel GetExternalLogByID(long id)
        {
            int userID = GetCurrentUserID();
            var query = ExternalApiRequestLogService.GetByID(id);
            return query;
        }



        [Route("records/ex-api-logs/RequestID/{requestID}")]
        [HttpGet]
        public ExternalApiRequestLogViewModel GetByRequestID(string requestID)
        {
            int userID = GetCurrentUserID();
            var query = ExternalApiRequestLogService.GetByRequestID(requestID);
            return query;
        }



    }
}
