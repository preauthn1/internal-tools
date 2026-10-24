using Newtonsoft.Json;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;

namespace target.Models.ViewModels.Base
{
    public class TransactionRecordExcelModel
    {
        [ExcelColumnSequence(1000)]
        [Description(i18nKeys.Common_Excel_TransactionRecordExcelModel_CreateTime)]
        public string CreateTime { get; set; }

        [ExcelColumnSequence(1002)]
        [Description(i18nKeys.Common_Excel_TransactionRecordExcelModel_MerchantName)]
        public string MerchantName { get; set; }

        [ExcelColumnSequence(1004)]
        [Description(i18nKeys.Common_Excel_TransactionRecordExcelModel_PlayerID)]
        public int PlayerID { get; set; }

        [ExcelColumnSequence(1006)]
        [Description(i18nKeys.Common_Excel_TransactionRecordExcelModel_PlayerName)]
        public string PlayerName { get; set; }

        [ExcelColumnSequence(1008)]
        [Description(i18nKeys.Common_Excel_TransactionRecordExcelModel_WalletMode)]
        public string WalletMode { get; set; }
		
        [ExcelColumnSequence(1010)]
        [Description(i18nKeys.Common_Excel_TransactionRecordExcelModel_GameName)]
        public string GameName { get; set; }

        [ExcelColumnSequence(1012)]
        [Description(i18nKeys.Common_Excel_TransactionRecordExcelModel_UpdateAmount)]
        public decimal UpdateAmount { get; set; }

        [ExcelColumnSequence(1014)]
        [Description(i18nKeys.Common_Excel_TransactionRecordExcelModel_CurrentAmount)]
        public decimal CurrentAmount { get; set; }

        [ExcelColumnSequence(1016)]
        [Description(i18nKeys.Common_Excel_TransactionRecordExcelModel_AfterAmount)]
        public decimal AfterAmount { get; set; }

        [ExcelColumnSequence(1018)]
        [Description(i18nKeys.Common_Excel_TransactionRecordExcelModel_IPAddress)]
        public string IPAddress { get; set; }
		
        [ExcelColumnSequence(1020)]
        [Description(i18nKeys.Common_Excel_TransactionRecordExcelModel_TransactionType)]
        public string TransactionType { get; set; }

    }
}
