using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class UserRequestModel
    {

        public string XXXXXXXXXXXXXXXXXXXXX { get; set; }

    }

    public class UserCreateRequestModel
    {
        public string UserName { get; set; }
        public string DisplayName { get; set; }
        public string Memo { get; set; }
        public int Status { get; set; }
        public int RoleID { get; set; }
        public string Password { get; set; }
        public string ConfirmPassword { get; set; }
        public int MerchantID { get; set; }


    }

    public class UserCreateUpdateModel
    {

        public int UserID { get; set; }


        public string DisplayName { get; set; }
        public int Status { get; set; }
        public string Password { get; set; }
        public string ConfirmPassword { get; set; }

    }

    public class UserChangeSelfPasswordModel
    {
        public string OldPassword { get; set; }

        public string NewPassword { get; set; }

        public string ConfirmPassword { get; set; }

    }

    public class UserChangePasswordModel
    {

        public int TargetUserID { get; set; }

        public string NewPassword { get; set; }

        public string ConfirmPassword { get; set; }

    }


}
