using System;
using System.Collections.Generic;
using System.Text;

namespace MSC.Models.ViewModels.Base
{
    public class PlayerTokenJwtModel : ITokenJwtModel
    {
        /// <summary>
        /// playerID，但是欄位命名請勿修改
        /// </summary>
        public int LoginUserID { get; set; }


        public string LoginName { get; set; }

        /// <summary>
        /// 語系
        /// </summary>
        public string Lang { get; set; }

        public long TokenExpiredTime { get; set; }

        public string RndID { get; set; }

        /// <summary>
        /// 加密驗證是否被偽造
        /// </summary>
        public string SignKey { get; set; }


        public int TokenVersion { get; set; }
    }





}
