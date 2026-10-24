using Microsoft.AspNetCore.Mvc.Filters;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using target.Models.ViewModels.Base;
using example.invalid.Http;
using System.Collections.Generic;
using System.Linq;

namespace target.Web.Modules.Filter
{
    /// <summary>
    /// 將api response格式 封裝成統一規格
    /// </summary>
    public class ApiResponseFormatterFilter : IActionFilter
    {

        private static List<string> skipFormatterApis = new List<string>
        {
            "target/validate-session"
        };


        public void OnActionExecuting(ActionExecutingContext context)
        {
        }


        public void OnActionExecuted(ActionExecutedContext context)
        {
            var url = context.HttpContext.Request.Path.ToString();

            //如果有指定特定的API 不採用包裝的格式
            if (skipFormatterApis.Any(keyword => url.Contains(keyword)))
                return;

            var emptyResult = context.Result as EmptyResult;

            if(emptyResult != null)
            {
                var newInstance = new ApiResponseModel
                {
                    Data = null
                };

                context.Result = new ObjectResult(newInstance)
                {
                    StatusCode = 200,
                    DeclaredType = typeof(ApiResponseModel)
                };

                return;
            }
                       
            //判斷是否為OKObject
            var okObject = context.Result as OkObjectResult;
            if(okObject != null)
            {
                return;
            }

            //多包裝一層data欄位
            var objectResult = context.Result as ObjectResult;
            if(objectResult != null && objectResult.DeclaredType != typeof(HttpResponseMessage))
            {
                objectResult.Value = new ApiResponseModel
                {
                    Data = objectResult.Value,
                };

                objectResult.DeclaredType = typeof(ApiResponseModel);
                return;
            }

            return;
        }

    }

}
