using game.Entity.Game;
using target.BL.Interface;
using target.BL.Services;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Models;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Models.ViewModels.GameVendor;
using target.Utility;
using target.Utility.Extensions;
using target.Utility.Extensions.Web;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Data;
using System.Data.SqlClient;
using System.Linq;

namespace target.BL.Services
{
    public class BetRecordService : BaseAdminSysService, IBetRecordService
    {

        public BetRecordRepository BetRecordRepository { get; set; }

        public BetRecordTransactionNoMappingRepository BetRecordTransactionNoMappingRepository { get; set; }

        public PerGameDetailLogRepository PerGameDetailLogRepository { get; set; }

        public PlayerWalletRepository PlayerWalletRepository { get; set; }

        public UserOperateRecordRepository UserOperateRecordRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public PerGameLogRepository PerGameLogRepository { get; set; }

        public UserRepository UserRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }

        public BetRecordMerchantNotifyQueueRepository BetRecordMerchantNotifyQueueRepository { get; set; }

        public SingleWalletMerchantApiRequestLogRepository SingleWalletMerchantApiRequestLogRepository { get; set; }



        public ISingleWalletMerchantBetProcessRefService SingleWalletMerchantBetProcessRefService { get; set; }

        public ISingleWalletPlatformApiRequestRefService SingleWalletPlatformApiRequestRefService { get; set; }



        /// <summary>
        /// 取得每一家投注資料 union起來的資料查詢
        /// </summary>
        public PagedList<UnionBetRecordViewModel> GetUnionBetRecordsByPagingFilter(int userID, QueryUnionBetRecordRequestViewModel filter)
        {
            filter.MerchantIDs = FilterAvailablePermissionMerchantIDs(userID, filter.MerchantIDs);

            //如果查單號 直接不篩選日期 (被日期誤導好多次了...)
            if (!string.IsNullOrWhiteSpace(filter.TransactionNo) ||
                !string.IsNullOrWhiteSpace(filter.TransactionNoPGS))
            {
                filter.CreateTimeBegin = null;
                filter.CreateTimeEnd = null;
            }

            var query = BetRecordRepository.GetUnionBetRecordsByPagingFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<UnionBetRecordQueryModel, UnionBetRecordViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            //額外欄位設置
            SetExtraProperty(viewModels.Items);

            return viewModels;
        }


        public BetRecordGeneralInfoViewModel GetByTransactionNo(int userID, BetRecordQueryRequestModel model)
        {
            model.RuleFor(z => z.TransactionNo).Rule_NotEmptyText();

            var currentUser = GetCurrentOwnPermissionByCache(userID);

            var queryMapping = BetRecordTransactionNoMappingRepository.GetByTransactionNo(model.TransactionNo);

            if (queryMapping == null)
                throw new TerproException(ApiErrorCode.NotFoundRecord);

            var queryGenralInfoByVendor = GetByPgsTransactionNo(queryMapping.GameVendorID, queryMapping.BetRecordNoPGS);

            if (queryGenralInfoByVendor == null)
                throw new TerproException(ApiErrorCode.NotFoundRecord);

            queryGenralInfoByVendor.GameName = DbMultiLanguageTextHelper.GetGameDisplayName(queryGenralInfoByVendor.GameID);

            if (!currentUser.IsSuperAdminOrg && !currentUser.AvailableMerchantIDs.Contains(queryGenralInfoByVendor.MerchantID))
                throw new TerproException(ApiErrorCode.NotFoundRecord);


            var viewModelConvert = queryGenralInfoByVendor.ReturnInjecteObject<BetRecordGeneralColumnViewModel>();
            viewModelConvert.MerchantName = MerchantRepository.GetSingleByCache(viewModelConvert.MerchantID)?.Name;

            if (viewModelConvert.ManualSettleUserID.HasValue)
            {
                viewModelConvert.IsManualSettlement = true;
                viewModelConvert.ManualSettlementUserName = UserRepository.GetSingleByCache(viewModelConvert.ManualSettleUserID.Value)?.UserName;
            }

            var queryGame = GameListRepository.GetSingleByCache(viewModelConvert.GameID);

            var viewModel = new BetRecordGeneralInfoViewModel
            {
                GameVendorID = queryMapping.GameVendorID,
                GameVendorName = ((GameVendorTypeEnum)queryMapping.GameVendorID).ToString(),
                GeneralInfo = viewModelConvert,
                IsSupportExternalLink = GameVendorDefinitionHelper.IsSupportExternalLinkUrl(queryGame.GameVendorType)
            };

            var queryPerGame = PerGameDetailLogRepository.GetByFilter(new QueryDetailLogRequestViewModel
                {
                    TransactionNo = queryMapping.BetRecordNoPGS,
                })
                .FirstOrDefault();

            if (queryPerGame != null)
            {
                viewModel.PerGameLogInfo = new PerGameDetailInfo
                {
                    PerGameLogID = queryPerGame.PerGameLogID,
                    PerGameDetailLogID = queryPerGame.Id,
                    ServerID = queryPerGame.ServerID
                };
            }

            return viewModel;
        }


