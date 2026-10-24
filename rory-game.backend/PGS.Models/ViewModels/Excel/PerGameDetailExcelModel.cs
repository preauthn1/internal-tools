using target.Models.CultureText;
using target.Utility.Helper;
using System;
using System.ComponentModel;

namespace target.Models.ViewModels.Base
{
    public class PerGameDetailExcelModel
    {

        [ExcelColumnSequence(1000)]
        [Description(i18nKeys.Common_Excel_PerGameDetailExcelModel_GameTIme)]
        public string GameTime { get; set; }
        

        [ExcelColumnSequence(1002)]
        [Description(i18nKeys.Common_Excel_PerGameDetailExcelModel_PerGameLogID)]
        public long PerGameLogID { get; set; }


        [ExcelColumnSequence(1004)]
        [Description(i18nKeys.Common_Excel_PerGameDetailExcelModel_ServerID)]
        public int ServerID { get; set; }


        [ExcelColumnSequence(1006)]
        [Description(i18nKeys.Common_Excel_PerGameDetailExcelModel_GameName)]
        public string GameName { get; set; }


        [ExcelColumnSequence(1008)]
        [Description(i18nKeys.Common_Excel_PerGameDetailExcelModel_MerchantName)]
        public string MerchantName { get; set; }



        [ExcelColumnSequence(1022)]
        [Description(i18nKeys.Common_Excel_PerGameDetailExcelModel_PlayerName)]
        public string PlayerName { get; set; }


        [ExcelColumnSequence(1024)]
        public decimal BetAmount { get; set; }

        
        [ExcelColumnSequence(1026)]
        public decimal PayoutAmount { get; set; }

        
        [ExcelColumnSequence(1028)]
        [Description(i18nKeys.Common_Excel_PerGameDetailExcelModel_UpdateAmount)]
        public decimal UpdateAmount { get; set; }
                

        [ExcelColumnSequence(1030)]
        [Description(i18nKeys.Common_Excel_PerGameDetailExcelModel_CurrentAmount)]
        public decimal CurrentAmount { get; set; }
                

        [ExcelColumnSequence(1032)]
        public decimal AfterAmount { get; set; }


    }


}
