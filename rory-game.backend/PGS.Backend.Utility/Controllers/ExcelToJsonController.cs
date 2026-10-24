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
    public class ExcelToJsonController : ControllerBase
    {

        private static string BREAK_LINE = Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine
                    + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine
                    + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine;

        /// <summary>
        /// 透過postman貼上 excel格式 可以轉成 json key value檔
        /// </summary>
        [Route("excel_to_json")]
        [HttpPost]
        public IActionResult ConvertExcelToJson([FromQuery]string langs = "en,tw,cn,ms" , bool isReplaceCharacter = true)
        {
            var requestBody = Request.GetRequestBodyString();

            if (string.IsNullOrWhiteSpace(requestBody))
                throw new Exception("requestBody empty");

            List<string> lines = requestBody.Split(new[] { Environment.NewLine , "\t\n" , "\n" }, StringSplitOptions.None)
                .Where(x => !string.IsNullOrWhiteSpace(x))
                .ToList();

            List<string> myLangs = langs.Split(',').ToList();
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
                    continue;

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
                var dic = new Dictionary<string, string>();

                foreach (var allText in allKeyValues.OrderBy(z => z.Key))
                {
                    var key = allText.Key;
                    var value = allText.Value[startIndex];

                    dic.Add(key, value);
                }

                string thisObjectJson = string.Empty;

                string doubleQuotes = @"""";

                foreach (var item in dic)
                {
                    thisObjectJson += $"{doubleQuotes}{item.Key}{doubleQuotes}: {doubleQuotes}{item.Value}{doubleQuotes},"+ Environment.NewLine;
                }

                responseText += thisObjectJson + BREAK_LINE;
                startIndex++;
            }

            if (responseText.Length <= 200)
                return Content(requestBody);

            if(isReplaceCharacter)
            {
                responseText = responseText.Replace("{", "");
                responseText = responseText.Replace("}", "");
            }

            return Content(responseText);
        }

        /// <summary>
        /// 透過postman貼上 excel格式 可以轉成 json key value檔
        /// </summary>
        [Route("excel_to_json_nested")]
        [HttpPost]
        public object ConvertExcelToNestedJson([FromQuery]string langs = "en,tw,cn,ms,it")
        {
            var requestBody = Request.GetRequestBodyString();

            if (string.IsNullOrWhiteSpace(requestBody))
                throw new Exception("requestBody empty");

            List<string> lines = requestBody.Split(new[] { Environment.NewLine , "\n" }, StringSplitOptions.None)
                .Where(x => !string.IsNullOrWhiteSpace(x))
                .ToList();

            List<string> myLangs = langs.Split(',').ToList();



            var allKeyValues = new Dictionary<string, object>();

            List<string> allParentObject = new List<string>();

            //先產生parent物件
            foreach (var eachColumn in lines)
            {
                List<string> split = eachColumn.Split('\t').ToList();

                string keyOfParent = split.FirstOrDefault();
                string keyOfProperty = split[1];

                if (allKeyValues.ContainsKey(keyOfParent))
                    continue;

                if(string.IsNullOrWhiteSpace(keyOfProperty))
                {
                    allKeyValues.Add(keyOfParent, new List<string>());
                }
                else
                {
                    allKeyValues.Add(keyOfParent, new Dictionary<string , List<string>>());
                }
            }

            foreach (var eachDBColumn in lines)
            {
                List<string> split = eachDBColumn.Split('\t').ToList();

                int expectCount = 1 + myLangs.Count();

                string parentName = split[0]; //第一欄是parent物件
                string keyName = split[1]; //第二欄開始才是 多語系鍵值
             
                if (split.Count() < expectCount)
                    continue;

                List<string> allText = new List<string>();

                foreach (var text in split.Skip(2))
                {
                    allText.Add(text);
                }

                if (!allKeyValues.ContainsKey(parentName))
                    continue;

                var parentObject = allKeyValues[parentName];

                //如果keyName空值 代表是第一層的多語系
                if(string.IsNullOrWhiteSpace(keyName))
                {
                    allKeyValues[parentName] = allText;
                    continue;
                }

                Dictionary<string, List<string>> parentObjectDic = parentObject as Dictionary<string, List<string>>;

                if (parentObjectDic.ContainsKey(keyName))
                    continue;

                parentObjectDic.Add(keyName, allText);
            }

            string responseText = string.Empty;


            Dictionary<string, object> listResult = new Dictionary<string, object>();
            int startIndex = 0;
            foreach (var lang in myLangs)
            {
                var jsonResultObject = new Dictionary<string, object>();

                foreach (var eachParentObject in allKeyValues.OrderByDescending(x=> IsListType(x.Value)).ThenBy(z => z.Key))
                {
                    var parentKey = eachParentObject.Key;
                    var parentValues = eachParentObject.Value;
                    bool isList = IsListType(parentValues);

                    if(isList)
                    {
                        List<string> values = (List<string>)parentValues;
                        string thisLangText = values[startIndex];
                        jsonResultObject.Add(parentKey, thisLangText);
                    }
                    else
                    {
                        Dictionary<string, List<string>> forceConvert = parentValues as Dictionary<string, List<string>>;

                        foreach (var eachProp in forceConvert.OrderBy(z => z.Key))
                        {
                            string key = eachProp.Key;
                            List<string> values = eachProp.Value;
                            string thisLangText = values[startIndex];

                            var parentObjectClone = InstanceNewOrGetExist(jsonResultObject, parentKey);
                            parentObjectClone.Add(key, thisLangText);
                        }
                    }
                }

                var jsonText = jsonResultObject.ConvertToJsonStringIndent();

                responseText += jsonText + BREAK_LINE;
                startIndex++;

                listResult.Add(lang, jsonResultObject);
            }

            if (responseText.Length <= 200)
                return Content(requestBody);

            return new
            {
                allLangs = myLangs,
                values = listResult
            };
        }



        private bool IsListType(object myValue)
        {
            if (myValue == null)
                return false;

            return myValue.GetType() == typeof(List<string>);
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
