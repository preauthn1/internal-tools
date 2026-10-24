namespace target.Models.ViewModels
{
    /// <summary>
    /// 請求參數
    /// </summary>
    public class SingleWalletGameRecordRequestModel 
    {
        public string MerchantSecretKey { get; set; }

        public int MerchantID { get; set; }

        public string TransactionNo { get; set; }
    }


}
