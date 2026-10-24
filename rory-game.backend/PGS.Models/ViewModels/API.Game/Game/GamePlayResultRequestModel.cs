using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    /// <summary>
    /// 寫分遊戲結果
    /// </summary>
    public class GamePlayResultRequestModel
    {
        public int PlayerID { get; set; }
        public int ServerID { get; set; }
        public long TradeID { get; set; }


        /// <summary>
        /// v1 => 輸贏金幣
        /// v2 => 派彩金幣
        /// </summary>
        public decimal UpdateAmount { get; set; }     

        public string IPAddress { get; set; }

        public SettlementBetInfoRequestModel BetInfo { get; set; }

        /// <summary>
        /// 有v1/v2版本的API 流程不太一樣
        /// </summary>
        [SwaggerIgnore]
        public bool IsPlaceBetProcess { get; set; }
    }


    public class SettlementBetInfoRequestModel
    {
        public string TransactionNo { get; set; }

        public decimal BetAmount { get; set; }

        public decimal ValidAmount { get; set; }

        public decimal PayoutAmount { get; set; }

        public decimal TaxAmount { get; set; }
    }


    /// <summary>
    /// 寫分遊戲結果 (回應參數)
    /// </summary>
    public class GamePlayResultViewModel
    {
        public decimal AfterAmount { get; set; }

    }


}
