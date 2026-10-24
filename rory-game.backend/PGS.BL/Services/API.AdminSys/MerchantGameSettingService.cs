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

    public class MerchantGameSettingService : BaseAdminSysService, IMerchantGameSettingService
    {
        public GameListRepository GameListRepository { get; set; }

        public MerchantGameSettingRepository MerchantGameSettingRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }


        public MerchantGameSettingViewModel GetGameSettingByMerchantID(int merchantID)
        {
            var querySettings = MerchantGameSettingRepository.GetByMerchantID(merchantID);

            var queryAllGames = GameListRepository.GetByFilter(new QueryGameListRequestViewModel
            {
                IsNotDelete = true
            });

            var viewModelItems = queryAllGames.Select(x =>
            {
                var convert = new MerchantGameSettingItemViewModel
                {
                    GameID = x.GameID,
                    GameName = DbMultiLanguageTextHelper.GetGameDisplayName(x.GameID),
                    Status = NormalActiveStatus.Disable.ToInteger()
                };

                var matchSetting = querySettings.FirstOrDefault(z => z.GameID == x.GameID);

                if (matchSetting != null)
                {
                    convert.Status = matchSetting.Status;
                }

                return convert;

            }).ToList();

            return new MerchantGameSettingViewModel
            {
                GameSettingItems = viewModelItems
            };
        }



        public void Update(int merchantID , MerchantGameSettingUpdateRequestModel model)
        {
            (new { merchantID }).RuleFor(x => x.merchantID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.Status).Rule_NotMinus();

            var queryMerchant = MerchantRepository.GetSingleByCache(merchantID);

            if (queryMerchant == null)
                throw new TerproException(ApiErrorCode.NotFoundMerchant);

            var queryGame = GameListRepository.GetSingleByCache(model.GameID);

            if (queryGame == null)
                throw new TerproException(ApiErrorCode.NotFoundGame);

            bool isSingleWalletMode = CommonLogicHelper.IsSingleWalletMode(queryMerchant.WalletMode);

            if(isSingleWalletMode && 
                model.Status == NormalActiveStatus.Enable.ToInteger() && 
                AppHelper.IsPgsSelfGame(model.GameID) && 
                queryGame.GameType != GameTypeEnum.Slot.ToInteger())
            {
                throw new TerproException(ApiErrorCode.BusinessLogicDeny);
            }

            var queryExist = MerchantGameSettingRepository.GetByMerchantIDAndGameID(merchantID , model.GameID);

            if(queryExist == null)
            {
                MerchantGameSettingRepository.Create(new MerchantGameSettingEntity
                {
                    GameID = model.GameID,
                    MerchantID = merchantID,
                    Status = model.Status,
                    Sort = 0,
                    UpdateBy = GetCurrentUser().LoginName,
                    UpdateTime = DateTimeHelper.GetCurrentTime(),
                });
            }
            else
            {
                queryExist.Status = model.Status;
                queryExist.UpdateBy = GetCurrentUser().LoginName;
                queryExist.UpdateTime = DateTimeHelper.GetCurrentTime();

                MerchantGameSettingRepository.Update(queryExist);
            }

            CacheRemoveHelper.RemoveCache_WhenUpdateMerchantGames(merchantID);
        }


    }
}
