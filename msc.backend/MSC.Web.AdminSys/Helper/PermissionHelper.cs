using MSC.Models;
using System;
using System.Collections.Generic;
using MSC.Utility.Extensions;
using System.Linq;
using MSC.Utility.Helper;
using MSC.Models.CultureText;
using MSC.Utility;
using Lck.Utility.Extensions.Helper;
using MSC.Models.Enums;
using MSC.Utility.Helper.HttpContext;
using MSC.Models.Models;

namespace MSC.Web.AdminSys.Models
{
    public class PermissionHelper
    {

        public static List<PermissionFunctionParentModel> AllPermissions = new List<PermissionFunctionParentModel>();


        static PermissionHelper()
        {
            InitializeAllPermissions();
        }

        public static string GetCurrentPageTile()
        {
            string url = MyHttpContext.Current.Request.Path;
            var match = AllPermissions.SelectMany(z => z.Childs).FirstOrDefault(x => x.Url == url);
            return match?.Name;
        }


        public static List<PermissionFunctionParentModel> GetMyPermissions()
        {
            var myPermissions = GetPermissions();

            //如果有開發中的功能要顯示在menu中 可在這邊加
            if (SettingHelper.IsDev)
            {
                string parent_test = "Test_Function";

                var devFunctions = new PermissionFunctionParentModel
                {
                    Name = "開發中功能",
                    Icon = "fas fa-tools",
                    Url = "",
                    Token = [REDACTED],
                    Childs = new List<PermissionFunctionChildModel>
                    {
                        //new PermissionFunctionChildModel
                        //{
                        //    ParentToken = [REDACTED],
                        //    Name = "系統參數設置",
                        //    Url = "/System/SystemSettingConfig",
                        //    Token = [REDACTED]",
                        //},
                    }
                };

                if(devFunctions.Childs.IsNotEmptyNull())
                {
                    myPermissions.Add(devFunctions);
                }
            }


            return myPermissions;
        }



        public static List<PermissionFunctionParentModel> GetPermissions()
        {
            var user = AuthTokenHelper.GetUserModelByCookieToken();

            var myPermissions = new List<PermissionFunctionParentModel>();

            if (user.Permissions.IsEmptyOrNull())
                throw new MSCException(ErrorCode.Forbidden_PagePermissionNotAllow);

            foreach (var item in AllPermissions)
            {
                bool isExistPermission = user.Permissions.Contains(item.Token);

                if (!isExistPermission)
                    continue;

                var newInstance = item.ReturnInjecteObject<PermissionFunctionParentModel>();
                newInstance.Childs = new List<PermissionFunctionChildModel>();

                foreach (var child in item.Childs)
                {
                    bool isExistPermissionChild = user.Permissions.Contains(child.Token);

                    if (!isExistPermissionChild)
                        continue;

                    newInstance.Childs.Add(child);
                }

                myPermissions.Add(newInstance);
            }

            return myPermissions;
        }


        public static bool PermissionExist(AdminPermissionTokenEnum permissionToken)
        {
            var user = AuthTokenHelper.GetUserModelByCookieToken();
            return user.Permissions.Any(x => x == permissionToken.ToString());
        }


        public static void CheckPermissionAvailable(AdminPermissionTokenEnum permissionToken)
        {
            var user = AuthTokenHelper.GetUserModelByCookieToken();
            bool isAny = user.Permissions.Any(x => x == permissionToken.ToString());

            if (!isAny)
            {
                throw new MSCException(ErrorCode.Forbidden_PagePermissionNotAllow);
            }
        }


