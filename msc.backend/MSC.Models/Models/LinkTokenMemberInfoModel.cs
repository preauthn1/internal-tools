using MSC.Models.Enums;
using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Text;

namespace MSC.Models.Models
{
    public class LinkTokenMemberInfoModel
    {
        [JsonProperty("id")]
        public string MSCRequestID { get; set; }

        [JsonProperty("mid")]
        public string MerchantRequestID { get; set; }

        [JsonProperty("mc")]
        public string MemberCode { get; set; }

        [JsonProperty("mi")]
        public int MerchantID { get; set; }

        [JsonProperty("t")]
        public DateTime ExpiredTime { get; set; }

        [JsonProperty("s")]
        public string SignKey { get; set; }

        [JsonProperty("l")]
        public LanguageEnum Lang { get; set; }

        [JsonProperty("aa")]
        public decimal AllowAmount { get; set; }

        [JsonProperty("rUrl")]
        public string ReturnUrl { get; set; }

    }
}
