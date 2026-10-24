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
using target.Models.Enums;
using target.Utility.Helper;

namespace target.Backend.Utility.Controllers
{
    public class CultureTextConvertController : Controller
    {
        private readonly ILogger<HomeController> _logger;

        private readonly IWebHostEnvironment _env;

        private readonly static string _sheetUrl = "https://example.invalid/spreadsheets/d/12NfEknrEBXXxA3x8eoaXdc111GGlVCN3r96_oyousbw/edit#gid=0";
        
        private static readonly int _startIndex = 2;


        public List<string> availableLanguage = new List<string>
        {
            "TW",
            "CN",
            "EN"
        };


        public CultureTextConvertController(IWebHostEnvironment env)
        {
            _env = env;
            GoogleSheetUtility.WebRootDirection = _env.ContentRootPath;
        }



        [HttpGet]
        public ActionResult Culture_To_Class(string filterClassName = "")
        {
            var query = GetOnlineExcelValueAndMergeCurrentClass(isMergeCurrentClass : false);

            if (query == null || query.Count == 0)
            {
                return Content("無解析任何內容或解析失敗");
            }

            List<ExcelConvertToClass> convertResults = new List<ExcelConvertToClass>();

            string resultContent = string.Empty;

            var allClass = new List<string>
            {
                "i18n"
            };


            //接著依序建立每個語系
            foreach (var eachClass in allClass)
            {
                int startIndex = _startIndex; //中文是在第幾個                

                foreach (var eachLang in availableLanguage)
                {
                    string className = eachClass + "_" + eachLang.ToUpper();

                    var myItem = new ExcelConvertToClass
                    {
                        ClassName = className,
                        Languages = eachLang,
                        PropertyValues = new List<ConvertClassPropertyValues>()
                    };

                    foreach (var eachText in query)
                    {
                        string propertyName = eachText.ColumnValues[1];

                        myItem.PropertyValues.Add(new ConvertClassPropertyValues
                        {
                            PropertyText = propertyName,
                            PropertyValue = eachText.ColumnValues[startIndex]
                        });
                    }

                    convertResults.Add(myItem);

                    //切換到下一個語言後 從第2欄 改成第3欄 持續往後+1
                    startIndex++;
                }
            }


            //還要額外先建立一個 i18nKeys物件
            var myInstance = new ExcelConvertToClass
            {
                ClassName = "i18nKeys",
                IsKeyObject = true,
                Languages = "",

                PropertyValues = new List<ConvertClassPropertyValues>()
            };

            int myStartIndex = _startIndex; //中文是在第幾個          

            foreach (var eachText in query)
            {
                string propertyName = eachText.ColumnValues[1];

                var refText = convertResults.SelectMany(x => x.PropertyValues)
                        .Where(x => x.PropertyText == propertyName)
                        .Select(x => x.PropertyValue)
                        .ToList();

                myInstance.PropertyValues.Add(new ConvertClassPropertyValues
                {
                    ReferenceText = refText,
                    PropertyText = propertyName,
                    PropertyValue = propertyName
                });
            }

            convertResults.Add(myInstance);

            if (!string.IsNullOrWhiteSpace(filterClassName))
            {
                convertResults = convertResults
                    .Where(x => x.ClassName.Contains(filterClassName))
                    .ToList();
            }


            convertResults = convertResults.OrderByDescending(x => x.ClassName == "i18nKeys").ToList();

            return PartialView(convertResults);
        }


        [HttpGet]
        public ActionResult Culture_To_ExcelForamt()
        {
            List<SheetReadModel> convertExcelValues = GetOnlineExcelValueAndMergeCurrentClass();

            //依序覆蓋 (必須以excel的為準)
            //若發現新增的 塞入
            return PartialView(convertExcelValues);
        }



        private List<SheetReadModel> GetOnlineExcelValueAndMergeCurrentClass(bool isMergeCurrentClass = true)
        {
            var originalExcelValues = GoogleSheetUtility.ReadBySheetUrl(_sheetUrl);

            RemoveHeaderLine(originalExcelValues);

            //讀取特定幾格class的檔案 並且轉成
            List<SheetReadModel> convertExcelValues = new List<SheetReadModel>();

            var file1 = ConvertClassToExcelForamt(new List<Type>
            {
                typeof(i18n_TW),
                typeof(i18n_CN),
                typeof(i18n_EN),
            });


            convertExcelValues.AddRange(file1);

            //補上新加的errorCode
            AddNewDefineErrorCode(convertExcelValues);

            //特定的enum自動補上新編號 (定義在這個檔案中)
            var allEnums = AppHelper.GetAutoGeneratei18nKeysEnums();

            foreach(Type enumType in allEnums)
            {
                AddNewDefineForSpecificEnum(convertExcelValues , enumType);
            }

            //補上i18n新加的keys
            if (isMergeCurrentClass)
            {
                AddNewDefineIn_i18nKeys(convertExcelValues);
            }


            //讀取線上的excel檔案 以線上的excel版本為主 覆蓋之
            foreach (var item in convertExcelValues)
            {
                var match = originalExcelValues.FirstOrDefault(x => x.ColumnValues[0] == item.ColumnValues[0] && x.ColumnValues[1] == item.ColumnValues[1]);

                if (match != null)
                {
                    int startIndex = _startIndex;

                    int maxLength = item.ColumnValues.Count() > match.ColumnValues.Count() ? match.ColumnValues.Count() : item.ColumnValues.Count();

                    for (int i = _startIndex; i < maxLength; i++)
                    {
                        item.ColumnValues[i] = match.ColumnValues[i];
                    }
                }
                else
                {
                    string date = DateTime.Now.ToShortDateString();
                    item.ColumnValues.Add(date + " **新增**");
                }
            }

            //讀取線上的excel檔案 找出是否有線上excel有 但class檔案無得
            foreach (var item in originalExcelValues)
            {
                var match = convertExcelValues.FirstOrDefault(x => x.ColumnValues[0] == item.ColumnValues[0] && x.ColumnValues[1] == item.ColumnValues[1]);

                if (match == null)
                {
                    convertExcelValues.Add(item);
                }
            }

            //重新整理排序
            convertExcelValues = convertExcelValues
                .OrderByDescending(x => x.ColumnValues[0] == "Common")
                .ThenByDescending(x => x.ColumnValues[0] == "Enums")
                .ThenByDescending(x => x.ColumnValues[0] == "ErrorCode")
                .ThenBy(x => x.ColumnValues[1])
                .ToList();

            return convertExcelValues;
        }

