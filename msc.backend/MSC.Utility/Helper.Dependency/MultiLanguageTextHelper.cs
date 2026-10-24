using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Collections.Specialized;
using MSC.Utility.Extensions;
using example.invalid;
using example.invalid;
using Microsoft.AspNetCore.Http;
using System.Reflection;
using System.Runtime.Versioning;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models;
using MSC.Models.CultureText;

namespace MSC.Utility.Helper
{

    /// <summary>
    /// 未來要多語系的話，比較容易統一在此處理
    /// </summary>
    public static class MultiLanguageTextHelper
    {

        private static Dictionary<LanguageEnum, Dictionary<string, string>> LanguageTextMappings;


        private static List<string> ColumnNameMappings = new List<string>
        {
            i18nKeys.Common_Label_GoogleAuthCode,
            i18nKeys.Common_Label_Account,
            i18nKeys.Common_Label_UserName,
            i18nKeys.Common_Label_DisplayName,
            i18nKeys.Common_Label_Password,
            i18nKeys.Common_Label_OldPassword,
            i18nKeys.Common_Label_NewPassword,
            i18nKeys.Common_Label_Status,
            i18nKeys.Common_Label_Amount,
            i18nKeys.Common_Label_Point,
        };
        

        static MultiLanguageTextHelper()
        {
            LanguageTextMappings = new Dictionary<LanguageEnum, Dictionary<string, string>>();

            //自幹的多語系處理(有點瞎)、系統啟動時，以static的方式存起來，可以減少舊版寫法的一些效能運算
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
        public static string GetErrorCodeMsg(ErrorCode errorCode, ErrorDetailModel errorDetail)
        {
            List<LanguageEnum> targetLanguages = GetMultiTargetLanguagesByHeaders();

            string fieldName = "ErrorCode_" + errorCode.ToString();

            string msg = GetCultureTextByLangClassAndColumnName(fieldName, targetLanguages);

            if (errorCode.ToString().Contains("InvalidParameter"))
            {
                if (!string.IsNullOrWhiteSpace(errorDetail?.ColumnName))
                {
                    string displayName = GetColumnDisplayName(errorDetail?.ColumnName);
                    msg = msg + $" ({displayName})";
                }
            }

            //如果有指定要強制覆蓋msg
            if (errorDetail != null)
            {
                if (!string.IsNullOrWhiteSpace(errorDetail.ReplaceMessage))
                {
                    msg = errorDetail.ReplaceMessage;
                }
            }

            if (string.IsNullOrWhiteSpace(msg))
                msg = errorCode.ToString() + $" ({errorCode.ToIntegerString()})";

            return msg;
        }


        public static string GetColumnDisplayName(string columnName)
        {
            if (string.IsNullOrWhiteSpace(columnName))
                return string.Empty;

            string toLower = columnName?.ToLower();

            if (toLower == "merchantid")
                return GetCommon_i18nText(i18nKeys.Common_Label_Merchant);

            //嘗試找出有定義好、要轉換名稱的欄位
            var match = ColumnNameMappings.FirstOrDefault(x => x.Replace("Common_Label_", "").ToLower() == toLower);

            if (match != null)
                return GetCommon_i18nText(match);
            
            return columnName;
        }

        public static string GetCommon_i18nText(string i18nKeys)
        {
            List<LanguageEnum> targetLanguages = GetMultiTargetLanguagesByHeaders();
            string text = GetCultureTextByLangClassAndColumnName(i18nKeys, targetLanguages);
            return text;
        }


        public static string GetEnums_i18nText(this System.Enum enumVal)
        {
            string fieldName = Geti18nFieldName(enumVal);
            List<LanguageEnum> targetLanguages = GetMultiTargetLanguagesByHeaders();

            string i18nKeys = "Enums_" + fieldName;
            string text = GetCultureTextByLangClassAndColumnName(i18nKeys, targetLanguages);

            if (string.IsNullOrWhiteSpace(text))
                return i18nKeys;

            return text;
        }
        

        private static string Geti18nFieldName(this System.Enum enumVal)
        {
            string enumName = enumVal.GetType().Name;
            string value = enumVal.ToString();
            string fieldName = $"{enumName}_{value}";
            return fieldName;
        }


        public static string DefaultErrorMsg(string rndText = "")
        {
            var langs = GetMultiTargetLanguagesByHeaders();

            if (!string.IsNullOrWhiteSpace(rndText))
            {
                string culText = GetCultureTextByLangClassAndColumnName(i18nKeys.Common_DefaultErrorMsg, langs);
                rndText = $"({culText}: {rndText}) ";
            }

            string templateText = GetCultureTextByLangClassAndColumnName(i18nKeys.Common_DefaultErrorMsg, langs);
            return templateText + rndText;
        }        


        private static List<LanguageEnum> GetMultiTargetLanguagesByHeaders()
        {
            List<LanguageEnum> targetLanguages = new List<LanguageEnum>();

            var currentLang = HeaderValueHelper.GetCurrentLanguage();

            targetLanguages.Add(currentLang);

            LanguageEnum systemDefaultLang = SystemConstant.DefaultLanguageEnum;

            //若主要語系不為簡體中文，會設第二順位的語系
            if (currentLang != systemDefaultLang)
            {
                targetLanguages.Add(systemDefaultLang);
            }

            return targetLanguages;
        }


        private static string GetCultureTextByLangClassAndColumnName(string fieldName, List<LanguageEnum> languages)
        {
            foreach (LanguageEnum targetLang in languages)
            {
                Dictionary<string, string> textDic = DeterminTextInstanceByLanguage(targetLang);

                //從指定語系的class抓
                string text = textDic.GetDictionaryOrEmptyValue(fieldName);

                //(防呆) 如果有設置文案，返回。若無設置文案 嘗試抓下一個語系的資料
                if (!string.IsNullOrWhiteSpace(text))
                {
                    return text;
                }
            }

            return string.Empty;
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


        private static Dictionary<string, string> DeterminTextInstanceByLanguage(LanguageEnum targetLang)
        {
            if (LanguageTextMappings.ContainsKey(targetLang))
                return LanguageTextMappings[targetLang];

            //預設語系
            return LanguageTextMappings[LanguageEnum.en_us];
        }


    }
}
