namespace target.Models.ViewModels
{
    public class GamePlaceBetRequestModel
    {
        public int ServerID { get; set; }
        public long TradeID { get; set; }
        public decimal BetAmount { get; set; }
        public string IPAddress { get; set; }
        public string Memo { get; set; }



        [SwaggerIgnore]
        public int PlayerID { get; set; }

        /// <summary>
        /// 另外傳入的
        /// </summary>
        [SwaggerIgnore]
        public string AccessJwtToken { get; set; }

    }



    public class GamePlaceBetResultViewModel
    {
        /// <summary>
        /// 更新後金額 (遊戲幣**** 因為Gs已經對接好了，想改欄位名稱也來不及了)
        /// </summary>
        public decimal AfterAmount { get; set; }

        /// <summary>
        /// 投注單號
        /// </summary>
        public string TransactionNo { get; set; }

    }



}
