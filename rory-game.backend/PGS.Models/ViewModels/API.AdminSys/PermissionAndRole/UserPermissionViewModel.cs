using Newtonsoft.Json;
using System.Collections.Generic;

namespace target.Models.ViewModels
{

    public class UserPermissionViewModel
    {        
        public List<UserPermissionItemViewModel> NestedPermissions { get; set; }


        public List<PermissionTokenBriefViewModel> AllPermissions { get; set; }

    }
    


    public class UserPermissionItemViewModel
    {
        public int PermissionID { get; set; }


        [JsonIgnore]
        public int RootPermissionID { get; set; }

        public int ParentPermissionID { get; set; }

        public string DisplayName { get; set; }

        public string PageUrl { get; set; }

        public int Sort { get; set; }

        public bool IsShowInMenu { get; set; }

        /// <summary>
        /// 參考 PermissionTypeEnum
        /// </summary>
        public int PermissionType { get; set; }

        public string Token { get; set; }

        public List<UserPermissionItemViewModel> ChildPermissions { get; set; }


    }


    public class PermissionTokenBriefViewModel
    {
        public int PermissionID { get; set; }

        public int ParentPermissionID { get; set; }

        public string Token { get; set; }

        public string PageUrl { get; set; }

    }



}
