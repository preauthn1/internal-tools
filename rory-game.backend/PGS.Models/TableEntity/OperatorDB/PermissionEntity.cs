using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class PermissionEntity
    {
        [Key]
        public int PermissionID { get; set; }
        public string PermissionName { get; set; }
        public int ParentPermissionID { get; set; }
        public string PermissionToken { get; set; }
        public string PageUrl { get; set; }
        public int PermissionType { get; set; }
        public int Sort { get; set; }
        public bool IsShowInMenu { get; set; }
        public string IconClass { get; set; }

    }
}
