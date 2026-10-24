using System;
using System.Web;
using example.invalid;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using ClosedXML.Excel;
using System.Reflection;
using game.Entity;
using System.Linq;
using target.Utility.Extensions;

namespace target.Utility.Helper
{
    public class SimpleExcelHelper
    {

        /// <summary>
        /// 將 list model轉成excel格式
        /// </summary>
        /// <typeparam name="T"></typeparam>
        /// <param name="data"></param>
        /// <returns></returns>
        public static Stream ConvertModelToExcelFileStream<T>(IList<T> data) where T : class
        {
            var singleSheet = new ExcelExportParameterModel
            {
                SheetName = "Sheet1",
                ListData = data
            };

            var result = ConvertMultiModelToExcelFileStream(new List<ExcelExportParameterModel> { singleSheet });
            return result;
        }


        /// <summary>
        /// 將 list model轉成excel格式 (有好幾個分頁的格式)
        /// </summary>
        /// <typeparam name="T"></typeparam>
        /// <param name="data"></param>
        /// <returns></returns>
        public static Stream ConvertMultiModelToExcelFileStream(List<ExcelExportParameterModel> eachSheetData)
        {
            List<DataTable> dataTables = new List<DataTable>();

            foreach (var item in eachSheetData)
            {
                dynamic list = item.ListData;
                var table = ConvertToDataTable(list, item);
                dataTables.Add(table);
            }

            using (XLWorkbook wb = new XLWorkbook())
            {
                foreach (var ds in dataTables)
                {
                    wb.Worksheets.Add(ds);
                }

                wb.Style.Alignment.Horizontal = XLAlignmentHorizontalValues.Center;
                wb.Style.Font.Bold = true;

                using (MemoryStream myMemoryStream = new MemoryStream())
                {
                    wb.SaveAs(myMemoryStream);

                    // return memory stream as byte array
                    var result = myMemoryStream.ToArray();
                    return new MemoryStream(result);
                }
            }
        }


