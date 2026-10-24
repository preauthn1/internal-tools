using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class UserLoginViewModel
    {

        public int UserID { get; set; }

        public string Token { get; set; }

        /// <summary>
        /// 是否尚未綁定google auth
        /// </summary>
        public bool IsGoogleAuthNotInitial { get; set; }


        public FirstTimeGoogleAuthModel GoogleAuth { get; set; }

    }

    

    public class FirstTimeGoogleAuthModel
    {
        public string Base64 { get; set; }
    }






}
