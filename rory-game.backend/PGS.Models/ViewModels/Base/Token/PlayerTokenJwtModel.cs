using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.ViewModels.Base
{
    public class PlayerTokenJwtModel : ITokenJwtModel
    {
        /// <summary>
        /// playerID，但是欄位命名請勿修改
        /// </summary>
        public int LoginUserID { get; set; }


        public string LoginName { get; set; }


        public int MerchantID { get; set; }



        /// <summary>
        /// 語系
        /// </summary>
        public string Lang { get; set; }

        public long TokenExpiredTime { get; set; }

        public string RndID { get; set; }

        /// <summary>
        /// 商戶端自訂的sessionToken
        /// </summary>
        public string SessionTokenRaw { get; set; }

        /// <summary>
        /// PGS格式的token
        /// </summary>
        public string SessionTokenPGS { get; set; }


        /// <summary>
        /// 加密驗證是否被偽造
        /// </summary>
        public string SignKey { get; set; }


        public int TokenVersion { get; set; }


        public int JwtTokenType { get; set; } 

    }

    public class PlayerJwtGenerateResultModel 
    {

        public string JwtToken { get; set; }


        public PlayerTokenJwtModel JwtModel { get; set; }

    }



}
