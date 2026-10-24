using System;
using System.Collections.Generic;
using target.BL.Interface;
using target.BL.Services;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Models.Models;
using System.Linq;
using target.DA.Repository;
using target.Models.Enums;
using target.Utility.Helper;
using Lck.Utility.Extensions.Helper;

namespace target.BL.Services
{

    public class GameInfoService : BaseVendorService, IGameInfoService
    {

        public GameListRepository GameListRepository { get; set; }

        public MerchantGameSettingRepository MerchantGameSettingRepository { get; set; }

        public SystemConfigSettingRepository SystemConfigSettingRepository { get; set; }



        /// <summary>
        /// 取得所有遊戲列表
        /// </summary>
        public List<GameInfoViewModel> GetAllGames(GameInfoRequestModel model)
        {          
            //檢查商戶狀態合法
            this.GetMerchantAndCheckIsValid(model.MerchantID , model.MerchantSecretKey);
            
            string cacheKey = CacheKey.VendorAPI_GetMerchantEnableGames(model.MerchantID);

            var queryEnableGames = CacheHelper.GetCacheOrQuery(cacheKey, () =>
               {
                   //查遊戲資料
                   var queryAllGames = GameListRepository.GetAllStatusEnableGames();

                   var queryGameSettings = MerchantGameSettingRepository.GetByMerchantID(model.MerchantID, NormalActiveStatus.Enable);

                   List<int> merchantAvailableGames = queryGameSettings.Select(x => x.GameID).ToList();

                   queryAllGames = queryAllGames
                       .Where(x => merchantAvailableGames.Contains(x.GameID))
                       .ToList();

                   return queryAllGames;
               });

            var imageSetting = SystemConfigSettingRepository.GetParseModelByCache(SystemConfigSettingEnum.LobbyPageGameImageUrl);
            string imageDomain = imageSetting.SettingValue;

            return queryEnableGames
                .Select(x => new GameInfoViewModel
                {
                    GameID = x.GameID,
                    GameCode = x.GameCode,
                    GameName = DbMultiLanguageTextHelper.GetGameDisplayName(x.GameID),
                    GameType = x.GameType,
                    GameTypeName = MultiLanguageTextHelper.GetEnumsCultureText((GameTypeEnum)x.GameType),
                    ImageUrl = AppHelper.GenerateGameImageUrl(imageDomain, x.GameCode),
                })
                .ToList();

        }


    }
}
