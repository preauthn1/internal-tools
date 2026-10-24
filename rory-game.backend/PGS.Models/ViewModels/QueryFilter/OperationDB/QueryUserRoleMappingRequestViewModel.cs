using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryUserRoleMappingRequestViewModel : PaginationParameter
    {

        public List<int> UserIDs { get; set; }

        public List<int> RoleIDs { get; set; }

    }


}
