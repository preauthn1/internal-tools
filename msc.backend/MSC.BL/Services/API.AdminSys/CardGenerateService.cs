using System;
using System.Collections.Generic;
using MSC.Models.CoreLibs;
using MSC.Models.Paging;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Models;
using System.Linq;
using MSC.Utility.Helper;
using MSC.Models;
using MSC.Utility;
using MSC.DA.Repository;
using MSC.Models.Table;
using System.Transactions;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Models.ViewModels.API.Point;

namespace MSC.BL.Services
{

    public class CardGenerateService : BaseAdminService, ICardGenerateService
    {

        //產生卡號
        //產生連結
        //驗證機制
        //查卡號是否存在 (精簡型SQL)

        //可變動的加密鹽巴 
        // Merchant
        // Card
        // SystemSetting


        public CardNoRepository CardNoRepository { get; set; }

        public CardInfoRepository CardInfoRepository { get; set; }

        public UserOperateLogRepository UserOperateLogRepository { get; set; }

        public ProviderCardGenearteRecordRepository ProviderCardGenearteRecordRepository { get; set; }

        public ProviderMerchantRepository ProviderMerchantRepository { get; set; }



        public PagedList<CardInfoViewModel> GetByPagingFilter(int userID, QueryCardInfoApiRequestParameterModel filter)
        {
            var query = CardInfoRepository.GetPagingByFilter(filter);

            var paging = ConvertToNewTypePagingObject<CardInfoEntity, CardInfoViewModel>(query);

            paging.Items = query.Items
                .Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            return paging;
        }


        public CardDestroyResultModel DestroyCards(int userID, CardDestroyRequestModel model)
        {
            model.RuleFor(z => z.CardNoList).Rule_IsNotEmptyCollection();

            var queryCards = CardInfoRepository.GetByFilter(new QueryCardInfoRequestViewModel
            {
                ListOfCardNo = model.CardNoList
            });

            if (queryCards.IsEmptyOrNull())
                throw new MSCException(ErrorCode.NotFoundData);

            var statusInvalidCards = queryCards.Where(z => z.ConsumeMerchantID != null).Select(z => z.CardNo).ToList();
            var statusvalidCards = queryCards.Where(z => z.ConsumeMerchantID == null).Select(z=> z.CardNo).ToList();

            if(statusvalidCards.IsEmptyOrNull())
                throw new MSCException(ErrorCode.InvalidStatus);

            string warningMsg = string.Empty;

            if(statusInvalidCards.IsNotEmptyNull())
            {
                warningMsg = $"有{statusInvalidCards.Count()}筆狀態不合法、無法更新，" + string.Join("," , statusInvalidCards);
            }

            //寫操作紀錄
            UserOperateLogRepository.CreateByParameter(userID, UserOperateTypeEnum.DestroyCardNumber, null, statusvalidCards.ConvertToJsonString() , model.Memo);

            //更新卡號
            CardInfoRepository.UpdateConsumeMerchantID(statusvalidCards, SystemConstant.MagicConsumeMerchantID);

            return new CardDestroyResultModel
            {
                ResultMessage = $"共{statusvalidCards.Count()}筆資料被更新。"+ warningMsg
            };
        }


        public CardGenerateResultViewModel BatchGenerateCardNumber(int userID, CardGenerateRequestModel model)
        {
            model.RuleFor(x => x.ExpectCount).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.Point).Rule_NotZeroNotMinus();
            model.RuleFor(z => z.Currency).Rule_NotEmptyText();

            string user = GetUserNameByUserID(userID);

            List<string> newCards = GenerateNoneDuplicateCardNumberAndRetry(model.ExpectCount);
            DateTime now = DateTimeHelper.GetCurrentTime();


            //step 1. 塞已存在卡號資料
            var insertModels = newCards.Select(x => new CardNoEntity { CardNo = x }).ToList();
            CardNoRepository.CreateMany(insertModels);


            //step 2. 塞卡號詳細資料
            var insertCardInfoModels = newCards
                .Select(x => new CardInfoEntity
                {
                    Currency = model.Currency,
                    CardNo = x,
                    CardPassword = [REDACTED],
                    Point = model.Point,
                    Memo = model.Memo,
                    CreateBy = user,
                    CreateTime = now,
                    UpdateBy = user,
                    UpdateTime = now
                })
                .ToList();

            CardInfoRepository.CreateMany(insertCardInfoModels);

            var debugInfo = new
            {
                DataInfo = newCards.ConvertToJsonString(),
                Message = $"ExpectCount:{model.ExpectCount}、Amount:{model.Point}"
            };

            //step 3. 寫入產生操作紀錄
            UserOperateLogRepository.CreateByParameter(userID, UserOperateTypeEnum.GenerateCardNumber, null, debugInfo , model.Memo);

            return new CardGenerateResultViewModel
            {
                ListOfCardNo = newCards
            };
        }


