using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.ViewModels.API.AdminSys
{
    public class DropdownValueViewModel
    {
        public string Text
        {
            get
            {
                return this.Text_CN;
            }
        }

        public string Text_CN { get; set; }


        public string Text_EN { get; set; }


        public string Value { get; set; }

        /// <summary>
        /// 如果某些資料 需要額外欄位的話
        /// </summary>
        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        [SwaggerIgnore]
        public object ExtraProperty { get; set; }

    }
}
