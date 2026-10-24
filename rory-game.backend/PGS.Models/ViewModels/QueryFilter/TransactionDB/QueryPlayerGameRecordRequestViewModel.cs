using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryPlayerGameRecordApiRequestParameterModel : QueryPlayerGameRecordRequestViewModel
    {

        public int? PlayerID { get; set; }

        public int? GameID { get; set; }

        public int? ServerID { get; set; }
    }



    public class QueryPlayerGameRecordRequestViewModel : PaginationParameter
    {

        [SwaggerIgnore]
        public List<int> PlayerIDs { get; set; }

        [SwaggerIgnore]
        public List<int> GameIDs { get; set; }

        [SwaggerIgnore]
        public List<int> ServerIDs { get; set; }

        [SwaggerIgnore]
        public List<int> MerchantIDs { get; set; }

        public string PlayerName { get; set; }

        public string IPAddress { get; set; }

        public long? TradeID { get; set; }

        public DateTime? LoginTimeBegin { get; set; }

        public DateTime? LoginTimeEnd { get; set; }

        public DateTime? LogoutTimeBegin { get; set; }

        public DateTime? LogoutTimeEnd { get; set; }

    }


}
