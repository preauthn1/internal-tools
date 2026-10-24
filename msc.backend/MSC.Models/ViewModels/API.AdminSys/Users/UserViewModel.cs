using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class UserViewModel
    {
        public int UserID { get; set; }
        public int OrganizationID { get; set; }
        public string UserName { get; set; }
        public string DisplayName { get; set; }
        public string Memo { get; set; }
        public int Status { get; set; }
        public string StatusText { get; set; }
        public int RoleID { get; set; }
        public string RoleName { get; set; }
        public string AvailableMerchantName { get; set; }
        public int AvailableMerchantCount { get; set; }
        public List<int> AvailableMerchantIDs { get; set; }
        public List<string> AvailableMerchantNames { get; set; }
        public DateTime CreateTime { get; set; }
        public string CreateBy { get; set; }
        public DateTime UpdateTime { get; set; }
        public string UpdateBy { get; set; }

        public bool IsBindGoogleAuth { get; set; }
        public bool IsAuthLoginEnable { get; set; }




    }



}
