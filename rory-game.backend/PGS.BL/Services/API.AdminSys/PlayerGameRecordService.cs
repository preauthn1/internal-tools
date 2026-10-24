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

    public class PlayerGameRecordService : BaseAdminSysService, IPlayerGameRecordService
    {
        public PlayerGameRecordRepository PlayerGameRecordRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public ServerListRepository ServerListRepository { get; set; }


        public PagedList<PlayerGameRecordViewModel> GetByPagingFilter(int userID , QueryPlayerGameRecordApiRequestParameterModel filter)
        {
            var currentUser = GetCurrentOwnPermissionByCache(userID);

            filter.MerchantIDs = currentUser.SuperAdminOrAvailableMerchantIDs;
            filter.PlayerIDs = filter.PlayerID.ToIntegerListCollection();
            filter.GameIDs = filter.GameID.ToIntegerListCollection();
            filter.ServerIDs = filter.ServerID.ToIntegerListCollection();

            var query = PlayerGameRecordRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<PlayerGameRecordEntity, PlayerGameRecordViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            return viewModels;
        }


        private PlayerGameRecordViewModel ConvertToViewModels(PlayerGameRecordEntity queryPlayerGameRecord)
        {
            var convert = queryPlayerGameRecord.ReturnInjecteObject<PlayerGameRecordViewModel>();

            var queryMerchat = MerchantRepository.GetSingleByCache(convert.MerchantID);
            convert.ServerName = ServerListRepository.GetBriefByCache(convert.ServerID)?.ServerName;
            convert.GameName = DbMultiLanguageTextHelper.GetGameDisplayName(convert.GameID);

            if (queryMerchat != null)
            {
                convert.MerchantName = queryMerchat.Name;
                convert.WalletMode = queryMerchat.WalletMode;
            }

            return convert;
        }

    }
}
