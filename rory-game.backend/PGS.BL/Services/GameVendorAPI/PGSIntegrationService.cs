using target.BL.Interface;
using target.DA.Repository;
using target.Models.Enums;
using target.Models.Models;
using target.Models.Table;
using target.Models.ViewModels;
using target.Models.ViewModels.Base;
using target.Models.ViewModels.GameVendor;
using target.Utility.Extensions;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Transactions;

namespace target.BL.Services
{

    public class PGSIntegrationService : BaseGameVendorIntegrationService, IGameVendorIntegrationService
    {

        public override GameVendorTypeEnum GameVendor
        {
            get
            {
                return GameVendorTypeEnum.target;
            }
        }


        public override string GameVendorPrefix
        {
            get
            {
                return "target";
            }
        }

        public override double TimeZone
        {
            get
            {
                return 0;
            }
        }


        public BetRecordRepository BetRecordRepository { get; set; }


        public SingleWalletExceptionHandleModel FormatParsing_ParseExceptionModelToVendorApiResponse(ApiResponseErrorModel model)
        {
            throw new TerproException(ApiErrorCode.NotImplement);
        }

        public object FormatParsing_ConvertToPlaceBetResponseResult(SingleWalletPlaceBetDebitAmountResponseParameter placeBetInfo)
        {
            throw new TerproException(ApiErrorCode.NotImplement);
        }


        public object FormatParsing_ConvertToSettleBetResponseResult(SettlementConvertRequiredInfoModel model)
        {
            throw new TerproException(ApiErrorCode.NotImplement);
        }


