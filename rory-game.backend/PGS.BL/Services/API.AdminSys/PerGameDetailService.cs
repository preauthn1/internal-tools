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
using example.invalid;
using target.Utility.Extensions.Web;

namespace target.BL.Services
{

    public class PerGameDetailService : BaseAdminSysService, IPerGameDetailService
    {
        public PerGameLogRepository PerGameLogRepository { get; set; }

        public PerGameDetailLogRepository PerGameDetailLogRepository { get; set; }

        public SlotGameHistoryLogRepository SlotGameHistoryLogRepository { get; set; }

        public ServerListRepository ServerListRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public PlayerRepository PlayerRepository { get; set; }

        public GamePayTableImageRepository GamePayTableImageRepository { get; set; }

        public BetRecordTransactionNoMappingRepository BetRecordTransactionNoMappingRepository { get; set; }



        public PagedList<PerGameDetailViewModel> GetByPagingFilter(int currentUserID, QueryPerGameDetailLogApiRequestParameterModel filter)
        {
            var currentUser = GetCurrentOwnPermissionByCache(currentUserID);

            filter.MerchantIDs = filter.MerchantID.ToIntegerListCollection();
            filter.MerchantIDs = FilterAvailablePermissionMerchantIDs(currentUserID, filter.MerchantIDs);

            //如果搜尋playerName，會先直接從player查詢 (效能考量)
            if (!string.IsNullOrWhiteSpace(filter.PlayerName))
            {
                var queryPlayers = PlayerRepository.GetByFilter(new QueryPlayerRequestViewModel { PlayerName = filter.PlayerName });
                filter.PlayerIDs = queryPlayers.Select(x => x.PlayerID).ToList();
            }

            //如果查單號 直接不篩選日期 (被日期誤導好多次了...)
            if(!string.IsNullOrWhiteSpace(filter.TransactionNo))
            {
                filter.GameStartTimeBegin = null;
                filter.GameStartTimeEnd = null;
            }

            //如果搜尋 單錢包投注單號
            bool isSearchTransactionNo = !string.IsNullOrWhiteSpace(filter.TransactionNo);
            if (isSearchTransactionNo)
            {
                var queryMapping = BetRecordTransactionNoMappingRepository.GetByTransactionNo(filter.TransactionNo);

                if (queryMapping == null)
                {
                    return new PagedList<PerGameDetailViewModel>();
                }
            }

            filter.GameIDs = filter.GameID.ToIntegerListCollection();
            filter.ServerIDs = filter.ServerID.ToIntegerListCollection();
            filter.PlayerIDs = filter.PlayerID.ToIntegerListCollection(filter.PlayerIDs);
            filter.PerGameLogIDs = filter.PerGameLogID.ToIntegerListCollection();

            var query = PerGameDetailLogRepository.GetPagingByFilter(filter);

            //若查得到mapping資訊 但perGameLog查無，跳出提示錯誤訊息
            if (isSearchTransactionNo && query.Items.IsEmptyOrNull())
            {
                throw new TerproException(ApiErrorCode.NotFoundPerGameLogButBetRecordExist);
            }

            var viewModels = ConvertToNewTypePagingObject<PerGameDetailLogEntity, PerGameDetailViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            return viewModels;
        }


        public PerGameSlotHistoryDetailInfoViewModel GetSlotHistoryByExternalLinkQuery(ExternalLinkQueryPerGameRequestModel model)
        {
            model.RuleFor(z => z.Token).Rule_NotEmptyText();
            model.RuleFor(z => z.TransactionNo).Rule_NotEmptyText();

            var json = model.Token.Base64Decode();

            //檢查token是否合法
            CheckIsValidSimpleToken(json , model.TransactionNo);

            var queryPerGame = PerGameDetailLogRepository.GetByFilter(new QueryDetailLogRequestViewModel
            {
                TransactionNo = model.TransactionNo

            }).FirstOrDefault();

            if (queryPerGame == null)
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);