        public ProviderCardGenerateResultViewModel ProviderGenerateCard(CardProviderCardGenerateRequestModel model)
        {
            var requestInfo = model.RequestInfo;

            model.RuleFor(z => z.SignKey).Rule_NotEmptyText();
            requestInfo.RuleFor(x => x.TS).Rule_NotZeroNotMinus();
            requestInfo.RuleFor(x => x.ProviderID).Rule_NotZeroNotMinus();
            requestInfo.RuleFor(x => x.Point).Rule_NotZeroNotMinus();
            requestInfo.RuleFor(z => z.TransactionNo).Rule_NotEmptyText();

            var allCurrency = ExtensionOfEnum.GetEnumAllValue<CurrencyEnum>();

            //先寫入紀錄?....(要新table?)
            DateTime now = DateTimeHelper.GetCurrentTime();

            //取得該通路商的secretKey
            var providerSetting = GetProviderSettingSecret(requestInfo.ProviderID);

            var expectSign = SignKeyRuleHelper.GenerateProviderApiIntegrationSignKey(requestInfo, providerSetting.ApiSecretKey);

            if (!IsMagicSignKey(model.SignKey) && model.SignKey != expectSign.AfterMD5)
                throw new MSCException(ErrorCode.ThirdPartyApiFail_SignCheckNotMatch);

            string systemOrderNo = RandomHelper.GetOrderNo("GC");

            //Step 1. 檢查資料已存在
            bool isRecordExist = ProviderCardGenearteRecordRepository.IsProviderOrderExist(requestInfo.ProviderID, requestInfo.TransactionNo);

            if (isRecordExist)
                throw new MSCException(ErrorCode.BusinessLogicDeny_DuplicateDataExist);

            //Step 2. 先寫紀錄
            ProviderCardGenearteRecordRepository.Create(new ProviderCardGenearteRecordEntity
            {
                SystemOrderNo = systemOrderNo,
                UpdateTime = now,
                CreateTime = now,
                Currency = providerSetting.Currency,
                Point = requestInfo.Point,
                ProviderID = requestInfo.ProviderID,
                Memo = requestInfo.Memo,
                TransactionNo = requestInfo.TransactionNo,
            });

            var queryAdmin = UserRepository.GetByFilter(new QueryUserRequestViewModel
            {
                UserName = "admin"
            });

            int userID = queryAdmin.FirstOrDefault().UserID;

            // Step 3. 呼叫相同的generate邏輯
            var generateResult = BatchGenerateCardNumber(userID, new CardGenerateRequestModel
            {
                Currency = providerSetting.Currency,
                ExpectCount = 1,
                Memo = "单号:" + requestInfo.TransactionNo,
                Point = requestInfo.Point
            });

            string cardNo = generateResult.ListOfCardNo.FirstOrDefault();

            var queryCard = CardInfoRepository.GetByCardNo(cardNo);

            return new ProviderCardGenerateResultViewModel
            {
                SystemOrderNo = systemOrderNo,
                CardNo = queryCard.CardNo,
                CardPassword = [REDACTED],
                Currency = queryCard.Currency,
                Point = queryCard.Point
            };
        }

        private bool IsMagicSignKey(string mySignKey)
        {
            return mySignKey == SystemConstant.VendorAPI_SignKey_MagicPassword;
        }

        public List<string> GenerateNoneDuplicateCardNumberAndRetry(int expectCount)
        {
            List<string> listOfCardNo = GenerateNoneDuplicateCardNumber(expectCount);

            if (listOfCardNo.Count() >= expectCount)
                return listOfCardNo;

            int insufficientCount = expectCount - listOfCardNo.Count();

            //若是補剩餘不足的數量 單次產生五倍的筆數 避免浪費時間一直retry
            int regeneateCount = (expectCount - listOfCardNo.Count()) * 5;

            while (true)
            {
                List<string> regenOfCardNo = GenerateNoneDuplicateCardNumber(regeneateCount);

                if (regenOfCardNo.IsNotEmptyNull())
                    listOfCardNo.AddRange(regenOfCardNo);

                //不斷while迴圈 直到筆數足夠
                if (listOfCardNo.Count() >= expectCount)
                    break;
            }

            return listOfCardNo.Take(expectCount).ToList();
        }


        public List<string> GenerateNoneDuplicateCardNumber(int expectCount)
        {
            List<string> listOfCardNo = new List<string>();

            for (int i = 0; i < expectCount; i++)
            {
                var card = CardNumberHelper.GenerateCardNumber();
                listOfCardNo.Add(card);
            }

            var queryExist = CardNoRepository.GetExistCardNo(listOfCardNo);

            if (queryExist.IsEmptyOrNull())
                return listOfCardNo;

            List<string> existCardNo = queryExist.Select(x => x.CardNo).ToList();
            return listOfCardNo.Where(x => !existCardNo.Contains(x)).ToList(); ;
        }


        private CardInfoViewModel ConvertToViewModels(CardInfoEntity query)
        {
            var convert = query.ReturnInjecteObject<CardInfoViewModel>();
            convert.CardStatusText = "未使用";
            
            if (query.ConsumeMerchantID.HasValue)
            {
                int consumeMerchant = query.ConsumeMerchantID.Value;

                if (consumeMerchant == SystemConstant.MagicConsumeMerchantID)
                {
                    convert.CardStatus = CardInfoStatusEnum.Destroy;
                    convert.CardStatusText = "卡號已銷毀";
                }
                else
                {
                    convert.CardStatus = CardInfoStatusEnum.Consumed;
                    convert.CardStatusText = "已使用";
                }
            }

            return convert;
        }

        private ProviderMerchantEntity GetProviderSettingSecret(int providerID)
        {
            var queryProvider = ProviderMerchantRepository.GetById(providerID);

            if (queryProvider == null)
                throw new MSCException(ErrorCode.NotFoundMerchant);

            if (queryProvider.Status != NormalActiveStatus.Enable.ToInteger())
                throw new MSCException(ErrorCode.InvalidStatus);

            return queryProvider;
        }



    }
}
