using target.Models.Enums;
using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class PlayerWalletTransferModel
    {
        public WalletTransferTypeEnum TranferType { get; set; }

        public int PlayerID { get; set; }

        public decimal Amount { get; set; }

        public string OrderNo_PGS { get; set; }

        public string OrderNo_Merchant { get; set; }

    }

    public class PlayerWalletTransferResultModel
    {

        public bool IsSuccess { get; set; }
        public string Msg { get; set; }
        public decimal BeforeAmount { get; set; }
        public decimal AfterAmount { get; set; }

    }



}
