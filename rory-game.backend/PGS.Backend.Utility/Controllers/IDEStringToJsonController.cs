using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using target.Backend.Utility.Models;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using example.invalid.Http;
using example.invalid.Http.Headers;
using Microsoft.AspNetCore.Mvc;
using target.Utility.Extensions;
using example.invalid;


namespace target.Backend.Utility.Controllers
{
    public class IDEStringToJsonController : Controller
    {

        /// <summary>
        /// 將這樣的格式轉成json (沒有雙引號包覆、沒有逗號隔開)
        // agent_username:""
        // cash_flow_id:""
        // source:0
        // filter:"1"
        // deposit_to:0
        // start_date:"2021-02-01 00:00:00"
        // end_date:"2021-02-28 23:59:59"
        // page:1
        // per_page:10
        /// </summary>
        [Route("ide_text/to_json")]
        [HttpPost]
        public IActionResult DbQueryToExcel()
        {
            var requestBody = Request.GetRequestBodyString();

            if (string.IsNullOrWhiteSpace(requestBody))
                throw new Exception("requestBody empty");

            List<string> lines = requestBody.Split(new[] { Environment.NewLine, "\t\n", "\n" }, StringSplitOptions.None)
                .Where(x => !string.IsNullOrWhiteSpace(x))
                .ToList();
         
            Dictionary<string, string> mappingColumnName = new Dictionary<string, string>();

            int rowIndex = 1;

            foreach (var eachDBColumn in lines)
            {
                bool isFirstRow = rowIndex == 1;

                List<string> splitEachColumns = eachDBColumn.Split(':' , 2).ToList();

                string columnName = splitEachColumns[0];
                string columnValue = splitEachColumns[1];

                columnValue = columnValue.Replace(@"""", "");

                mappingColumnName.Add(columnName, columnValue);
                rowIndex++;
            }

            return Json(mappingColumnName);
        }

    }
}
