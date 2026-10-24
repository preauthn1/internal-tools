using System;
using System.Collections.Generic;
using MSC.Models.CoreLibs;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Models;
using System.Linq;
using MSC.Utility.Helper;
using MSC.Models;
using MSC.DA.Repository;
using MSC.Utility;
using Lck.Utility.Extensions.Helper;

namespace MSC.BL.Services
{

    public class CardInputService : BaseService, ICardInputService
    {

        public ITokenCheckRefService TokenCheckRefService { get; set; }

        public INotifyCardDepositService NotifyCardDepositService { get; set; }

        public IRequestForbiddenRefService RequestForbiddenRefService { get; set; }


        public CardInfoRepository CardInfoRepository { get; set; }

        public CardLinkRequestRecordRepository CardLinkRequestRecordRepository { get; set; }

        public CardConsumeRecordRepository CardConsumeRecordRepository { get; set; }

        public SystemSettingConfigRepository SystemSettingConfigRepository { get; set; }

        public RequestForbiddenLogRepository RequestForbiddenLogRepository { get; set; }

        public MerchantWalletRepository MerchantWalletRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }



        public CardInputViewModel InputCard(CardInputRequestModel model)
        {
            model.RuleFor(x => x.MerchantToken).Rule_NotEmptyText().Rule_MinLengthEnough(30);
            model.RuleFor(x => x.SignKey).Rule_NotEmptyText();
            model.RuleFor(x => x.TS).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.CardNo).Rule_NotEmptyText();
            model.RuleFor(x => x.CardPassword).Rule_NotEmptyText().Rule_MinLengthEnough(16).Rule_IsLetterAndNumberOnly();

            AddIPAddressCounterTimes();

            AddTokenCounterTimes(model.MerchantToken);

            //檢查token(格式的部分)
            LinkTokenMemberInfoModel tokenParsing = TokenCheckRefService.CheckIsTokenFormatCorrect(model.MerchantToken);

            //檢查IP 存取太多次
            CheckIPAddressTooManyTimes(tokenParsing);

            //檢查token 存取太多次
            CheckTokenHashTooManyTimes(tokenParsing, model.MerchantToken);


            //檢查卡號
            CardNumberHelper.CheckCardIsValid(model.CardNo);

            //檢查signKey加密
            CheckSignKeyIsValid(model, tokenParsing);


            //以下開始是需要透過 [DB] 檢查------------------------------     
            //以下開始是需要透過 [DB] 檢查------------------------------     
            //以下開始是需要透過 [DB] 檢查------------------------------     

            //系統維護開關
            CheckSystemEnableFlag();

            //檢查商戶狀態開關
            CheckMerchantIsValid(tokenParsing.MerchantID);

            //檢查token有效 (db、cache層檢查)
            CardLinkRequestRecordEntity queryRequestInfo = CheckMerchantRequestTokenIsValid(model);

