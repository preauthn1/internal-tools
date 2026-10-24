using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class TransferRecordViewModel
    {

        public string OrderNo { get; set; }
        public string MerchantOrderNo { get; set; }
        public int TransferType { get; set; }
        public string TransferTypeText { get; set; }
        public int PlayerID { get; set; }
        public string PlayerName { get; set; }
        public int MerchantID { get; set; }
        public string MerchantName { get; set; }


        public decimal Amount { get; set; }
        public decimal BeforeAmount { get; set; }
        public decimal AfterAmount { get; set; }
        public string Memo { get; set; }
        public DateTime CreateTime { get; set; }

    }

    public class ExcelDownloadResult
    {

        public string DownloadUrl { get; set; }

    }



}
