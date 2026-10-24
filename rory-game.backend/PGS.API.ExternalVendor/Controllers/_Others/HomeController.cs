using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.ComponentModel.DataAnnotations;
using example.invalid;
using System.Linq;
using System.Reflection;
using System.Runtime.Versioning;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using target.Models;
using target.Models.ViewModels.Base;
using target.Utility;
using target.Web.Modules.Filter;

namespace target.API.ExternalVendor.Controllers
{
    public class HomeController : BaseController
    {
        [ApiExplorerSettings(IgnoreApi = true)]
        [ApiConfig(IsLogApiEnable = false)]
        [Route("")]
        [HttpGet]
        public object HomePage(string password = "")
        {
            bool isPasswordIncorrect = password != SystemConstant.PublicAPI_SignKey_MagicPassword;
            bool isNotComapnyIP = IPHelper.GetIP() != SystemConstant.CompanyIPAddress;

            //externalAPI [總是] 保持password保護 (因為廠商會看到)
            if (isPasswordIncorrect && isNotComapnyIP)
            {
                //production環境會隱藏某些資訊
                return ConvertToGeneralFormat(RuntimeHelper.GetHideInformationForProduction());
            }

            return ConvertToGeneralFormat(RuntimeHelper.GetHomePageDebugInformation());
        }

       
        /// <summary>
        /// 讀取指定的config
        /// </summary>
        [ApiExplorerSettings(IgnoreApi = true)]
        [ApiConfig(IsLogApiEnable = false)]
        [Route("config/read")]
        public object GetConfigValue(string key = "")
        {
            return BaseSettingHelper.GetSettingValue(key);
        }

        /// <summary>
        /// ExternalAPI 沒有套用 ApiResponseFormatterFilter 所以格式會略有落差
        /// </summary>
        private ApiResponseModel ConvertToGeneralFormat(object model)
        {
            return new ApiResponseModel
            {
                Code = 0,
                Data = model
            };
        }

    }
}
