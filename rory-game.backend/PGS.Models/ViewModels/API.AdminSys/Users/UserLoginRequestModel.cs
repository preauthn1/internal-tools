using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    /// <summary>
    /// 請求參數
    /// </summary>
    public class UserLoginRequestModel
    {

        public string Account { get; set; }


        public string Password { get; set; }

    }


    /// <summary>
    /// 回應參數
    /// </summary>
    public class UserLoginViewModel
    {
        public string Access_Token { get; set; }

    }


}
