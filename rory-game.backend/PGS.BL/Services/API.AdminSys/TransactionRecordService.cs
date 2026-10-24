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
using example.invalid;
using target.Models.Models.StoreProcedures;

namespace target.BL.Services
{

    public class TransactionRecordService : BaseAdminSysService, ITransactionRecordService
    {
        public TransactionRepository TransactionRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }

        public ServerListRepository ServerListRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public PlayerRepository PlayerRepository { get; set; }


        public PagedList<TransactionRecordViewModel> GetPagingByFilter(int userID , QueryTransactionRecordRequestViewModel model)
        {
            var currentUser = GetCurrentOwnPermissionByCache(userID);
            var filter = model.ReturnInjecteObject<QueryTransactionRecordRequestViewModel>();

            if(model.PlayerName.IsNotNullEmptyText())
            {
                var queryPlayers = PlayerRepository.GetByPlayerName(model.PlayerName);

                //如果查得到唯一的playerID，用playerID來查 (因為有設index)
                if(queryPlayers.IsNotEmptyNull() && queryPlayers.Count() == 1)
                {
                    filter.PlayerID = queryPlayers.FirstOrDefault().PlayerID;
                    filter.PlayerName = string.Empty;
                }
            }

            filter.MerchantIDs = filter.MerchantID.ToIntegerListCollection();
            filter.MerchantIDs = FilterAvailablePermissionMerchantIDs(userID, filter.MerchantIDs);

            PagedList<TransactionEntity> query = TransactionRepository.GetPagingByFilter(filter);

            var paging = ConvertToNewTypePagingObject<TransactionEntity, TransactionRecordViewModel>(query);
            paging.Items = query.Items.Select(x =>
            {
                var convert = ConvertToViewModels(x);
                if (convert.GameID > 0)
                {
                    convert.GameName = DbMultiLanguageTextHelper.GetGameDisplayName(convert.GameID);
                }
                return convert;
            }).ToList();

            return paging;
        }
              
        private TransactionRecordViewModel ConvertToViewModels(TransactionEntity query)
        {
            var convert = query.ReturnInjecteObject<TransactionRecordViewModel>();

            var queryMerchant = MerchantRepository.GetSingleByCache(query.MerchantID);
            convert.MerchantName = queryMerchant.Name;
            convert.WalletMode = queryMerchant.WalletMode;
            convert.UpdateAmount = query.Amount;
            convert.CurrentAmount = query.BeforeAmount;
            convert.TransactionTypeText = MultiLanguageTextHelper.GetEnumsCultureText((TransactionTypeEnum)query.TransactionType);

            return convert;
        }
    }
}
