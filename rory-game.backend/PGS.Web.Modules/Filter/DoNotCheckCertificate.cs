using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using target.Models.ViewModels.Base;
using example.invalid.Http;
using System.Collections.Generic;
using System.Linq;
using System;

namespace target.Web.Modules.Filter
{

    [AttributeUsage(AttributeTargets.Class)]
    public class DoNotCheckCertificate : System.Attribute
    {


    }

}
