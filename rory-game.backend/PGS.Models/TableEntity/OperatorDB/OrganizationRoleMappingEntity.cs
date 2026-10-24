using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class OrganizationRoleMappingEntity
    {
        [Key]
        public int ID { get; set; }
        public int OrganizationID { get; set; }
        public int RoleID { get; set; }

    }
}
