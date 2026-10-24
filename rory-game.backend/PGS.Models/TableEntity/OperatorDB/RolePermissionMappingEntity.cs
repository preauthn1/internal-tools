using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class RolePermissionMappingEntity
    {
        [Key]
        public int ID { get; set; }
        public int RoleID { get; set; }
        public int PermissionID { get; set; }

    }
}
