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

    public class TokenCheckRefService : BaseVendorService, ITokenCheckRefService
    {

        public CardLinkRequestRecordRepository CardLinkRequestRecordRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }

        public RequestForbiddenLogRepository RequestForbiddenLogRepository { get; set; }


        private static string AES_KEY = SystemConstant.API_Vendor_AES_KEY;



        public LinkTokenMemberInfoModel CheckIsValidToken(string token)
        {
            var parseToken = [REDACTED];

            CheckIsDbExistToken(token);

            return parseToken;
        }


        /// <summary>
        /// token格式檢查 (不查DB)
        /// </summary>
        public LinkTokenMemberInfoModel CheckIsTokenFormatCorrect(string token)
        {
            if (string.IsNullOrWhiteSpace(token))
                throw new MSCException(ErrorCode.InvalidParameter);

            string decrypt = string.Empty;

            try
            {
                //解密
                decrypt = EncryptAES.Decrypt(token.Base64Decode() , AES_KEY);
            }
            catch (Exception ex)
            {
                throw new MSCException(ErrorCode.ParsingError_EncryptOrDecryptError);
            }

            LinkTokenMemberInfoModel tokenModel = null;

            try
            {
                //解析json
                tokenModel = decrypt.Base64Decode().JsonConvertToModel<LinkTokenMemberInfoModel>();
            }
            catch (Exception ex)
            {
                throw new MSCException(ErrorCode.ParsingError_ParseTokenObjectFail);
            }
            
            if(tokenModel == null)
            {
                throw new MSCException(ErrorCode.ParsingError_ParseTokenObjectFail);
            }

            string expectSignKey = SignKeyRuleHelper.GeneratePointTokenSignKey(tokenModel.MerchantID, tokenModel.MemberCode);

            //檢查 SignKey
            if (tokenModel.SignKey != expectSignKey)
            {
                throw new MSCException(ErrorCode.ThirdPartyApiFail_SignCheckNotMatch);
            }

            var now = DateTimeHelper.GetCurrentTime();

            //檢查過期
            if(tokenModel.ExpiredTime <= now)
            {
                throw new MSCException(ErrorCode.ThirdPartyApiFail_TokenExpired , new ErrorDetailModel
                {
                    ExtraData = tokenModel.ExpiredTime.ToString(),
                });
            }

            return tokenModel;
        }


        public void CheckIsDbExistToken(string token)
        {
            string tokenHash = token.MD5Hash();

            // cache檢查1. token是否已失效
            CheckTokenDisabledCache(tokenHash);

            bool? cacheOfTokenExist = GetTokenExistCacheAndCheck(tokenHash);

            // cache檢查2. token存在與否
            if (cacheOfTokenExist.HasValue && cacheOfTokenExist.Value == true)
                return;
                       
            //檢查DB是否有資料 (含 IsDisabled 確認 )
            var isExist = CardLinkRequestRecordRepository.IsTokenHashExist(tokenHash);

            string cacheKeyOfTokenExist = CacheKey.PointAPI_TokenExist(tokenHash);

            CacheHelper.SetCacheValue(cacheKeyOfTokenExist, isExist, 5);

            if (!isExist)
                throw new MSCException(ErrorCode.ThirdPartyApiFail_TokenExpired);

        }

        public void CheckTokenDisabledCache(string tokenHash)
        {
            string cacheKey = CacheKey.PointAPI_TokenExpired(tokenHash);

            var queryDisabledCache = CacheHelper.GetCacheValue<string>(cacheKey);

            if (queryDisabledCache != null)
                throw new MSCException(ErrorCode.ThirdPartyApiFail_TokenExpired);

        }

        private bool? GetTokenExistCacheAndCheck(string tokenHash)
        {
            string cacheKey = CacheKey.PointAPI_TokenExist(tokenHash);

            bool? queryTokenExistdCache = CacheHelper.GetCacheValue<bool?>(cacheKey);

            if (queryTokenExistdCache.HasValue && queryTokenExistdCache.Value == false)
                throw new MSCException(ErrorCode.Authorize_AccessToken_NotValid);

            return queryTokenExistdCache;
        }



    }
}
