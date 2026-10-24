using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{
    public class RolePermissionEntity
    {
        [Key]
        public int Id { get; set; }
        public int RoleID { get; set; }
        public string PermissionKey { get; set; }

    }
}
