using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class PlayerManualDepositRequestModel
    {
        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public decimal Amount { get; set; }

        public string Memo { get; set; }
    }


    public class PlayerManualDepositResultViewModel
    {
        public decimal BeforeAmount { get; set; }

        public decimal AfterAmount { get; set; }

    }


}
