using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;

namespace MSC.Models.Enums
{
   
    public enum RequestForbiddenTypeEnum
    {

        /// <summary>
        /// 商戶 請求支付
        /// </summary>
        Request_CardInput_WebsiteLink = 3501,

        /// <summary>
        /// 卡號輸入 禁止特定 IPAddress
        /// </summary>
        Request_CardInput_CheckCardNumber_IPAddress = 5001,

        /// <summary>
        /// 卡號輸入 禁止特定 Token
        /// </summary>
        Request_CardInput_CheckCardNumber_TokenHash = 5002,
    }
}
