using Microsoft.AspNetCore.Mvc;
using MSC.Models.ViewModels.Base;
using System;
using System.Collections.Generic;
using System.Text;
using System.Linq;
using MSC.Utility.Extensions;
using MSC.Models.Table;
using MSC.Utility;
using MSC.Utility.IOC;
using MSC.BL.Interface;
using MSC.Models.Models;

namespace MSC.Web.Modules.Filter
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
