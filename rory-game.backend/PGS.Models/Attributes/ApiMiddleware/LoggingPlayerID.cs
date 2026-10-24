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
    /// 如果API有標記此，寫API log時，會嘗試解析 requestBody / queryString 找出是否有playerID參數 做為log記錄用
    /// </summary>
    [AttributeUsage(AttributeTargets.Method, Inherited = false, AllowMultiple = true)]
    public class LoggingPlayerID : System.Attribute
    {
        public bool IsEnable { get; set; }

        public LoggingPlayerID()
        {
            IsEnable = true;
        }

        public LoggingPlayerID(bool isEnable)
        {
            IsEnable = isEnable;
        }
    }

}
