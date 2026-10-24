using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using target.Models.Paging;
using target.Web.Modules.Filter;
using target.Models.Enums;

namespace target.API.AdminSys.Controllers
{
    public class PerGameDetailController : BaseController
    {
        public IPerGameDetailService PerGameDetailService { get; set; }

        /// <summary>
        /// 游戏历程查询
        /// </summary>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.GamesManagementRecord)]
        [Route("perGame/player-details")]
        [HttpGet]
        public PagedList<PerGameDetailViewModel> GetByPagingFilter([FromQuery]QueryPerGameDetailLogApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = PerGameDetailService.GetByPagingFilter(userID , filter);
            return query;
        }


        /// <summary>
        /// 游戏历程查询 (供外部查詢、不需要登入token)
        /// </summary>
        [Route("perGame/slot-result/transactionNo")]
        [HttpGet]
        public PerGameSlotHistoryDetailInfoViewModel GetSlotHistoryByExternalLinkQuery([FromQuery] ExternalLinkQueryPerGameRequestModel model)
        {
            var query = PerGameDetailService.GetSlotHistoryByExternalLinkQuery(model);
            return query;
        }


        /// <summary>
        /// 游戏历程查询
        /// </summary>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.GamesManagementRecord)]
        [Route("perGame/slot-result/{perGameDetailLogID}")]
        [HttpGet]
        public PerGameSlotHistoryDetailInfoViewModel GetSlotHistoryFromPerGameLog(long perGameDetailLogID)
        {
            int userID = GetCurrentUserID();
            var query = PerGameDetailService.GetSlotHistoryFromPerGameLog(userID, perGameDetailLogID);
            return query;
        }


        /// <summary>
        /// 游戏历程查询-详情
        /// </summary>
        /// <param name="perGameLogID"></param>
        /// <returns></returns>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.GamesManagementRecord)]
        [Route("perGame/{perGameLogID}/full-info")]
        [HttpGet]
        public PerGameLogFullInformationViewModel GetFullInformationInPerGame(int perGameLogID)
        {
            int userID = GetCurrentUserID();
            var query = PerGameDetailService.GetFullInformationInPerGame(userID , perGameLogID);
            return query;
        }


    }
}
