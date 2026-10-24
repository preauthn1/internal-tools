using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{


    public class QueryCardConsumeMerchantNotifyQueueApiRequestParameterModel : QueryCardConsumeMerchantNotifyQueueRequestViewModel
    {


        public int? MerchantID { get; set; }

    }


    public class QueryCardConsumeMerchantNotifyQueueRequestViewModel : PaginationParameter
    {
        [SwaggerIgnore]
        public List<int> MerchantIDs { get; set; }

        public string CardNo { get; set; }

        public string ReferenceOrderNo { get; set; }

        public string MerchantRequestID { get; set; }

        public List<string> MerchantRequestIDs { get; set; }

        public DateTime? CreateTimeBegin { get; set; }

        public DateTime? CreateTimeEnd { get; set; }

        public int? Status { get; set; }
    }


}
