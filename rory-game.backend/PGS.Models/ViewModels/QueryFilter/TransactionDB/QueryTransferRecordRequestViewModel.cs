using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryTransferRecordApiRequestParameterModel : QueryTransferRecordRequestViewModel
    {

        public int? PlayerID { get; set; }


        public int? MerchantID { get; set; }




    }


    public class QueryTransferRecordRequestViewModel : PaginationParameter
    {
        public string PlayerName { get; set; }

        public string OrderNo { get; set; }

        public string MerchantOrderNo { get; set; }

        public int? TransferType { get; set; }


        public DateTime? CreateTimeBegin { get; set; }


        public DateTime? CreateTimeEnd { get; set; }



        [SwaggerIgnore]
        public List<int> TransferRecordIDs { get; set; }

        [SwaggerIgnore]
        public List<int> PlayerIDs { get; set; }

        [SwaggerIgnore]
        public List<int> MerchantIDs { get; set; }

    }


}
