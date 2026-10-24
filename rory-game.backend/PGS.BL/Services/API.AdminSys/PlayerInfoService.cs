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
using game.Entity.Game;
using target.Utility;
using Lck.Utility.Extensions.Helper;

namespace target.BL.Services
{

    public class PlayerInfoService : BaseAdminSysService, IPlayerInfoService
    {
        public PlayerWalletRepository PlayerWalletRepository { get; set; }

        public PlayerRepository PlayerRepository { get; set; }

        public PlayerGameSummaryRepository PlayerGameSummaryRepository { get; set; }

        public PlayerTotalSummaryRepository PlayerTotalSummaryRepository { get; set; }

        public PlayerCurrentServerRepository PlayerCurrentServerRepository { get; set; }

        public ServerListRepository ServerListRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public UserOperateRecordRepository UserOperateRecordRepository { get; set; }

        public SystemConfigSettingRepository SystemConfigSettingRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }

        public IPlayerCreateRefService PlayerCreateRefService { get; set; }


        public PagedList<PlayerInfoViewModel> GetByPagingFilter(int userID , QueryPlayerApiRequestParameterModel filter)
        {
            filter.MerchantIDs = filter.MerchantID.ToIntegerListCollection();
            filter.MerchantIDs = FilterAvailablePermissionMerchantIDs(userID, filter.MerchantIDs);
            filter.PlayerIDs = filter.PlayerID.ToIntegerListCollection();
            filter.StatusValues = filter.Status.ToIntegerListCollection();

            var query = PlayerRepository.GetPagingByFilter(filter);            

            var paging = ConvertToNewTypePagingObject<PlayerEntity, PlayerInfoViewModel>(query);

            paging.Items = query.Items
                .Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();


            SetPlayerWalletAmount(paging.Items);
            return paging;
        }


        public void CreateInernalTestAccount(int userID, CreateInternalTestAccountRequestModel model)
        {
            model.RuleFor(z => z.PlayerName).Rule_NotEmptyText().Rule_IsValidAccountFormat();
            model.RuleFor(z => z.MerchantID).Rule_NotZeroNotMinus();

            var queryMerchant = MerchantRepository.GetSingleByCache(model.MerchantID);

            int playerID = PlayerCreateRefService.CreatePlayerForMerchantRequest(new PlayerCreateParameterModel
            {
                Currency = queryMerchant.Currency,
                MerchantID = queryMerchant.MerchantID,
                PlayerName = model.PlayerName
            });

            if(model.IsAddToInternalTestAccount)
            {
                var querySetting = SystemConfigSettingRepository.GetById(SystemConfigSettingEnum.InternalTestAccountPlayerIDs.ToInteger());
                querySetting.SettingValue = querySetting.SettingValue + "," + playerID.ToString();

                //加入內部測試用帳號
                SystemConfigSettingRepository.Update(querySetting);
            }

            return;
        }


        public PlayerManualDepositResultViewModel InternalTestAccountWalletDepositAmount(int userID, InternalTestAccountPlusWalletRequestModel model)
        {
            model.RuleFor(z => z.PlayerID).Rule_NotZeroNotMinus();

            var queryPlayer = PlayerRepository.GetSingleByCache(model.PlayerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            bool isTestAccountOrMerchant = IsInternalTestAccountPlayerOrMerchantPlayer(model.PlayerID);

            if (!isTestAccountOrMerchant)
            {
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_SettingDataNotCompleted, new ErrorDetailModel("InternalTestAccountPlayerIDs"));
            }

            var result = ManualDepositPlayerWallet(userID, new PlayerManualDepositRequestModel
            {
                Amount = model.Amount,
                Memo = "內部測試帳號加值",
                PlayerID = queryPlayer.PlayerID,
                PlayerName = queryPlayer.PlayerName
            });

            return result;
        }


