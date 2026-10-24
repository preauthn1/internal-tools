using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using target.Models.ViewModels.Base;
using example.invalid.Http;
using System.Collections.Generic;
using System.Linq;
using System;
using target.Models.Enums;
using Newtonsoft.Json;

namespace target.Web.Modules.Filter
{

    /// <summary>
    /// 建議欄位名稱/型別  都定義跟 ApiConfig 一樣
    /// </summary>
    public class ApiConfigAnalyzeModel
    {

        public bool IsLogApiEnable { get; set; }

        public bool IsRequestBodyLogEnable { get; set; }

        public bool IsResponseLogEnable { get; set; }

        public bool IsRedirectToErrorPage { get; set; }

        public bool IsHeaderLogEnable { get; set; }


        public ApiConfigCategoryEnum ApiCategory { get; set; }


        /// <summary>
        /// 對應 LoggingPlayerID attribute
        /// </summary>
        public bool IsLoggingPlayerID { get; set; }

    }

}
