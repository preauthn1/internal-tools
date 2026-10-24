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
using target.DA.Repository;
using System.Linq;

namespace target.BL.Services
{
    /// <summary>
    /// 如果常常用到 又需要存取db的cache 可以考慮放這邊統一呼叫
    /// </summary>
    public static class DbValueCacheHelper
    {
        private static GameListRepository GameListRepository { get; set; }

        private static LanguageTextRepository LanguageTextRepository { get; set; }
        

        static DbValueCacheHelper()
        {
            GameListRepository = new GameListRepository();
            LanguageTextRepository = new LanguageTextRepository();
        }

        public static GameListEntity GetGame(int gameID)
        {
            return GameListRepository.GetSingleByCache(gameID);
        }


        public static List<LanguageTextEntity> GetLanguageTextsByCache(string languageCode)
        {
            return LanguageTextRepository.GetByLanguageCodeByCache(languageCode);
        }


        /// <summary>
        /// 抓出遊戲的顯示名稱 (db存的)
        /// </summary>
        public static string GetGameNameFromLanguageText(int gameID)
        {
            if (gameID <= 0)
                return string.Empty;
           
            var queryGame = GetGame(gameID);

            if (queryGame == null)
                return $"undefined_{gameID}";

            string languageCode = queryGame.LanguageCode;

            if (string.IsNullOrWhiteSpace(languageCode))
                return queryGame.GameName;

            //從 LanguageText 資料表抓資料
            string gameNameText = GetTargetTextByLanguageTextTable(languageCode);

            if (string.IsNullOrWhiteSpace(gameNameText))
                return queryGame.GameName;

            return gameNameText;
        }


        private static string GetTargetTextByLanguageTextTable(string languageCode)
        {
            List<LanguageEnum> languages = MultiLanguageTextHelper.GetRequiredLanguages();

            //查db的部分
            var queryTexts = GetLanguageTextsByCache(languageCode);

            foreach (var lan in languages)
            {
                string targetLang = lan.ToString();
                var targetWord = queryTexts.FirstOrDefault(x => x.Language == targetLang);

                if (targetWord != null && targetWord.Text.IsNotNullEmptyText())
                {
                    return targetWord.Text;
                }
            }

            return string.Empty;
        }


    }
}
