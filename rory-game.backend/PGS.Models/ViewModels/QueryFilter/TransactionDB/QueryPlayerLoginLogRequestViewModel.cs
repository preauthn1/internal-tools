using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QueryPlayerLoginLogApiRequestParameterModel : QueryPlayerLoginLogRequestViewModel
    {
        
        public int? PlayerID { get; set; }



    }

    

    public class QueryPlayerLoginLogRequestViewModel : PaginationParameter
    {

        public DateTime? CreateTimeBegin { get; set; }

        public DateTime? CreateTimeEnd { get; set; }

        public string PlayerName { get; set; }

        public int? VendorID { get; set; }




        [SwaggerIgnore]
        public List<int> PlayerIDs { get; set; }

        [SwaggerIgnore]
        public List<int> MerchantIDs { get; set; }


        public bool IsThirdPartyGame { get; set; }



    }


}
