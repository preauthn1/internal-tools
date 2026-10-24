namespace target.Models.ViewModels
{

    public class ImageUploadRequestModel
    {
        /// <summary>
        /// 傳進來的字串 請勿包含 data:image/png;base64,
        /// </summary>
        public string Base64 { get; set; }
    }


}
