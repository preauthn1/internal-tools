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
using target.Models.Models;
using target.Models.ViewModels;
using target.Models.ViewModels.Base;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.Utility.IOC;

namespace target.API.AdminSys.Controllers
{
    [Consumes("application/json")]
    [Produces("application/json")]
    [ApiController]
    public class BaseController : ControllerBase
    {
        public BaseController()
        {
            InstanceAllServiceImplement(this);
        }


        /// <summary>
        /// 根據當前登入的用戶jwt，設定可以篩選哪些merchantIDs
        /// </summary>
        protected void SetFilterMerchantIDs<T>(T model) where T: IMerchantFilterModel
        {
            var userJwtModel = GetAdminUserModelFromJwtHeader();

            if (userJwtModel == null)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);

            //若為最高admin(或組織) 不需篩選merchantID
            if (userJwtModel.IsSuperAdminOrg)
                return;

            model.MerchantIDs = userJwtModel.MerchantIDs ?? new List<int>();
        }



        protected List<int> GetFilterMerchantIDsOrEmpty()
        {
            var userJwtModel = GetAdminUserModelFromJwtHeader();

            if (userJwtModel == null)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);

            //若為最高admin(或組織) 不需篩選merchantID
            if (userJwtModel.IsSuperAdminOrg)
                return new List<int>();

            //若非管理員組織，不允許組織為空
            if (userJwtModel.MerchantIDs.IsEmptyOrNull())
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_SettingDataNotCompleted);

            return userJwtModel.MerchantIDs;
        }


        protected List<int> GetFilterOrganizationIDsOrEmpty()
        {
            var userJwtModel = GetAdminUserModelFromJwtHeader();

            if (userJwtModel == null)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);

            //若為最高admin(或組織) 不需篩選merchantID
            if (userJwtModel.IsSuperAdminOrg)
                return new List<int>();

            return userJwtModel.OrganizationID.ToListCollection();
        }



        protected int GetCurrentUserID()
        {
            var userJwtModel = GetAdminUserModelFromJwtHeader();

            if (userJwtModel == null)
                return 0;

            return userJwtModel.LoginUserID;
        }



        /// <summary>
        /// 直接呼叫此沒有檢查 過期/反查DB 的機制  請記得自行標上 AdminTokenrAuthorizeFilter 檢查
        /// </summary>
        /// <returns></returns>
        protected AdminUserTokenJwtModel GetAdminUserModelFromJwtHeader()
        {
            string jwt = HttpContext.Request.GetHeaderValueByKey(SystemConstant.HeaderKey_AccessToken);
            AdminUserTokenJwtModel userJwtModel = MyJwtHelper.GetAdminUserModelFromJwtString(jwt);
            return userJwtModel;
        }


        protected AdminUserTokenJwtModel GetAdminUserModelFromJwtAndCheckException(string jwt)
        {
            if (string.IsNullOrWhiteSpace(jwt))
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);

            AdminUserTokenJwtModel userJwtModel = MyJwtHelper.GetAdminUserModelFromJwtString(jwt);

            if (userJwtModel == null)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);

            return userJwtModel;
        }

        

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
