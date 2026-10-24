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
        public static string VendorAPI_GetMerchantEnableGames(int merchantID)
        {
            string suffix = merchantID == 0 ? "" : merchantID.ToString();
            return $"VendorAPI_GetMerchantEnableGames_{suffix}";
        }

        public static string AdminSys_OrganizationsList()
        {
            return $"AdminSys_OrganizationsList";
        }

        public static string AdminSys_UserOwnPermission(int userID)
        {
            return $"AdminSys_UserOwnPermission_{userID}";
        }

        public static string AdminSys_RolePermission(string roleIdsText)
        {
            return $"AdminSys_RolePermission{roleIdsText}";
        }

        public static string AdminSys_UserInfo(int userID)
        {
            return $"AdminSys_UserInfo{userID}";
        }

        public static string AdminSys_MerchantInfo(int merchantId)
        {
            return $"AdminSys_MerchantInfo{merchantId}";
        }

       
        public static string AdminSys_PlayerBriefInfo(int playerID)
        {
            return $"AdminSys_PlayerBrief{playerID}";
        }

        public static string AdminSys_ServerBriefInfo(int serverID)
        {
            return $"AdminSys_ServerBriefInfo{serverID}";
        }

        public static string DbEntity_VendorSessionToken(string token)
        {
            return $"DbEntity_VendorSessionToken_{token}";
        }

        public static string DbEntity_Player(int playerID)
        {
            return $"DbEntity_Player{playerID}";
        }

        public static string DbEntity_ParseModel_MerchantCustomizedSetting(int merchantID , int settingKey)
        {
            return $"DbEntity_ParseModel_MerchantCustomizedSetting{merchantID}_{settingKey}";
        }

        public static string DbEntity_GameList(int gameID)
        {
            return $"DbEntity_GameList{gameID}";
        }

        public static string DbEntity_GameList_GameVendorCode(int vendorID , string thirdGameCode)
        {
            return $"DbEntity_GameList_GameVendorCode{vendorID}_{thirdGameCode}";
        }

        public static string DbEntity_ParseModel_SystemConfigSetting(int settingID)
        {
            return $"DbEntity_SystemConfigSetting{settingID}";
        }

        public static string DbEntity_LanguageTexts(string languageCode)
        {
            return $"DbEntity_LanguageTexts{languageCode}";
        }

        public static string DbEntity_VendorMaintain_VendorID(int vendorID)
        {
            return $"DbEntity_VendorMaintain_VendorID{vendorID}";
        }


        public static string ExternalVendor_VendorConfigSettings(int merchantID , int vendorID)
        {
            return $"ExternalVendor_VendorConfigSettings_{merchantID}_{vendorID}";
        }




        public static string Job_JobSettingControl(int jobType)
        {
            return $"Job_JobSettingControl_{jobType}";
        }


    }


}
