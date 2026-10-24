namespace target.Models.ViewModels.Base
{
    public interface ITokenJwtModel
    {
        int LoginUserID { get; set; }

        string LoginName { get; set; }

        string RndID { get; set; }

        long TokenExpiredTime { get; set; }

        /// <summary>
        /// 加密驗證是否被偽造
        /// </summary>
        string SignKey { get; set; }

        int TokenVersion { get; set; }

        /// <summary>
        /// 對應 JwtTokenTypeEnum
        /// </summary>
        int JwtTokenType { get; set; }
    }





}
