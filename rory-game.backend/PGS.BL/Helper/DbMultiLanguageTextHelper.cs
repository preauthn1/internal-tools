using target.Models;
using System;
using System.Collections.Generic;
using System.Text;
using target.Utility.Extensions;
using target.Models.Models;
using target.Models.Enums;
using target.Models.Table;
using target.BL.Interface;
using target.Utility.Helper;
using Lck.Utility.Extensions.Helper;

namespace target.BL.Services
{
    public class DbMultiLanguageTextHelper
    {
      
        public static string GetGameDisplayName(int gameID)
        {
            var sourceResult = GetGameDisplayNameAndGetSource(gameID);

            if (sourceResult.IsLanguageResourceDefined)
                return sourceResult.GameName;

            return sourceResult.DBGameName;
        }


        /// <summary>
        /// 1.  從多語系檔取得 遊戲名稱
        /// 2.  從cache(DB設置) 取得遊戲名稱
        /// </summary>
        public static (string GameName , string DBGameName , bool IsLanguageResourceDefined) GetGameDisplayNameAndGetSource(int gameID)
        {
            var gameEnum = (GameListEnum)gameID;

            var languageText = MultiLanguageTextHelper.GetEnumsCultureText(gameEnum);

            string fieldName = MultiLanguageTextHelper.GetEnums_i18n_FieldName(gameEnum);

            bool isResourceTextUndefined = fieldName == languageText || string.IsNullOrWhiteSpace(languageText);

            if (!isResourceTextUndefined)
                return (languageText , "" , true);

            //嘗試從DB抓資料
            string gameDisplayName = DbValueCacheHelper.GetGameNameFromLanguageText(gameID);
            return ("" , gameDisplayName , false);
        }



    }
}
