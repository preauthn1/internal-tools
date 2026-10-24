using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class UserLoginRequestModel
    {

        public string Account { get; set; }

        public string Password { get; set; }

        public string GoogleAuthCode { get; set; }

    }


}
