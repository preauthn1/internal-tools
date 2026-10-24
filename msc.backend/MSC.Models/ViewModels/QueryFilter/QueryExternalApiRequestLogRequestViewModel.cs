using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class QueryExternalApiRequestLogApiRequestParameterModel : QueryExternalApiRequestLogRequestViewModel
    {

    }
    

    public class QueryExternalApiRequestLogRequestViewModel : PaginationParameter
    {

        [SwaggerIgnore]
        public List<int> IDs { get; set; }

        public string ApiCategory { get; set; }

        public string Group { get; set; }

    }


}
