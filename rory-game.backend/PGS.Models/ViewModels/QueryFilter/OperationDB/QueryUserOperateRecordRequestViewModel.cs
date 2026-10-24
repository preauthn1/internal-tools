using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryUserOperateRecordApiRequestParameterModel : QueryUserOperateRecordRequestViewModel
    {
        public int? MerchantID { get; set; }

        public int? PlayerID { get; set; }

        public string PlayerName { get; set; }

    }





    public class QueryUserOperateRecordRequestViewModel : PaginationParameter
    {

        [SwaggerIgnore]
        public List<int> OperateTypes { get; set; }

        [SwaggerIgnore]
        public List<int> MerchantIDs { get; set; }

        [SwaggerIgnore]
        public List<int> TargetIDs { get; set; }

        public DateTime? CreateTimeBegin { get; set; }

        public DateTime? CreateTimeEnd { get; set; }

        public string Memo { get; set; }


    }


}
