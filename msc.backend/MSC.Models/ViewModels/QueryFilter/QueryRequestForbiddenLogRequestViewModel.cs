using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class QueryRequestForbiddenLogApiRequestParameterModel : QueryRequestForbiddenLogRequestViewModel
    {
        public int? Type { get; set; }
    }


    public class QueryRequestForbiddenLogRequestViewModel : PaginationParameter
    {

        [SwaggerIgnore]
        public List<int> Types { get; set; }

        public string MemberCode { get; set; }

        public DateTime? CreateTimeBegin { get; set; }

        public DateTime? CreateTimeEnd { get; set; }


        [SwaggerIgnore]
        public List<int> IDs { get; set; }

    }


}
