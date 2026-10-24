namespace target.Models.ViewModels
{
    /// <summary>
    /// 產生後台 盤面結果 的網址
    /// </summary>
    public class BetRecordDetailInfoLinkGenerateRequestModel
    {
        public string MerchantSecretKey { get; set; }

        public int MerchantID { get; set; }

        public string TransactionNo { get; set; }
    }



    /// <summary>
    /// 產生後台 盤面結果 的網址
    /// </summary>
    public class BetRecordDetailInfoLinkGenerateViewModel
    {
        public string Url { get; set; }

    }




}
