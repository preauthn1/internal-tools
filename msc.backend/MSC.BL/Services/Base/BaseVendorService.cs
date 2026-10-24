using MSC.Models.Paging;
using MSC.Utility.IOC;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using MSC.Core.Repository;
using MSC.Models;
using MSC.Models.ViewModels;
using MSC.Models.Models;
using System.Text;
using MSC.Models.Table;
using MSC.Models.Enums;
using MSC.DA.Repository;

namespace MSC.BL.Services
{
    public class BaseVendorService : BaseService
    {

        public MerchantRepository MerchantRepository { get; set; }


        protected MerchantEntity GetAndCheckIsValidMerchant<T>(T model) where T : BaseVendorAPIRequestModel
        {
            var queryMerchant = MerchantRepository.GetById(model.MerchantID);

            if (queryMerchant == null)
                throw new MSCException(ErrorCode.NotFoundMerchant);

            if (queryMerchant.Status != NormalActiveStatus.Enable.ToInteger())
                throw new MSCException(ErrorCode.InvalidStatus);

            return queryMerchant;
        }


        protected void CheckSignKeyIsValid<T>(string signKey, T model , string merchantSecret) where T: BaseVendorAPIRequestModel
        {
            //若使用神奇密碼 不檢查
            if (signKey == SystemConstant.API_Vendor_MagicCode)
            {
                return;
            }

            var now_UTCTime = DateTimeHelper.GetCurrentTime();
            var var_UTCtime = model.TS.UnixTimeStampToDateTime();          

            var signKeyResult = SignKeyRuleHelper.GenerateVendorApiIntegrationSignKey(model, merchantSecret);
            string expectSignKey = signKeyResult.AfterMD5;

            if (expectSignKey != signKey)
                throw new MSCException(ErrorCode.ThirdPartyApiFail_SignCheckNotMatch);
        }


    }
}
