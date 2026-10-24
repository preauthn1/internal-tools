namespace target.Models.ViewModels
{
    public class BetRecordManualSettleRequestModel
    {
        /// <summary>
        /// PGS單號
        /// </summary>
        public string TransactionNo { get; set; }

        public decimal PayoutAmount { get; set; }
        public decimal ValidAmount { get; set; }
        public string Memo { get; set; }
    }


}
