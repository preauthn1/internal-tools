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

namespace target.BL.Services
{

    public class PlayerLoginLogService : BaseAdminSysService, IPlayerLoginLogService
    {
        public PlayerLoginLogRepository PlayerLoginLogRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }


        public PagedList<PlayerLoginLogViewModel> GetByPagingFilter(int userID , QueryPlayerLoginLogApiRequestParameterModel filter)
        {
            var cureentUser = GetCurrentOwnPermissionByCache(userID);

            filter.MerchantIDs = cureentUser.SuperAdminOrAvailableMerchantIDs;
            filter.PlayerIDs = filter.PlayerID.ToIntegerListCollection();

            var query = PlayerLoginLogRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<PlayerLoginLogEntity, PlayerLoginLogViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            return viewModels;
        }


        private PlayerLoginLogViewModel ConvertToViewModels(PlayerLoginLogEntity query)
        {
            var convert = query.ReturnInjecteObject<PlayerLoginLogViewModel>();
            convert.MerchantName = MerchantRepository.GetSingleByCache(query.MerchantID)?.Name;
            convert.HDType = HDTypeEnum.PC.ToInteger();
            convert.GameName = DbMultiLanguageTextHelper.GetGameDisplayName(query.GameID);
            convert.StatusText = MultiLanguageTextHelper.GetEnumsCultureText((NormalSuccessStatusEnum)query.Status);
            convert.RequestUUID = AppHelper.GenerateGameLaunchRequestUUID(query.Id);

            var gameCache = DbValueCacheHelper.GetGame(convert.GameID);

            if(gameCache != null)
            {
                convert.VendorID = gameCache.GameVendorType;
                convert.VendorName = MultiLanguageTextHelper.GetGameVendorDisplayName(gameCache.GameVendorType);
                convert.IsLaunchRecordExist = GameMerchantCustomizeHelper.IsGameLaunchRecordExist(gameCache.GameVendorType);
            }

            return convert;
        }

    }
}
