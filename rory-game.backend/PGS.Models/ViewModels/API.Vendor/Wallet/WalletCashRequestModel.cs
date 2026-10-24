using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{




    /// <summary>
    /// 請求參數
    /// </summary>
    public class WalletCashPlayerAmountRequestModel : BaseVendorRequestParameter
    {

        public string PlayerName { get; set; }
    }


    /// <summary>
    /// 回應參數
    /// </summary>
    public class WalletCashPlayerAmountViewModel
    {

        public string Currency { get; set; }

        public decimal Amount { get; set; }

    }

}
