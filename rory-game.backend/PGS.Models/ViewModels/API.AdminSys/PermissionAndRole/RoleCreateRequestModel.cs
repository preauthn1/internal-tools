using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class RoleCreateRequestModel
    {    
        public string RoleName { get; set; }


        public int? TemplateID { get; set; }

    }


    public class RoleUpdateRequestModel
    {
        public int RoleID { get; set; }
        public string RoleName { get; set; }

    }



}
