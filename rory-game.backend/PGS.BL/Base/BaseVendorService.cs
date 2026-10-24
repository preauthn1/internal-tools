using target.Core.Repository;
using target.Models.Enums;
using target.Models.Models;
using target.Models.ViewModels;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using target.Utility.Extensions;
using target.DA.Repository;
using target.Models.Table;
using Microsoft.AspNetCore.Http;
using target.Models;
using target.Utility.Helper.HttpContext;

namespace target.BL.Services
{
    public class BaseVendorService : BaseService
    {

        public MerchantRepository MerchantRepository { get; set; }


        public MerchantEntity GetMerchantAndCheckIsValid(int merchantID, string secret)
        {
            var validModel = new BaseVendorRequestParameter
            {
                MerchantID = merchantID,
                MerchantSecretKey = secret
            };

            validModel.RuleFor(x => x.MerchantID).Rule_NotZero().Rule_NotMinus();
            validModel.RuleFor(x => x.MerchantSecretKey).Rule_NotEmptyText();

            //用id + secretKey 查出供應商資料
            var queryMerchant = MerchantRepository.GetById(merchantID);

            if (queryMerchant == null)
                throw new TerproException(ApiErrorCode.NotFoundMerchant);

            //透過此方式 動態改變廠商端API要顯示的語系
            if (!MyHttpContext.Current.Request.Headers.ContainsKey(SystemConstant.HeaderKey_Lang))
            {
                MyHttpContext.Current.Request.Headers.Add(SystemConstant.HeaderKey_Lang, queryMerchant.DefaultLanguage);
            }

            if (queryMerchant.ApiSecretKey != secret)
                throw new TerproException(ApiErrorCode.NotFoundMerchant);

            if (queryMerchant.Status != (int)NormalActiveStatus.Enable)
                throw new TerproException(ApiErrorCode.InvalidStatus);

            AddMercahntInfoToHeader(queryMerchant.MerchantID);

            return queryMerchant;
        }


        /// <summary>
        /// 透過此方式，ApiLog會獨立某個欄位記錄是哪個merchantID
        /// </summary>
        /// <param name="merchantID"></param>
        private void AddMercahntInfoToHeader(int merchantID)
        {
            if (MyHttpContext.Current.Request.Headers.ContainsKey(SystemConstant.HeaderKey_IdentityMemo))
                return;

            MyHttpContext.Current.Request.Headers.Add(SystemConstant.HeaderKey_IdentityMemo, merchantID.ToString());
        }



      
    }
}
