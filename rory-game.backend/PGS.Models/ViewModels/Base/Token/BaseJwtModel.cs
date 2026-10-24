namespace target.Models.ViewModels.Base
{
    public class BaseJwtModel : ITokenJwtModel
    {
        public int LoginUserID { get; set; }

        public string LoginName { get; set; }

        public string RndID { get; set; }

        public long TokenExpiredTime { get; set; }

        /// <summary>
        /// 加密驗證是否被偽造
        /// </summary>
        public string SignKey { get; set; }

        public int TokenVersion { get; set; }

        public int JwtTokenType { get; set; }
    }





}