        private static DataTable ConvertToDataTable<T>(IEnumerable<T> data, ExcelExportParameterModel setting) where T : new()
        {
            int startIndexValue = 10;
            PropertyInfo[] properties = (typeof(T))
                .GetProperties()
                .OrderBy(x =>
                {
                    startIndexValue++;
                    var attributes = x.GetCustomAttributes(typeof(ExcelColumnSequenceAttribute), false);

                    if (attributes == null || attributes.Length == 0)
                        return startIndexValue;

                    return ((ExcelColumnSequenceAttribute)attributes.FirstOrDefault()).Sequence;
                })
                .ToArray();

            DataTable table = new DataTable(setting.SheetName);

            foreach (PropertyInfo prop in properties)
            {
                string columnName = GetDescriptionOrColumnName(prop);

                Type type = Nullable.GetUnderlyingType(prop.PropertyType) ?? prop.PropertyType;

                //如果欄位有標記 特定的Attribute 屬性，型別強制轉為string
                if (IsAttributeExistAndNumberFormat<ToCashMoneyFormatAttribute>(prop) ||
                    IsAttributeExistAndNumberFormat<ToMoneyFormatAttribute>(prop) ||
                    IsAttributeExistAndNumberFormat<DisplayEmptyWhenZero>(prop) ||
                    prop.PropertyType == typeof(decimal)
                   )
                {
                    type = typeof(string);
                }

                table.Columns.Add(columnName, type);
            }

            List<T> dataToList = data.ToList();

            //如果有指定要自動產生加總欄位
            if (setting.IsAutoGenerateTotalSumRow && dataToList.IsNotEmptyNull())
            {
                var totalSumRow = dataToList.CalculateTotalSumMergeToOneObject();
                dataToList.Add(totalSumRow);
            }


            var propertyAttributeCheck = new Dictionary<string, PropertyAttributeCheckResult>();
            
            //開始塞資料前 先檢查每個欄位是否有定義特殊attribute
            foreach (PropertyInfo prop in properties)
            {
                var newInstance = new PropertyAttributeCheckResult
                {
                    ColumnDisplayName = GetDescriptionOrColumnName(prop)
                };

                if (prop.PropertyType == typeof(decimal))
                {
                    newInstance.isDecimal = true;
                }

                //如果欄位有標記 ToCashAttribute 屬性
                if (IsAttributeExistAndNumberFormat<ToCashMoneyFormatAttribute>(prop))
                {
                    newInstance.IsToCashMoneyFormat = true;
                }

                //如果欄位有標記 ToCashAttribute 屬性
                if (IsAttributeExistAndNumberFormat<ToMoneyFormatAttribute>(prop))
                {
                    newInstance.IsToMoneyFormat = true;
                }

                //如果有指定數值為0 改顯示空字串
                if (IsAttributeExistAndNumberFormat<DisplayEmptyWhenZero>(prop))
                {
                    newInstance.IsDisplayEmptyWhenZero = true;
                }

                propertyAttributeCheck.Add(prop.Name, newInstance);
            }


            foreach (T item in dataToList)
            {
                DataRow row = table.NewRow();
                foreach (PropertyInfo prop in properties)
                {
                    var matchProprty = propertyAttributeCheck[prop.Name];

                    string columnName = matchProprty.ColumnDisplayName;
                    var value = prop.GetValue(item) ?? DBNull.Value;

                    if (matchProprty.IsDisplayEmptyWhenZero)
                    {
                        //如果有指定數值為0 改顯示空字串
                        string valueToString = value.ToString();

                        if (valueToString == "0")
                        {
                            value = string.Empty;
                        }
                        else
                        {
                            value = valueToString;
                        }
                    }
                    else if (matchProprty.IsToCashMoneyFormat)
                    {     
                        //如果欄位有標記 ToCashAttribute 屬性
                        value = ((long)value).To_Cash();
                    }
                    else if (matchProprty.IsToMoneyFormat)
                    {
                        //如果欄位有標記 ToCashAttribute 屬性
                        value = ((long)value).ToString("N2");
                    }                                      
                    else if(matchProprty.isDecimal)
                    {
                        value = ((decimal)value).ToString("N2");
                    }

                   

                    row[columnName] = value;
                }
                table.Rows.Add(row);
            }

            return table;
        }


        private static bool IsAttributeExistAndNumberFormat<T>(PropertyInfo prop)
        {
            var attributes = prop.GetCustomAttributes(typeof(T), false);

            if (attributes == null || attributes.Length == 0)
                return false;

            if (prop.PropertyType == typeof(long) || prop.PropertyType == typeof(int))
                return true;

            return false;
        }


        private static string GetDescriptionOrColumnName(PropertyInfo prop)
        {
            try
            {
                var attributes = prop.GetCustomAttributes(typeof(DescriptionAttribute), false);

                if (attributes == null || attributes.Length == 0)
                    return prop.Name;

                string multiKey = (attributes[0] as DescriptionAttribute).Description;

                var text = MultiLanguageTextHelper.GetCommonCultureTextByClassAndColumnName(multiKey);

                if (string.IsNullOrWhiteSpace(text))
                    return prop.Name;

                return text;
            }
            catch
            {

            }

            return prop.Name;
        }

    }


    public class ExcelExportParameterModel
    {

        /// <summary>
        /// 若傳入多個sheetName不得相同 (請至少用流水號拆開)
        /// </summary>
        public string SheetName { get; set; }


        /// <summary>
        /// 請 務必傳入 IList、List、IEnumerable  
        /// 泛型不太好寫 只好傳入
        /// </summary>
        public object ListData { get; set; }



        /// <summary>
        /// 如果指定為true 會把數字欄位都加總 變成最後一列
        /// </summary>
        public bool IsAutoGenerateTotalSumRow { get; set; }

    }


    public class PropertyAttributeCheckResult
    {
        public string ColumnDisplayName { get; set; }

        public bool isDecimal { get; set; }

        public bool IsToCashMoneyFormat { get; set; }

        public bool IsToMoneyFormat { get; set; }

        public bool IsDisplayEmptyWhenZero { get; set; }

    }

}
