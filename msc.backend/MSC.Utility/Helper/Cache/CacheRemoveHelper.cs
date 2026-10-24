using Lck.Utility.Extensions.Helper;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Lck.Utility.Extensions.Helper
{
    public class CacheRemoveHelper
    {

        public static void WhenLoginAdminSys(int userID , int roleID , List<string> tokens)
        {
            string cacheKey1 = CacheKey.DB_Data_User(userID);
            string cacheKey2 = CacheKey.DB_Data_RolePermissions(roleID);
            string cacheKey3 = CacheKey.AdminSysAPI_UserMerchantIDs(userID);

            CacheHelper.RemoveCacheByKey(cacheKey1);
            CacheHelper.RemoveCacheByKey(cacheKey2);
            CacheHelper.RemoveCacheByKey(cacheKey3);

            foreach (var tokenHash in tokens)
            {
                string cacheKeyToken = [REDACTED];
                CacheHelper.RemoveCacheByKey(cacheKeyToken);
            }
        }


        public static void WhenUpdateUser(int userID)
        {
            string cacheKey1 = CacheKey.DB_Data_User(userID);

            CacheHelper.RemoveCacheByKey(cacheKey1);
        }

        public static void WhenCreateMerchant(int userID)
        {
            List<string> cacheKeys = new List<string>
            {                  
                CacheKey.AdminSysAPI_UserMerchantIDs(userID)
            };

            RemoveCacheByKeys(cacheKeys);
        }


        public static void WhenUpdateMerchant(int merchantID)
        {
            List<string> cacheKeys = new List<string>
            {
                CacheKey.DB_Data_Merchant(merchantID),
                CacheKey.PointAPI_PageLayoutSetting(merchantID)
            };

            RemoveCacheByKeys(cacheKeys);
        }


        public static void WhenUpdateUserMerchantMappings(int userID)
        {
            List<string> cacheKeys = new List<string>
            {
                CacheKey.AdminSysAPI_UserMerchantIDs(userID),
            };

            RemoveCacheByKeys(cacheKeys);
        }



        public static void WhenUpdateSystemSettingConfig(int settingKey)
        {
            List<string> cacheKeys = new List<string>
            {
                CacheKey.PointAPI_ShoppingIframeUrl(), //需要順手刪除這個
                CacheKey.PointAPI_ShoppingStore(),     //需要順手刪除這個
                CacheKey.DB_Data_SystemSettingConfig(settingKey),
            };

            RemoveCacheByKeys(cacheKeys);
        }



        private static void RemoveCacheByKeys(List<string> cacheKeys)
        {
            foreach(var key in cacheKeys)
            {
                CacheHelper.RemoveCacheByKey(key);
            }
        }


    }

}
