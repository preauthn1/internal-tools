using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using Lck.GoogleApi;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using target.Backend.Utility.Models;
using target.Models.CultureText;
using System.Reflection;
using target.Utility.Extensions;
using target.Models.Models;

namespace target.Backend.Utility.Controllers
{
    public class ExcelToSqlInsertController : Controller
    {        

        /// <summary>
        /// 透過postman送出excel格式 轉成insert into 語法
        /// </summary>
        [Route("post/excel_to_sql_insert")]
        [HttpPost]
        public IActionResult Post_ExcelFormat_ConvertTo_InsertSQL()
        {           
            var requestBody = Request.GetRequestBodyString();

            List<string> lines = requestBody.Split(new[] { Environment.NewLine, "\t\n", "\n" }, StringSplitOptions.None)
             .Where(x => !string.IsNullOrWhiteSpace(x))
             .ToList();                       

            if(lines.IsEmptyOrNull() || lines.Count() <= 1)
                return Content("解析失敗");

            var headerLine = lines.FirstOrDefault();

            List<string> columns = headerLine.Split('\t').ToList();


            string sql_columns = string.Join(",", columns);
            string sql_cmd = $"INSERT INTO XXXXXXXXXXX ({sql_columns}) VALUES " ;
            

            foreach (var eachLine in lines.Skip(1))
            {
                List<string> value_eachColumns = eachLine.Split('\t').ToList();
                
                string sql_value_columns = string.Join(",", value_eachColumns.Select(x=> ParseToSpecialValue(x)));
                sql_cmd += $" {Environment.NewLine} ( {sql_value_columns} ) ,";
            }

            sql_cmd = sql_cmd.TruncateLastWordIfMatch(',');

            return Content(sql_cmd);
        }


        private string ParseToSpecialValue(string value)
        {
            if (value.Contains("getDate"))
                return "GETDATE()";

            return $"'{value.Replace("'" , " ")}'";
        }



    }
}
