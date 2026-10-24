using Microsoft.AspNetCore.Mvc;
using target.Models.ViewModels.Base;
using System;
using System.Collections.Generic;
using System.Text;
using System.Linq;
using target.Utility.Extensions;
using target.Utility.Helper.Dependency;
using target.Models.Table;
using target.Utility;
using target.Utility.IOC;
using target.BL.Interface;
using target.Models.Models;

namespace target.Web.Modules.Filter
{
    public class SetErrorInfoHelper
    {

        public static void SetDebugInfoBeforeWriteLog(ApiRequestLogEntity logModel , ApiResponseErrorModel error)
        {
            logModel.ResponseStatusCode = 400;
            logModel.DebugInfo = error.Code.ToString();
            logModel.ResponseText = error.ConvertToJsonString();
        }

    }
}
