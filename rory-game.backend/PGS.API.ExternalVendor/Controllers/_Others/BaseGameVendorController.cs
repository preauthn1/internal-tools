using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.ComponentModel.DataAnnotations;
using example.invalid;
using System.Linq;
using System.Reflection;
using System.Runtime.Versioning;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using target.BL.Interface;
using target.BL.Services;
using target.DA.Repository;
using target.Models.Enums;
using target.Models.Models;
using target.Utility;
using target.Utility.IOC;

namespace target.API.ExternalVendor.Controllers
{
    [ApiController]
    public abstract class BaseGameVendorController : BaseController
    {

        public IGameVendorSingleWalletService SingleWalletVendorService { get; set; }

        protected abstract GameVendorTypeEnum _gameVendor { get; }


        /// <summary>
        /// 根據配置 模擬錯誤發生 (建議新對接一家遊戲商，要測試下對方有無retry機致)
        /// </summary>
        protected void Develop_Simulate_DelayTimeout(SystemConfigSettingEnum type)
        {
            if (SettingHelper.IsProduction)
                return;

            var rep = new SystemConfigSettingRepository();

            var query = rep.GetParseModelByCache(type);

            if (query == null)
                return;

            var delaySeconds = query.SettingValueOfInt;

            if (delaySeconds <= 0)
                return;

            Thread.Sleep(1000 * delaySeconds);
            return;
        }


    }
}
