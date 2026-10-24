using Lck.Utility.Extensions.Helper;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Lck.Utility.Extensions.Helper
{
    public class CacheKey
    {
        public static string DB_Data_Merchant(int merchantID)
        {
            return $"DB_Data_Merchant{merchantID}";
        }

        public static string DB_Data_MerchantCustomizedSetting(int merchantID, int key)
        {
            return $"DB_Data_MerchantCustomizedSetting{merchantID}_{key}";
        }

        public static string DB_Data_SystemSettingConfig(int settingKey)
        {
            return $"DB_Data_SystemSettingConfig{settingKey}";
        }

        public static string DB_Data_User(int userID)
        {
            return $"DB_Data_UserID{userID}";
        }

        public static string DB_Data_UserToken(string tokenHash)
        {
            return $"DB_Data_UserToken{tokenHash}";
        }


        public static string DB_Data_RolePermissions(int roleID)
        {
            return $"DB_Data_RolePermissions{roleID}";
        }



        public static string AdminSysAPI_UserMerchantIDs(int userID)
        {
            return $"AdminSysAPI_UserMerchantIDs{userID}";
        }

        public static string PointAPI_TokenExpired(string tokenHash)
        {
            return $"PointAPI_TokenExpired{tokenHash}";
        }

        public static string PointAPI_ShoppingStore()
        {
            return $"PointAPI_ShoppingStore";
        }

        public static string PointAPI_PageLayoutSetting(int merchantID)
        {
            return $"PointAPI_PageLayoutSetting_{merchantID}";
        }

        public static string PointAPI_ShoppingIframeUrl()
        {
            return $"PointAPI_ShoppingIframeUrl";
        }


        public static string PointAPI_TokenExist(string tokenHash)
        {
            return $"PointAPI_TokenExist{tokenHash}";
        }

        public static string VendorAPI_RequestForbidden(int merchantID , string memberCode)
        {
            return $"VendorAPI_RequestForbidden_{merchantID}_{memberCode}";
        }



    }


}
