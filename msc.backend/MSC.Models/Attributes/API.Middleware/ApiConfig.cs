using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using MSC.Models.ViewModels.Base;
using example.invalid.Http;
using System.Collections.Generic;
using System.Linq;
using System;
using MSC.Models.Enums;
using Newtonsoft.Json;

namespace MSC.Models.API.Middleware
{
    [AttributeUsage(AttributeTargets.Method, Inherited = false, AllowMultiple = true)]
    public class ApiConfig : System.Attribute
    {
        public bool IsLogApiEnable { get; set; }

        public bool IsRequestBodyLogEnable { get; set; }

        public bool IsResponseLogEnable { get; set; }

        public bool IsHeaderLogEnable { get; set; }


        public ApiConfig()
        {
            IsLogApiEnable = true;
            IsResponseLogEnable = true;
            IsRequestBodyLogEnable = true;
            IsHeaderLogEnable = true;
        }     
    }


}
