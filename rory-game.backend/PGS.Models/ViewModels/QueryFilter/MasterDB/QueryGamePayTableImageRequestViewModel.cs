using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QueryGamePayTableImageApiRequestParameterModel : QueryGamePayTableImageRequestViewModel
    {
        public int? GameID { get; set; }
    }


    public class QueryGamePayTableImageRequestViewModel : PaginationParameter
    {

        public List<int> IDs { get; set; }

        public List<int> GameIDs { get; set; }

        public string Code { get; set; }

    }


}
