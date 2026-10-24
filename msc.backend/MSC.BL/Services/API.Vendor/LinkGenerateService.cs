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

    public class LinkGenerateService : BaseVendorService, ILinkGenerateService
    {

        public CardLinkRequestRecordRepository CardLinkRequestRecordRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }

        public RequestForbiddenLogRepository RequestForbiddenLogRepository { get; set; }

        public SystemSettingConfigRepository SystemSettingConfigRepository { get; set; }


        private static string AES_KEY = SystemConstant.API_Vendor_AES_KEY;

        private static int DEFAULT_MAX_TIMES = 20;


        /// <summary>
        /// (測試用) 產生連結
        /// </summary>
        public LinkGenerateViewModel GenerateTestOrderLinkByMerchantID(int merchantID)
        {
            var queryMerchant = MerchantRepository.GetById(merchantID);

            var model = new LinkGenerateRequestModel
            {
                SignKey = SystemConstant.API_Vendor_MagicCode,
                RequestInfo = new LinkGenerateRequestInfoModel
                {
                    TS = DateTimeHelper.GetCurrentTimeTimestamp(),
                    Amount = 150,
                    Lang = LanguageEnum.zh_cn.ToString(),
                    MemberCode = "admin_test_01",
                    MerchantID = queryMerchant.MerchantID,
                    NotifyUrl = "http://0.0.0.0",
                    ReturnUrl = "http://0.0.0.0:11449/home/Success",
                    RequestID = Guid.NewGuid().ToString(),
                    RandomID = Guid.NewGuid().ToString()
                }
            };

            var result = GenerateCardLink(model);
            return result;
        }




        public LinkGenerateViewModel GenerateCardLink(LinkGenerateRequestModel model)
        {
            model.RuleFor(x => x.SignKey).Rule_NotEmptyText();
            model.RuleFor(x => x.RequestInfo).Rule_ModelNotNull();

            var requestInfo = model.RequestInfo;

            requestInfo.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            requestInfo.RuleFor(x => x.TS).Rule_NotZeroNotMinus().Rule_IsValidTimestamp();
            requestInfo.RuleFor(x => x.RandomID).Rule_NotEmptyText();
            requestInfo.RuleFor(x => x.RequestID).Rule_NotEmptyText();
            requestInfo.RuleFor(x => x.MemberCode).Rule_NotEmptyText();
            requestInfo.RuleFor(x => x.NotifyUrl).Rule_NotEmptyText();
            requestInfo.RuleFor(x => x.Amount).Rule_NotMinus();


            var queryMerchant = GetAndCheckIsValidMerchant(model.RequestInfo);

            CheckSignKeyIsValid(model.SignKey, model.RequestInfo, queryMerchant.ApiSecretKey);

            //檢查請求次數超過上限
            CheckRequestTooManyTimes(requestInfo);

            //檢查系統維護中
            CheckSystemMaintain();

            if (requestInfo.Amount > 0)
            {
                List<decimal> amounts = GetAvailableAmounts(queryMerchant.MerchantID);

                if(!amounts.Contains(requestInfo.Amount))
                {
                    throw new MSCException(ErrorCode.Forbidden_OrderAmount_NotAllow);
                }
            }

            //檢查商戶請求號重複
            CheckMerchantRequestDuplicate(requestInfo);

            //如果有配置要強制語系
            var langSetting = MerchantCustomizedSettingRepository.GetByMerchantSettingKey(requestInfo.MerchantID, MerchantCustomizedSettingKeyEnum.ForceChangeLang);

            if(langSetting != null && !string.IsNullOrWhiteSpace(langSetting.SettingValue))
            {
                requestInfo.Lang = langSetting.SettingValue;
            }

            LanguageEnum lang = HeaderValueHelper.ConvertLangCodeToEnum(requestInfo.Lang);

            string systemRequestID = RandomHelper.GetOrderNo("LK");
            var tokenInfo = new LinkTokenMemberInfoModel
            {
                MSCRequestID = systemRequestID,
                MerchantRequestID = requestInfo.RequestID,
                MerchantID = requestInfo.MerchantID,
                SignKey = SignKeyRuleHelper.GeneratePointTokenSignKey(requestInfo.MerchantID, requestInfo.MemberCode),
                MemberCode = requestInfo.MemberCode,
                ExpiredTime = DateTimeHelper.GetCurrentTime().AddHours(1),
                Lang = lang,
                AllowAmount = model.RequestInfo.Amount,
                ReturnUrl = requestInfo.ReturnUrl
            };

            string originalJson = tokenInfo.ConvertToJsonString();
            var token = [REDACTED], AES_KEY).Base64Encode();

            //寫入一筆請求紀錄
            CardLinkRequestRecordRepository.Create(new CardLinkRequestRecordEntity
            {
                CreateTime = DateTimeHelper.GetCurrentTime(),
                IPAddress = IPHelper.GetIP(),
                MemberCode = requestInfo.MemberCode,
                MerchantRequestID = requestInfo.RequestID,
                SystemRequestID = systemRequestID,
                MerchantID = requestInfo.MerchantID,
                TokenHashValue = token.MD5Hash(),
                ReturnUrl = requestInfo.ReturnUrl,
                NotifyUrl = requestInfo.NotifyUrl,
                Amount = model.RequestInfo.Amount
            });

            string tokenUrlEncode = token.UrlEncode();

            var domainSetting = MerchantCustomizedSettingRepository.GetByMerchantSettingKey(requestInfo.MerchantID, MerchantCustomizedSettingKeyEnum.Basic_CardInputPage_DomainUrl);

            if (domainSetting == null)
                throw new MSCException(ErrorCode.BusinessLogicDeny_SettingDataNotCompleted);

            string domain = domainSetting.SettingValue.ReplaceLastWordIfMatch('/');

            return new LinkGenerateViewModel
            {
                Link = string.Format("{0}/card/input?token=[REDACTED]", domain, tokenUrlEncode)
            };
        }

        private void CheckSystemMaintain()
        {
            var querySystemEnableSetting = SystemSettingConfigRepository.GetSingleByCache(SystemSettingConfigEnum.System_SystemEnableFlag, 1);

            if (querySystemEnableSetting == null)
                return;

            bool isEnable = querySystemEnableSetting.SettingValueOfBoolean;

            if (isEnable)
                return;

            throw new MSCException(ErrorCode.System_SystemMaintain);
        }


        private void CheckMerchantRequestDuplicate(LinkGenerateRequestInfoModel requestInfo)
        {
            var queryExist = CardLinkRequestRecordRepository.GetByFilter(new QueryCardLinkRequestRecordRequestViewModel
            {
                MerchantRequestID = requestInfo.RequestID
            });

            //若請求號已存在 拋錯誤
            if (queryExist.IsNotEmptyNull())
                throw new MSCException(ErrorCode.BusinessLogicDeny_DuplicateDataExist, new ErrorDetailModel { ColumnName = "RequestID" });
        }

        private void CheckRequestTooManyTimes(LinkGenerateRequestInfoModel requestInfo)
        {
            string cacheKey = CacheKey.VendorAPI_RequestForbidden(requestInfo.MerchantID, requestInfo.MemberCode);

            var query = CacheHelper.GetCacheValue<string>(cacheKey);

            //若cache存活 拋錯誤
            if (query != null)
                throw new MSCException(ErrorCode.BusinessLogicDeny_RequestTooManyTimes);

            DateTime previous_hour = DateTimeHelper.GetCurrentTime().AddHours(-1);

            // 分類依據: memberCode + merchantID
            int totalCount = CardLinkRequestRecordRepository.GetTotalCountForTheDuration(requestInfo.MerchantID, requestInfo.MemberCode, previous_hour, DateTimeHelper.GetCurrentTime());

            //超過五次才會去做更深入的檢查
            if (totalCount >= 5)
            {
                var setting = MerchantCustomizedSettingRepository.GetByMerchantSettingKey(requestInfo.MerchantID, MerchantCustomizedSettingKeyEnum.MaxTimesOfRequestLinkPerMember);

                int maxTimes = DEFAULT_MAX_TIMES;

                if (setting != null)
                {
                    maxTimes = setting.SettingValueOfInt;
                }

                if (totalCount >= maxTimes)
                {
                    //加入cache
                    CacheHelper.SetCacheValue(cacheKey, "true", new TimeSpan(0, 30, 0));

                    //寫入紀錄
                    RequestForbiddenLogRepository.Create(new RequestForbiddenLogEntity
                    {
                        BlockKey = $"{requestInfo.MerchantID}_{requestInfo.MemberCode}",
                        CreateTime = DateTimeHelper.GetCurrentTime(),
                        MemberCode = requestInfo.MemberCode,
                        MerchantID = requestInfo.MerchantID,
                        Type = RequestForbiddenTypeEnum.Request_CardInput_WebsiteLink.ToInteger()
                    });

                    throw new MSCException(ErrorCode.BusinessLogicDeny_RequestTooManyTimes);
                }
            }
        }


        private List<decimal> GetAvailableAmounts(int merchantID)
        {
            var querySetting = MerchantCustomizedSettingRepository.GetByMerchantSettingKey(merchantID, MerchantCustomizedSettingKeyEnum.Basic_AvailableAmounts);

            if (querySetting == null)
                throw new MSCException(ErrorCode.BusinessLogicDeny_SettingDataNotCompleted);

            return querySetting.SettingValueListDecimal;
        }


        private string GenerateMerchantRequestSignKey(int merchantID, string requestID, long timestamp)
        {
            return (merchantID.ToString() + requestID + timestamp.ToString()).MD5Hash();
        }


    }
}
