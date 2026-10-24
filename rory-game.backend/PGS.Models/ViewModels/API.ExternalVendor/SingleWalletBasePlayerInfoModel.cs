namespace target.Models.ViewModels
{

    /// <summary>
    /// 繼承用 裏面包含了 GameVendorRawRequestInfo 
    /// 可用來操作原始遊戲商請求的物件
    /// </summary>
    public class SingleWalletBasePlayerInfoModel
    {
        /// <summary>
        /// 部分遊戲商 可能會沒有玩家ID之類的參考數據 可指定跳過檢查
        /// </summary>
        public bool IsSkipCheckPlayerData { get; set; }

        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public string SessionToken { get; set; }

        public EachGameVendorRawRequestInfoModel GameVendorRawRequestInfo { get; set; }

    }





}
