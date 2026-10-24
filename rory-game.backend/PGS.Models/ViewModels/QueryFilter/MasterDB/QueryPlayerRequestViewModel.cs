using target.Models.Paging;
using target.Models.ViewModels;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace game.Entity.Game
{


    public class QueryPlayerApiRequestParameterModel : QueryPlayerRequestViewModel
    {

        public int? PlayerID { get; set; }

        public int? Status { get; set; }

        public int? MerchantID { get; set; }

    }


    public class CreateInternalTestAccountRequestModel
    {
        public string PlayerName { get; set; }

        public int MerchantID { get; set; }

        public bool IsAddToInternalTestAccount { get; set; }
    }


    public class InternalTestAccountPlusWalletRequestModel
    {
        public int PlayerID { get; set; }

        public decimal Amount { get; set; }
    }


    public class QueryPlayerRequestViewModel : PaginationParameter
    {
        public DateTime? CreateTimeBegin { get; set; }

        public DateTime? CreateTimeEnd { get; set; }

        public DateTime? LastLoginTimeBegin { get; set; }

        public DateTime? LastLoginTimeEnd { get; set; }

        public string PlayerName { get; set; }

        public string LastIPAddress { get; set; }

        public string CreateIPAddress { get; set; }



        [SwaggerIgnore]
        public List<int> PlayerIDs { get; set; }


        [SwaggerIgnore]
        public List<string> PlayerNames { get; set; }


        [SwaggerIgnore]
        public List<int> StatusValues { get; set; }


        [SwaggerIgnore]
        public string Token { get; set; }


        [SwaggerIgnore]
        public List<int> MerchantIDs { get; set; }


        [SwaggerIgnore]
        public bool IsNoLockQuery { get; set; }


    }


}
