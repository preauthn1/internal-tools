using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class OrganizationEntity
    {
        [Key]
        public int OrganizationID { get; set; }
        public string OrganizationName { get; set; }
        public int OrganizationType { get; set; }
        public string Lang { get; set; }
        public int Status { get; set; }
        public string ShortName { get; set; }
        public string OrganizationCode { get; set; }


        public string Memo { get; set; }
        public string CreateBy { get; set; }
        public DateTime CreateTime { get; set; }
        public string UpdateBy { get; set; }
        public DateTime UpdateTime { get; set; }

    }
}
