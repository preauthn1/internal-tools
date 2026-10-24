using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;
using example.invalid.Http;
using example.invalid.Http.Headers;
using Microsoft.AspNetCore.Mvc;
using target.Utility.Extensions;
using example.invalid;
using Lck.GoogleApi;
using System.Text.RegularExpressions;

namespace target.Backend.Utility.Controllers
{
    [ApiExplorerSettings(IgnoreApi = true)]
    public class DBQueryToJsonController : ControllerBase
    {

        private static string BREAK_LINE = Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine
                    + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine
                    + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine + Environment.NewLine;



        /// <summary>
        /// 透過postman貼上 excel格式 可以再解析裡面json物件
        /// </summary>
        [Route("db_query_to_excel")]
        [HttpPost]
        public IActionResult DbQueryToExcel([FromQuery]string jsonColumns = "")
        {
            var requestBody = Request.GetRequestBodyString();

            if (string.IsNullOrWhiteSpace(requestBody))
                throw new Exception("requestBody empty");

            List<string> lines = requestBody.Split(new[] { Environment.NewLine, "\t\n", "\n" }, StringSplitOptions.None)
                .Where(x => !string.IsNullOrWhiteSpace(x))
                .ToList();

            List<SheetReadModel> eachExcelRows = new List<SheetReadModel>();

            List<string> targetJsonColumns = jsonColumns.Split(',').ToList();

            Dictionary<int, string> positionMappingColumnName = new Dictionary<int, string>();

            int rowIndex = 1;
            int columnIndex = 1;

            foreach (var eachDBColumn in lines)
            {
                bool isFirstRow = rowIndex == 1;

                List<string> splitEachColumns = eachDBColumn.Split('\t').ToList();

                var newInstance = new SheetReadModel();

                foreach (string text in splitEachColumns)
                {
                    if (IsJson(text))
                    {
                        var parseToJson = text.TryJsonConvertToModel<Dictionary<string, string>>();

                        if (parseToJson != null)
                        {
                            //如果包含 * 字串 全部印出來
                            if (jsonColumns == "*")
                            {
                                foreach (var item in parseToJson)
                                {
                                    string columnName = item.Key;
                                    newInstance.ColumnValues.Add($"{item.Value}");

                                    if(isFirstRow)
                                    {
                                        AssignColumnNameOnPoistion(positionMappingColumnName, ref columnIndex, columnName);
                                    }
                                }
                            }
                            else
                            {
                                foreach (var column in targetJsonColumns)
                                {
                                    if (parseToJson.ContainsKey(column))
                                    {
                                        string columnName = column;
                                        string myText = parseToJson[column];
                                        newInstance.ColumnValues.Add($"{myText}");

                                        if (isFirstRow)
                                        {
                                            AssignColumnNameOnPoistion(positionMappingColumnName, ref columnIndex, columnName);
                                        }
                                    }
                                    else
                                    {
                                        string columnName = $"column{columnIndex}";
                                        newInstance.ColumnValues.Add("");

                                        if (isFirstRow)
                                        {
                                            AssignColumnNameOnPoistion(positionMappingColumnName, ref columnIndex, columnName);
                                        }
                                    }
                                }
                            }
                        }
                    }
                    else
                    {
                        string columnName = $"column{columnIndex}";
                        newInstance.ColumnValues.Add(text);

                        if (isFirstRow)
                        {
                            AssignColumnNameOnPoistion(positionMappingColumnName, ref columnIndex, columnName);
                        }
                    }
                }

                eachExcelRows.Add(newInstance);
                rowIndex++;

            }



            string allText = "";

            //表頭列
            foreach (var column in positionMappingColumnName.OrderBy(x=> x.Key))
            {
                allText += column.Value + "\t";
            }

            allText += Environment.NewLine;
            

            //實際的數值
            foreach (var row in eachExcelRows)
            {
                allText += "";

                foreach (var text in row.ColumnValues)
                {
                    allText += text + "\t";
                }

                allText += Environment.NewLine;
            }

            return Content(allText);

        }

        private void AssignColumnNameOnPoistion(Dictionary<int , string> dic , ref int columnIndex , string columnName)
        {
            dic.Add(columnIndex, columnName);
            columnIndex++;
        }


        private bool IsJson(string text)
        {
            if (string.IsNullOrWhiteSpace(text))
                return false;

            if (text.Length <= 5)
                return false;

            bool isKeywordExsit = text.Contains("{") && text.Contains(":") && text.Contains("}");
            return isKeywordExsit;
        }


    }

}