            var result = GetSlotHistoryFromPerGameLog(SystemConstant.SkipValidUserID, queryPerGame.Id);
            return result;
        }

      
        public PerGameSlotHistoryDetailInfoViewModel GetSlotHistoryFromPerGameLog(int currentUserID, long perGameLogDetailID)
        {
            var queryDetail = PerGameDetailLogRepository.GetById(perGameLogDetailID);

            if (queryDetail == null)
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);

            if(currentUserID != SystemConstant.SkipValidUserID)
            {
                this.CheckViewMerchantPermission(currentUserID, queryDetail.MerchantID);
            }

            var queryMerchant = MerchantRepository.GetSingleByCache(queryDetail.MerchantID);

            var querySlotHistory = SlotGameHistoryLogRepository.GetByFilter(new QuerySlotGameHistoryLogRequestViewModel
            {
                GameID = queryDetail.GameID,
                PlayerID = queryDetail.PlayerID,
                GameTimeBegin = queryDetail.GameTime.AddSeconds(-3),
                GameTimeEnd = queryDetail.GameTime.AddSeconds(3)
            });

            if (querySlotHistory.IsEmptyOrNull())
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);

            //部分欄位沒有設index，程式端篩選
            var matchAmounts = querySlotHistory
                .Where(x => (x.WinAmount - x.BetAmount) == queryDetail.UpdateAmount)
                .OrderBy(x => ((long)(queryDetail.GameTime - x.CreateTime).TotalMilliseconds).ToAbs()) //時間最接近的一筆
                .ToList();

            List<SlotGameHistoryLogEntity> slotsResultList = new List<SlotGameHistoryLogEntity>();

            SlotGameHistoryLogEntity firstMatch = null;

            //如果比對不到資料，可能原因之一: 有多個盤面結果
            if (matchAmounts.IsEmptyOrNull())
            {
                slotsResultList = GetRelativeHistoryByRawDataCompare(queryDetail, querySlotHistory);
                firstMatch = slotsResultList.OrderByDescending(z => z.GameTime).FirstOrDefault();
            }
            else
            {
                firstMatch = matchAmounts.FirstOrDefault();
            }

            if (firstMatch == null)
                throw new TerproException(ApiErrorCode.NotFoundRecord);

            List<List<string>> slotLineResult = ParseSlotResultToArray(firstMatch.GameData);

            var lineImages = ConvertCodeToImages(firstMatch.GameID, slotLineResult);

            //找出這批結果 個別出現過的圖片
            var distinctIcons = lineImages.SelectMany(z => z).GroupBy(z => z.Code).Select(x => x.First()).OrderBy(z => z.Code).ToList();

            return new PerGameSlotHistoryDetailInfoViewModel
            {
                isMultipleSlotsResult = slotsResultList.Count() > 1,
                BeforeAmount = queryDetail.CurrentAmount,
                AfterAmount = queryDetail.CurrentAmount + queryDetail.UpdateAmount,
                LineResult = lineImages,
                DistinctIcons = distinctIcons,
                BetAmount = firstMatch.BetAmount,
                GameID = firstMatch.GameID,
                GameTime = firstMatch.GameTime,
                PlayerID = firstMatch.PlayerID,
                PlayerName = PlayerRepository.GetSingleByCache(firstMatch.PlayerID)?.PlayerName,
                PerGameDetailLogID = perGameLogDetailID,
                SlotGameHistoryLogID = firstMatch.ID,
                PerGameLogID = queryDetail.PerGameLogID,
                IsSingleWalletMode = CommonLogicHelper.IsSingleWalletMode(queryMerchant.WalletMode),
                PerGameTransactionNo = queryDetail.TransactionNo,
                TaxAmount = queryDetail.TaxAmount,
                UpdateAmount = queryDetail.UpdateAmount,
                GameName = DbMultiLanguageTextHelper.GetGameDisplayName(firstMatch.GameID)
            };
        }

        private List<SlotGameHistoryLogEntity> GetRelativeHistoryByRawDataCompare(PerGameDetailLogEntity queryDetail, List<SlotGameHistoryLogEntity> querySlotHistory)
        {
            List<SlotGameHistoryLogEntity> slotsResultList = new List<SlotGameHistoryLogEntity>();

            var queryPerGame = PerGameLogRepository.GetById(queryDetail.PerGameLogID);

            //解析找出有盤面結果的字串行
            List<string> lines = queryPerGame.GameData.Split(new[] { Environment.NewLine, "\t\n", "\n" }, StringSplitOptions.None)
                  .Where(x => !string.IsNullOrWhiteSpace(x))
                  .Where(x => x.Contains("#") && x.Contains("|"))
                  .ToList();

            foreach (var text in lines)
            {
                //perGame存的是 橫式array結果
                var slotsResultArray = ParseSlotResultToArray_HorizonMode(text);

                string toCompareText = ToArrayCompareString(slotsResultArray);

                foreach (var slotHistoryItem in querySlotHistory)
                {
                    var myResultArray = ParseSlotResultToArray(slotHistoryItem.GameData);

                    string compareTextOfSlot = ToArrayCompareString(myResultArray);

                    bool isEqual = toCompareText == compareTextOfSlot;

                    if (isEqual)
                    {
                        slotsResultList.Add(slotHistoryItem);
                    }
                }
            }

            return slotsResultList;
        }

        private string ToArrayCompareString(List<List<string>> array)
        {
            string result = string.Empty;

            foreach (var innerArray in array)
            {
                foreach (var text in innerArray)
                {
                    result += text;
                }
            }

            return result;
        }


        private string SpliTextByHashSymbol(string gameData)
        {
            var splitByHash = gameData.Split('#').ToList();

            string spinLineResult = splitByHash[0];

            if (splitByHash.Count() > 1)
            {
                spinLineResult = splitByHash[1];
            }

            return spinLineResult;
        }


        private List<List<string>> ParseSlotResultToArray(string gameData)
        {
            //先除掉hash符號
            string spinLineResult = SpliTextByHashSymbol(gameData);

            var allLines = spinLineResult.Split('|');

            List<List<string>> eachLineIconCode = new List<List<string>>();

            foreach (var item in allLines)
            {
                List<string> slotIcons = item.Split(',').Where(x => !string.IsNullOrEmpty(x)).Select(z => z.Trim()).ToList();
                eachLineIconCode.Add(slotIcons);
            }

            return eachLineIconCode;
        }


        /// <summary>
        /// 有些地方是存成 直式(先上下、再左右) / 橫式(先左右、再上下)
        /// 這邊會統一返回成直式的array
        /// </summary>
        /// <param name="gameData"></param>
        /// <returns></returns>
        private List<List<string>> ParseSlotResultToArray_HorizonMode(string gameData)
        {
            //先除掉hash符號
            string spinLineResult = SpliTextByHashSymbol(gameData);

            var allHorizonLines = spinLineResult.Split('|');

            List<List<string>> eachLineIconCode = new List<List<string>>();

            int howManyColumns = 0;

            //先計算總共有幾columns (幾排)
            foreach (var item in allHorizonLines)
            {
                List<string> slotIcons = item.Split(',').Where(x => !string.IsNullOrEmpty(x)).Select(z => z.Trim()).ToList();
                if (howManyColumns == 0)
                {
                    howManyColumns = slotIcons.Count();
                    break;
                }
            }

            //建好N排 空array
            for (int i = 0; i < howManyColumns; i++)
            {
                eachLineIconCode.Add(new List<string>());
            }

            //依序產生直式的array
            foreach (var item in allHorizonLines)
            {
                List<string> slotIcons = item.Split(',').Where(x => !string.IsNullOrEmpty(x)).Select(z => z.Trim()).ToList();

                int loopIndex = 0;
                foreach (string eachIcon in slotIcons)
                {
                    eachLineIconCode[loopIndex].Add(eachIcon);
                    loopIndex++;
                }
            }

            return eachLineIconCode;
        }




        private List<List<SlotHistoryResultIconViewModel>> ConvertCodeToImages(int gameID, List<List<string>> codes)
        {
            var queryImagesOfGame = GamePayTableImageRepository.GetByFilter(new QueryGamePayTableImageRequestViewModel
            {
                GameIDs = gameID.ToListCollection()
            });

            List<List<SlotHistoryResultIconViewModel>> lineImages = new List<List<SlotHistoryResultIconViewModel>>();

            foreach (var each_line in codes)
            {
                var items = new List<SlotHistoryResultIconViewModel>();

                foreach (var iconCode in each_line)
                {
                    var match = queryImagesOfGame.FirstOrDefault(x => x.Code == iconCode);

                    items.Add(new SlotHistoryResultIconViewModel
                    {
                        Code = iconCode,
                        ImageUrl = match?.ImageUrl
                    });
                }

                lineImages.Add(items);
            }

            return lineImages;
        }


        public PerGameLogFullInformationViewModel GetFullInformationInPerGame(int currentUserID, long perGameLogID)
        {
            var query = PerGameLogRepository.GetById(perGameLogID);

            var viewModel = query.ReturnInjecteObject<PerGameLogFullInformationViewModel>();

            var queryDetailLogs = PerGameDetailLogRepository.GetByFilter(new QueryDetailLogRequestViewModel
            {
                PerGameLogIDs = perGameLogID.ToListCollection()
            });

            viewModel.PlayerDetails = queryDetailLogs.Select(x => ConvertToViewModels(x)).ToList();

            var serverInfo = ServerListRepository.GetBriefByCache(query.ServerID);

            if (serverInfo != null)
            {
                viewModel.ServerName = serverInfo.ServerName;
            }

            viewModel.GameName = DbMultiLanguageTextHelper.GetGameDisplayName(query.GameID);
            return viewModel;
        }




        private PerGameDetailViewModel ConvertToViewModels(PerGameDetailLogEntity query)
        {
            var convert = query.ReturnInjecteObject<PerGameDetailViewModel>();

            var serverInfo = ServerListRepository.GetBriefByCache(convert.ServerID);

            if (serverInfo != null)
            {
                convert.ServerName = serverInfo.ServerName;
            }

            var matchMerchant = MerchantRepository.GetSingleByCache(query.MerchantID);

            convert.PerGameDetailLogID = query.Id;
            convert.GameName = DbMultiLanguageTextHelper.GetGameDisplayName(query.GameID);
            convert.MerchantName = matchMerchant.Name;
            convert.IsExternalVendorGame = AppHelper.IsExternalVendorGame(query.GameID);

            var queryGame = DbValueCacheHelper.GetGame(query.GameID);

            if (queryGame != null)
            {
                convert.VendorName = MultiLanguageTextHelper.GetGameVendorDisplayName(queryGame.GameVendorType);
                convert.ThirdPartyGameCode = queryGame.ThirdPartyGameCode;
                convert.IsSlotGame = queryGame.GameType == GameTypeEnum.Slot.ToInteger();
                convert.IsHideAfterAmountInfo = queryGame.IsHideAfterAmountInfo;
            }

            if (CommonLogicHelper.IsSingleWalletMode(matchMerchant.WalletMode) || queryGame.GameType == (int)GameTypeEnum.Fishing)
            {
                convert.IsHideAfterAmountInfo = true;
            }

            convert.AfterAmount = query.CurrentAmount + query.UpdateAmount;
            convert.IsSupportExternalLink = GameVendorDefinitionHelper.IsSupportExternalLinkUrl(queryGame.GameVendorType);

            return convert;
        }


        private static void CheckIsValidSimpleToken(string json , string transactionNo)
        {
            var convert = json.JsonConvertToModel<SimpleExternalTokenModel>();

            if (convert == null)
                throw new TerproException(ApiErrorCode.InvalidParameter_ParsingError);

            //禁止使用其他組token 查詢其他注單號
            //if(convert.TransactionNo != transactionNo)
            //    throw new TerproException(ApiErrorCode.InvalidParameter_ParsingError);

            var sinKey = AppHelper.GenerateSimpleTokenModelSignKey(convert);

            if (convert.Sign != sinKey)
                throw new TerproException(ApiErrorCode.InvalidParameter_SignKeyError);

            if (convert.ExpiredTime <= DateTimeHelper.GetCurrentTime())
                throw new TerproException(ApiErrorCode.InvalidParameter_SignKeyError);
        }

    }
}
