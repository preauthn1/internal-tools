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
using Newtonsoft.Json.Linq;

namespace target.Backend.Utility.Controllers
{

    public class JsonToClassController : Controller
    {


        private static string ChangeLine =Environment.NewLine;

        Dictionary<string, string> eachClassName = new Dictionary<string, string>();

        public string ClassSuffixText = string.Empty;

        public string ClassPrefixText = string.Empty;


        /// <summary>
        /// 透過postman貼上 json 可以轉成C# class程式碼
        /// </summary>
        [Route("Json_to_Class")]
        [HttpPost]
        public IActionResult ConvertDataSchema_To_CSharpClass(string suffixText = "", string prefixText = "")
        {
            ClassSuffixText = suffixText;
            ClassPrefixText = prefixText;

            var requestBody = Request.GetRequestBodyString();

            if (string.IsNullOrWhiteSpace(requestBody))
                throw new Exception("requestBody empty");
            

            var jsonObject = requestBody.JsonConvertToModel<Dictionary<string , dynamic>>();

            GenerateClass("MasterClass" ,  jsonObject);

            string outputText = string.Empty;

            foreach(var item in eachClassName)
            {
                outputText += item.Value + ChangeLine;      
            }

            return Content(outputText);
        }
        
        
        private string GenerateColumnRow(string columnName , string columnType)
        {
            string get_set_text = "{get;set;}";
            return  $" public {columnType} {columnName} {get_set_text} {ChangeLine}";
        }



        private string GenerateClass( string targetClassName , Dictionary<string, dynamic> dic)
        {
            string classFileText = string.Empty;
            classFileText += " public class "+ targetClassName + " { " + ChangeLine;
            
            foreach (var eachProp in dic)
            {
                string columnName = eachProp.Key;
                dynamic propValue = eachProp.Value;
                string columnType = DetermineColumnType(columnName, propValue);
                string generateText = GenerateColumnRow(columnName, columnType);
                classFileText += generateText;
            }

            classFileText += " } ";
            eachClassName.Add(targetClassName, classFileText);
            return string.Empty;
        }


        
        private string DetermineColumnType(string propertyName , object value)
        {
            if (value == null)
                return "object";

            var myType = value.GetType();

            if(myType == typeof(JObject))
            {
                string targetClassName = ClassPrefixText + FirstCharToUpper(propertyName) + ClassSuffixText;

                if(!eachClassName.ContainsKey(targetClassName))
                {
                    var parsingValue = value.ConvertToJsonString().JsonConvertToModel<Dictionary<string, dynamic>>();
                    GenerateClass(targetClassName, parsingValue);
                }

                return targetClassName;
            }

            if (myType == typeof(JArray))
            {
                var firstValue = value.ConvertToJsonString().JsonConvertToModel<List<dynamic>>().FirstOrDefault();
                var arrayDataType = DetermineColumnType(propertyName , firstValue);
                return $" List<{arrayDataType}> ";
            }

            var valueText = value.ToString().ToLower();

            if (myType == typeof(string))
                return "string";

            if (myType == typeof(Boolean))
                return "bool";

            if (myType == typeof(DateTime))
                return "DateTime";

            if (myType == typeof(int))
                return "int";

            if (myType == typeof(long))
            {
                if (valueText.Length <= 6)
                    return "int";

                return "long";
            }

            if (valueText.Contains("."))
            {
                decimal tryValue_Decimal = 0;
                var isDecimal = decimal.TryParse(valueText, out tryValue_Decimal);

                if(isDecimal)
                {
                    return "decimal";
                }
            }

            return "string";

        }

        public static string FirstCharToUpper(string input)
        {
            if (string.IsNullOrWhiteSpace(input))
                return string.Empty;

            return input.First().ToString().ToUpper() + input.Substring(1);
        }

    }
}
