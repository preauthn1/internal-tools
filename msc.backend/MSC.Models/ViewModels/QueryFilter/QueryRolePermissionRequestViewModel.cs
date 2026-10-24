using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class QueryRolePermissionRequestViewModel : PaginationParameter
    {

        public List<int> IDs { get; set; }

        public List<int> RoleIDs { get; set; }
    }


}
