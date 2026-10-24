using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QueryExternalApiRequestLogApiRequestParameterModel : QueryExternalApiRequestLogRequestViewModel
    {
        public long? Id { get; set; }

    }

    public class QueryExternalApiRequestLogRequestViewModel : PaginationParameter
    {

        public List<long> IDs { get; set; }

        public string RequestID { get; set; }

        public string Category { get; set; }

        public string SubCategory { get; set; }

        public DateTime? CreateTimeBegin { get; set; }

        public DateTime? CreateTimeEnd { get; set; }

    }


}
