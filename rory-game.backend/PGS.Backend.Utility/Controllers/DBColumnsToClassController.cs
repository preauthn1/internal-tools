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
    [ApiExplorerSettings(IgnoreApi = true)]
    public class DBColumnsToClassController : ControllerBase
    {
        /// <summary>
        /// 透過postman貼上 SSMS複製的schema格式 可以轉成C# class程式碼
        /// </summary>
        [Route("DB_to_Class/{tableName}")]
        [HttpPost]
        public IActionResult ConvertDataSchema_To_CSharpClass([FromQuery]string tableName)
        {
            var requestBody = Request.GetRequestBodyString();

            if (string.IsNullOrWhiteSpace(requestBody))
                throw new Exception("requestBody empty");

            List<string> lines = requestBody.Split(new[] { Environment.NewLine }, StringSplitOptions.None).Where(x=> !string.IsNullOrWhiteSpace(x)).ToList();

            int totalColumns = 0;
            string yourClassText = "";

            foreach (var eachDBColumn in lines)
            {
                List<string> split = eachDBColumn.Split('\t').Where(x=> !string.IsNullOrWhiteSpace(x)).ToList();

                if (split.Count() != 3)
                    continue;

                var columName = split[0];
                var columType = split[1].Trim();
                var isNullChecked = split[2].Trim();

                var csharpType = "string";

                if (columType == "int")
                    csharpType = "int";
                if (columType == "bit")
                    csharpType = "bool";
                else if (columType == "bigint")
                    csharpType = "long";
                else if (columType == "tinyint")
                    csharpType = "byte";
                else if (columType == "smallint")
                    csharpType = "short";
                else if (columType == "float")
                    csharpType = "float";        
                else if (columType == "datetime")
                    csharpType = "DateTime";
                else if (columType == "date")
                    csharpType = "DateTime";
                else if (columType == "uniqueidentifier")
                    csharpType = "Guid";                
                else if (columType.Contains("decimal"))
                    csharpType = "decimal";
                else if (columType.Contains("char"))
                    csharpType = "string";

                string get_set_text = "{get;set;}";
                string isNullableText = "";

                if(csharpType != "string")
                {
                    if(isNullChecked == "Checked")
                    {
                        isNullableText = "?";
                    }
                }
                

                yourClassText += $"　public {csharpType}{isNullableText} {columName} {get_set_text}" + Environment.NewLine;
                totalColumns++;
            }

            return Content(yourClassText);
        }


    }

}
