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
using MSC.Utility.Helper.Dependency;

namespace MSC.Web.Modules.Filter
{
    public class ModalValidCustomizedResponse
    {

        public static BadRequestObjectResult ResponseCustomizedErrorResponse(ActionContext context)
        {
            var invalidColumns = context.ModelState.Where(x => x.Value.ValidationState == Microsoft.AspNetCore.Mvc.ModelBinding.ModelValidationState.Invalid);
            var columns = context.ModelState.Keys.Where(z=> !string.IsNullOrWhiteSpace(z)).Select(x => x).ToList();
            string msg = "您所带入的表单请求参数格式不正确";

            if(invalidColumns.IsNotEmptyNull())
            {
                var allColumns = string.Join("、", invalidColumns.Select(z=> z.Key));
                msg = $"{allColumns} 非合法格式，请检查您填写的参数";
            }
            else
            {
                var error = context.ModelState.Root.Errors.FirstOrDefault();
                if(error.ErrorMessage.Contains("non-empty"))
                {
                    string extraMsg = string.Empty;
                    if(context.HttpContext.Request.Method == "GET")
                    {
                        extraMsg = "（检查您代码中是否正确标记 [FromQuery]？）";
                    }

                    msg = $"不允許请求参数為空。{extraMsg}";
                }
            }

            var res = new ApiResponseErrorModel
            {
                Code = ErrorCode.InvalidParameter_ParsingError.ToInteger(),
                Message = msg,
                Data = null
            };
            
            var logModel = new ApiRequestLogEntity();
            HttpRequestLogHelper.SetRequestInfo(context.HttpContext.Request, logModel);
            SetErrorInfoHelper.SetDebugInfoBeforeWriteLog(logModel, res);
            logModel.ResponseStatusCode = 400;

            //errorLog紀錄
            string tableName = SettingHelper.ApiRequestErrorLogTableName;
            var service = IOCHelper.Resolve<IApiRequestLogService>();
            service.CreateTryCatch(tableName, logModel);
            
            return new BadRequestObjectResult(res)
            {
                ContentTypes = { "application/json" }
            };
        }

    }
}
