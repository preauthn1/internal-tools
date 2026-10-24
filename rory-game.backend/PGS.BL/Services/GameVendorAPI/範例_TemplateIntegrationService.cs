using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Models.ViewModels.Base;
using target.Utility.Utility;
using target.Utility.Extensions;
using System.Linq;
using target.Models.Models;
using target.Utility;
using target.DA.Repository;
using target.BL.Interface;
using target.Models.ViewModels.GameVendor;
using target.Utility.Helper;
using target.Utility.Helper.HttpContext;
using target.Models.Attributes;

namespace target.BL.Services
{

    /// <summary>
    /// 範例
    /// </summary>
    public class Sample_TemplateIntegrationService : BaseGameVendorIntegrationService, IGameVendorIntegrationService
    {

        /// <summary>
        /// 遊戲商編號
        /// </summary>
        public override GameVendorTypeEnum GameVendor
        {
            get
            {
                return GameVendorTypeEnum.Funky;
            }
        }


        /// <summary>
        /// 遊戲商前綴
        /// </summary>
        public override string GameVendorPrefix
        {
            get
            {
                return "xxxx";
            }
        }

        /// <summary>
        /// 對方API的時區
        /// </summary>
        public override double TimeZone 
        {
            get
            {
                return 0;
            }
        }


        //參考完寫法 麻煩刪掉這個
        private FunkyIntegrationService _this_is_sample = new FunkyIntegrationService();


        //--------------------------------------
        //--------------------------------------
        // 建議按照 標記優先度 1、2、3、4 的順序實作
        // 建議按照 標記優先度 1、2、3、4 的順序實作
        // 建議按照 標記優先度 1、2、3、4 的順序實作
        // 建議按照 標記優先度 1、2、3、4 的順序實作
        //--------------------------------------
        //--------------------------------------

        //優先度：１. 啟動遊戲、進入遊戲 
        //優先度：２. 投注、結算、寫入pgs資料庫
        //優先度：３. 補齊一些完整性
        //優先度：４. 後台報表、對帳數據 



        [DocumentMemo("優先度：３")]
        public SingleWalletExceptionHandleModel FormatParsing_ParseExceptionModelToVendorApiResponse(ApiResponseErrorModel model)
        {
            //將 PGS系統的 apiResponse 轉成對方要求的API格式

            // (看完請刪掉這段**) 可參考這邊的寫法
            _this_is_sample.FormatParsing_ParseExceptionModelToVendorApiResponse(model);

            return new SingleWalletExceptionHandleModel
            {
                StatusCode = 400,
                GameVendorCustomizedResponse = new
                {
                    code = model.Code,
                    message = model.Message
                }                
            };
        }

        [DocumentMemo("優先度：２")]
        public object FormatParsing_ConvertToPlaceBetResponseResult(SingleWalletPlaceBetDebitAmountResponseParameter placeBetInfo)
        {
            // [投注API] 專用: 轉成對方文件要求的API回應格式

            // (看完請刪掉這段**) 可參考這邊的寫法
            _this_is_sample.FormatParsing_ConvertToPlaceBetResponseResult(placeBetInfo);

            return new 
            {
                code = 0,
                message = "OK",
                data = 123
            };
        }


        [DocumentMemo("優先度：２")]
        public object FormatParsing_ConvertToSettleBetResponseResult(SettlementConvertRequiredInfoModel model)
        {
            // [結算API] 專用: 轉成對方文件要求的 API回應格式

            // (看完請刪掉這段**) 可參考這邊的寫法
            _this_is_sample.FormatParsing_ConvertToSettleBetResponseResult(model);

            return new
            {
                code = 0,
                message = "OK",
                data = 123
            };
        }


        [DocumentMemo("優先度：２")]
        public PerGameDetailLogEntity FormatParsing_ConvertBetRecordToPerGameDetailLog(string betRecordNo)
        {
            // [結算API] 專用，將三方的投注資料 轉成 PGS的perGameDetailLog
            // **注意 未結算的注單 請不要寫perGameLog
            throw new TerproException(ApiErrorCode.NotImplement);

            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            _this_is_sample.FormatParsing_ConvertBetRecordToPerGameDetailLog(betRecordNo);
        }



        [DocumentMemo("優先度：２")]
        public void WriteRawBetRecordWhenPlaceBet(InsertRawBetRecordRequiredInfoModel model, EachGameVendorRawRequestInfoModel rawInfo)
        { 
            // [投注API] 專用: 投注後，要寫入一筆 raw bet record資料 (三方的投注資料原始資料)
            rawInfo.RuleFor(x => x.Funky).Rule_ModelNotNull();
            throw new TerproException(ApiErrorCode.NotImplement);


            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            _this_is_sample.WriteRawBetRecordWhenPlaceBet(model, rawInfo);
        }


        [DocumentMemo("優先度：2")]
        public void CheckIsValidRequestInfo(ExternalVendorInputValidateRequestParameter model)
        {
            // [投注API] 檢查是否為合法三方的請求 (包含token、加密...)
            throw new TerproException(ApiErrorCode.NotImplement);
        }


        [DocumentMemo("優先度：4")]
        public BetRecordGeneralColumnModel GetBetRecordByNo(string betRecordNo, string gameCode)
        {
            // [查帳API] 使用，提供單號、讓廠商查帳
            throw new TerproException(ApiErrorCode.NotImplement);
        }