        /// <summary>
        /// 手動取消注單
        /// </summary>
        public void ManualRefundBetRecord(int userID, BetRecordManualRefundRequestModel model)
        {
            model.RuleFor(z => z.TransactionNo).Rule_NotEmptyText();

            var currentUser = GetCurrentOwnPermissionByCache(userID);

            var queryMapping = BetRecordTransactionNoMappingRepository.GetByTransactionNo(model.TransactionNo);

            if (queryMapping == null)
                throw new TerproException(ApiErrorCode.NotFoundRecord);

            var queryGenralInfoByVendor = GetByPgsTransactionNo(queryMapping.GameVendorID, queryMapping.BetRecordNoPGS);

            if (queryGenralInfoByVendor == null)
                throw new TerproException(ApiErrorCode.NotFoundRecord);

            //只有pending中的狀態可以操作退款
            if (queryGenralInfoByVendor.Status != (int)GameVendorBetRecordStatusEnum.Pending)
            {
                throw new TerproException(ApiErrorCode.InvalidStatus);
            }

            var queryMerchant = MerchantRepository.GetById(queryGenralInfoByVendor.MerchantID);


            string ip = IPHelper.GetIP();


            // Step 1. 修改狀態
            var vendorService = ServiceInstanceHelper.InstanceVendorService(queryMapping.GameVendorID);
            vendorService.CancelBet(queryMapping.BetRecordNoPGS, model.Memo);


            // Step 2. 寫入一筆操作紀錄
            UserOperateRecordRepository.Create(new UserOperateRecordEntity
            {
                UserID = userID,
                OperateType = UserOperateTypeEnum.SingleWallet_CancelRefundWalletManual.ToInteger(),
                CreateIPAddress = ip,
                CreateTime = DateTimeHelper.GetCurrentTime(),
                TargetID = queryGenralInfoByVendor.Id,
                ChangeLogInfo = queryGenralInfoByVendor.ConvertToJsonString(),
                Memo = model.Memo
            });


            // Step 3. 若 [單錢包模式] 
            if (CommonLogicHelper.IsSingleWalletMode(queryMerchant.WalletMode))
            {
                //呼叫對方API
                SingleWalletMerchantBetProcessRefService.CancelProcessPassToMerchantApi(new game.Entity.Game.SingleWalletModeCancelApiRequestModel
                {
                    TransactionNo = queryGenralInfoByVendor.BetRecordNoPGS,
                    GameID = queryGenralInfoByVendor.GameID,
                    IsExternalGameVendor = true,
                    MerchantID = queryGenralInfoByVendor.MerchantID
                });
            }
            else
            {
                //更新PGS錢包
                PlayerWalletRepository.UpdateWalletAmountForbiddenMinus(
                    queryGenralInfoByVendor.PlayerID,
                    queryGenralInfoByVendor.GameID,
                    queryGenralInfoByVendor.BetAmount,
                    TransactionTypeEnum.SingleWallet_CancelRefundWalletManual,
                    queryGenralInfoByVendor.BetRecordNoPGS,
                    ip
                );
            }


            return;
        }


