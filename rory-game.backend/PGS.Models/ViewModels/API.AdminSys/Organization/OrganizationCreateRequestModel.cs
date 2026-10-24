using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class OrganizationCreateRequestModel
    {

        public string OrganizationName { get; set; }

        public string OrganizationCode { get; set; }

        public string Memo { get; set; }

        public string Lang { get; set; }

        public int Status { get; set; }
       
        public List<int> Roles { get; set; }

        public List<int> Merchants { get; set; }


    }


    public class OrganizationUpdateRequestModel
    {
        public int OrganizationID { get; set; }

        public string OrganizationName { get; set; }

        public string Memo { get; set; }

        public string Lang { get; set; }

        public int Status { get; set; }

        public int OrganizationType { get; set; }

        public List<int> Roles { get; set; }

        public List<int> Merchants { get; set; }

    }



}
