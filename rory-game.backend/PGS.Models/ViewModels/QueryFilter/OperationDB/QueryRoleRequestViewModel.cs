using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    
    public class QueryRoleApiRequestParameterModel : QueryRoleRequestViewModel
    {


    }




    public class QueryRoleRequestViewModel : PaginationParameter
    {

        [SwaggerIgnore]
        public List<int> RoleIDs { get; set; }

    }


}
