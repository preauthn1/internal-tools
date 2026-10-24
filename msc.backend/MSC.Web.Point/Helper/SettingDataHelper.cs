using Lck.Utility.Extensions.Helper;
using MSC.Models.Enums;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.Base;
using MSC.Utility;
using MSC.Utility.Utility;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace MSC.Web.Point.Helper
{
    public class SettingDataHelper
    {

        public static string GetShoppingIframeUrl()
        {
            string cacheKey = CacheKey.PointAPI_ShoppingIframeUrl();

            var queryCache = CacheHelper.GetCacheValue<string>(cacheKey);

            if (queryCache != null)
                return queryCache;

            var query = GetShoppingIframeUrlByApi();

            CacheHelper.SetCacheValue(cacheKey, query, 15);
            return query;
        }


        public static List<ShoppingStoreViewModel> GetShoppingStoreByCache()
        {
            string cacheKey = CacheKey.PointAPI_ShoppingStore();

            var queryCache = CacheHelper.GetCacheValue<List<ShoppingStoreViewModel>>(cacheKey);

            if (queryCache != null)
                return queryCache;

            var query = GetShoppingStore();

            CacheHelper.SetCacheValue(cacheKey, query, 15);
            return query;
        }


        public static CardInputPageSettingViewModel GetPageLayoutSettings(int merchantID)
        {
            string cacheKey = CacheKey.PointAPI_PageLayoutSetting(merchantID);

            var queryCache = CacheHelper.GetCacheValue<CardInputPageSettingViewModel>(cacheKey);

            if (queryCache != null)
                return queryCache;

            var query = GetPageLayoutSetting(merchantID);

            CacheHelper.SetCacheValue(cacheKey, query, 15);
            return query;
        }




        public static List<ShoppingStoreViewModel> GetShoppingStore()
        {
            //呼叫API檢查
            string baseurl = SettingHelper.WebApplicationSetting.F2E_Point_API_Doamin_Internal;
            var res = APIRequestHelper.RequestApiResult<ApiResponseModel<List<ShoppingStoreViewModel>>>(new RequestParameters
            {
                BaseUrl = baseurl,
                ApiUrl = "settings/shopping-store",
                Method = HttpMethodType.HttpGet

            });

            if (res.Exception == null && res.ResultObject.Data != null)
            {
                //只顯示enable的
                return res.ResultObject.Data
                    .Where(z => z.IsEnable)
                    .ToList(); ;
            }

            return new List<ShoppingStoreViewModel>();
        }



        public static CardInputPageSettingViewModel GetPageLayoutSetting(int merchantID)
        {
            //呼叫API檢查
            string baseurl = SettingHelper.WebApplicationSetting.F2E_Point_API_Doamin_Internal;
            var res = APIRequestHelper.RequestApiResult<ApiResponseModel<CardInputPageSettingViewModel>>(new RequestParameters
            {
                BaseUrl = baseurl,
                ApiUrl = $"settings/card-page/{merchantID}",
                Method = HttpMethodType.HttpGet

            });

            if (res.Exception == null && res.ResultObject.Data != null)
            {
                //只顯示enable的
                return res.ResultObject.Data;
            }

            return new CardInputPageSettingViewModel
            {
                WebSiteBrandingType = (int)WebSiteBrandingTypeEnum.MSCard,
                IsBackButtonVisible = true
            };
        }




        public static string GetShoppingIframeUrlByApi()
        {
            //呼叫API檢查
            string baseurl = SettingHelper.WebApplicationSetting.F2E_Point_API_Doamin_Internal;
            var res = APIRequestHelper.RequestApiResult<ApiResponseModel<string>>(new RequestParameters
            {
                BaseUrl = baseurl,
                ApiUrl = "settings/shopping-iframe-url",
                Method = HttpMethodType.HttpGet

            });

            if (res.Exception == null && res.ResultObject.Data != null)
            {
                return res.ResultObject.Data ?? "";
            }

            return string.Empty;
        }



    }
}
