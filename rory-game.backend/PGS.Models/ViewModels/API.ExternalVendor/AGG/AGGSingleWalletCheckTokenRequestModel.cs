using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class AGGSingleWalletCheckTokenRequestModel
    {
       
        public string UserToken { get; set; }

    }

    public class AGGSingleWalletAPICheckTokenResponseDataModel
    {
        public int UserID { get; set; }

        public List<AGGSingleWalletEachBalanceItemViewModel> Balances { get; set; }

        public string StatusCode { get; set; }

    }

}
