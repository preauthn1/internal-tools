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
    [AdminTokenrAuthorizeFilter]
    public class PlayerInfoController : BaseController
    {
        public IPlayerInfoService PlayerInfoService { get; set; }

        public IViewPermissionCheckService IViewPermissionCheckService { get; set; }


        /// <summary>
        /// 用户列表
        /// </summary>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.MembersManagementList_Search)]
        [Route("players")]
        [HttpGet]
        public PagedList<PlayerInfoViewModel> GetByPagingFilter([FromQuery]QueryPlayerApiRequestParameterModel model)
        {
            int userID = GetCurrentUserID();
            var query = PlayerInfoService.GetByPagingFilter(userID , model);
            return query;
        }


        /// <summary>
        /// 用戶補錢 (內部測試帳號)
        /// </summary>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.MembersManagementList_Search)]
        [Route("players/internal-account/manual-deposit")]
        [HttpPost]
        public PlayerManualDepositResultViewModel InternalTestAccountWalletDepositAmount([FromBody] InternalTestAccountPlusWalletRequestModel model)
        {
            int userID = GetCurrentUserID();
            var result = PlayerInfoService.InternalTestAccountWalletDepositAmount(userID, model);
            return result;
        }


        /// <summary>
        /// 產生測試用帳號
        /// </summary>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.MembersManagementList_Search)]
        [Route("players/internal-account/create")]
        [HttpPost]
        public void CreateInernalTestAccount([FromBody] CreateInternalTestAccountRequestModel model)
        {
            int userID = GetCurrentUserID();
            PlayerInfoService.CreateInernalTestAccount(userID, model);
            return;
        }


        /// <summary>
        /// 用戶補錢 (目前僅供使用API觸發)
        /// </summary>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.MembersManagementList_Search)]
        [Route("players/manual-deposit")]
        [HttpPost]
        public PlayerManualDepositResultViewModel ManualDepositPlayerWallet([FromBody] PlayerManualDepositRequestModel model)
        {
            int userID = GetCurrentUserID();
            var result = PlayerInfoService.ManualDepositPlayerWallet(userID, model);
            return result;
        }


        /// <summary>
        /// 用户查询
        /// </summary>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.MembersManagementList_Search , AdminSysPermissionKey.MembersManagementDetail_Search)]
        [Route("players/{playerID}")]
        [HttpGet]
        public PlayerInfoSummaryViewModel GetPlayerSummaryInfo(int playerID)
        {
            int userID = GetCurrentUserID();
            IViewPermissionCheckService.CheckViewPlayerPermission(userID , playerID);
            var query = PlayerInfoService.GetPlayerSummaryInfo(playerID);
            return query;
        }


        /// <summary>
        /// 用户查询
        /// </summary>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.MembersManagementList_Search, AdminSysPermissionKey.MembersManagementDetail_Search)]
        [Route("players/detail")]
        [HttpGet]
        public PlayerInfoSummaryViewModel GetPlayerSummaryInfoByNameOrID(string playerNameOrID)
        {
            int userID = GetCurrentUserID();

            int playerID = PlayerInfoService.GetPlayerIdByPlayerNameOrID(userID , playerNameOrID);
            IViewPermissionCheckService.CheckViewPlayerPermission(userID, playerID);
            var query = PlayerInfoService.GetPlayerSummaryInfo(playerID);
            return query;
        }


        /// <summary>
        /// 解除卡房
        /// </summary>
        [AdminPermissionAuthorizeFilter(AdminSysPermissionKey.MembersManagementDetail_Lock)]
        [Route("players/{playerID}/room-lock/remove")]
        [HttpDelete]
        public void RemovePlayerRoomLock(int playerID)
        {
            int userID = GetCurrentUserID();
            IViewPermissionCheckService.CheckViewPlayerPermission(userID, playerID);
            PlayerInfoService.RemovePlayerRoomLock(userID , playerID);
            return;
        }



        /// <summary>
        /// 用户锁定/解锁
        /// </summary>
        /// <param name="playerID"></param>
        /// <param name="model"></param>
        [Route("players/{playerID}/status")]
        [HttpPut]
        public void UpdatePlayerStatus(int playerID , [FromBody]PlayerUpdateStatusRequestModel model)
        {
            int userID = GetCurrentUserID();
            IViewPermissionCheckService.CheckViewPlayerPermission(userID, playerID);

            model.PlayerID = playerID;
            PlayerInfoService.UpdatePlayerStatus(userID , model);
        }

        /// <summary>
        /// 用户批量锁定
        /// </summary>
        /// <param name="model"></param>
        [Route("players/batch-lock")]
        [HttpPut]
        public void UpdateMultiplePlayerStatus([FromBody]PlayerMultipleUpdateStatusRequestModel model)
        {
            int userID = GetCurrentUserID();
            model.Status = (int)PlayerStatusEnum.Disable;
            PlayerInfoService.UpdateMultiplePlayerStatus(userID, model);
        }


    }
}