        public void ManualSettleBetRecord(int userID, BetRecordManualSettleRequestModel model)
        {
            model.RuleFor(z => z.TransactionNo).Rule_NotEmptyText();
            model.RuleFor(z => z.PayoutAmount).Rule_NotMinus();
            model.RuleFor(z => z.ValidAmount).Rule_NotMinus();
            model.RuleFor(z => z.Memo).Rule_NotEmptyText();

            var currentUser = GetCurrentOwnPermissionByCache(userID);

            var queryMapping = BetRecordTransactionNoMappingRepository.GetByTransactionNo(model.TransactionNo);

            if (queryMapping == null)
                throw new TerproException(ApiErrorCode.NotFoundRecord);

            var queryGenralInfoByVendor = GetByPgsTransactionNo(queryMapping.GameVendorID, queryMapping.BetRecordNoPGS);

            if (queryGenralInfoByVendor == null)
                throw new TerproException(ApiErrorCode.NotFoundRecord);

            var queryMerchant = MerchantRepository.GetById(queryGenralInfoByVendor.MerchantID);

            string ip = IPHelper.GetIP();
            string processID = Guid.NewGuid().ToString();

            decimal payoutAmount = model.PayoutAmount;
            decimal validAmount = model.ValidAmount;

            // Step 1. 修改狀態
            var vendorService = ServiceInstanceHelper.InstanceVendorService(queryMapping.GameVendorID);
            vendorService.UpdateRawBetRecordAfterSettlement(new SettlementBetRecordRequiredInfoModel
            {
                BetRecordNoPGS = queryMapping.BetRecordNoPGS,
                PayoutAmount = payoutAmount,
                ValidAmount = validAmount,
                IsAdminSysManualOperate = true,
                DebugExtraMemo = processID,
                UserID = userID
            });

            // Step 2. 寫入一筆操作紀錄
            UserOperateRecordRepository.Create(new UserOperateRecordEntity
            {
                UserID = userID,
                OperateType = UserOperateTypeEnum.SingleWallet_CancelRefundWalletManual.ToInteger(),
                CreateIPAddress = ip,
                CreateTime = DateTimeHelper.GetCurrentTime(),
                TargetID = queryGenralInfoByVendor.Id,
                ChangeLogInfo = queryGenralInfoByVendor.ConvertToJsonString(),
                Memo = model.Memo
            });

            // step 3. 轉成perGameLog格式
            var convertPerGame = vendorService.FormatParsing_ConvertBetRecordToPerGameDetailLog(queryGenralInfoByVendor.BetRecordNoPGS);

            UpdateBalanceSPResult updateResult = null;

            // step 4. 更新錢包金額
            if (CommonLogicHelper.IsSingleWalletMode(queryMerchant.WalletMode))
            {
                //單錢包模式
                var swResult = SingleWalletMerchantBetProcessRefService.SettlemenetProcessPassToMerchantApi(new SingleWalletModeSettleApiParameter
                {
                    PlayerID = queryGenralInfoByVendor.PlayerID,
                    BetRecordNoPGS = queryMapping.BetRecordNoPGS,
                    SettlementTime = DateTimeHelper.GetCurrentTime(),
                    PayoutAmount = convertPerGame.PayoutAmount,
                    BetAmount = convertPerGame.BetAmount,
                    GameID = convertPerGame.GameID,
                    MerchantID = convertPerGame.MerchantID,
                    TaxAmount = convertPerGame.TaxAmount,
                    ValidAmount = convertPerGame.ValidAmount,
                    DebugProcessUUID = processID
                });

                updateResult = new UpdateBalanceSPResult
                {
                    AfterAmount = swResult.Balance,
                    Code = ApiErrorCode.Success
                };
            }
            else
            {
                //轉帳模式
                var transactionType = TransactionTypeEnum.SingleWallet_BetSettlement;
                updateResult = PlayerWalletRepository.UpdateWalletAmountForbiddenMinus(queryGenralInfoByVendor.PlayerID, queryGenralInfoByVendor.GameID, payoutAmount, transactionType, queryGenralInfoByVendor.BetRecordNoPGS, "");
            }

            // step 5. 寫入perGameLog
            convertPerGame.CurrentAmount = updateResult.BeforeAmount;
            WritePerGameLog(convertPerGame);


            return;
        }


