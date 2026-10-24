using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class UpdateCurrentUserPasswordRequestModel
    {
        public string OriginalPassword { get; set; }

        public string NewPassword { get; set; }

        public string ConfirmPassword { get; set; }
    }


    public class UpdateCurrentUserLangdRequestModel
    {
        public string Lang { get; set; }

    }


}
