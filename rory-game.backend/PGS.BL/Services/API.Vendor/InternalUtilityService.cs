using System;
using System.Linq;
using System.Collections.Generic;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.BL.Services;
using target.Models.Models;
using target.Models.ViewModels.Base;
using target.Utility.Helper.Dependency;

namespace target.BL.Interface
{

    public class InternalUtilityService : BaseGameService, IInternalUtilityService
    {

        public ServerListRepository ServerListRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }

        public MerchantGameSettingRepository MerchantGameSettingRepository { get; set; }

        public ServerMerchantMappingRepository ServerMerchantMappingRepository { get; set; }

        public SystemConfigSettingRepository SystemConfigSettingRepository { get; set; }

        public IPlayerWalletRefService PlayerWalletRefService { get; set; }

        public IPlayerJwtCheckService PlayerJwtCheckService { get; set; }


        public List<DemoGameCategoryViewModels> GetMerchantServers(int merchantID)
        {
            var gameCategories = ExtensionOfEnum.GetEnumAllValue<GameTypeEnum>();

            var queryAllGames = GameListRepository.GetByFilter(new QueryGameListRequestViewModel
            {
                Status = NormalActiveStatus.Enable.ToInteger()
            });

            var queryServers = ServerListRepository.GetByFilter(new QueryServerListRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection()
            });

            var convert = queryServers.Select(x =>
            {
                var parse = x.ReturnInjecteObject<DemoServerItemViewModels>();
                parse.GameType = queryAllGames.FirstOrDefault(z => z.GameID == x.GameID)?.GameType ?? 0;
                return parse;
            });


            var gameTypeGroup = convert.GroupBy(x => x.GameType).Select(z => z.Key).ToList();
            var gameIDGroup = convert.GroupBy(x => x.GameID).Select(z => z.Key).ToList();

            List<DemoGameItemViewModels> listOfEachGames = new List<DemoGameItemViewModels>();
            List<DemoGameCategoryViewModels> listResult = new List<DemoGameCategoryViewModels>();


            foreach (int gameID in gameIDGroup)
            {
                var matchGame = queryAllGames.FirstOrDefault(x => x.GameID == gameID);
                var instance = new DemoGameItemViewModels
                {
                    GameID = gameID,
                    GameName = DbMultiLanguageTextHelper.GetGameDisplayName(gameID),
                    GameType = matchGame?.GameType ?? 0,
                    GameCode = matchGame?.GameCode ?? "",
                    Servers = convert.Where(x => x.GameID == gameID).ToList()
                };

                listOfEachGames.Add(instance);
            }


            foreach (int gameType in gameTypeGroup)
            {
                var instance = new DemoGameCategoryViewModels
                {
                    GameType = gameType,
                    GameTypeName = MultiLanguageTextHelper.GetEnumsCultureText((GameTypeEnum)gameType),
                    Games = listOfEachGames.Where(x => x.GameType == gameType).ToList()
                };

                listResult.Add(instance);
            }


            //嘗試取得三方遊戲
            var queryEnableGames = MerchantGameSettingRepository.GetByMerchantID(merchantID , NormalActiveStatus.Enable);

            List<int> enableGameIDs = queryEnableGames.Select(z => z.GameID).ToList();

            var thirdPartyVendorGames = queryAllGames
                .Where(z => AppHelper.IsExternalVendor(z.GameVendorType))
                .Where(x => enableGameIDs.Contains(x.GameID))
                .ToList();

            if (thirdPartyVendorGames.IsNotEmptyNull())
            {
                var thirdGames = thirdPartyVendorGames.Select(z => new DemoGameItemViewModels
                    {
                        GameID = z.GameID,
                        GameCode = z.GameCode,
                        GameName = DbMultiLanguageTextHelper.GetGameDisplayName(z.GameID),
                        GameType = z.GameType,
                        Servers = new List<DemoServerItemViewModels>
                            {
                                new DemoServerItemViewModels
                                {
                                    ServerID = 0,
                                    ServerName = "Demo",
                                    FlagLevel = 1,
                                    GameID = z.GameID,
                                    GameType = z.GameType
                                }
                            }
                    })
                    .ToList();


                var allGameTypes = thirdPartyVendorGames.Select(z => z.GameType).Distinct();

                //補上 listResult 沒有支援的遊戲分類
                foreach (int gameType in allGameTypes)
                {
                    bool isGameTypeExist = listResult.Any(x => x.GameType == gameType);

                    if(!isGameTypeExist)
                    {
                        listResult.Add(new DemoGameCategoryViewModels
                        {
                            Games = new List<DemoGameItemViewModels>(),
                            GameType = gameType,
                            GameTypeName = MultiLanguageTextHelper.GetEnumsCultureText((GameTypeEnum)gameType)
                        });
                    }
                }

                foreach (var item in listResult)
                {
                    var matchGameTypes = thirdGames.Where(z => z.GameType == item.GameType).ToList();

                    if (matchGameTypes.IsNotEmptyNull())
                    {
                        item.Games.AddRange(matchGameTypes);
                    }
                }
              
            }

            return listResult;

        }



        /// <summary>
        /// 目前是Lobby頁面 才會用到
        /// </summary>
        public LobbyApiGetGamesInfoViewModel GetLobbyGames(LobbyApiGetGamesRequestModel model)
        {
            model.RuleFor(z => z.Ts).Rule_NotZeroNotMinus();
            model.RuleFor(z => z.PlayerName).Rule_NotEmptyText();
            model.RuleFor(z => z.SignKey).Rule_NotEmptyText();

            SignKeyRuleHelper.CheckIsValidPublicApiLobbySignkey(model.Ts, model.PlayerName, model.SignKey);

            string token = [REDACTED];

            var parseModel = GetTokenAndCheckValid(token);

            if (parseModel.LoginName != model.PlayerName)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_NotValid);

            var queryEnableGames = MerchantGameSettingRepository.GetByFilter(new QueryMerchantGameSettingRequestViewModel
            {
                MerchantIDs = parseModel.MerchantID.ToListCollection(),
                Status = NormalActiveStatus.Enable.ToInteger(),
                EnableGameOnly = true
            });

            var imageSetting = SystemConfigSettingRepository.GetParseModelByCache(SystemConfigSettingEnum.LobbyPageGameImageUrl);
            string imageDomain = imageSetting.SettingValue;

            var games = queryEnableGames.Select(z =>
            {
                var matchGame = GameListRepository.GetSingleByCache(z.GameID);

                var convert = new LobbyGamesViewModel
                {
                    GameID = z.GameID,
                    GameImageUrl = AppHelper.GenerateGameImageUrl(imageDomain , matchGame.GameCode),
                    GameName = DbMultiLanguageTextHelper.GetGameDisplayName(matchGame.GameID)
                };

                return convert;

            }).ToList();

            var balanceResult = PlayerWalletRefService.GetWalletAmountByMerchantWalletMode(parseModel.LoginUserID, parseModel.SessionTokenPGS);

            decimal balanceAmount = balanceResult.Amount;
            decimal balancePoint = ConvertMoneyToGamePoint(parseModel.MerchantID , balanceAmount);

            return new LobbyApiGetGamesInfoViewModel
            {
                Games = games,
                PlayerName = parseModel.LoginName,
                Balance = balancePoint
            };
        }


        private PlayerTokenJwtModel GetTokenAndCheckValid(string token)
        {
            var parseModel = PlayerJwtCheckService.ParsingLobbyToken(token);
            return parseModel;
        }


    }
}
