using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Lck.Utility.Extensions.Helper;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models.Models;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using MSC.Web.Modules.Filter;
using Newtonsoft.Json.Linq;

namespace MSC.API.Game.Controllers
{
    [ApiExplorerSettings(IgnoreApi = true)]
    public class DebugController : BaseController
    {       
        public ITokenCheckRefService TokenCheckRefService { get; set; }


        [Route("debug/is-valid-card")]
        [HttpGet , HttpPost]
        public object CheckIsValidToken(string token)
        {
            var result = TokenCheckRefService.CheckIsValidToken(token);
            return result;
        }


        [DisabledWhenProduction]
        [Route("debug/signkey/preview")]
        [HttpPost]
        public SignKeyResultModel SignKeyPreview([FromBody]SignKeyPreviewRequestModel model)
        {
            model.RuleFor(x => x.RequestInfo).Rule_ModelNotNull();

            Dictionary<string, object> convert = ConvertToDic(model.RequestInfo);

            var result = SignKeyRuleHelper.GenerateVendorApiIntegrationSignKey(convert, model.MyMerchantSecret);
            return result;
        }



        private static Dictionary<string, object> ConvertToDic(object requestInfo)
        {
            var convertText = requestInfo.ConvertToJsonString().JsonConvertToModel<Dictionary<string, string>>();
            var convert = requestInfo.ConvertToJsonString().JsonConvertToModel<Dictionary<string, object>>();
            List<string> decimalColumns = new List<string>();
            List<string> dateColumns = new List<string>();

            foreach (var item in convertText)
            {
                var value = item.Value;

                if (value == null)
                    continue;

                if (DateTime.TryParse(value , out DateTime dateValue))
                {
                    dateColumns.Add(item.Key);
                }

                if (value.Contains(".") && decimal.TryParse(value, out decimal decimalValue))
                {
                    decimalColumns.Add(item.Key);
                }
            }
            
            //如果是時間欄位 要轉成iso 8601
            foreach (var item in dateColumns)
            {
                var value = convertText[item];
                convert[item] = DateTimeHelper.ToISO_8601(DateTime.Parse(value).ToUniversalTime());
            }

            //如果是浮點數欄位 要轉成
            foreach (var item in decimalColumns)
            {
                var value = convertText[item];
                convert[item] = decimal.Parse(value).ToString("#0.00");
            }

            return convert;
        }






    }
}
