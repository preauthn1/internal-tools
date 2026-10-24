using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    /// <summary>
    /// 用戶列表
    /// </summary>
    public class UsersViewModel
    {
        public int UserID { get; set; }


        public string UserName { get; set; }


        public string DisplayName { get; set; }


        public int OrganizationID { get; set; }


        public string OrganizationName { get; set; }


        public string RoleName_CN { get; set; }


        public string RoleName_EN { get; set; }


        public int Status { get; set; }
        

        public int TimeZoneAreaID { get; set; }


        public string TimeZoneAreaName { get; set; }


        public DateTime CreateTime { get; set; }


        public DateTime LastLogonTime { get; set; }


        public string LastLogonIPAddress { get; set; }


    }


    public class UsersUpdatePageViewModel : UsersViewModel
    {
   
        public int RoleID { get; set; }

        public List<int> Merchants { get; set; }

    }



    public class UserCreateModel
    {
        public string UserName { get; set; }

        public string DisplayName { get; set; }

        public int OrganizationID { get; set; }

        public int RoleID { get; set; }

        public string Password { get; set; }

        public string ConfirmPassword { get; set; }

        public List<int> Merchants { get; set; }

        public int TimeZoneAreaID { get; set; }

    }



    public class UserUpdateModel
    {
        public int UserID { get; set; }

        public string DisplayName { get; set; }

        public int Status { get; set; }

        public int RoleID { get; set; }

        public string Password { get; set; }

        public string ConfirmPassword { get; set; }

        public List<int> Merchants { get; set; }

        public int TimeZoneAreaID { get; set; }
    }




}
