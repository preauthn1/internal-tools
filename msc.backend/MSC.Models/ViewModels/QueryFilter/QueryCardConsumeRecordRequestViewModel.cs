using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class QueryCardConsumeRecordApiRequestParameterModel : QueryCardConsumeRecordRequestViewModel
    {


        public int? MerchantID { get; set; }

    }




    public class QueryCardConsumeRecordRequestViewModel : PaginationParameter
    {

        [SwaggerIgnore]
        public List<long> IDs { get; set; }

        [SwaggerIgnore]
        public List<string> ListOfCardNo { get; set; }

        [SwaggerIgnore]
        public List<int> MerchantIDs { get; set; }

        public string MerchantRequestID { get; set; }

        public string CardNo { get; set; }
        
        public string OrderNo { get; set; }

        public string MemberCode { get; set; }

        public DateTime? CreateTimeBegin { get; set; }

        public DateTime? CreateTimeEnd { get; set; }
    }


}
