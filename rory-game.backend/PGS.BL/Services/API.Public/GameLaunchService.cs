using System;
using System.Collections.Generic;
using target.BL.Interface;
using target.BL.Services;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.Models.Models;
using System.Linq;
using target.DA.Repository;
using target.Models;
using target.Utility;
using target.Models.Enums;
using target.Models.ViewModels.Base;
using target.Models.ViewModels.GameVendor;
using target.Utility.Utility.ThirdParty;

namespace target.BL.Services
{

    public partial class GameLaunchService : BaseVendorService, IGameLaunchService
    {
        public PlayerRepository PlayerRepository { get; set; }

        public PlayerLoginLogRepository PlayerLoginLogRepository { get; set; }

        public PlayerTotalSummaryRepository PlayerTotalSummaryRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }

        public MerchantGameSettingRepository MerchantGameSettingRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }

        public ThirdPartySlot88MD5RecordRepository ThirdPartySlot88MD5RecordRepository { get; set; }

        public SystemConfigSettingRepository SystemConfigSettingRepository { get; set; }

        public ServerListRepository ServerListRepository { get; set; }

        public VendorSessionTokenRepository VendorSessionTokenRepository { get; set; }

        public VendorMaintenanceRepository VendorMaintenanceRepository { get; set; }

        public PlayerCurrentServerRepository PlayerCurrentServerRepository { get; set; }

        public SingleWalletMerchantSessionTokenRepository SingleWalletMerchantSessionTokenRepository { get; set; }

        public LobbyTokenLogRepository LobbyTokenLogRepository { get; set; }


        public IPlayerWalletRefService PlayerWalletRefService { get; set; }

        public IPlayerJwtCheckService PlayerJwtCheckService { get; set; }

        public IMerchantIntergratedExternalApiService MerchantIntergratedExternalApiService { get; set; }

        public IPlayerCreateRefService PlayerCreateRefService { get; set; }

        public IGameServerDistributeRefService GameServerDistributeRefService { get; set; }


        /// <summary>
        /// 內部測試專用的API
        /// </summary>
        public string LaunchGameAndGetGamePageUrlByPlayerID(int playerID, int gameID)
        {
            var queryPlayer = PlayerRepository.GetById(playerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            string resultUrl = LaunchGameAndGetGamePageUrl(new GameLaunchRequestModel
            {
                SessionToken = [REDACTED],
                GameID = gameID,
                PlayerName = queryPlayer.PlayerName,
                MerchantID = queryPlayer.MerchantID,
            });

            return resultUrl;
        }


        public string LaunchGameAndGetGamePageUrl(GameLaunchRequestModel model)
        {
            //最主要的邏輯處理在此
            var result = LaunchGameAndCheckPlayerExist(model);

            if (result == null)
                throw new TerproException(ApiErrorCode.BusinessLogicDeny);

            // Case 1. Lobby網址
            if (result.IsLobbyLink)
            {
                return GenerateLobbyLink(model);
            }

            // Case 2. 外部遊戲商網址
            if (result.ExternalVendorInfo != null)
            {
                return result.ExternalVendorInfo.ExternalVendorUrl;
            }

            // Case 3. PGS自製遊戲 網址
            return GenerateRedirectUrl(model.MerchantID, result);
        }


        public GameLaunchViewModel LaunchGameAndCheckPlayerExist(GameLaunchRequestModel param)
        {
            param.RuleFor(x => x.PlayerName).Rule_NotEmptyText().Rule_IsValidAccountFormat();
            param.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            param.RuleFor(x => x.GameID).Rule_NotMinus();

            MerchantEntity queryMerchant = MerchantRepository.GetSingleByCache(param.MerchantID);

            if (queryMerchant == null)
                throw new TerproException(ApiErrorCode.NotFoundMerchant);

            bool isSingleWalletMode = queryMerchant.WalletMode == MerchantWalletModeEnum.SingleWalletMode.ToInteger();

            //如果是單錢包模式 有額外參數要檢查
            if (queryMerchant.WalletMode == MerchantWalletModeEnum.SingleWalletMode.ToInteger())
            {
                param.RuleFor(z => z.SessionToken).Rule_NotEmptyText();
                param.RuleFor(z => z.SessionToken).Rule_MinLengthEnough(10);
                param.RuleFor(z => z.SessionToken).Rule_MaxLengthOK(80);

                if (queryMerchant.Code.IsNullOrWhiteSpace())
                {
                    throw new TerproException(ApiErrorCode.BusinessLogicDeny_SettingDataNotCompleted, new ErrorDetailModel("Code"));
                }
            }

            //查詢供應商資料、檢查狀態
            this.GetMerchantAndCheckIsValid(param.MerchantID, queryMerchant.ApiSecretKey);

            //內部測試用商戶 可以跳過檢查
            if (!MerchantCustomizedSettingRepository.IsInternalTestMerchant(param.MerchantID))
            {
                //檢查系統維護
                SystemConfigSettingRepository.CheckSystemMaintainAndThrowException();
            }

            bool isLobbyLink = param.GameID == 0;

            //如果判斷是要導向 大廳Lobby
            if (isLobbyLink)
            {
                return new GameLaunchViewModel
                {
                    IsLobbyLink = true
                };
            }

            GameListEntity queryGame = GetGameAndCheckStatus(param);

            //檢查是否有開通遊戲
            var isEnableGame = MerchantGameSettingRepository.IsMerchantGameEnableInSystem(param.MerchantID, param.GameID);

            if (!isEnableGame)
                throw new TerproException(ApiErrorCode.InvalidStatus_GameStatus);

            if (!RegexRuleHelper.IsSupportLanguage(param.Lang))
                param.Lang = SystemConstant.DefaultLang;

            bool isExternalVendorMode = AppHelper.IsExternalVendor(queryGame.GameVendorType);

            if (isExternalVendorMode)
            {
                // 判斷為外部遊戲商遊戲
                // 判斷為外部遊戲商遊戲
                // 判斷為外部遊戲商遊戲
                // 判斷為外部遊戲商遊戲
                // 判斷為外部遊戲商遊戲
                // 判斷為外部遊戲商遊戲
                // 判斷為外部遊戲商遊戲
                //執行 [外部遊戲商] 的啟動模式 (打對方的API)
                return LaunchGameExternalVendorMode(param, queryGame, queryMerchant);
            }

            //若是第一次launch、新建帳號、相關資料，並返回playerID
            int playerID = FirstTimeLaunchThenCreateNewPlayerData(param, queryMerchant);
                
            string pgsSessionToken = [REDACTED];

            if (isSingleWalletMode)
            {
                CreateSingleWalletMerchantToken(param, param.MerchantID, playerID, pgsSessionToken);
            }

            //由後端處理 房間分配策略
            var chooseGameServer = GameServerDistributeRefService.ChooseAvailableGameServer(playerID, param.MerchantID, param.GameID, new ServerDistributeOptionalParameter
            {
                RoomFlagLevelID = param.RoomLevelFlagID,
                SecretParameter_ServerID = param.SecretParameter_ServerID
            });

            //寫入一筆登入紀錄
            string myIP = IPHelper.GetIP();

            PlayerLoginLogRepository.Create(new PlayerLoginLogEntity
            {
                PlayerID = playerID,
                MerchantID = param.MerchantID,
                MerchantSession = "",
                PlayerName = param.PlayerName,
                CreateIPAddress = myIP,
                CreateTime = DateTimeHelper.GetCurrentTime(),
                HDType = 0,
                HDCode = "",
                UserAgent = HeaderValueHelper.GetUserAgent()
            });


            //產生jwtToken 供H5使用
            var jwtResult = MyJwtHelper.GeneratePlayerJwt(playerID, param.MerchantID, param.PlayerName, param.Lang , pgsSessionToken , param.SessionToken);

            string jwtHash = jwtResult.JwtToken.MD5Hash();

            //讓所有三方session失效
            ResetAllPlayerSessionToken(playerID);

            //更新token (讓前次的token自動過期)
            PlayerRepository.UpdateTokenAndLastTime(playerID, jwtHash, myIP);

            //更新登錄次數
            PlayerTotalSummaryRepository.UpdateLoginTimes(playerID);


            //部分設置是存在 CustomizedSetting 這邊
            List<MerchantCustomizedSettingModel> customizedSettings = GetMerchantCustomizedSettings(param.MerchantID);

            var resultViewModel = new GameLaunchViewModel
            {
                PlayerJwtModel = jwtResult.JwtModel,
                IsMultipleServers = chooseGameServer.IsRedirectToHall,
                Info = new GameLaunchRequiredInfoViewModel
                {
                    ServerID = chooseGameServer.ServerID,
                    ServerConnection = chooseGameServer.ServerIPAddress,
                    GameID = queryGame.GameID,
                    Lang = param.Lang,
                    PlayerID = playerID,
                    PlayerName = param.PlayerName,
                    GameCode = queryGame.GameCode?.ToLower(),
                    Access_Token = [REDACTED],
                }
            };

            //客製customized設置
            AssignCustomizedSettingsValueToInfo(customizedSettings, example.invalid);

            //slot88特殊參數設置
            SetSlot88Information(resultViewModel, queryGame);

            return resultViewModel;
        }

        /// <summary>
        /// 上游遊戲商的啟用遊戲模式
        /// </summary>
        public GameLaunchViewModel LaunchGameExternalVendorMode(GameLaunchRequestModel model, GameListEntity queryGame, MerchantEntity queryMerchant)
        {
            var queryMaintain = VendorMaintenanceRepository.GetVendorMaintenance(queryGame.GameVendorType);

            if (queryMaintain != null)
                throw new TerproException(ApiErrorCode.GameUnderMaintain);

            //若是第一次launch、新建帳號、相關資料，並返回playerID
            int playerID = FirstTimeLaunchThenCreateNewPlayerData(model, queryMerchant);

            //檢查卡房
            var queryCurrentServer = PlayerCurrentServerRepository.GetByPlayerID(playerID);

            if (queryCurrentServer != null)
            {
                throw new TerproException(ApiErrorCode.Forbidden_DuplicateEnterGameServer);
            }

            //寫入一筆登入紀錄
            string myIP = IPHelper.GetIP();

            DateTime now = DateTimeHelper.GetCurrentTime();
            DateTime afterOneDay = DateTimeHelper.GetCurrentTime().AddDays(1);

            bool isSingleWalletMode = queryMerchant.WalletMode == MerchantWalletModeEnum.SingleWalletMode.ToInteger();
            string pgsSessionToken = [REDACTED];

            //若平台方 採用 [單錢包模式]
            if (isSingleWalletMode)
            {
                CreateSingleWalletMerchantToken(model, queryMerchant.MerchantID, playerID , pgsSessionToken);
            }

            //讓所有三方session失效
            ResetAllPlayerSessionToken(playerID);

            VendorSessionTokenRepository.Create(new VendorSessionTokenEntity
            {
                CreateTime = now,
                UpdateTime = now,
                ExpiredTime = afterOneDay,
                GameID = queryGame.GameID,
                PlayerID = playerID,
                SessionToken = [REDACTED],
                VendorID = queryGame.GameVendorType
            });

            long resultID = PlayerLoginLogRepository.Create(new PlayerLoginLogEntity
            {
                PlayerID = playerID,
                MerchantID = model.MerchantID,
                MerchantSession = "",
                PlayerName = model.PlayerName,
                CreateIPAddress = myIP,
                CreateTime = DateTimeHelper.GetCurrentTime(),
                HDType = 0,
                HDCode = "",
                UserAgent = HeaderValueHelper.GetUserAgent(),
                GameID = queryGame.GameID,
                Memo = pgsSessionToken,
                Status = NormalSuccessStatusEnum.Success.ToInteger()
            });

            //根據不同遊戲商 實例化不同service實作
            IGameVendorIntegrationService vendorService = ServiceInstanceHelper.InstanceVendorService(queryGame.GameVendorType);

            bool isTrialMode = IsAlwaysTrialModeEnable(queryMerchant.MerchantID);

            string launchUrl = string.Empty;

            try
            {
                // step 1. 先產生三方廠商的遊戲帳號
                vendorService.CreateVendorPlayerAccountOrNot(new ExternalVendorCreateAccountRequestParameter
                {
                    PlayerName = model.PlayerName,
                    Currency = queryMerchant.Currency
                });

                string currencyParse = AppHelper.ConvertCurrencyBETA_To_CNY(queryMerchant.Currency);

                var targetLang = AppHelper.ConvertLangToEnum(model.Lang);
                
                string lobbyUrl = GenerateLobbyLink(model);

                // step 2. 取得遊戲網址
                var launchResult = vendorService.GetGameLaunchUrl(new ExternalVendorLaunchUrlRequestModel
                {
                    Language = targetLang,
                    Currency = currencyParse,
                    ThirdPartyGameCode = queryGame.ThirdPartyGameCode,
                    MerchantID = model.MerchantID,
                    PlayerName = model.PlayerName,
                    PlayerID = playerID,
                    ThirdPartyLaunchSessionToken = [REDACTED],
                    IsTrialMode = isTrialMode,
                    BackUrl = lobbyUrl,
                    DebugReferenceUUID = AppHelper.GenerateGameLaunchRequestUUID(resultID)
                });

                launchUrl = launchResult.Url;
            }
            catch (Exception ex)
            {
                PlayerLoginLogRepository.UpdateStatus(resultID, NormalSuccessStatusEnum.Fail);
                throw ex;
            }

            //產生jwtToken 供H5使用
            var jwtResult = MyJwtHelper.GeneratePlayerJwt(playerID, model.MerchantID, model.PlayerName, model.Lang , pgsSessionToken , model.SessionToken);

            string jwtHash = jwtResult.JwtToken.MD5Hash();

            //更新token (讓前次的token自動過期)
            PlayerRepository.UpdateTokenAndLastTime(playerID, jwtHash, myIP);

            //更新登錄次數
            PlayerTotalSummaryRepository.UpdateLoginTimes(playerID);

            return new GameLaunchViewModel
            {
                ExternalVendorInfo = new ExternalVendorLaunchInfoViewModel
                {
                    ExternalVendorUrl = launchUrl
                }
            };
        }

        private void CreateSingleWalletMerchantToken(GameLaunchRequestModel model, int merchantID, int playerID, string pgsSessionToken)
        {
            //只有在Lobby模式下 才允許使用同一組sessionToken
            if (!model.IsAllowDuplicateSessionToken)
            {
                var queryExist = SingleWalletMerchantSessionTokenRepository.GetByFilter(new QuerySingleWalletMerchantSessionTokenRequestViewModel
                {
                    MerchantIDs = merchantID.ToListCollection(),
                    SessionTokenRaw = model.SessionToken
                });

                if (queryExist.IsNotEmptyNull())
                {
                    throw new TerproException(ApiErrorCode.BusinessLogicDeny_DuplicateDataExist);
                }
            }

            DateTime now = DateTimeHelper.GetCurrentTime();
            SingleWalletMerchantSessionTokenRepository.Create(new SingleWalletMerchantSessionTokenEntity
            {
                SessionTokenPGS = pgsSessionToken,
                SessionTokenRaw = model.SessionToken,
                CreateTime = now,
                PlayerDisplayName = model.PlayerName,
                GameID = model.GameID,
                MerchantID = model.MerchantID,
                PlayerID = playerID,
                UpdateTime = now,
            });
        }

        private string GenerateLobbyLink(GameLaunchRequestModel param)
        {
            var queryMerchant = MerchantRepository.GetSingleByCache(param.MerchantID);

            //若是第一次launch、新建帳號、相關資料，並返回playerID
            int playerID = FirstTimeLaunchThenCreateNewPlayerData(param, queryMerchant);

            //產生jwtToken 供H5使用
            var jwtResult = MyJwtHelper.GeneratePlayerJwt(playerID, param.MerchantID, param.PlayerName, param.Lang, "" , param.SessionToken);

            string jwtHash = jwtResult.JwtToken.MD5Hash();

            LobbyTokenLogRepository.ResetAllPlayerToken(playerID);

            Guid guid = Guid.NewGuid();

            LobbyTokenLogRepository.Create(new LobbyTokenLogEntity
            {
                CreateTime = DateTimeHelper.GetCurrentTime(),
                ExpiredTime = DateTimeHelper.GetCurrentTime().AddDays(2),
                MerchantID = param.MerchantID,
                PlayerID = playerID,
                IdentityHash = jwtHash,
                IdentityGuid = guid,
                RawJwtToken = [REDACTED]
            });

            var queryDomain = SystemConfigSettingRepository.GetParseModelByCache(SystemConfigSettingEnum.LobbyPageUrl);

            if (queryDomain == null || queryDomain.SettingValue.IsNullOrWhiteSpace())
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_SettingDataNotCompleted, new ErrorDetailModel { ColumnName = "LobbyPageDomain" });
            
            //舊版寫法 有max length 2048的可能性
            //string base64Token = [REDACTED]; 

            //因為之前舊版有跟前端約定好會做base64Encode，新版需保留一樣的模式
            string base64Token = [REDACTED];            
            return $"{queryDomain.SettingValue}?playerName={param.PlayerName}&token=[REDACTED]";
        }


        private bool IsAlwaysTrialModeEnable(int merchantID)
        {
            var settintType = MerchantCustomizedSettingKeyEnum.System_IsAlwaysTrialMode;

            var query = MerchantCustomizedSettingRepository.GetByMerchantSettingKey(merchantID, settintType);

            if (query != null)
                return query.SettingValueOfBoolean;

            var defaultValue = SettingKeyEnumLogicHelper.GetMerchantSettingDefaultValueByType(settintType);

            var parsingModel = new MerchantCustomizedSettingModel
            {
                SettingValue = defaultValue,
                CustomizedSettingKey = settintType.ToInteger()
            };

            return parsingModel.SettingValueOfBoolean;
        }


        private void ResetAllPlayerSessionToken(int playerID)
        {
            VendorSessionTokenRepository.ResetAllPlayerSessionToken(playerID);
        }


        private GameListEntity GetGameAndCheckStatus(GameLaunchRequestModel model)
        {
            var queryGame = GameListRepository.GetById(model.GameID);

            if (queryGame == null)
                throw new TerproException(ApiErrorCode.NotFoundGame);

            if (queryGame.Status != (int)NormalActiveStatus.Enable)
                throw new TerproException(ApiErrorCode.InvalidStatus_GameStatus);

            return queryGame;
        }


        private int FirstTimeLaunchThenCreateNewPlayerData(GameLaunchRequestModel model, MerchantEntity queryMerchant)
        {
            var newUserModel = new PlayerCreateParameterModel
            {
                PlayerName = model.PlayerName.Trim(),
                MerchantID = model.MerchantID,
                Currency = queryMerchant.Currency
            };

            //新建新用戶或返回已存在playerID
            int playerID = PlayerCreateRefService.CreatePlayerOrReturnPlayerID(newUserModel);
            return playerID;
        }


        private void AssignCustomizedSettingsValueToInfo(List<MerchantCustomizedSettingModel> customizedSettings, GameLaunchRequiredInfoViewModel info)
        {
            var styleSetting = customizedSettings.FirstOrDefault(x => x.CustomizedSettingKey == (int)MerchantCustomizedSettingKeyEnum.Style_LoadingStyleMode);
            var clientLogSetting = customizedSettings.FirstOrDefault(x => x.CustomizedSettingKey == (int)MerchantCustomizedSettingKeyEnum.System_ClientLogEnable);

            info.StyleMode = styleSetting.SettingValueOfInt;
            info.Log = clientLogSetting.SettingValueOfInt;
        }



        private List<MerchantCustomizedSettingModel> GetMerchantCustomizedSettings(int merchantID)
        {
            List<MerchantCustomizedSettingKeyEnum> settings = new List<MerchantCustomizedSettingKeyEnum>
            {
                MerchantCustomizedSettingKeyEnum.Style_LoadingStyleMode,
                MerchantCustomizedSettingKeyEnum.System_ClientLogEnable
            };

            var queryStyleSetting = MerchantCustomizedSettingRepository.GetMultipleMerchantSettings(merchantID, settings);


            List<MerchantCustomizedSettingModel> results = new List<MerchantCustomizedSettingModel>();


            foreach (var eachEnum in settings)
            {
                var matchSetting = queryStyleSetting.FirstOrDefault(x => x.CustomizedSettingKey == eachEnum.ToInteger());
                var defaultValueResult = SettingKeyEnumLogicHelper.GetMerchantSettingOrDefaultValue(eachEnum, matchSetting);
                results.Add(defaultValueResult);
            }

            return results;
        }


        /// <summary>
        /// 針對特定三方遊戲商有特殊欄位要帶入
        /// </summary>
        /// <param name="resultViewModel"></param>
        /// <param name="queryGame"></param>
        private void SetSlot88Information(GameLaunchViewModel resultViewModel, GameListEntity queryGame)
        {
            var vendorType = (GameVendorTypeEnum)queryGame.GameVendorType;

            if (vendorType == GameVendorTypeEnum.target)
                return;

            if (vendorType == GameVendorTypeEnum.Slot88)
            {
                string guid = RandomHelper.GetGuidWithoutDash();


                var launchInfo = example.invalid;

                var info = new Slot88MD5InfoResultModel
                {
                    PlayerInfo = new Slot88MD5PlayerInfoModel
                    {
                        Access_Token = [REDACTED],
                        Lang = launchInfo.Lang,
                        PlayerID = launchInfo.PlayerID,
                        PlayerName = launchInfo.PlayerName
                    },
                    ServerInfo = new Slot88MDServerInfoModel
                    {
                        ServerConnection = launchInfo.ServerConnection,
                        ServerID = launchInfo.ServerID,
                        GameID = queryGame.GameID,
                        GameVendorType = queryGame.GameVendorType,
                        ConnType = "Websocket"
                    },
                    ExpiredTime = DateTimeHelper.GetCurrentTime().AddMinutes(SettingHelper.PublicApiSetting.Slot88MD5ExpiredMinutes),
                    RandomID = guid
                };

                string beforeMD5 = (guid + SystemConstant.ThirdPartyMD5_SaltKey);
                string MD5 = beforeMD5.MD5Hash();

                ThirdPartySlot88MD5RecordRepository.Create(new ThirdPartySlot88MD5RecordEntity
                {
                    CreateTime = DateTimeHelper.GetCurrentTime(),
                    GameID = queryGame.GameID,
                    MD5 = MD5,
                    OriginalJson = info.ConvertToJsonString(),
                    PlayerID = launchInfo.PlayerID,
                });

                launchInfo.ThirdPartyInfo = new ThirdPartyExtraInfoModel
                {
                    GameVendorID = queryGame.GameVendorType,
                    MD5Text = MD5
                };
            }
        }




        private string GetGameLinkDomain(int merchantID)
        {
            //也許未來會根據 merchantID 不同 分配不同domain (避免域名劫持?)
            var setting = SystemConfigSettingRepository.GetEntityByEnum(SystemConfigSettingEnum.GamePlayDomain);

            if (setting == null || setting.SettingValue.IsEmptyOrNull())
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_SettingDataNotCompleted);

            return setting.SettingValue;
        }



        private string GenerateRedirectUrl(int merchantID, GameLaunchViewModel result)
        {
            string domain = GetGameLinkDomain(merchantID);
            string redirectUrl = GenerateRedirectUrl(domain, result);
            return redirectUrl;
        }


        private string GenerateRedirectUrl(string settingUrl, GameLaunchViewModel result)
        {
            //實際上只有info這個物件會包給client & servers使用
            var launchInfo = example.invalid;

            var myObject = new
            {
                data = launchInfo
            };

            string json64Encode = myObject.ConvertToJsonStringCamelCase().Base64Encode();
            string domain = settingUrl.TruncateLastWordIfMatch('/');
            string langCode = launchInfo.Lang?.ToUpper();
            string folderName = SettingHelper.IsDev ? "web-mobile-dev" : "web-mobile"; //懶得多加config設定 所以這樣處理

            //導到選廳頁
            if (result.IsMultipleServers)
                return $"{domain}/{folderName}/hall/{langCode}/?gameID={launchInfo.GameID}&styleMode={launchInfo.StyleMode}&token=[REDACTED]";

            //直接導入遊戲頁
            string redirectUrl = $"{domain}/{folderName}/{launchInfo.GameCode.ToLower()}/{langCode}/?{json64Encode}";
            return redirectUrl;
        }




    }
}
