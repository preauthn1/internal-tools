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
using MSC.Utility.Helper;
using MSC.Utility.IOC;
using MSC.Utility.Extensions;
using MSC.Models;
using MSC.Models.Models;
using MSC.BL.Interface;

namespace MSC.API.Job.Controllers
{
    /// <summary>
    /// Base Controller
    /// </summary>
    [ApiController]
    [Consumes("application/json")]
    [Produces("application/json")]
    public class BaseController : ControllerBase
    {

        /// <summary>
        /// Base Controller
        /// </summary>
        public BaseController()
        {
            InstanceAllServiceImplement(this);
        }


        /// <summary>
        /// 取得玩家ID
        /// </summary>
        protected int GetCurrentUserID()
        {
            string jwtToken = [REDACTED];
            var service = IOCHelper.Resolve<IAdminSysTokenService>();
            var jwtModel = service.GetAndVerifyTokenIsValid(jwtToken);            
            return jwtModel.LoginUserID;
        }


        /// <summary>
        /// 實例所有服務
        /// </summary>
        protected void InstanceAllServiceImplement<T>(T controller) where T : class
        {
            var allService = this.GetType()
                .GetProperties(BindingFlags.NonPublic | BindingFlags.Instance | BindingFlags.Public)
                .Where(x => x.PropertyType.IsInterface)
                .Where(x => x.Name.ToLower().Contains("service"))
                .ToList();

            foreach (var item in allService)
            {
                var instanceValue = IOCHelper.Resolve(item.PropertyType);
                item.SetValue(this, instanceValue);
            }
        }

    }
}
