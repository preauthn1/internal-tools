using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class UserEntity
    {
        [Key]
        public int UserID { get; set; }
        public string UserName { get; set; }
        public int OrganizationID { get; set; }
        public string DisplayName { get; set; }
        public int Status { get; set; }
        public string Password { get; set; }
        public DateTime LastLogonTime { get; set; }
        public string LastLogonIPAddress { get; set; }
        public string CreateBy { get; set; }
        public DateTime CreateTime { get; set; }
        public string UpdateBy { get; set; }
        public DateTime UpdateTime { get; set; }
        public string DefaultLang { get; set; }
        public int TimeZoneAreaID { get; set; }
    }
}
