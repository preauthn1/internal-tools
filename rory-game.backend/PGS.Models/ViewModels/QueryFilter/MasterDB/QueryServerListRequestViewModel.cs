using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryServerListApiRequestParameterModel : QueryServerListRequestViewModel
    {

        public int? GameID { get; set; }


    }



    public class QueryServerListRequestViewModel : PaginationParameter
    {

        public string ConnIPKeyword { get; set; }

        [SwaggerIgnore]
        public string ServerNameFullMatch { get; set; }

        [SwaggerIgnore]
        public List<int> ServerIDs { get; set; }

        [SwaggerIgnore]
        public List<int> GameIDs { get; set; }

        [SwaggerIgnore]
        public List<int> MerchantIDs { get; set; }


        public string Currency { get; set; }



        /// <summary>
        /// 對應 ServerStatusEnum
        /// </summary>
        public int? Status { get; set; }

    }


}