        [DocumentMemo("優先度：4")]
        public BetRecordGeneralColumnModel GetBetRecordByPGSTransactionNo(string betRecordNoPGS)
        {
            // [查帳API] 使用，提供單號、讓廠商查帳
            throw new TerproException(ApiErrorCode.NotImplement);
        }


        [DocumentMemo("優先度：２")]
        public void UpdateRawBetRecordAfterSettlement(SettlementBetRecordRequiredInfoModel model, EachGameVendorRawRequestInfoModel rawInfo)
        {
            // [結算API] 專用: 將指定table的原始注單資料 更新 狀態/輸贏金額/有效流水/最後更新時間
            throw new TerproException(ApiErrorCode.NotImplement);

            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            _this_is_sample.UpdateRawBetRecordAfterSettlement(model, rawInfo);
        }


        [DocumentMemo("優先度：２")]
        public void UpdateRawBetRecordAfterSettlement(SettlementBetRecordRequiredInfoModel model)
        {
            // [後台手動結算／結算API] 專用: 將指定table的原始注單資料 更新 狀態/輸贏金額/有效流水/最後更新時間
            throw new TerproException(ApiErrorCode.NotImplement);

            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            _this_is_sample.UpdateRawBetRecordAfterSettlement(model);
        }


        [DocumentMemo("優先度：４")]
        public void UpdateBetRecordStatus(string transactionNoPGS, GameVendorBetRecordStatusEnum targetStatus, string memo)
        {
            // [後台API] 專用: 提供手動退還金幣
            throw new TerproException(ApiErrorCode.NotImplement);
        }



        [DocumentMemo("優先度：３")]
        public ExternalVendorCancelBetResultModel CancelBet(string betRecordNoPGS, string memo)
        {
            // [取消注單API] 專用，將三方的投注資料 狀態改成已取消

            //最後再轉成 對方文件 要求的格式
            return new ExternalVendorCancelBetResultModel
            {
                GameVendorCustomizedResponse = "xxxxxx"
            };

            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            _this_is_sample.CancelBet(betRecordNoPGS, memo);
        }



        [DocumentMemo("優先度：１")]
        public void CreateVendorPlayerAccountOrNot(ExternalVendorCreateAccountRequestParameter model)
        {
            // [啟動遊戲API] 專用，如果對方有 create account 的API，實作該行為(如果不需要 留空白即可)
            throw new TerproException(ApiErrorCode.NotImplement);
        }



        [DocumentMemo("優先度：１")]
        public ExternalVendorLaunchUrlViewModel GetGameLaunchUrl(ExternalVendorLaunchUrlRequestModel model)
        {
            // [啟動遊戲API] 專用，取得遊戲連結   
            throw new TerproException(ApiErrorCode.NotImplement);

            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            _this_is_sample.GetGameLaunchUrl(model);
        }



        [DocumentMemo("優先度：１")]
        public List<ExternalVendorGameListConvertModel> GetGames(int merchantID)
        {
            // [啟動遊戲API] 專用 取得三方遊戲的所有遊戲，並轉成PGS的table格式 (含多語系)
            throw new TerproException(ApiErrorCode.NotImplement);

            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            _this_is_sample.GetGames(merchantID);
        }


        [DocumentMemo("優先度：４")]
        public ThirdPartyGameVendorBetRecordDetailUrlViewModel GetExternalBetDetailInfo(ExternalVendorBetRecordDetailInfoRequestModel model)
        {
            // [後台報表API] 專用，取得注單詳細資訊(外部連結/json資訊)
            throw new TerproException(ApiErrorCode.NotImplement);

            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            _this_is_sample.GetExternalBetDetailInfo(model);
        }



        [DocumentMemo("優先度：４")]
        public ExternalVendorGetReportResultModel GetDbAndApiSummaryReport(ExternalVendorGetReportRequestParameter model)
        {
            // [對帳API]專用 取得統計數據、以及統計自己db原始投注資料的數據 並返回回去給service 比對判斷是否數據相等

            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            // (看完請刪掉**) 請自行參考這邊的寫法
            _this_is_sample.GetDbAndApiSummaryReport(model);


            throw new TerproException(ApiErrorCode.NotImplement);
        }



        private void SetRequestInfoByDBSetting(RequestParameters req, int merchantID, string apiCategory, string requestUUID = "")
        {
            req.LogAction = WriteRequestLog;
            req.AdvancedSetting = new AdvancedSetting
            {
                Log_ApiCategory = apiCategory,
                Log_SubCategory = GameVendor.ToString(),
                Log_RequestUUID = requestUUID
            };

            string apiDomain = GetVendorSettingByKey(merchantID, VendorConfigSettingEnum.ApiUrl)?.SettingValue;
            string opeatorID = GetVendorSettingByKey(merchantID, VendorConfigSettingEnum.MainAssignVendorID)?.SettingValue;
            string apiToken = [REDACTED], VendorConfigSettingEnum.MainSecretKeyOrApiToken)?.SettingValue;

            if (apiDomain.IsNullOrWhiteSpace() || opeatorID.IsNullOrWhiteSpace() || apiToken.IsNullOrWhiteSpace())
            {
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_SettingDataNotCompleted);
            }

            req.BaseUrl = apiDomain;
            req.UserAgent = opeatorID;
            req.AddHeaderData("Authentication", apiToken);
            req.AddHeaderData("X-Request-ID", Guid.NewGuid().ToString());
        }


    }
}
