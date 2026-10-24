using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{




    public class FunkySingleWalletAPIGetBalanceRequestModel
    {
        public string SessionId { get; set; }

        public int PlayerId { get; set; }

    }


    public class FunkySingleWalletAPIGetBalanceResponseDataModel
    {
        public decimal Balance { get; set; }

        public string Currency { get; set; }

    }


}
