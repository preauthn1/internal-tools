namespace target.Models.ViewModels.GameVendor
{
    public class ExternalVendorBetRecordDetailInfoRequestModel
    {
        /// <summary>
        /// 三方投注單號 (非PGS單號)
        /// </summary>
        public string BetRecordNo { get; set; }

        /// <summary>
        /// 有些遊戲 BetRecordNo + gameCode 或許才是唯一值
        /// </summary>
        public string ThirdPartyGameCode { get; set; }


        public int MerchantID { get; set; }

    }



}
