using Newtonsoft.Json;
using target.Models.CultureText;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;

namespace target.Models.ViewModels.Base
{
    public class TransferRecordExcelModel
    {

        [ExcelColumnSequence(1000)]
        [Description(i18nKeys.Common_Excel_TransferRecordExcelModel_CreateTime)]
        public string CreateTime { get; set; }

        [ExcelColumnSequence(1002)]
        [Description(i18nKeys.Common_Excel_TransferRecordExcelModel_MerchantName)]
        public string MerchantName { get; set; }

        [ExcelColumnSequence(1004)]
        [Description(i18nKeys.Common_Excel_TransferRecordExcelModel_PlayerID)]
        public int PlayerID { get; set; }

        [ExcelColumnSequence(1006)]
        [Description(i18nKeys.Common_Excel_TransferRecordExcelModel_PlayerName)]
        public string PlayerName { get; set; }

        [ExcelColumnSequence(1008)]
        [Description(i18nKeys.Common_Excel_TransferRecordExcelModel_TransferTypeName)]
        public string TransferTypeName { get; set; }

        [ExcelColumnSequence(1010)]
        [Description(i18nKeys.Common_Excel_TransferRecordExcelModel_OrderNo)]
        public string OrderNo { get; set; }

        [ExcelColumnSequence(1012)]
        [Description(i18nKeys.Common_Excel_TransferRecordExcelModel_MerchantOrderNo)]
        public string MerchantOrderNo { get; set; }

        [ExcelColumnSequence(1014)]
        [Description(i18nKeys.Common_Excel_TransferRecordExcelModel_BeforeAmount)]
        public decimal BeforeAmount { get; set; }

        [ExcelColumnSequence(1016)]
        [Description(i18nKeys.Common_Excel_TransferRecordExcelModel_Amount)]
        public decimal Amount { get; set; }

        [ExcelColumnSequence(1018)]
        [Description(i18nKeys.Common_Excel_TransferRecordExcelModel_AfterAmount)]
        public decimal AfterAmount { get; set; }

        [ExcelColumnSequence(1020)]
        [Description(i18nKeys.Common_Excel_TransferRecordExcelModel_Memo)]
        public string Memo { get; set; }




    }


}
