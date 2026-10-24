using Microsoft.AspNetCore.Mvc;
using target.Models.ViewModels.Base;
using System;
using System.Collections.Generic;
using System.Text;
using System.Linq;
using target.Utility.Extensions;
using target.Utility.Helper.Dependency;
using target.Models.Table;
using target.Utility;
using target.Models.Models;
using target.Models;
using target.Utility.Helper;
using System.Reflection;
using target.Web.Modules.Filter;
using target.Utility.Utility.ThirdParty;

namespace target.Utility.Extensions.Web
{
    public class ApiConfigDefineHelper
    {
        private static Dictionary<string, ApiConfigAnalyzeModel> _allApiActions = new Dictionary<string, ApiConfigAnalyzeModel>();

        private static List<string> _skipLogResponseApis = new List<string>();

        private static List<string> _skipLogRequestBodyApis = new List<string>();


        /// <summary>
        /// 專案啟動後，先解析所有API是否有標記apiConfig或某些特殊attribute
        /// </summary>
        /// <param name="targetProject"></param>
        public static void AnalyzeAllApiConfigDefines(Assembly targetProject)
        {
            var allController = targetProject.GetTypes()
                     .Where(type => typeof(ControllerBase).IsAssignableFrom(type))
                     .ToList();

            foreach(var controller in allController)
            {
                var methods = controller.GetMethods(BindingFlags.Instance | BindingFlags.DeclaredOnly | BindingFlags.Public)
                    .Where(m => !m.GetCustomAttributes(typeof(System.Runtime.CompilerServices.CompilerGeneratedAttribute), true).Any())
                    .ToList();

                foreach(var eachAction in methods)
                {
                    string controllerName = eachAction.DeclaringType.Name.Replace("Controller", "");
                    string actionUniqueName = $"{controllerName}_{eachAction.Name}";

                    var loggingAttributes = (LoggingPlayerID)eachAction.GetCustomAttributes(typeof(LoggingPlayerID), false).FirstOrDefault();

                    var configAttributes = (ApiConfig)eachAction.GetCustomAttributes(typeof(ApiConfig), false).FirstOrDefault();

                    ApiConfigAnalyzeModel apiConfig = InstanceDefaultApiConfigValue();

                    if (configAttributes != null)
                    {
                        //若有自訂配置 以自訂的attribute設定為優先
                        apiConfig = configAttributes.ReturnInjecteObject<ApiConfigAnalyzeModel>();
                    }

                    if (loggingAttributes != null)
                        apiConfig.IsLoggingPlayerID = loggingAttributes.IsEnable;


                    try
                    {
                        //如果你看到這邊拋錯誤 代表有API的method名稱一樣 請改成不一樣的名稱
                        _allApiActions.Add(actionUniqueName, apiConfig);
                    }
                    catch(Exception ex)
                    {
                        //好多次了，發個通知吧
                        MyTelegramUtility.SendMessageToChannelAsync(new TelegramMessageParameter
                        {
                            TelegramChannelEnum = TelegramChannelEnum.pgs_common_error,
                            Message = "AnalyzeAllApiConfigDefines_" + ex.Message
                        });

                        throw ex;
                    }


                }
            }

            return;
        }


        public static ApiConfigAnalyzeModel GetApiConfigSettingByCurrentRequest()
        {
            string uniqueRouteName = HeaderValueHelper.HeaderTemp_ApiUniqueRouteName();

            return GetApiConfigSettingByUniqueRouteName(uniqueRouteName);
        }


        public static ApiConfigAnalyzeModel GetApiConfigSettingByUniqueRouteName(string uniqueRouteName)
        {
            if (string.IsNullOrWhiteSpace(uniqueRouteName))
                return InstanceDefaultApiConfigValue();

            if (!_allApiActions.ContainsKey(uniqueRouteName))
                return InstanceDefaultApiConfigValue();

            return _allApiActions[uniqueRouteName];
        }


        private static ApiConfigAnalyzeModel InstanceDefaultApiConfigValue()
        {
            //這邊有定義 apiConfig各種預設值
            var defaultSetting = new ApiConfig();

            var convert = defaultSetting.ReturnInjecteObject<ApiConfigAnalyzeModel>();
            convert.IsLoggingPlayerID = false;

            return convert;
        }




    }
}
