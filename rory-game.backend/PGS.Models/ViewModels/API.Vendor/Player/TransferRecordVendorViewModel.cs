using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryTransferRecordRequestModel : BaseVendorRequestParameter
    {
     

        /// <summary>
        /// PGS方的轉帳單號
        /// </summary>
        public string PGSTransactionNo { get; set; }


        /// <summary>
        /// 對接商戶的轉帳單號
        /// </summary>
        public string MerchantTransactionNo { get; set; }
    }



    public class TransferRecordVendorViewModel
    {
        public string PGSTransactionNo { get; set; }
        public string MerchantTransactionNo { get; set; }
        public int TransferType { get; set; }
        public int PlayerID { get; set; }
        public string PlayerName { get; set; }
        public decimal Amount { get; set; }
        public decimal BeforeAmount { get; set; }
        public decimal AfterAmount { get; set; }
        public DateTime CreateTime { get; set; }

    }



}
