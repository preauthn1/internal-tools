using System.Collections.Generic;

namespace target.Models.ViewModels
{
    public class RolePermissionViewModel
    {
        public int RoleID { get; set; }


        public List<int> Permissions { get; set; }

    }

    public class RolePermissionUpdateModel
    {
        public int RoleID { get; set; }


        public List<int> Permissions { get; set; }

    }



}
