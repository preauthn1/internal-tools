using Google.Apis.Auth.OAuth2;
using Google.Apis.Services;
using Google.Apis.Sheets.v4;
using Google.Apis.Sheets.v4.Data;
using Google.Apis.Util.Store;
using System;
using System.Collections.Generic;
using example.invalid;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Web;

namespace Lck.GoogleApi
{
    public static class GoogleSheetUtility
    {

        static string[] Scopes =
        {
            SheetsService.Scope.SpreadsheetsReadonly,
        };


        public static string WebRootDirection = "";


        //-------- 限制(注意事項) --------
        //限制不得超過26個columns


        //應用程式的名字需要英文
        static string ApplicationName = "PGS_GoogleSheet_Read";



        public static List<SheetReadModel> ReadBySheetUrl(string sheetUrlOr)
        {
            var result = Read(new SheetReqeustParameter
            {
                SheetIdOrSheetUrl = sheetUrlOr
            });

            return result;
        }


        public static SheetCreateResultModel CreateNew(string title)
        {
            SheetsService service = InstanceGoogleSheetService();
            
            var myNewSheet = new Google.Apis.Sheets.v4.Data.Spreadsheet();
            myNewSheet.Properties = new SpreadsheetProperties();
            myNewSheet.Properties.Title = title;
            
            var newSheet = service.Spreadsheets.Create(myNewSheet).Execute();
            
            return new SheetCreateResultModel
            {
                SheetId = newSheet.SpreadsheetId,
                SheetUrl = newSheet.SpreadsheetUrl
            };
        }



        private static SheetSummaryModel GetSheetSummary(string sheetIdOrSheetUrl)
        {
            SheetsService service = InstanceGoogleSheetService();

            //指定google sheet Id
            string spreadsheetId = ParseSheetId(sheetIdOrSheetUrl);

            var sheetSummary = service.Spreadsheets.Get(spreadsheetId).Execute();

            //取出所有分頁
            List<string> sheetList = new List<string>();

            foreach (Sheet sheet in sheetSummary.Sheets)
            {
                sheetList.Add(sheet.Properties.Title);
            }

            //range (第幾格到第幾格)
            string sheetPageTitle = sheetList.FirstOrDefault();
            
            string range = string.Format("{0}!A1:Z", sheetPageTitle);
            
            int maxRows = service.Spreadsheets.Values.Get(spreadsheetId, range).Execute().Values.Count();

            return new SheetSummaryModel
            {
                SheetId = spreadsheetId,
                MaxRows = maxRows,
                FirstPageName = sheetPageTitle
            };
        }



        public static List<SheetReadModel> Read(SheetReqeustParameter model)
        {
            SheetsService service = InstanceGoogleSheetService();

            //指定google sheet Id
            string spreadsheetId = ParseSheetId(model.SheetIdOrSheetUrl);

            var sheetSummary = GetSheetSummary(spreadsheetId);
            
            //range (第幾格到第幾格)
            string sheetPageName = sheetSummary.FirstPageName;
            string range = string.Format("{0}!A1:Z", sheetPageName); //目前這樣的寫法 超過26 columns會出問題


            SpreadsheetsResource.ValuesResource.GetRequest request =
                    service.Spreadsheets.Values.Get(spreadsheetId, range);
            

            var listResult = new List<SheetReadModel>();

            ValueRange response = request.Execute();
            IList<IList<Object>> values = response.Values;
            if (values != null && values.Count > 0)
            {
                int rowIndex = 1;
                int? columns = null;
                foreach (var row in values)
                {
                    var newRowValue = new SheetReadModel
                    {
                        RowIndex = rowIndex,
                        ColumnValues = new List<string>()
                    };


                    //用第一row 來決定有多少columns
                    if (!columns.HasValue)
                    {
                        columns = row.Count();
                    }

                    for (int i = 0; i < columns; i++)
                    {
                        newRowValue.ColumnValues.Add((GetValueByIndex(row, i) ?? string.Empty).ToString());
                    }

                    listResult.Add(newRowValue);
                    rowIndex++;
                }
            }
            else
            {
                return listResult;
            }

            return listResult;
        }



        public static void InsertNewRowsValue(SheetUpdateModel model)
        {
            SheetsService sheetsService = InstanceGoogleSheetService();
            
            var summary = GetSheetSummary(model.SheetId);
            
            UpdateSheetValue(model, summary.MaxRows + 1, summary.FirstPageName);
        }


        public static void UpdateSheetValue(SheetUpdateModel model , int targetRowIndex , string targetPageName)
        {
            SheetsService sheetsService = InstanceGoogleSheetService();
            
            string spreadsheetId = model.SheetId;
            string range = $"{targetPageName}!A{targetRowIndex}";  // update cell F5 

            ValueRange valueRange = new ValueRange();
            valueRange.MajorDimension = "ROWS";//"ROWS" 橫向逐筆寫入  ;//COLUMNS  直向逐筆寫入

            var oblist = model.UpdateCellValues.Select(x => (object)x).ToList();
            valueRange.Values = new List<IList<object>> { oblist };

            SpreadsheetsResource.ValuesResource.UpdateRequest update = sheetsService.Spreadsheets.Values.Update(valueRange, spreadsheetId, range);
            update.ValueInputOption = SpreadsheetsResource.ValuesResource.UpdateRequest.ValueInputOptionEnum.RAW;
            UpdateValuesResponse result = update.Execute();
            return;
        }

     

        private static string ParseSheetId(string sheetId_or_sheetUrl)
        {
            if (string.IsNullOrWhiteSpace(sheetId_or_sheetUrl))
                throw new System.Exception("invalid parameter sheetId_or_sheetUrl");

            if (!sheetId_or_sheetUrl.Contains("http"))
                return sheetId_or_sheetUrl;

            //用斜線符號拆解後 取字串最長的那段當作sheetId
            var maxLengthText = sheetId_or_sheetUrl.Split('/')
                .Where(x => !x.Contains("google"))
                .OrderByDescending(z => z.Length)
                .FirstOrDefault();

            return maxLengthText;
        }



        private static object GetValueByIndex(IList<object> values, int index)
        {
            if (values.Count() <= index)
                return null;

            return values[index];
        }



        private static SheetsService InstanceGoogleSheetService()
        {           
            GoogleCredential credential;
            
            //請使用服務帳號的方式 產生google token，並放在專案跟目錄下
            string rootPath = WebRootDirection;

            string tokenPath = rootPath + @"\myGoogleToken.json";

            using (var stream =
                new FileStream(tokenPath, FileMode.Open, FileAccess.Read))
            {
                credential = GoogleCredential.FromStream(stream).CreateScoped(Scopes);
            }
            
            // Create Google Sheets API service.
            var service = new SheetsService(new BaseClientService.Initializer()
            {
                HttpClientInitializer = credential,
                ApplicationName = ApplicationName,
            });

            return service;
        }


    }
}
