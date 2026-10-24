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
    public class ExcelToPHPArrayController : ControllerBase
    {

        private static string BREAK_LINE = Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine
                    + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine
                    + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine;

        /// <summary>
        /// 透過postman貼上 excel格式 可以轉成 json key value檔
        /// </summary>
        [Route("excel_to_php_array")]
        [HttpPost]
        public IActionResult ConvertExcelToJson([FromQuery]string langs = "en,tw,cn,ms" , bool isReplaceCharacter = true)
        {
            var requestBody = Request.GetRequestBodyString();

            if (string.IsNullOrWhiteSpace(requestBody))
                throw new Exception("requestBody empty");

            List<string> lines = requestBody.Split(new[] { Environment.NewLine , "\t\n" , "\n" }, StringSplitOptions.None)
                .Where(x => !string.IsNullOrWhiteSpace(x))
                .ToList();


            bool isUnexpectCount = false;

            List<string> myLangs = langs.Split(',').ToList();

            //裡面存的是某一組 key 對應的所有語系文字(list<string>)
            Dictionary<string, List<string>> allKeyValues = new Dictionary<string, List<string>>();

            foreach (var eachDBColumn in lines)
            {
                List<string> split = eachDBColumn.Split('\t').ToList();

                int expectCount = 1 + myLangs.Count();

                string keyName = split[0];

                if (string.IsNullOrWhiteSpace(keyName))
                    continue;

                //如果沒有填特定語系 這邊的解析就會跳出
                if (split.Count() < expectCount)
                {
                    isUnexpectCount = true;
                    break;
                }

                List<string> allText = new List<string>();

                foreach (var text in split.Skip(1))
                {
                    allText.Add(text);
                }

                if (allKeyValues.ContainsKey(keyName))
                    continue;

                allKeyValues.Add(keyName, allText);
            }

            string responseText = string.Empty;

            var debug1 = allKeyValues.FirstOrDefault(x => x.Key == "Agent_Table_Column_Level");
            var debug2 = allKeyValues.FirstOrDefault(x => x.Key == "Msg_ApplySuccess");

            int startIndex = 0;
            foreach (var lang in myLangs)
            {
                var allTextOfThisLang = new Dictionary<string, string>();

                foreach (var allText in allKeyValues.OrderBy(z => z.Key))
                {
                    var key = allText.Key;
                    var value = allText.Value[startIndex];

                    allTextOfThisLang.Add(key, value);
                }

                string thisObjectJson = string.Empty;

                string doubleQuotes = @"""";

                foreach (var item in allTextOfThisLang)
                {
                    thisObjectJson += $"{doubleQuotes}{item.Key}{doubleQuotes} => {doubleQuotes}{item.Value}{doubleQuotes},"+ Environment.NewLine;
                }

                responseText += thisObjectJson + BREAK_LINE;
                startIndex++;
            }

            if(isUnexpectCount)
                return Content($"langs語系數量 ({myLangs.Count()}) 與 excel欄位數不一致，無法解析");

            if (responseText.Length <= 200)
                return Content(requestBody);

            if(isReplaceCharacter)
            {
                responseText = responseText.Replace("{", "");
                responseText = responseText.Replace("}", "");
            }

            return Content(responseText);
        }



        private string ReplaceFirstWordIfMatch(string source, char keyword)
        {
            char firstWord = source[0];

            if (firstWord == keyword)
            {
                return source.Remove(0, 1);
            }

            return source;
        }


        private Dictionary<string , string> InstanceNewOrGetExist(Dictionary<string, object> dic, string key)
        {
            if (dic.ContainsKey(key))
                return dic[key] as Dictionary<string, string>;

            var newInstance = new Dictionary<string, string>();

            dic.Add(key, newInstance);
            return newInstance;
        }












    }

}
