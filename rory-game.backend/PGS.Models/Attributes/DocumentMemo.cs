using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using target.Models.ViewModels.Base;
using example.invalid.Http;
using System.Collections.Generic;
using System.Linq;
using System;

namespace target.Models.Attributes
{

    [AttributeUsage(AttributeTargets.Method)]
    public class DocumentMemo : System.Attribute
    {
        public DocumentMemo(string memo)
        {
            this.Memo = memo;
        }

        public string Memo { get; set; }
    }

}
