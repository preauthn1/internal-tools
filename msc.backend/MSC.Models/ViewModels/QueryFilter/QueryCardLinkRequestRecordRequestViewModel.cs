using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class QueryCardLinkRequestRecordApiRequestParameterModel : QueryCardLinkRequestRecordRequestViewModel
    {
        public int? MerchantID { get; set; }

    }


    public class QueryCardLinkRequestRecordRequestViewModel : PaginationParameter
    {
        [SwaggerIgnore]
        public List<int> MerchantIDs { get; set; }

        public string MerchantRequestID { get; set; }

        public string SystemRequestID { get; set; }

        public string TokenHash { get; set; }

        public string MemberCode { get; set; }

        public DateTime? CreateTimeBegin { get; set; }

        public DateTime? CreateTimeEnd { get; set; }

    }


}