        public PerGameDetailLogEntity FormatParsing_ConvertBetRecordToPerGameDetailLog(string betRecordNo)
        {
            var queryBetRecord = BetRecordRepository.GetByTransactionNo(betRecordNo);

            if (queryBetRecord == null)
            {
                throw new TerproException(ApiErrorCode.NotFoundData);
            }

            if (queryBetRecord.Status != (int)GameVendorBetRecordStatusEnum.SettleComplete)
            {
                throw new TerproException(ApiErrorCode.NotFoundData);
            }

            var queryPlayer = PlayerRepository.GetBriefByCache(queryBetRecord.PlayerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            return new PerGameDetailLogEntity
            {
                BetAmount = queryBetRecord.BetAmount,
                ValidAmount = queryBetRecord.ValidAmount,
                UpdateAmount = queryBetRecord.WinLossAmount,
                PayoutAmount = queryBetRecord.PayoutAmount,
                TaxAmount = 0,
                ServerID = 0,
                GameID = queryBetRecord.GameID,
                CreateTime = DateTimeHelper.GetCurrentTime(),
                CurrentAmount = 0, //後半段的code再補上
                GameTime = queryBetRecord.SettlementTime.Value,
                MerchantID = queryBetRecord.MerchantID,
                PlayerID = queryBetRecord.PlayerID,
                PlayerName = queryPlayer.PlayerName,
                TradeID = 0,
                ThirdPartyBetRecordNo = queryBetRecord.BetRecordNoPGS,
                TransactionNo = queryBetRecord.BetRecordNoPGS
            };
        }


        public BetRecordGeneralColumnModel GetBetRecordByNo(string betRecordNo, string gameCode)
        {
            var queryRecord = BetRecordRepository.GetByTransactionNo(betRecordNo);
            return ConvertToGeneralModel(queryRecord);
        }


        public BetRecordGeneralColumnModel GetBetRecordByPGSTransactionNo(string betRecordNoPGS)
        {
            var queryRecord = BetRecordRepository.GetByTransactionNo(betRecordNoPGS);
            return ConvertToGeneralModel(queryRecord);
        }


        public void WriteRawBetRecordWhenPlaceBet(InsertRawBetRecordRequiredInfoModel model, EachGameVendorRawRequestInfoModel rawInfo)
        {
            throw new TerproException(ApiErrorCode.NotImplement);
        }




        public void UpdateRawBetRecordAfterSettlement(SettlementBetRecordRequiredInfoModel model, EachGameVendorRawRequestInfoModel rawInfo)
        {
            throw new TerproException(ApiErrorCode.NotImplement);
        }


        public void UpdateRawBetRecordAfterSettlement(SettlementBetRecordRequiredInfoModel model)
        {
            model.RuleFor(z => z.BetRecordNoPGS).Rule_NotEmptyText();

            DateTime now = DateTimeHelper.GetCurrentTime();
            string memo = $"processID:{model.DebugExtraMemo}";

            using (var scope = new TransactionScope())
            {
                var queryExist = BetRecordRepository.GetByTransactionNo(model.BetRecordNoPGS);

                //檢查狀態合法
                if (queryExist == null)
                    throw new TerproException(ApiErrorCode.NotFoundBetRecord);

                //如果狀態一樣 禁止操作(可能同時操作了)
                if (queryExist.Status == (int)GameVendorBetRecordStatusEnum.SettleComplete.ToInteger())
                {
                    throw new TerproException(ApiErrorCode.InvalidStatus);
                }

                queryExist.UpdateTime = now;
                queryExist.DebugInfo += RepositoryHelper.GenerateDebugInfo(queryExist.Status, memo);
                queryExist.Status = GameVendorBetRecordStatusEnum.SettleComplete.ToInteger();
                queryExist.PayoutAmount = model.PayoutAmount;
                queryExist.ValidAmount = model.ValidAmount;
                queryExist.WinLossAmount = model.PayoutAmount - queryExist.BetAmount;
                queryExist.SettlementTime = now;
                queryExist.ManualSettleUserID = model.UserID;

                BetRecordRepository.Update(queryExist);
                scope.Complete();
            }
        }

        public void UpdateBetRecordStatus(string transactionNoPGS, GameVendorBetRecordStatusEnum targetStatus, string memo)
        {
            DateTime now = DateTimeHelper.GetCurrentTime();
            using (var scope = new TransactionScope())
            {
                var queryExist = BetRecordRepository.GetByTransactionNo(transactionNoPGS);

                //檢查狀態合法
                if (queryExist == null)
                    throw new TerproException(ApiErrorCode.NotFoundBetRecord);

                //如果狀態一樣 禁止操作(可能同時操作了)
                if (queryExist.Status == (int)targetStatus)
                {
                    throw new TerproException(ApiErrorCode.InvalidStatus);
                }

                queryExist.UpdateTime = now;
                queryExist.DebugInfo += RepositoryHelper.GenerateDebugInfo(queryExist.Status, memo);
                queryExist.Status = targetStatus.ToInteger();

                BetRecordRepository.Update(queryExist);
                scope.Complete();
            }
        }


        public ExternalVendorCancelBetResultModel CancelBet(string betRecordNoPGS, string memo)
        {
            DateTime now = DateTimeHelper.GetCurrentTime();
            using (var scope = new TransactionScope())
            {
                var queryExist = BetRecordRepository.GetByTransactionNo(betRecordNoPGS);

                //檢查狀態合法
                RepositoryHelper.CheckBetRecordStatusIsWhiteList(queryExist, GameVendorBetRecordStatusEnum.Pending);

                queryExist.UpdateTime = now;
                queryExist.DebugInfo += RepositoryHelper.GenerateDebugInfo(queryExist.Status, memo);
                queryExist.Status = GameVendorBetRecordStatusEnum.Cancel.ToInteger();

                BetRecordRepository.Update(queryExist);
                scope.Complete();
            }

            //PGS不需實作return
            return new ExternalVendorCancelBetResultModel();
        }


        public void CheckIsValidRequestInfo(ExternalVendorInputValidateRequestParameter model)
        {
            throw new TerproException(ApiErrorCode.NotImplement);
        }


        public void CreateVendorPlayerAccountOrNot(ExternalVendorCreateAccountRequestParameter model)
        {
            throw new TerproException(ApiErrorCode.NotImplement);
        }


        public string GetSingleBetRecordByApi(ExternalVendorGetSingleBetRecordRequestParameter model)
        {
            throw new TerproException(ApiErrorCode.NotImplement);
        }



        public ExternalVendorLaunchUrlViewModel GetGameLaunchUrl(ExternalVendorLaunchUrlRequestModel model)
        {
            throw new TerproException(ApiErrorCode.NotImplement);
        }


        public ThirdPartyGameVendorBetRecordDetailUrlViewModel GetExternalBetDetailInfo(ExternalVendorBetRecordDetailInfoRequestModel model)
        {
            throw new TerproException(ApiErrorCode.NotImplement);
        }


        public List<ExternalVendorGameListConvertModel> GetGames(int merchantID)
        {
            throw new TerproException(ApiErrorCode.NotImplement);
        }


        public ExternalVendorGetReportResultModel GetDbAndApiSummaryReport(ExternalVendorGetReportRequestParameter model)
        {
            throw new TerproException(ApiErrorCode.NotImplement);
        }


        private BetRecordGeneralColumnModel ConvertToGeneralModel(BetRecordEntity queryRecord)
        {
            var convert = queryRecord.ReturnInjecteObject<BetRecordGeneralColumnModel>();
            return convert;
        }

    }
}
