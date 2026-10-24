using target.Models;
using System;
using System.Collections.Generic;
using System.Text;
using target.Utility.Extensions;
using target.Models.Enums;
using target.Models.Attributes;
using example.invalid;
using target.Models.ViewModels;
using System.Linq;
using target.Models.ViewModels.Base;
using target.Models.Models.StoreProcedures;

namespace target.Utility.Helper
{
    public class ExcelDataConvertHelper
    {

        public static Stream ConvertToGameUserGameInfoExcelStream(List<TransferRecordViewModel> list)
        {
            var excelSheets = new List<ExcelExportParameterModel>();

            excelSheets.Add(new ExcelExportParameterModel
            {
                SheetName = "Sheet1",
                ListData = list.Select(x => 
                {
                    var convert = x.ReturnInjecteObject<TransferRecordExcelModel>();
                    convert.CreateTime = ToExcelTimeFormat(x.CreateTime);
                    convert.TransferTypeName = MultiLanguageTextHelper.GetEnumsCultureText((WalletTransferTypeEnum)x.TransferType);

                    return convert;
                })
                .ToList()
            });

            Stream fileStream = SimpleExcelHelper.ConvertMultiModelToExcelFileStream(excelSheets);
            return fileStream;
        }


        public static Stream ConvertToPerGameDetailExcelStream(List<PerGameDetailViewModel> list)
        {
            var excelSheets = new List<ExcelExportParameterModel>();

            excelSheets.Add(new ExcelExportParameterModel
            {
                SheetName = "Sheet1",
                ListData = list.Select(x =>
                {
                    var convert = x.ReturnInjecteObject<PerGameDetailExcelModel>();
                    convert.GameTime = ToExcelTimeFormat(x.GameTime);
                    return convert;
                })
                .ToList()
            });

            Stream fileStream = SimpleExcelHelper.ConvertMultiModelToExcelFileStream(excelSheets);
            return fileStream;
        }

        public static Stream ConvertToTransactionRecordExcelStream(List<TransactionRecordViewModel> list)
        {
            var excelSheets = new List<ExcelExportParameterModel>();

            excelSheets.Add(new ExcelExportParameterModel
            {
                SheetName = "Sheet1",
                ListData = list.Select(x =>
                {
                    var convert = x.ReturnInjecteObject<TransactionRecordExcelModel>();
                    convert.CreateTime = ToExcelTimeFormat(x.CreateTime);
                    convert.WalletMode = ((WalletTypeEnum)x.WalletMode).GetDescription();
                    convert.TransactionType = ((TransactionTypeEnum)x.TransactionType).GetDescription();
                    return convert;
                })
                .ToList()
            });

            Stream fileStream = SimpleExcelHelper.ConvertMultiModelToExcelFileStream(excelSheets);
            return fileStream;
        }



        public static Stream ConvertToMerchantDailySummaryExcelStream(List<MerchantDailySummaryStatReportViewModel> list)
        {
            var excelSheets = new List<ExcelExportParameterModel>();

            excelSheets.Add(new ExcelExportParameterModel
            {
                SheetName = "Sheet1",
                ListData = list.Select(x =>
                {
                    var convert = x.ReturnInjecteObject<MerchantDailySummaryExcelModel>();
                    return convert;
                })
                .ToList()
            });

            Stream fileStream = SimpleExcelHelper.ConvertMultiModelToExcelFileStream(excelSheets);
            return fileStream;
        }






        private static string ToExcelTimeFormat(DateTime time)
        {
            return time.ToString("yyyy-MM-dd HH:mm:ss");
        }


    }
}
