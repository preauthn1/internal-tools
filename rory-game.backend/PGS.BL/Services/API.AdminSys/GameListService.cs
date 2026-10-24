using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.BL.Interface;
using target.Models.Enums;
using target.Models.Models;
using System.Linq;
using target.Utility.Helper;
using target.Models;
using target.DA.Repository;
using Lck.Utility.Extensions.Helper;

namespace target.BL.Services
{

    public class GameListService : BaseAdminSysService, IGameListService
    {
        public GameListRepository GameListRepository { get; set; }

        public LanguageTextRepository LanguageTextRepository { get; set; }

        public SystemConfigSettingRepository SystemConfigSettingRepository { get; set; }


        public PagedList<GameListViewModel> GetByPagingFilter(QueryGameListApiRequestParameter filter)
        {
            filter.GameVendorTypes = filter.GameVendorType.ToIntegerListCollection();
            filter.IsNotDelete = true;

            var query = GameListRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<GameListEntity, GameListViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            SetExtraProperty(viewModels.Items);

            //如果是一次撈超過100筆，採用複雜order by 條件 (用底層的 orderByColumn 不太好寫)
            if (viewModels.Items.IsNotEmptyNull() && viewModels.Items.Count() > 100)
            {
                viewModels.Items = viewModels.Items
                    .OrderBy(x => AppHelper.IsExternalVendor(x.GameVendorType) ? 999 : 1)
                    .ThenByDescending(x => x.GameID)
                    .ToList();
            }

            return viewModels;
        }



        public GameListViewModel GetByID(int targetID)
        {
            var queryGameList = GameListRepository.GetById(targetID);

            if (queryGameList == null)
                throw new TerproException(ApiErrorCode.NotFoundGame);

            var convert = ConvertToViewModels(queryGameList);
            return convert;
        }


        public void Create(GameListCreateRequestModel model)
        {
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameCode).Rule_NotEmptyText();
            model.RuleFor(x => x.GameName).Rule_NotEmptyText();

            var isExist = GameListRepository.IsExist(model.GameID);

            if (isExist)
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_DuplicateDataExist);

            //檢查gameCode
            CheckGameCodeExist(model.GameCode);

            string languageCode = AppHelper.GetGameListLanguageCode(model.GameCode);


            var newEntity = model.ReturnInjecteObject<GameListEntity>();
            newEntity.CreateTime = DateTimeHelper.GetCurrentTime();
            newEntity.LanguageCode = languageCode;


            var newGameListID = GameListRepository.Create(newEntity);


            //產生預設幾種多語系
            List<LanguageEnum> targetLanguages = new List<LanguageEnum>
            {
                LanguageEnum.zh_cn,
                LanguageEnum.zh_tw,
                LanguageEnum.en_us
            };

            foreach (var lang in targetLanguages)
            {
                string gameName = model.GameName;

                if (lang == LanguageEnum.zh_cn)
                    gameName = AppHelper.TraditionalToSimplified(model.GameName);
                else
                    gameName = AppHelper.SimplifiedToTraditional(model.GameName);

                var instance = new LanguageTextEntity
                {
                    LanguageCode = languageCode,
                    Text = gameName,
                    Language = lang.ToString()
                };

                LanguageTextRepository.Create(instance);
            }

        }

        private void CheckGameCodeExist(string gameCode)
        {
            var queryGameCode = GameListRepository.GetByFilter(new QueryGameListRequestViewModel
            {
                GameCode = gameCode
            });

            if (queryGameCode.IsNotEmptyNull())
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_DuplicateDataExist);
        }


        public void Update(GameListUpdateRequestModel model)
        {
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameCode).Rule_NotEmptyText();
            model.RuleFor(x => x.GameName).Rule_NotEmptyText();

            var queryExist = GameListRepository.GetById(model.GameID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundGame);

            queryExist.GameName = model.GameName;
            queryExist.GameCode = model.GameCode;
            queryExist.Status = model.Status;
            queryExist.GameType = model.GameType;
            queryExist.GameVendorType = model.GameVendorType;

            GameListRepository.Update(queryExist);

            CacheRemoveHelper.RemoveCache_WhenUpdateGame(queryExist.GameID , queryExist.GameVendorType , queryExist.ThirdPartyGameCode);
        }


        public void Delete(int targetID)
        {
            var queryExist = GameListRepository.GetById(targetID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundGame);

            queryExist.Status = NormalActiveStatus.PermanentDelete.ToInteger();
            GameListRepository.Update(queryExist);

            CacheRemoveHelper.RemoveCache_WhenUpdateGame(queryExist.GameID, queryExist.GameVendorType, queryExist.ThirdPartyGameCode);
        }


        private GameListViewModel ConvertToViewModels(GameListEntity queryGame)
        {
            var convert = queryGame.ReturnInjecteObject<GameListViewModel>();
            convert.GameTypeText = MultiLanguageTextHelper.GetEnumsCultureText((GameTypeEnum)convert.GameType);
            convert.GameVendorTypeText = MultiLanguageTextHelper.GetEnumsCultureText((GameVendorTypeEnum)convert.GameVendorType);
            convert.IsMultiLanguageDefined = true;
            convert.IsGameListEnumDefined = true; //三方遊戲不需要定義

            var isThirdPartyGame = AppHelper.IsExternalVendorGame(queryGame.GameID);

            convert.IsThirdPartyGame = isThirdPartyGame;

            if (!isThirdPartyGame)
            {
                var i18nResource = DbMultiLanguageTextHelper.GetGameDisplayNameAndGetSource(queryGame.GameID);
                convert.IsMultiLanguageDefined = i18nResource.IsLanguageResourceDefined;
                convert.IsGameListEnumDefined = ExtensionOfEnum.IsDefinedValue((GameListEnum)queryGame.GameID);
            }

            return convert;
        }

        private void SetExtraProperty(List<GameListViewModel> items)
        {
            var imageSetting = SystemConfigSettingRepository.GetParseModelByCache(SystemConfigSettingEnum.LobbyPageGameImageUrl);
            string imageDomain = imageSetting.SettingValue;

            foreach (var item in items)
            {
                item.GameImageUrl = AppHelper.GenerateGameImageUrl(imageDomain, item.GameCode);
            }
        }

    }
}
