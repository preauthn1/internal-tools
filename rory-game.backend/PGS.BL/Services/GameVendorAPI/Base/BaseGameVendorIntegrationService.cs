using System;
using System.Collections.Generic;
using Lck.Utility.Extensions.Helper;
using target.BL.Interface;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.IOC;
using target.Utility.Utility;
using target.Utility.Extensions;
using System.Linq;
using target.Models.Models;
using target.Utility.Helper;
using target.Models.ViewModels.GameVendor;

namespace target.BL.Services
{
    public class BaseGameVendorIntegrationService : BaseService
    {
        /// <summary>
        /// 需實作項目，對應哪個vendor列舉
        /// </summary>
        public virtual GameVendorTypeEnum GameVendor { get; }

        /// <summary>
        /// 需實作項目，遊戲商前綴英文 (建議2~4碼 小寫英文)
        /// </summary>
        public virtual string GameVendorPrefix { get; }

        /// <summary>
        /// 報表API的時區依據
        /// </summary>
        public virtual double TimeZone { get; }



        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }

        public VendorSessionTokenRepository VendorSessionTokenRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }

        public PlayerRepository PlayerRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public int GameVendorID
        {
            get
            {
                return this.GameVendor.ToInteger();
            }
        }


        public VendorConfigSettingRepository VendorConfigSettingRepository { get; set; }

        public IExternalApiRequestLogRefService ExternalApiRequestLogRefService { get; set; }


        protected void WriteRequestLog(ApiRequestLogTemporaryInfoModel model)
        {
            ExternalApiRequestLogRefService.Create(model);
        }

        public double GetVendorTimeZone()
        {
            return this.TimeZone;
        }


        public int GetGameIdByGameCode(string gameCode)
        {
            var gameVendor = this.GameVendor.ToInteger();

            var queryGame = GameListRepository.GetByFilter(new QueryGameListRequestViewModel
            {
                ThirdPartyGameCode = gameCode,
                GameVendorTypes = gameVendor.ToListCollection()
            });

            //如果遊戲編號未定義
            if (queryGame.IsEmptyOrNull())
            {
                queryGame = GameListRepository.GetByFilter(new QueryGameListRequestViewModel
                {
                    ThirdPartyGameCode = "DEFAULT",
                    GameVendorTypes = gameVendor.ToListCollection()
                });
            }

            if (queryGame.IsEmptyOrNull())
                throw new TerproException(ApiErrorCode.NotFoundGame);


            return queryGame.FirstOrDefault().GameID;
        }


        protected MerchantVendorSettingModel GetVendorSettingByKey(int merchantID, VendorConfigSettingEnum settingKey)
        {
            string cacheKey = CacheKey.ExternalVendor_VendorConfigSettings(merchantID, this.GameVendor.ToInteger());

            var settingCache = CacheHelper.GetCacheOrQuery(cacheKey, () =>
            {
                var querySettings = VendorConfigSettingRepository.GetTargetMerchantVendorSettings(merchantID, this.GameVendor);
                return querySettings;

            }, 10);

            var match = settingCache.FirstOrDefault(x => x.SettingKeyEnum == settingKey);

            if (match == null)
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_SettingDataNotCompleted, new ErrorDetailModel(settingKey.ToString()));

            return match;
        }

        protected (bool IsDefine, decimal ConvertRate) GetMerchantConvertRate(int merchantID, GameVendorTypeEnum gameVendor)
        {
            var settingKeyResult = AppHelper.GetConvertRateCustomizedSettingKey(gameVendor);

            if (!settingKeyResult.IsDefine)
                return (false, 1.0M);

            decimal rate = MerchantCustomizedSettingRepository.GetMerchantGamePointConvertRateByCache(merchantID, settingKeyResult.SettingKey);
            return (true, rate);
        }


        protected bool IsBalanceNotEnughError(ApiErrorCode errorCode)
        {
            if (errorCode == ApiErrorCode.BetFail_WalletAmountNotEnough ||
               errorCode == ApiErrorCode.Forbidden_EnterGameError_WalletAmountNotEnough ||
               errorCode == ApiErrorCode.WalletTranferFail_WalletAmountNotEnough
                )
            {
                return true;
            }

            return false;
        }

        protected bool IsAuthenrizeError(ApiErrorCode errorCode)
        {
            if (errorCode == ApiErrorCode.NotFoundSessionToken ||
               errorCode == ApiErrorCode.Authorize_AccessToken_Expired ||
               errorCode == ApiErrorCode.Authorize_AccessToken_NotValid ||
               errorCode == ApiErrorCode.Authorize_TokenHasChanged ||
               errorCode == ApiErrorCode.NotFoundPlayer
                )
            {
                return true;
            }

            return false;
        }

        protected decimal ConvertGamePointToMoney(int merchantID, decimal point)
        {
            var rate = GetMerchantConvertRate(merchantID, GameVendor);
            return point / rate.ConvertRate;
        }


        protected bool IsNotFoundBetRecordError(ApiErrorCode errorCode)
        {
            if (errorCode == ApiErrorCode.NotFoundBetRecord || errorCode == ApiErrorCode.NotFoundRecord
                )
            {
                return true;
            }

            return false;

        }
        protected bool IsInvalidStatusError(ApiErrorCode errorCode)
        {
            if (errorCode == ApiErrorCode.InvalidStatus || errorCode == ApiErrorCode.InvalidStatus_DisableNow)
            {
                return true;
            }

            return false;

        }


        /// <summary>
        /// 檢查sessionToken是否合法 並返回playerID
        /// </summary>
        protected int CheckSessionTokenAndPlayer(ExternalVendorInputValidateRequestParameter model)
        {
            //預設要檢查 除非有指定不檢查
            if (!model.IsSkipCheckSessionToken)
            {
                model.RuleFor(x => x.SessionOrToken).Rule_NotEmptyText();

                var query = VendorSessionTokenRepository.GetTokenByCache(model.SessionOrToken);

                if (query == null)
                {
                    throw new TerproException(ApiErrorCode.NotFoundSessionToken);
                }

                if (model.PlayerID != 0 && query.PlayerID != model.PlayerID)
                {
                    throw new TerproException(ApiErrorCode.NotFoundSessionToken);
                }

                if (query.ExpiredTime <= DateTimeHelper.GetCurrentTime())
                {
                    throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);
                }

                return query.PlayerID;
            }

            return model.PlayerID;
        }
    }
}