        public BetRecordNotifyQueueDetailInfoViewModel GetNotifyQueueInfo(int userID, long queueID)
        {
            BetRecordMerchantNotifyQueueEntity query = BetRecordMerchantNotifyQueueRepository.GetById(queueID);

            if (query == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            CheckViewMerchantPermission(userID, query.MerchantID);

            //會有多筆一樣的
            var queryLogs = SingleWalletMerchantApiRequestLogRepository.GetByFilter(new QuerySingleWalletMerchantApiRequestLogRequestViewModel
            {
                RequestID = query.LastRequestUUID
            });

            var viewModel = query.ReturnInjecteObject<BetRecordNotifyQueueDetailInfoViewModel>();
            viewModel.TypeName = ((BetRecordMerchantNotifyQueueTypeEnum)query.Type).ToString();

            viewModel.Logs = queryLogs.OrderByDescending(z=> z.CreateTime).ToList();

            return viewModel;
        }


        private void SetExtraProperty(List<UnionBetRecordViewModel> models)
        {
            List<string> transactionNoList = models.Select(z => z.BetRecordNoPGS).ToList();

            var queryQueue = BetRecordMerchantNotifyQueueRepository.GetByFilter(new QueryBetRecordMerchantNotifyQueueRequestViewModel
            {
                ListOfBetRecordNoPGS = transactionNoList
            });

            foreach (var item in models)
            {
                var match = queryQueue.FirstOrDefault(z => z.BetRecordNoPGS == item.BetRecordNoPGS);

                if (match != null)
                {
                    item.MerchantNotifyQueueID = match.Id;
                }
            }
        }

        private BetRecordGeneralColumnModel GetByPgsTransactionNo(int gameVendorID, string pgsTransactionNo)
        {
            var service = ServiceInstanceHelper.InstanceVendorService(gameVendorID);
            var queryByVendor = service.GetBetRecordByPGSTransactionNo(pgsTransactionNo);
            queryByVendor.StatusText = MultiLanguageTextHelper.GetEnumsCultureText((GameVendorBetRecordStatusEnum)queryByVendor.Status);
            return queryByVendor;
        }


        private UnionBetRecordViewModel ConvertToViewModels(UnionBetRecordQueryModel query)
        {
            var convert = query.ReturnInjecteObject<UnionBetRecordViewModel>();
            convert.GameName = DbMultiLanguageTextHelper.GetGameDisplayName(query.GameID);
            convert.GameVendorName = MultiLanguageTextHelper.GetGameVendorDisplayName(query.GameVendor);
            convert.StatusText = MultiLanguageTextHelper.GetEnumsCultureText((GameVendorBetRecordStatusEnum)query.Status);
            convert.MerchantName = MerchantRepository.GetSingleByCache(query.MerchantID)?.Name;
            convert.TypeName = MultiLanguageTextHelper.GetEnumsCultureText((BetRecordTypeEnum)query.Type);
            return convert;
        }

        private long WritePerGameLog(PerGameDetailLogEntity convertPerGame)
        {
            //step 1. 寫入perGame
            long perGameID = PerGameLogRepository.Create(new PerGameLogEntity
            {
                BeginTime = convertPerGame.GameTime,
                EndTime = convertPerGame.GameTime,
                GameData = "",
                ServerID = convertPerGame.ServerID,
                GameID = convertPerGame.GameID,
                TaxAmount = 0,
                CreateTime = DateTimeHelper.GetCurrentTime()
            });

            convertPerGame.PerGameLogID = perGameID;

            //step 2. 寫入perGameDetail
            long resultPerGameDetailID = PerGameDetailLogRepository.Create(convertPerGame);
            return resultPerGameDetailID;
        }

    }

}
