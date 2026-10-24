using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryServerRobotConfigApiRequestParameterModel : QueryServerRobotConfigRequestViewModel
    {

        public int? GameID { get; set; }

    }






    public class QueryServerRobotConfigRequestViewModel : PaginationParameter
    {

        [SwaggerIgnore]
        public List<int> GameIDs { get; set; }


        [SwaggerIgnore]
        public List<int> ServerIDs { get; set; }

    }


}
