using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.API.Game.Controllers
{
    public class SystemSettingController : BaseController
    {
        public ISystemSettingService SystemSettingService { get; set; }


        /// <summary>
        /// 取得黑名單IP，Gs如果防護配置有打開才會呼叫
        /// </summary>
        [Route("system/setting/deny-ip")]
        [HttpGet]
        public List<DenyIpAddressSettingViewModel> GetDenyIpAddressList()
        {
            var query = SystemSettingService.GetDenyIpAddressList();
            return query;
        }

    }
}
