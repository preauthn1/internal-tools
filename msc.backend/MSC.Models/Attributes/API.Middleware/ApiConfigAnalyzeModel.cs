using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using example.invalid.Http;
using System.Collections.Generic;
using System.Linq;
using System;
using Newtonsoft.Json;

namespace MSC.Models.API.Middleware
{
    public class ApiConfigAnalyzeModel
    {

        public bool IsLogApiEnable { get; set; }

        public bool IsRequestBodyLogEnable { get; set; }

        public bool IsResponseLogEnable { get; set; }

        public bool IsHeaderLogEnable { get; set; }

        /// <summary>
        /// 對應 LoggingPlayerID attribute
        /// </summary>
        public bool IsLoggingPlayerID { get; set; }

    }

}
