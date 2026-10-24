namespace target.Models.ViewModels.GameVendor
{
    public class ExternalVendorInputValidateRequestParameter
    {
        /// <summary>
        /// 加密簽 (如果有的話)
        /// </summary>
        public string EncryptSign { get; set; }

        public string SessionOrToken { get; set; }

        public int PlayerID { get; set; }

        public bool IsSkipCheckSessionToken { get; set; }


        /// <summary>
        /// 每家遊戲商自己帶的請求參數
        /// </summary>
        public EachGameVendorRawRequestInfoModel GameVendorRawRequestInfo { get; set; }


    }


}
