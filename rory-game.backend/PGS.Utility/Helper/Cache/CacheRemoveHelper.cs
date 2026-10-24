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

        /// <summary>
        /// 後台修改用戶狀態
        /// </summary>
        public static void RemoveCache_WhenUpdatePlayer(int playerID)
        {
            CacheHelper.RemoveCacheByKey(CacheKey.AdminSys_PlayerBriefInfo(playerID));
        }


        /// <summary>
        /// 後台用戶每次登入 刷新某些cache
        /// </summary>
        public static void RemoveCache_WhenUpdateSystemSetting(int settingKey)
        {
            CacheHelper.RemoveCacheByKey(CacheKey.DbEntity_ParseModel_SystemConfigSetting(settingKey));
        }


        /// <summary>
        /// 後台用戶每次登入 刷新某些cache
        /// </summary>
        public static void RemoveCache_WhenUserLogin(int userID)
        {
            CacheHelper.RemoveCacheByKey(CacheKey.AdminSys_UserOwnPermission(userID));
        }


        /// <summary>
        /// 後台用戶資訊 被修改
        /// </summary>
        public static void RemoveCache_WhenUpdateUser(int userID)
        {
            CacheHelper.RemoveCacheByKey(CacheKey.AdminSys_UserOwnPermission(userID));
            CacheHelper.RemoveCacheByKey(CacheKey.AdminSys_UserInfo(userID));
        }
        
        /// <summary>
        /// 修改角色權限時
        /// </summary>
        public static void RemoveCache_WhenRolePermissionUpdate(int roleID)
        {
            CacheHelper.RemoveCacheByKey(CacheKey.AdminSys_RolePermission(roleID.ToString()));
        }


        /// <summary>
        /// 修改 多語系文字
        /// </summary>
        public static void RemoveCache_WhenUpdateLanguageCode(string languageCode)
        {
            CacheHelper.RemoveCacheByKey(CacheKey.DbEntity_LanguageTexts(languageCode));
        }


        /// <summary>
        /// 修改伺服器
        /// </summary>
        public static void RemoveCache_WhenUpdateServer(int serverID)
        {
            CacheHelper.RemoveCacheByKey(CacheKey.AdminSys_ServerBriefInfo(serverID));
        }


        /// <summary>
        /// 修改遊戲
        /// </summary>
        public static void RemoveCache_WhenUpdateGame(int gameID , int vendorID , string thirdGameCode)
        {
            CacheHelper.RemoveCacheByKey(CacheKey.DbEntity_GameList(gameID));
            CacheHelper.RemoveCacheByKey(CacheKey.DbEntity_GameList_GameVendorCode(vendorID , thirdGameCode));

            //遊戲開關狀態可能會有有異動 刪除 VendorAPI
            string keyword = CacheKey.VendorAPI_GetMerchantEnableGames(0);

            var queryKeys = CacheHelper.GetCacheKeysByKeyword(keyword);

            foreach(string key in queryKeys)
            {
                CacheHelper.RemoveCacheByKey(key);
            }
        }


        /// <summary>
        /// 修改商戶
        /// </summary>
        public static void RemoveCache_WhenUpdateMerchant(int merchantID)
        {
            CacheHelper.RemoveCacheByKey(CacheKey.AdminSys_MerchantInfo(merchantID));
        }


        /// <summary>
        /// 修改商戶開通遊戲
        /// </summary>
        public static void RemoveCache_WhenUpdateMerchantGames(int merchantID)
        {
            CacheHelper.RemoveCacheByKey(CacheKey.VendorAPI_GetMerchantEnableGames(merchantID));
        }


        /// <summary>
        /// 修改商戶
        /// </summary>
        public static void RemoveCache_WhenUpdateMerchantSetting(int merchantID , int settingKey)
        {
            CacheHelper.RemoveCacheByKey(CacheKey.DbEntity_ParseModel_MerchantCustomizedSetting(merchantID , settingKey));
        }


        /// <summary>
        /// 修改組織
        /// </summary>
        public static void RemoveCache_WhenUpdateOrganization()
        {
            CacheHelper.RemoveCacheByKey(CacheKey.AdminSys_OrganizationsList());
        }


    }

}
