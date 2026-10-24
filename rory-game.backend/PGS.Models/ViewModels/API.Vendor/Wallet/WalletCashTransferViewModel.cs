using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    /// <summary>
    /// 轉帳 (轉入/轉出)
    /// </summary>
    public class WalletTransferRequestModel : BaseVendorRequestParameter
    {
        public decimal Amount { get; set; }

        public string MerchantOrderNo { get; set; }

        public string PlayerName { get; set; }
    }




    /// <summary>
    /// 轉帳結果
    /// </summary>
    public class WalletTransferResultViewModel
    {

        /// <summary>
        /// 我方系統的轉帳單號
        /// </summary>
        public string TransactionNo { get; set; }


        /// <summary>
        /// 本次交易金額
        /// </summary>
        public decimal Amount { get; set; }


        /// <summary>
        /// (也許不需要)
        /// </summary>
        public decimal BeforeAmount { get; set; }


        /// <summary>
        /// (也許不需要)
        /// </summary>
        public decimal AfterAmount { get; set; }
    }

}
