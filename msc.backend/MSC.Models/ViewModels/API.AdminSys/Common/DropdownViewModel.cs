using MSC.Models.Table;
using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class DropdownSelectValueViewModel
    {

        public string Text { get; set; }

        public string Value { get; set; }


        /// <summary>
        /// 如果某些資料 需要額外欄位的話
        /// </summary>
        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        [SwaggerIgnore]
        public object ExtraProperty { get; set; }


    }



}
