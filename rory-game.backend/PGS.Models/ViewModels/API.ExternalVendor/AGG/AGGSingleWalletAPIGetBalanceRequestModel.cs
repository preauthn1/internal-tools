using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class AGGSingleWalletAPIGetBalanceRequestModel
    {
        /// <summary>
        /// 對應我方 playerID
        /// </summary>
        public int UserID { get; set; }

    }




    public class AGGSingleWalletAPIGetBalanceResponseDataModel
    {
        public List<AGGSingleWalletEachBalanceItemViewModel> Balances { get; set; }


        public string StatusCode { get; set; }

    }


    public class AGGSingleWalletEachBalanceItemViewModel
    {
        public string CurrencyCode { get; set; }

        public decimal Amount { get; set; }

    }


}