        private static void InitializeAllPermissions()
        {
            var parents = new List<PermissionFunctionParentModel>();


            //卡號管理
            parents.Add(new PermissionFunctionParentModel
            {
                Name = t.Text(i18nKeys.Common_Permission_Card_Management),
                Icon = "fas fa-book",
                Url = "",
                Token = [REDACTED],
                Childs = new List<PermissionFunctionChildModel>
                    {
                        InstancePermission(AdminPermissionTokenEnum.Card_Management.ToString() , "/Card/List" , AdminPermissionTokenEnum.Card_Card_List.ToString()),
                        InstancePermission(AdminPermissionTokenEnum.Card_Management.ToString() , "/Card/GenerateCardNo" , AdminPermissionTokenEnum.Card_Card_GenerateCardNo.ToString()),
                    }
            });


            //商戶管理
            string parent_merchant = AdminPermissionTokenEnum.Merchant_Management.ToString();
            parents.Add(new PermissionFunctionParentModel
            {
                Name = t.Text(i18nKeys.Common_Permission_Merchant_Management),
                Icon = "fas fa-user-friends",
                Url = "",
                Token = [REDACTED],
                Childs = new List<PermissionFunctionChildModel>
                    {
                        InstancePermission(parent_merchant , "/Merchant/List" , AdminPermissionTokenEnum.System_Merchant_List.ToString()),
                        InstancePermission(parent_merchant , "/Merchant/Summary" , AdminPermissionTokenEnum.Merchant_MerchantSummary.ToString()),
                        InstancePermission(parent_merchant , "/Record/MerchantWalletLogs" , AdminPermissionTokenEnum.Record_MerchantWalletLog_List.ToString()),
                    }
            });


            //歷史紀錄查看
            string parent_history = AdminPermissionTokenEnum.History_Management.ToString();
            parents.Add(new PermissionFunctionParentModel
            {
                Name = t.Text(i18nKeys.Common_Permission_History_Management),
                Icon = "fas fa-history",
                Url = "",
                Token = [REDACTED],
                Childs = new List<PermissionFunctionChildModel>
                {
                     InstancePermission(parent_history , "/Record/CardLinkRequestRecordList" , AdminPermissionTokenEnum.Record_MerchantLinkRequest_List.ToString()),
                     InstancePermission(parent_history , "/Record/CardConsumeRecordList" , AdminPermissionTokenEnum.Record_CardConsume_List.ToString()),
                     InstancePermission(parent_history , "/Record/CardConsumeMerchantNotifyQueueList" , AdminPermissionTokenEnum.Record_CardConsumeMerchantNotifyQueue_List.ToString()),
                     InstancePermission(parent_history , "/Record/UserOperateLogs" , AdminPermissionTokenEnum.Record_UserOperateLogs_List.ToString()),
                     InstancePermission(parent_history , "/Record/RequestForbiddenLogList" , AdminPermissionTokenEnum.Record_RequestForbiddenLog_List.ToString()),
                     InstancePermission(parent_history , "/Merchant/Withdrawal" , AdminPermissionTokenEnum.Merchant_Withdrawal_List.ToString()),
                }
            });



            //報表類
            string parent_report = AdminPermissionTokenEnum.Report_Management.ToString();
            parents.Add(new PermissionFunctionParentModel
            {
                Name = t.Text(i18nKeys.Common_Permission_Report_Management),
                Icon = "fas fa-chart-bar",
                Url = "",
                Token = [REDACTED],
                Childs = new List<PermissionFunctionChildModel>
                    {
                        InstancePermission(parent_report , "/Report/DailySummary" , AdminPermissionTokenEnum.Report_DailySummary_List.ToString()),
                    }
            });



            //系統設置類
            string parent_system = AdminPermissionTokenEnum.System_Management.ToString();
            parents.Add(new PermissionFunctionParentModel
            {
                Name = t.Text(i18nKeys.Common_Permission_System_Management),
                Icon = "fas fa-cog",
                Url = "",
                Token = [REDACTED],
                Childs = new List<PermissionFunctionChildModel>
                    {
                        InstancePermission(parent_system , "/User/List" , AdminPermissionTokenEnum.System_User_List.ToString()),
                        InstancePermission(parent_system , "/InternalOperation/ExternalLink" , AdminPermissionTokenEnum.System_InternalUtility.ToString()),
                        InstancePermission(parent_system , "/System/SystemSettingConfig" , AdminPermissionTokenEnum.System_SystemSettingConfig_List.ToString()),
                    }
            });


            AllPermissions = parents;
        }


        private static PermissionFunctionChildModel InstancePermission(string parent, string url, string permissionToken)
        {
            string i18nKeys = "Common_Permission_" + permissionToken;
            string i18nText = MultiLanguageTextHelper.GetCommon_i18nText(i18nKeys);

            if (string.IsNullOrWhiteSpace(i18nText))
                i18nText = i18nKeys;

            return new PermissionFunctionChildModel
            {
                ParentToken = [REDACTED],
                Name = i18nText,
                Url = url,
                Token = [REDACTED],
            };
        }


    }


    public class PermissionFunctionParentModel
    {
        public string Icon { get; set; }

        public string Token { get; set; }

        public string Name { get; set; }

        public string Url { get; set; }

        public List<PermissionFunctionChildModel> Childs { get; set; }

    }

    public class PermissionFunctionChildModel
    {
        public string ParentToken { get; set; }

        public string Token { get; set; }

        public string Name { get; set; }

        public string Url { get; set; }
    }






}
