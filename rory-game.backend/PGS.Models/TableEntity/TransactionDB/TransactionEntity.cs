using target.Models.Enums;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class TransactionEntity
    {
        [Key]
        public long TransactionID { get; set; }
        public int PlayerID { get; set; }
        public string PlayerName { get; set; }
        public int MerchantID { get; set; }

        /// <summary>
        /// 請參考 TransactionTypeEnum (請勿亂改型別)
        /// </summary>
        public int TransactionType { get; set; }
        public decimal Amount { get; set; }
        public decimal BeforeAmount { get; set; }
        public decimal AfterAmount { get; set; }
        public DateTime CreateTime { get; set; }
        public long ReferenceID { get; set; }
        public int GameID { get; set; }
        public decimal TaxAmount { get; set; }
        public decimal FlowAmount { get; set; }
        public string IPAddress { get; set; }
    }
}
