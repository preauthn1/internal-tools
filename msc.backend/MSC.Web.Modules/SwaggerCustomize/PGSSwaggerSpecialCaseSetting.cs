using Microsoft.OpenApi.Any;
using Microsoft.OpenApi.Models;
using MSC.Models.CoreLibs;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using Swashbuckle.AspNetCore.SwaggerGen;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;


namespace MSC.Web.Modules
{

    public static class PGSSwaggerSpecialCaseSetting
    {

        /// <summary>
        /// 請求參數，若名稱符合、swagger UI上會顯示此處定義的欄位說明
        /// </summary>
        public static Dictionary<string, string> QueryParameterSampleValueDefined = new Dictionary<string, string>()
        {
            { "pageIndex" , "查第幾頁 (從1開始)" },
            { "pageSize" , "每頁幾筆" },
            { "orderByColumn" , "用哪個欄位排序" },
            { "isASC" , "正排序、倒排序" },
            { "signKey" , "加密驗簽，需將Ts以及部分欄位加密，以檢查防止API偽造" },
            { "ts" , "時間戳記timestamp" },
        };

        /// <summary>
        /// 回應參數，若名稱符合、swagger UI上會顯示此處定義的範例值
        /// </summary>
        public static Dictionary<string, string> OutputParameterSampleValueDefined = new Dictionary<string, string>()
        {
            { "lang" , "zh-tw" },
            { "defaultLanguage" , "zh-tw" },
            { "defaultLang" , "zh-tw" },
            { "currency" , "CNY" },
            { "memo" , "備註欄位" },
            { "createBy" , "adminUser" },
            { "updateBy" , "adminUser" },
            { "userName" , "用戶名／用戶帳號" },
            { "displayName" , "顯示名稱" },
            { "password" , "輸入您的密碼" },
            { "confirmPassword" , "再次確認您的密碼" },
            { "apiSecretKey" , "8acc72e29907" },
            { "apiDoaminUrl" , "http://0.0.0.0" },
            { "apiValidatePlayerSessionUrl" , "target/validatePlayerSession" },
            { "ipAddress" , "0.0.0.0" },
            { "userGameDetail" , "10001:2000:-10:20180401,10002:1000:-10:20180402" },
            { "serverName" , "多福多財伺服器_1.0.1" },
            { "merchantName" , "PGS娛樂投注平台" },
            { "organizationName" , "PGS天馬團隊" },
            { "gameName" , "DFDZ多福多財" },
            { "orderNo" , "2019010100751241" },
            { "merchantOrderNo" , "AZY201906553641" },
            { "playerName" , "player_account" },
        };


        /// <summary>
        /// 回應參數，若名稱符合、swagger UI上會顯示此處定義的範例值
        /// </summary>
        public static Dictionary<string, long> OutputParameterSampleNumberValueDefined = new Dictionary<string, long>()
        {
            { "gameID" , 25 },
            { "serverID" , 1 },
            { "tradeID" , 20191200000095331},
            { "playerID" , 1000003},
            { "faceID" , 15},
            { "gameTimeSecond" , 360},     
        };


        /// <summary>
        /// 回應參數，若名稱符合、swagger UI上會顯示此處定義的範例值
        /// </summary>
        public static Dictionary<string, decimal> OutputParameterSampleDecimalValueDefined = new Dictionary<string, decimal>()
        {
            { "currentAmount" , 5000.75M},
            { "updateAmount" , 2000.10M},
            { "flowAmount" , 100.50M},
            { "taxAmount" , 5.5M}
        };


        /// <summary>
        /// 這些變數不會產生文件 因為會有無窮迴圈
        /// </summary>
        public static HashSet<string> OutputParameterIgnoreOutputDocumentColumns = new HashSet<string>()
        {
            "childPermissions"
        };


        public static Dictionary<string, string> SkipOutputResponseColumnWhenPagedList = new Dictionary<string, string>();


        public static void AddWhenNotExist(string name)
        {
            name = name.ToLower();

            if (SkipOutputResponseColumnWhenPagedList.ContainsKey(name))
                return;

            SkipOutputResponseColumnWhenPagedList.Add(name, "");
        }



        public static bool IsKeyExist(string className , string columnName)
        {
            string keyName = $"{className.ToLower()}_{columnName.ToLower()}";

            if (SkipOutputResponseColumnWhenPagedList.ContainsKey(keyName))
                return true;

            return false;
        }

        public static bool IsOutputNumberColumnsExist(string columnName)
        {
            if (OutputParameterSampleNumberValueDefined.ContainsKey(columnName))
                return true;

            return false;
        }

        public static bool IsOutputDecimalColumnsExist(string columnName)
        {
            if (OutputParameterSampleDecimalValueDefined.ContainsKey(columnName))
                return true;

            return false;
        }



        public static bool IsOutputKeyExist(string columnName)
        {
            if (OutputParameterSampleValueDefined.ContainsKey(columnName))
                return true;

            return false;
        }

    }


}
