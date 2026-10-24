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
    public class ModalValidCustomizedResponse
    {

        public static BadRequestObjectResult ResponseCustomizedErrorResponse(ActionContext context)
        {
            var invalidColumns = context.ModelState.Where(x => x.Value.ValidationState == Microsoft.AspNetCore.Mvc.ModelBinding.ModelValidationState.Invalid);
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
                else
                {
                    msg = msg + error.ErrorMessage;
                }
            }

            var res = new ApiResponseErrorModel
            {
                Code = ApiErrorCode.InvalidParameter_ParsingError.ToInteger(),
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