        public PlayerManualDepositResultViewModel ManualDepositPlayerWallet(int userID, PlayerManualDepositRequestModel model)
        {
            model.RuleFor(z => z.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(z => z.PlayerName).Rule_NotEmptyText();
            model.RuleFor(z => z.Amount).Rule_NotZeroNotMinus();

            var currentUser = GetCurrentOwnPermissionByCache(userID);

            if (!currentUser.IsSuperAdminOrg)
                throw new TerproException(ApiErrorCode.Forbidden_PagePermissionNotAllow);

            var queryPlayer = PlayerRepository.GetById(model.PlayerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            if (queryPlayer.PlayerName != model.PlayerName)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            string ip = IPHelper.GetIP();

            var queryMerchant = MerchantRepository.GetById(queryPlayer.MerchantID);

            if (queryMerchant.WalletMode != MerchantWalletModeEnum.TransferMode.ToInteger())
                throw new TerproException(ApiErrorCode.BusinessLogicDeny, "手動補錢，僅針對轉帳模式玩家");

            // Step 1. 更新錢包
            var spResult = PlayerWalletRepository.UpdateWalletAmountForbiddenMinus(queryPlayer.PlayerID, 0, model.Amount, TransactionTypeEnum.ManualDepositWalletAmount, "", ip);
                        
            UserOperateRecordEntity operateLog = InstanceOperateLog(userID, UserOperateTypeEnum.ManualDepositWalletAmount, queryPlayer.PlayerID);
            operateLog.MerchantID = queryMerchant.MerchantID;
            operateLog.Memo = "update amount:" + model.Amount + ". " + model.Memo;

            // Step 2. 寫操作紀錄
            UserOperateRecordRepository.Create(operateLog);

            return new PlayerManualDepositResultViewModel
            {
                AfterAmount = spResult.AfterAmount,
                BeforeAmount = spResult.BeforeAmount
            };
        }


        public PlayerInfoSummaryViewModel GetPlayerSummaryInfoByFilter(PlayerInfoSummaryRequestModel model)
        {
            if(model.PlayerID > 0)
            {
                var query = GetPlayerSummaryInfo(model.PlayerID);
                return query;
            }
            else
            {
                return null;
            }
        }

        public int GetPlayerIdByPlayerNameOrID(int userID , string playerNameOrID)
        {
            if (string.IsNullOrWhiteSpace(playerNameOrID))
                throw new TerproException(ApiErrorCode.InvalidParameter);

            //非數字
            if(!playerNameOrID.IsNumeric())
            {
                return FindPlayerByString(userID, playerNameOrID);
            }

            int parseToNumber = Int32.Parse(playerNameOrID);

            if(parseToNumber > 1000000 && parseToNumber < 1300000)
            {
                //優先從playerID去查
                var query = PlayerRepository.GetById(parseToNumber);

                if (query != null)
                    return query.PlayerID;

                return FindPlayerByString(userID , playerNameOrID);
            }
            else
            {
                return FindPlayerByString(userID, playerNameOrID);
            }
        }

     
        public PlayerInfoSummaryViewModel GetPlayerSummaryInfo(int playerID)
        {
            var queryPlayer = PlayerRepository.GetById(playerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            var convert = queryPlayer.ReturnInjecteObject<PlayerInfoSummaryViewModel>();
            var queryMerchant = MerchantRepository.GetSingleByCache(queryPlayer.MerchantID);

            if(queryMerchant != null)
            {
                convert.MerchantName = queryMerchant.Name;
                convert.WalletMode = queryMerchant.WalletMode;
            }

            convert.LastLoginTime = queryPlayer.LastLoginTime.HasValue ? queryPlayer.LastLoginTime.Value : queryPlayer.CreateTime;
            var queryTotalSummary = PlayerTotalSummaryRepository.GetById(playerID);

            var queryGameSummary = PlayerGameSummaryRepository.GetByFilter(new QueryPlayerGameSummaryRequestViewModel 
            { 
                PlayerIDs = playerID.ToListCollection()
            });

            convert.TotalSummary = new PlayerTotalSummaryViewModel
            {
                TotalLoginTimes = queryTotalSummary?.SumLoginTimes ?? 0
            };

            convert.PlayerGameSummaries = queryGameSummary.Select(x =>
                {
                    var convertTo = x.ReturnInjecteObject<PlayerGameSummaryViewModel>();
                    convertTo.GameName = DbMultiLanguageTextHelper.GetGameDisplayName(x.GameID);
                    convertTo.RTP = CommonLogicHelper.CalculateRTP_Percentage(x.BetAmount, x.PayoutAmount);

                    return convertTo;
                })
                .ToList();

            convert.PlayerTypeText = MultiLanguageTextHelper.GetEnumsCultureText((PlayerTypeEnum)queryPlayer.PlayerType);
            convert.WalletAmount = PlayerWalletRepository.GetPlayerWalletByType(playerID, WalletTypeEnum.GameWallet).Amount;
            convert.PlayerGameSummaryRowTotal = convert.PlayerGameSummaries.CalculateTotalSumMergeToOneObject();

            var queryMyCurrentServer = PlayerCurrentServerRepository.GetByPlayerID(queryPlayer.PlayerID);

            if(queryMyCurrentServer != null)
            {
                var queryServer = ServerListRepository.GetBriefByCache(queryMyCurrentServer.ServerID);
                convert.PlayerCurrentGameGame = DbMultiLanguageTextHelper.GetGameDisplayName(queryServer.GameID);
                convert.PlayerCurrentServerID = queryMyCurrentServer.ServerID;
            }

            convert.IsTestAccountOrMerchant = IsInternalTestAccountPlayerOrMerchantPlayer(queryPlayer.PlayerID);

            //有投注的遊戲才會RTP算進平均值
            convert.PlayerGameSummaryRowTotal.RTP = convert.PlayerGameSummaries
                .Where(x => x.BetAmount > 0)
                .Select(x=> x.RTP)
                .DefaultIfEmpty()
                .Average();

            return convert;
        }


        
        public void UpdatePlayerStatus(int userID , PlayerUpdateStatusRequestModel model)
        {
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.Status).Rule_NotMinus();

            UpdateMultiplePlayerStatus(userID, new PlayerMultipleUpdateStatusRequestModel
            {
                PlayerIDs = model.PlayerID.ToListCollection(),
                Status = model.Status,
                Memo = model.Memo
            });

        }


        public void UpdateMultiplePlayerStatus(int userID , PlayerMultipleUpdateStatusRequestModel model)
        {
            if (model.PlayerIDs.IsEmptyOrNull())
                return;

            int targetStatus = model.Status;

            var queryPlayers = PlayerRepository.GetByIds(model.PlayerIDs);

            if (queryPlayers.IsEmptyOrNull())
                throw new TerproException(ApiErrorCode.NotFoundPlayer);            

            var currentUser = GetCurrentOwnPermissionByCache(userID);

            foreach (var eachPlayer in queryPlayers)
            {
                //如果有權限不符合的(update別家商戶) 跳過不處理
                if (!currentUser.IsSuperAdminOrg && 
                    !currentUser.AvailableMerchantIDs.Contains(eachPlayer.MerchantID))
                    continue;

                var changeLog = new UserOperateChangeLogModel
                {
                    BeforeObject = new { Status = eachPlayer.Status },
                    AfterObject = new { Status = targetStatus }
                };

                eachPlayer.Status = targetStatus;

                //Step 1. 更新狀態
                PlayerRepository.UpdateStatus(eachPlayer.PlayerID , targetStatus);

                //Step 2. 寫異動紀錄
                var operatorType = (targetStatus == PlayerStatusEnum.Enable.ToInteger() ? UserOperateTypeEnum.Player_UnlockStatus : UserOperateTypeEnum.Player_LockStatus);
                UserOperateRecordEntity operateLog = InstanceOperateLog(userID, operatorType, eachPlayer.PlayerID, changeLog);
                operateLog.MerchantID = eachPlayer.MerchantID;
                operateLog.Memo = string.IsNullOrWhiteSpace(model.Memo) ? operateLog.Memo : model.Memo;

                UserOperateRecordRepository.Create(operateLog);

                CacheRemoveHelper.RemoveCache_WhenUpdatePlayer(eachPlayer.PlayerID);
            }
        }

        public void RemovePlayerRoomLock(int userID , int playerID)
        {
            var query = PlayerCurrentServerRepository.GetById(playerID);

            if (query == null)
                return;

            //Step 1. 刪除資料
            PlayerCurrentServerRepository.Delete(playerID);

            //Step 2. 寫異動紀錄
            var changeLog = new UserOperateChangeLogModel
            {
                BeforeObject = query,
                AfterObject = null
            };

            UserOperateRecordEntity operateLog = InstanceOperateLog(userID, UserOperateTypeEnum.Player_RoomLockRemove, playerID , changeLog);
            UserOperateRecordRepository.Create(operateLog);
        }

        private int FindPlayerByString(int userID, string playerNameOrID)
        {
            var currentUser = GetCurrentOwnPermissionByCache(userID);
            var query = PlayerRepository.GetByFilter(new QueryPlayerRequestViewModel
            {
                PlayerName = playerNameOrID.ToLower(),
                MerchantIDs = currentUser.SuperAdminOrAvailableMerchantIDs
            });

            if (query.IsEmptyOrNull())
            {
                throw new TerproException(ApiErrorCode.NotFoundPlayer);
            }

            if (query.Count() > 1)
            {
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_MatchDataIsNotSingleOne);
            }

            return query.FirstOrDefault().PlayerID;
        }

        private void SetPlayerWalletAmount(List<PlayerInfoViewModel> models)
        {
            List<int> playerIDs = models.Select(x => x.PlayerID).ToList();

            var query = PlayerWalletRepository.GetByFilter(new QueryPlayerWalletRequestViewModel
            {
                PlayerIDs = playerIDs
            });

            foreach (var item in models)
            {
                var match = query.FirstOrDefault(x => x.PlayerID == item.PlayerID);

                if (match != null)
                {
                    item.WalletAmount = match.Amount;
                }
            }
        }
        

        private PlayerInfoViewModel ConvertToViewModels(PlayerEntity query)
        {
            var convert = query.ReturnInjecteObject<PlayerInfoViewModel>();
            var queryMerchant = MerchantRepository.GetSingleByCache(query.MerchantID);

            if(queryMerchant != null)
            {
                convert.MerchantName = queryMerchant.Name;
                convert.WalletMode = queryMerchant.WalletMode;
            }

            return convert;
        }


    }
}
