using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class RolePermissionViewModel
    {

        public int RoleID { get; set; }

        public string RoleName { get; set; }

        public int TotalPermissionCount { get; set; }

        public List<string> Permissions { get; set; }

        /// <summary>
        /// 沒有分配到的權限
        /// </summary>
        public List<string> UnassignPermissions { get; set; }

    }



}
