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
using MSC.Models;
using MSC.Models.API.Middleware;
using MSC.Utility;
using MSC.Utility.Helper;
using MSC.Web.Modules.Filter;

namespace MSC.API.Job.Controllers
{

    public class HomeController : BaseController
    {

        [ApiExplorerSettings(IgnoreApi = true)]
        [ApiConfig(IsLogApiEnable = false)]
        [Route("")]
        [HttpGet]
        public object HomePage(string password = "")
        {
            bool isPasswordIncorrect = password != SystemConstant.VendorAPI_SignKey_MagicPassword;
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