            //主要商業邏輯處理在此
            //主要商業邏輯處理在此
            //主要商業邏輯處理在此
            //主要商業邏輯處理在此
            //主要商業邏輯處理在此
            var result = InputCardBusinessProcess(model, queryRequestInfo);
            return result;
        }


        /// <summary>
        /// 內部測試使用，隨機選一張卡輸入 (正式環境禁止使用)
        /// </summary>
        public CardInputViewModel FakeTestRandomInputCard(CardInputFakeTestRequestModel model)
        {
         
            //從db直接撈指定的 [支付請求單號] 做為測試使用 (非正式流程)
            var queryRequestInfo = CardLinkRequestRecordRepository.GetByMerchantIDAndRequestID(model.MerchantID, model.MerchantRequestID);

            if (queryRequestInfo == null)
                throw new MSCException(ErrorCode.NotFoundRequestData);

            //token已過期/禁用
            if (queryRequestInfo.IsDisabled)
                throw new MSCException(ErrorCode.ThirdPartyApiFail_TokenExpired);

            var queryRandomCards = CardInfoRepository.GetPagingByFilter(new QueryCardInfoRequestViewModel
            {
                Point = queryRequestInfo.Amount > 0 ? (decimal?)queryRequestInfo.Amount : null,
                CardStatus = CardInfoStatusEnum.Pending,
                PageSize = 100
            });

            //隨便自動使用一張卡號
            var randomCard = queryRandomCards.Items.OrderBy(x => Guid.NewGuid()).FirstOrDefault();

            var fakeModel = new CardInputRequestModel
            {
                TS = DateTimeHelper.GetCurrentTimeTimestamp(),
                CardNo = randomCard.CardNo,
                CardPassword = [REDACTED],
                MerchantToken = [REDACTED]",
                SignKey = "XXXXXXXX",
            };

            var result = InputCardBusinessProcess(fakeModel, queryRequestInfo);
            return result;
        }


        /// <summary>
        /// 內部測試使用，輸入一張指定的卡號使用 (跟正式API差別在於 merchantToken是直接用ID指定)
        /// </summary>
        public CardInputViewModel FakeTestInputTargetCardNumber(CardInputFakeTestTargetCardNumberRequestModel model)
        {
            model.RuleFor(z => z.CardNo).Rule_NotEmptyText();
            model.RuleFor(z => z.CardPassword).Rule_NotEmptyText();
            model.RuleFor(z => z.MagicPassword).Rule_NotEmptyText();

            if (model.MagicPassword != SystemConstant.VendorAPI_SignKey_MagicPassword)
                throw new MSCException(ErrorCode.BusinessLogicDeny);

            var queryCard = CardInfoRepository.GetByCardNo(model.CardNo);

            if (queryCard == null)
                throw new MSCException(ErrorCode.NotFoundCardNo_FormatIncorrect);

            if (queryCard.Point >= 500) 
                throw new MSCException(ErrorCode.BusinessLogicDeny , "Fake API 卡片面額不得輸入太大");

            //從db直接撈指定的 [支付請求單號] 做為測試使用 (非正式流程)
            var queryRequestInfo = CardLinkRequestRecordRepository.GetByMerchantIDAndRequestID(model.MerchantID, model.MerchantRequestID);

            if (queryRequestInfo == null)
                throw new MSCException(ErrorCode.NotFoundRequestData);

            //token已過期/禁用
            if (queryRequestInfo.IsDisabled)
                throw new MSCException(ErrorCode.ThirdPartyApiFail_TokenExpired);

            var fakeModel = new CardInputRequestModel
            {
                TS = DateTimeHelper.GetCurrentTimeTimestamp(),
                CardNo = model.CardNo,
                CardPassword = [REDACTED],
                MerchantToken = [REDACTED]",
                SignKey = "XXXXXXXX",
            };

            var result = InputCardBusinessProcess(fakeModel, queryRequestInfo);
            return result;
        }

        private CardInputViewModel InputCardBusinessProcess(CardInputRequestModel model, CardLinkRequestRecordEntity queryRequestInfo)
        {
            //檢查卡號存在、是否已使用
            var result = CardInfoRepository.CheckIsValidCard(model.CardNo, queryRequestInfo.MerchantID);

            if (result.Code != ErrorCode.Success)
            {
                throw new MSCException(result.Code);
            }

            var queryCard = CardInfoRepository.GetByCardNo(model.CardNo);

            //檢查password
            if (queryCard.CardPassword != model.CardPassword)
                throw new MSCException(ErrorCode.NotFoundCardNo_OrPassword);

            if (queryCard.ConsumeMerchantID == SystemConstant.MagicConsumeMerchantID)
                throw new MSCException(ErrorCode.Forbidden_Card_CardStatusIsDestroy);

            var queryMerchant = MerchantRepository.GetSingleByCache(queryRequestInfo.MerchantID);

            if(queryMerchant.Currency != queryCard.Currency)
                throw new MSCException(ErrorCode.Forbidden_Card_CardCurrencyNotSupport);

            //如果訂單有金額，輸入的卡號面額要檢查是否一致
            CheckCardAmountAndOrderAmountMatch(queryRequestInfo, queryCard);

            string orderNo = RandomHelper.GetOrderNo("CI");

            DateTime now = DateTimeHelper.GetCurrentTime();

            //step 1. 寫入一筆紀錄
            //step 2. 更新卡號已被使用
            var resultOfConsume = CardConsumeRecordRepository.TransactionWriteConsumeRecord(new CardConsumeRecordEntity
            {
                DeviceID = HeaderValueHelper.GetDeviceID(),
                CardNo = model.CardNo,
                CreateTime = now,
                IPAddress = IPHelper.GetIP(),
                MemberCode = queryRequestInfo.MemberCode,
                MerchantID = queryRequestInfo.MerchantID,
                MerchantRequestID = queryRequestInfo.MerchantRequestID,
                OrderNo = orderNo,
                Point = queryCard.Point
            });

            if(resultOfConsume.Code != ErrorCode.Success)
            {
                throw new MSCException(resultOfConsume.Code);
            }

            var queryNewRecord = CardConsumeRecordRepository.GetByFilter(new QueryCardConsumeRecordRequestViewModel
                {
                    CardNo = model.CardNo
                })
                .FirstOrDefault();

            long resultID = queryNewRecord.ID;

            //step 3 更新商戶資金
            MerchantWalletRepository.UpdateWalletAndWriteLog(new UpdateMerchantWalletSPRequestParameter
            {
                MerchantID = queryRequestInfo.MerchantID,
                TransactionType = TransactionTypeEnum.UserCardDeposit,
                FeeRate = GetFeeRate(queryRequestInfo.MerchantID),
                RefNo = orderNo,
                SourceAmount = queryCard.Point
            });

            //step 4. token禁止再次使用
            CardLinkRequestRecordRepository.UpdateToDisabled(queryRequestInfo.Id);

            //step 5. 通知到帳成功
            NotifyCardDepositService.WriteQueueAndTriggerApiImmediately(new NotifyCardDepositRequestModel
            {
                CardConsumeRecordID = resultID,
                NotifyApiUrl = queryRequestInfo.NotifyUrl
            });

            //step 6. 紀錄過期的token
            string cacheKey = CacheKey.PointAPI_TokenExpired(model.MerchantToken.MD5Hash());
            CacheHelper.SetCacheValue(cacheKey, "true", 120);


            return new CardInputViewModel
            {
                OrderNo = orderNo,
                ReturnUrl = queryRequestInfo.ReturnUrl
            };
        }

        private CardLinkRequestRecordEntity CheckMerchantRequestTokenIsValid(CardInputRequestModel model)
        {
            string tokenHash = model.MerchantToken.MD5Hash();

            //檢查token失效 (從cache檢查)
            TokenCheckRefService.CheckTokenDisabledCache(tokenHash);

            //檢查db存在
            var queryRequestInfo = CardLinkRequestRecordRepository.GetByTokenHash(model.MerchantToken.MD5Hash());

            if (queryRequestInfo == null)
                throw new MSCException(ErrorCode.NotFoundRequestData);

            //token已過期/禁用
            if (queryRequestInfo.IsDisabled)
                throw new MSCException(ErrorCode.ThirdPartyApiFail_TokenExpired);

            return queryRequestInfo;
        }

        private void CheckSystemEnableFlag()
        {
            var querySetting = SystemSettingConfigRepository.GetSingleByCache(SystemSettingConfigEnum.System_SystemEnableFlag);

            if (!querySetting.SettingValueOfBoolean)
                throw new MSCException(ErrorCode.System_SystemMaintain);
        }


        private void CheckMerchantIsValid(int merchantID)
        {
            var queryMerchant = MerchantRepository.GetSingleByCache(merchantID);

            if (queryMerchant == null)
                throw new MSCException(ErrorCode.NotFoundMerchant);

            if (queryMerchant.Status != NormalActiveStatus.Enable.ToInteger())
                throw new MSCException(ErrorCode.InvalidStatus);
        }



        private decimal GetFeeRate(int merchantID)
        {
            var setting = MerchantCustomizedSettingRepository.GetByMerchantSettingKey(merchantID, MerchantCustomizedSettingKeyEnum.DefaultTransactionFeeRate);

            if (setting == null)
                return 0.15M;

            return setting.SettingValueOfDecimal;
        }



        private static void CheckCardAmountAndOrderAmountMatch(CardLinkRequestRecordEntity queryRequestInfo, CardInfoEntity queryCard)
        {
            bool isRestrictAmount = queryRequestInfo.Amount > 0;

            if (isRestrictAmount)
            {
                if (queryCard.Point != queryRequestInfo.Amount)
                {
                    throw new MSCException(ErrorCode.Forbidden_CardPoint_OrderAmount_NotMatch);
                }
            }
        }

        private void CheckSignKeyIsValid(CardInputRequestModel model, LinkTokenMemberInfoModel token)
        {
            var expect_signKey = (model.CardNo + model.TS.ToString() + token.MSCRequestID).MD5Hash();

            //如果採用萬用密碼
            if (model.SignKey == SystemConstant.API_Point_MagicCode)
            {
                return;
            }

            //前端傳來的ts可能並非精確時間(或偽造時間) 因此不檢查ts時間過久

            if (model.SignKey != expect_signKey)
            {
                throw new MSCException(ErrorCode.InvalidParameter_SignKeyError);
            }
        }


        private void AddIPAddressCounterTimes()
        {
            string myIP = IPHelper.GetIP();
            RequestForbiddenRefService.AddForbiddenCounterTimes(RequestForbiddenTypeEnum.Request_CardInput_CheckCardNumber_IPAddress, myIP);
        }


        private void AddTokenCounterTimes(string merchantToken)
        {
            string tokenHash = merchantToken.MD5Hash();
            RequestForbiddenRefService.AddForbiddenCounterTimes(RequestForbiddenTypeEnum.Request_CardInput_CheckCardNumber_TokenHash, tokenHash);
        }


        private void CheckIPAddressTooManyTimes(LinkTokenMemberInfoModel token)
        {
            RequestForbiddenTypeEnum forbiddenType = RequestForbiddenTypeEnum.Request_CardInput_CheckCardNumber_IPAddress;

            string myIP = IPHelper.GetIP();

            CheckForbiddenCounterAndException(token, forbiddenType, myIP);
        }


        private void CheckTokenHashTooManyTimes(LinkTokenMemberInfoModel token , string merchantToken)
        {
            RequestForbiddenTypeEnum forbiddenType = RequestForbiddenTypeEnum.Request_CardInput_CheckCardNumber_TokenHash;

            string tokenHash = merchantToken.MD5Hash();

            CheckForbiddenCounterAndException(token, forbiddenType, tokenHash);
        }


        private void CheckForbiddenCounterAndException(LinkTokenMemberInfoModel token, RequestForbiddenTypeEnum forbiddenType, string forbiddenKeyOrIP)
        {
            int maxTimes = GetForbiddenMaxTimesSetting();

            var checkResult = RequestForbiddenRefService.CheckForbiddenCounterTooManyTimes(forbiddenType, forbiddenKeyOrIP, maxTimes);

            if (checkResult.IsFobidden)
            {
                if (checkResult.IsNewForbiddenEvent)
                {
                    //首次被封鎖那次事件 寫入log (之後則是判斷cache存在 直接拋exception)
                    RequestForbiddenLogRepository.Create(new RequestForbiddenLogEntity
                    {
                        Type = forbiddenType.ToInteger(),
                        CreateTime = DateTimeHelper.GetCurrentTime(),
                        MemberCode = token.MemberCode,
                        MerchantID = token.MerchantID,
                        BlockKey = checkResult.ForbiddenCacheKey
                    });
                }

                throw new MSCException(ErrorCode.BusinessLogicDeny_RequestTooManyTimes);
            }
        }


        private int GetForbiddenMaxTimesSetting()
        {
            //查出DB設置 (超過多少次要拋exception)
            var querySetting = SystemSettingConfigRepository.GetSingleByCache(SystemSettingConfigEnum.MaxTimesOfRequestCardInputPerIPAddress, 60);

            int maxTimes = querySetting == null ? 30 : querySetting.SettingValueOfInt;

            return maxTimes;
        }

    }
}
