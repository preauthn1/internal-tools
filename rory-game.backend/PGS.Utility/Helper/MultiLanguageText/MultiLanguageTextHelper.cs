using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Collections.Specialized;
using target.Utility.Extensions;
using example.invalid;
using example.invalid;
using Microsoft.AspNetCore.Http;
using System.Reflection;
using System.Runtime.Versioning;
using target.Models.Enums;
using target.Models.Models;
using target.Models.CultureText;
using target.Models;

namespace target.Utility.Helper
{

    /// <summary>
    /// 未來要多語系的話，比較容易統一在此處理
    /// </summary>
    public static class MultiLanguageTextHelper
    {

        private static Dictionary<LanguageEnum, Dictionary<string, string>> LanguageTextMappings;


        static MultiLanguageTextHelper()
        {
            LanguageTextMappings = new Dictionary<LanguageEnum, Dictionary<string, string>>();

            //自幹的多語系處理(後來覺得有點瞎)、系統啟動時，以static的方式存起來，可以減少舊版寫法的一些效能運算
            LanguageTextMappings.Add(LanguageEnum.en_us, GetTargetLanguageAllTextWhenInitialize(typeof(i18n_EN)));
            LanguageTextMappings.Add(LanguageEnum.zh_tw, GetTargetLanguageAllTextWhenInitialize(typeof(i18n_TW)));
            LanguageTextMappings.Add(LanguageEnum.zh_cn, GetTargetLanguageAllTextWhenInitialize(typeof(i18n_CN)));
        }


        /// <summary>
        /// 取得api errorCode的錯誤碼
        /// </summary>
        /// <param name="errorCode"></param>
        /// <param name="errorDetail"></param>
        /// <returns></returns>
        public static string GetErrorCodeMsg(ApiErrorCode errorCode , ErrorDetailModel errorDetail)
        {
            List<LanguageEnum> targetLanguages = GetRequiredLanguages();

            string fieldName = "ErrorCode_" + errorCode.ToString();

            string msg = GetCultureTextByTargetClassAndColumnName(fieldName, targetLanguages);
                       
            //如果有指定columnName
            if (!string.IsNullOrWhiteSpace(errorDetail?.ColumnName))
            {
                string columnNameDisplayName = GetColumnDisplayName(errorDetail?.ColumnName);
                msg = msg + $" ({columnNameDisplayName})";
            }            

            //如果有指定要強制覆蓋msg
            if (errorDetail != null)
            {
                if (!string.IsNullOrWhiteSpace(errorDetail.ReplaceMessage))
                {
                    msg = errorDetail.ReplaceMessage;
                }
            }

            return msg;
        }


        private static string GetColumnDisplayName(string columnName)
        {
            columnName = columnName.ToLower();

            //TOCO 透過swtich case 轉成各欄位的顯示名稱 (如果直接顯示columnName(變數名) 給客戶看，不適當的話)
            return columnName;
        }


        public static string GetGameVendorDisplayName(int vendorID)
        {
            return ((GameVendorTypeEnum)vendorID).ToString();
        }


        public static string GetEnumsCultureText(System.Enum enumVal)
        {
            string fieldName = GetEnums_i18n_FieldName(enumVal);
            List<LanguageEnum> targetLanguages = GetRequiredLanguages();
            string text = GetCultureTextByTargetClassAndColumnName(fieldName, targetLanguages);
            return text;
        }


        public static string GetEnumsCultureTextByLanguage(System.Enum enumVal , LanguageEnum lang)
        {
            string fieldName = GetEnums_i18n_FieldName(enumVal);
            string text = GetCultureTextByTargetClassAndColumnName(fieldName, new List<LanguageEnum> { lang });
            return text;
        }

      
        public static string GetTimeZoneDisplayName(int timeZoneAreaID)
        {
            var name = GetEnumsCultureText((TimeZoneAreaEnum)timeZoneAreaID);
            return name;
        }

        public static string DefaultErrorMsg(string rndText = "")
        {
            if (!string.IsNullOrWhiteSpace(rndText))
            {
                string culText = GetCommonCultureTextByClassAndColumnName(i18nKeys.Common_System_ErrorMsg_ErrorSN);
                rndText = $"({culText}: {rndText}) ";
            }

            string templateText = GetCommonCultureTextByClassAndColumnName(i18nKeys.Common_System_ErrorMsg_Defaul);
            return templateText + rndText;
        }


        public static string GameServerReconnectMsg(string roomName)
        {
            string templateText = GetCommonCultureTextByClassAndColumnName(i18nKeys.Common_System_ErrorMsg_ReconnetToServerRoomIn);
            return string.Format(templateText, roomName);
        }


        public static string EnterGameRoom_AmountNotEnough(decimal limitAmount, decimal playerAmount)
        {
            string templateText = GetCommonCultureTextByClassAndColumnName(i18nKeys.Common_Global_ErrorMsg_MoneyNotEnough);
            return string.Format(templateText, limitAmount , playerAmount);
        }

        
        public static string GetCommonCultureTextByClassAndColumnName(string fieldName)
        {
            List<LanguageEnum> targetLanguages = GetRequiredLanguages();

            string text = GetCultureTextByTargetClassAndColumnName(fieldName, targetLanguages);

            return text;
        }


        public static string GetEnums_i18n_FieldName(this System.Enum enumVal)
        {
            string enumName = enumVal.GetType().Name;
            string value = enumVal.ToString();
            string fieldName = $"Enums_{enumName}_{value}";
            return fieldName;
        }

        

        public static List<LanguageEnum> GetRequiredLanguages()
        {
            var currentLang = HeaderValueHelper.GetCurrentLanguage();

            List<LanguageEnum> targetLanguages = new List<LanguageEnum>
            {
                currentLang
            };              
                        
            //若主要語系不為簡體中文，會設第二順位的語系
            if (currentLang != SystemConstant.DefaultLanguageEnum)
            {
                targetLanguages.Add(SystemConstant.DefaultLanguageEnum);
            }

            return targetLanguages;
        }


        private static Dictionary<string, string> GetTargetLanguageAllTextWhenInitialize(Type type)
        {
            Dictionary<string, string> eachi18nTexts = new Dictionary<string, string>();

            foreach (var eachField in type.GetFields())
            {
                string fieldName = eachField.Name;

                if (eachi18nTexts.ContainsKey(fieldName))
                {
                    continue;
                }
                var value = eachField.GetValue(null)?.ToString();
                eachi18nTexts.Add(fieldName, value);
            }

            return eachi18nTexts;
        }


        private static string GetCultureTextByTargetClassAndColumnName(string fieldName , List<LanguageEnum> languages) 
        {
            foreach(LanguageEnum targetLang in languages)
            {
                Dictionary<string ,string> textDic = DeterminTextInstanceByLanguage(targetLang);

                //從指定語系的class抓
                string text = textDic.GetDictionaryOrEmptyValue(fieldName);                

                //(防呆) 如果有設置文案，返回。若無設置文案 嘗試抓下一個語系的資料
                if (!string.IsNullOrWhiteSpace(text))
                {
                    return text;
                }
            }

            return fieldName;
        }


        private static Dictionary<string , string> DeterminTextInstanceByLanguage(LanguageEnum targetLang)
        {
            if (LanguageTextMappings.ContainsKey(targetLang))
                return LanguageTextMappings[targetLang];

            //預設語系
            return LanguageTextMappings[LanguageEnum.en_us];
        }


    }

}
