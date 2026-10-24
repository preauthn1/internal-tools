using Newtonsoft.Json.Linq;
using System.Collections.Generic;

namespace target.Models.ViewModels
{
    public class AGGSingleWalletAPITransactionRequestModel
    {

        /// <summary>
        /// 對應我方的playerID
        /// </summary>
        public int UserID { get; set; }

        public string TransactionId { get; set; }

        public string AccountToken { get; set; }

        public string GameSessionId { get; set; }

        public string GameId { get; set; }

        public string WagerId { get; set; }

        public List<AGGSingleWalletOperationRequestModel> Operations { get; set; }

    }


    public class AGGSingleWalletOperationRequestModel
    {
        public string OperationType { get; set; }

        public string TagType { get; set; }

        public string Action { get; set; }

        public decimal Withdraw { get; set; }

        public decimal Deposit { get; set; }

        public AggTagModel Tag { get; set; }

    }

    public class AggTagModel
    {
        public List<AggTagItemModel> Actions { get; set; }
    }


    public class AggTagItemModel
    {
        public string Action { get; set; }

        public decimal MoneyAmount { get; set; }
    }


    public class AGGSingleWalletAPIPlaceBetResponseViewModel
    {
        public List<AGGSingleWalletEachBalanceItemViewModel> Balances { get; set; }

        public string StatusCode { get; set; }

    }





}