        private List<SheetReadModel> ConvertClassToExcelForamt(List<Type> types)
        {
            //*請注意第一個一定要放簡體中文 CN的class
            var first = types.FirstOrDefault();

            List<string> columns = first.GetFields().Select(x => x.Name).ToList();

            List<SheetReadModel> listResult = new List<SheetReadModel>();

            foreach (var eachColumnKey in columns)
            {
                var newItem = new SheetReadModel
                {
                    ColumnValues = new List<string>(),
                };

                string categoryName = eachColumnKey.Split('_').FirstOrDefault() ?? "Common";

                newItem.ColumnValues.Add(categoryName);
                newItem.ColumnValues.Add(eachColumnKey);

                foreach (var eachType in types)
                {
                    var cultureText = eachType.GetField(eachColumnKey)?.GetValue(null)?.ToString();
                    newItem.ColumnValues.Add(cultureText);
                }

                listResult.Add(newItem);
            }

            return listResult;
        }
        private void AddNewDefineIn_i18nKeys(List<SheetReadModel> excelModels)
        {
            var allCodes = ExtensionOfEnum.GetEnumAllValue<ApiErrorCode>();

            foreach (var item in typeof(i18nKeys).GetFields())
            {
                var propName = item.Name;
                bool isExist = excelModels.Any(x => x.ColumnValues[1] == propName);

                if (!isExist)
                {
                    string prefix = propName.Split('_').FirstOrDefault() ?? "Common";
                    excelModels.Add(new SheetReadModel
                    {
                        ColumnValues = new List<string>
                        {
                            prefix,
                            propName,
                            "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
                            "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
                            "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX",                        }
                    });
                }
            }
        }


        private void AddNewDefineErrorCode(List<SheetReadModel> excelModels)
        {
            var allCodes = ExtensionOfEnum.GetEnumAllValue<ApiErrorCode>();
            string prefixName = "ErrorCode";

            foreach (ApiErrorCode item in allCodes)
            {
                string enumPropName = prefixName + "_" + item.ToString();
                var isExist = excelModels.Any(x => x.ColumnValues[1] == enumPropName);

                if (!isExist)
                {
                    excelModels.Add(new SheetReadModel
                    {
                        ColumnValues = new List<string>
                        {
                            prefixName,
                            enumPropName,
                            "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
                            "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
                            "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX",                        }
                    });
                }
            }
        }


        private void AddNewDefineForSpecificEnum<T>(List<SheetReadModel> excelModels) where T : IComparable, IFormattable, IConvertible
        {
            var allCodes = ExtensionOfEnum.GetEnumAllValue<T>();
            string prefixName = "Enums";

            foreach (var item in allCodes)
            {
                string enumPropName = prefixName + "_" + typeof(T).Name +"_" + item.ToString();
                var isExist = excelModels.Any(x => x.ColumnValues[1] == enumPropName);

                if (!isExist)
                {
                    excelModels.Add(new SheetReadModel
                    {
                        ColumnValues = new List<string>
                        {
                            prefixName,
                            enumPropName,
                            "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
                            "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
                            "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX",                        
                        }
                    });
                }
            }
        }


        private void AddNewDefineForSpecificEnum(List<SheetReadModel> excelModels , System.Type enumType)
        {
            List<object> enumsAll = new List<object>();

            foreach (var enumValue in System.Enum.GetValues(enumType))
            {
                enumsAll.Add(enumValue);
            }

            string prefixName = "Enums";

            foreach (var item in enumsAll)
            {
                string enumPropName = prefixName + "_" + enumType.Name + "_" + item.ToString();
                var isExist = excelModels.Any(x => x.ColumnValues[1] == enumPropName);

                if (!isExist)
                {
                    excelModels.Add(new SheetReadModel
                    {
                        ColumnValues = new List<string>
                        {
                            prefixName,
                            enumPropName,
                            "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
                            "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
                            "XXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
                        }
                    });
                }
            }
        }



        private void RemoveHeaderLine(List<SheetReadModel> query)
        {
            if (query == null || query.Count == 0)
            {
                return;
            }

            //表頭列 不處理、移除之
            query.RemoveAt(0);
        }



    }
}
