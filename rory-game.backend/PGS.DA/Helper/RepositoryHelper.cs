using target.Models.Models;
using System;
using System.Collections.Generic;
using System.Text;
using target.Utility.Extensions;
using target.Models.Enums;
using target.Models.Table;

namespace target.DA.Repository
{
    public class RepositoryHelper
    {
        public static void CheckBetRecordStatusIsWhiteList<T>(T data, GameVendorBetRecordStatusEnum statusRules) where T : IBetRecordEntity
        {
            CheckBetRecordStatusIsWhiteList(data , new List<int> { statusRules.ToInteger() });
        }

        public static void CheckBetRecordStatusIsBlackList<T>(T data, GameVendorBetRecordStatusEnum statusRules) where T : IBetRecordEntity
        {
            CheckBetRecordStatusIsBlackList(data, new List<int> { statusRules.ToInteger() });
        }


        public static void CheckBetRecordStatusIsWhiteList<T>(T data, List<int> statusRules) where T : IBetRecordEntity
        {
            //檢查是否為 [指定允許] 的狀態
            if (!statusRules.Contains(data.Status))
            {
                throw new TerproException(ApiErrorCode.InvalidStatus);
            }
        }

        public static void CheckBetRecordStatusIsBlackList<T>(T data, List<int> statusRules) where T : IBetRecordEntity
        {
            //檢查是否為 [禁止] 的狀態
            if (statusRules.Contains(data.Status))
            {
                throw new TerproException(ApiErrorCode.InvalidStatus);
            }
        }

        public static void CheckRecordNotFoundAndDuplicateStatus<T>(T data, GameVendorBetRecordStatusEnum targetStatus, decimal? targetPayoutAmount) where T: IBetRecordEntity
        {
            //傳入null 代表查無資料
            if (data == null)
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);

            //如果狀態一樣 禁止操作(可能同時操作了)
            if (data.Status == (int)targetStatus)
            {
                var error = new ErrorDetailModel
                {

                };

                //如果重複發通知了 但是金額不一樣
                if (targetPayoutAmount.HasValue && data.PayoutAmount != targetPayoutAmount)
                {
                    error.IsAlwaysNotify = true;
                    error.LogExtraInfo = $"當前 PayoutAmount: {data.PayoutAmount}，預計 PayoutAmount: {targetPayoutAmount}";
                }

                throw new TerproException(ApiErrorCode.InvalidStatus_BetRecordDuplicateSettlement, error);
            }                
        }

        public static string GenerateDebugInfo(int currentStatus, string inputMemo)
        {
            return string.Format(" Before status: {0} . {1}", currentStatus, inputMemo);
        }


    }
}
