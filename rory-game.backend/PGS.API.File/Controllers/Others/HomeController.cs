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
using target.API.File.Helper;
using target.Models;
using target.Utility;
using target.Web.Modules.Filter;

namespace target.API.File.Controllers
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
            bool isProduction = SettingHelper.IsProduction;
            bool isNotComapnyIP = IPHelper.GetIP() != SystemConstant.CompanyIPAddress;

            if (isProduction)
            {
                if (isPasswordIncorrect && isNotComapnyIP)
                {
                    //production環境會隱藏某些資訊
                    return RuntimeHelper.GetHideInformationForProduction();
                }
            }

            return RuntimeHelper.GetHomePageDebugInformation();
        }                
      
    }
}
