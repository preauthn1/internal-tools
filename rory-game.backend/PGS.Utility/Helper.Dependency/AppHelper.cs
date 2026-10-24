using target.Models;
using System;
using System.Collections.Generic;
using System.Text;
using target.Utility.Extensions;
using target.Models.Models;
using target.Models.Enums;
using target.Models.Table;
using Microsoft.International.Converters.TraditionalChineseToSimplifiedConverter;

namespace target.Utility.Helper
{
    public class AppHelper
    {

        /// <summary>
        /// 對接三方遊戲，並不支援BETA幣別，強制轉成某一種真實幣別
        /// </summary>
        /// <param name="currency"></param>
        /// <returns></returns>
        public static string ConvertCurrencyBETA_To_CNY(string currency)
        {
            if (currency == CurrencyEnum.BETA.ToString())
                return CurrencyEnum.CNY.ToString();

            return currency;
        }


        public static string GenerateSimpleTokenModelSignKey(SimpleExternalTokenModel model)
        {
            string beforeMD5 = $"{model.RequestID}_{SystemConstant.JwtToken_SaltKey}_{model.ExpiredTime.ToString("yyyy-MM-dd HH:mm")}_{model.TransactionNo}";
            return beforeMD5.MD5Hash();
        }

        public static bool IsChineseLang(LanguageEnum lang)
        {
            return lang == LanguageEnum.zh_cn || lang == LanguageEnum.zh_tw;
        }

        public static LanguageEnum ConvertLangToEnum(string lang , LanguageEnum defaultLang = LanguageEnum.en_us)
        {
            lang = lang.Replace("-", "_");

            return ExtensionOfEnum.GetEnumByTextValue<LanguageEnum>(lang, defaultLang);
        }

        public static string GenerateGameLaunchRequestUUID(long id)
        {
            return $"GL_" + id.ToString();
        }


        public static string GetGameListLanguageCode(string gameCode)
        {
            return $"GameList_{gameCode}";
        }


        public static string GeneratePasswordHash(string password)
        {
            return $"{SystemConstant.Password_SaltKey}_{password}".MD5Hash();
        }

        public static string GenerateGameImageUrl(string imageDomain , string pgs_game_code)
        {
            //imageDomain 是透過 SystemConfigSettingEnum.LobbyPageGameImageUrl 配置拿到
            return ($"{imageDomain}/{pgs_game_code.ToLower()}.png").ReplaceDuplicateUrlChar();
        }

        /// <summary>
        /// 轉成 時間號的table e.g. ApiRequestLog_20200101
        /// </summary>
        public static string ConvertToDateSequenceTableName(string logTableName)
        {
            string replaceTableName = logTableName.TruncateLastWordIfMatch(']');
            string tableName = replaceTableName + "_" + DateTimeHelper.GetCurrentTime().ToString("yyyyMMdd") + "]";
            return tableName;
        }


        /// <summary>
        /// 判斷是否為 外部對接的遊戲商遊戲
        /// </summary>
        public static bool IsExternalVendorGame(int gameID)
        {
            return gameID > 99999;
        }


        /// <summary>
        /// 判斷是否為 PGS自製的遊戲
        /// </summary>
        public static bool IsPgsSelfGame(int gameID)
        {
            return !IsExternalVendorGame(gameID);
        }


        /// <summary>
        /// 判斷是否為 外部對接的遊戲商
        /// </summary>
        public static bool IsExternalVendor(int vendorID)
        {
            if (vendorID == GameVendorTypeEnum.target.ToInteger() ||
                vendorID == GameVendorTypeEnum.Slot88.ToInteger())
            {
                return false;
            }

            return true;
        }


        /// <summary>
        /// 取得轉換比例的 配置enum編號
        /// </summary>
        public static (bool IsDefine , MerchantCustomizedSettingKeyEnum SettingKey) GetConvertRateCustomizedSettingKey(GameVendorTypeEnum vendorID)
        {
            if (vendorID == GameVendorTypeEnum.target)
                return (true , MerchantCustomizedSettingKeyEnum.GameVendor_ConvertRate_PGS);

            //編號的規則是 9000000 + 遊戲商ID
            int expectSettingID = 9000000 + vendorID.ToInteger();
            var convertToEnum = (MerchantCustomizedSettingKeyEnum)expectSettingID;

            if(!convertToEnum.IsDefinedValue())
                return (false, MerchantCustomizedSettingKeyEnum.GameVendor_ConvertRate_PGS);

            return (true, convertToEnum);
        }

        /// <summary>
        /// 簡體字轉繁體
        /// </summary>
        public static string SimplifiedToTraditional(string text)
        {
            return ChineseConverter.Convert(text, ChineseConversionDirection.SimplifiedToTraditional);//转繁体
        }

        /// <summary>
        /// 繁體轉簡體
        /// </summary>
        public static string TraditionalToSimplified(string text)
        {
            return ChineseConverter.Convert(text, ChineseConversionDirection.TraditionalToSimplified);//转繁体
        }


        /// <summary>
        /// 定義了 target.Backend.Utility專案中  哪些enums要自動產生多語系
        /// </summary>
        public static List<Type> GetAutoGeneratei18nKeysEnums()
        {
            return new List<Type>
            {
                typeof(MerchantCustomizedSettingKeyEnum),
                typeof(MerchantLoadingStyleModeEnum),
                typeof(MerchantTypeEnum),
                typeof(SystemConfigSettingEnum),
                typeof(SystemMaintainStatusEnum),
                typeof(DataCategoryEnum),
                typeof(UserOperateTypeEnum),
                typeof(NormalSuccessStatusEnum),
                typeof(VendorConfigSettingEnum),
                typeof(GameListEnum),
                typeof(GameTypeEnum),
                typeof(GameVendorTypeEnum),
                typeof(GameVendorBetRecordStatusEnum),
                typeof(TransactionTypeEnum),
            };
        }




    }
}
