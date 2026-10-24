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

namespace target.API.Game.Controllers
{
    /// <summary>
    /// Server Setting Controller
    /// </summary>
    public class ServerSettingController : BaseController
    {
        /// <summary>
        /// Game Server Setting Service
        /// </summary>
        public IGameServerSettingService GameServerSettingService { get; set; }


        /// <summary>
        /// 取得伺服器配置設定、開放給哪些供應商使用
        /// </summary>
        [Route("system/server/settings")]
        [HttpGet]
        public List<ServerSettingViewModel> GetServerListAndSetting([FromQuery]ServerListAndSettingRequestModel model)
        {
            var result = GameServerSettingService.GetServerListAndSetting(model);
            return result;
        }


        /// <summary>
        /// 伺服器「啟動」 相關註冊與初始化
        /// </summary>
        [Route("system/server/start-process")]
        [HttpPost]
        public ServerStartProcessViewModel RunServerStartProcess(ServerStartProcessRequestModel model)
        {
            var result = GameServerSettingService.RunServerStartProcess(model);
            return result;
        }


        /// <summary>
        /// 伺服器「關閉」 相關程序
        /// </summary>
        [Route("system/server/close-process")]
        [HttpPost]
        public void RunServerCloseProcess(ServerCloseProcessRequestModel model)
        {
            GameServerSettingService.RunServerCloseProcess(model);
        }



        /// <summary>
        /// GameServer的資料 可能是透過這隻建的
        /// Server (Gs) 啟動的時候 之後每30秒服務器註冊
        /// </summary>
        [Route("system/server/re-register")]
        [ApiConfig(IsHeaderLogEnable = false)]
        [HttpPost]
        public void ServerReRegister(SystemReregisterRequestModel model)
        {
            GameServerSettingService.ServerReRegister(model);
        }
        
    }
}
