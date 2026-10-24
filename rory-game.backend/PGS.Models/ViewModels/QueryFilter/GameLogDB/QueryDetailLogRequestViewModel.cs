using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QueryPerGameDetailLogApiRequestParameterModel : QueryDetailLogRequestViewModel
    {

        public int? PlayerID { get; set; }

        public long? PerGameLogID { get; set; }

        public int? GameID { get; set; }

        public int? ServerID { get; set; }

        public int? MerchantID { get; set; }


        [SwaggerIgnore]
        public string Lang { get; set; }


        [SwaggerIgnore]
        public string AccessToken { get; set; }

    }



    public class QueryDetailLogRequestViewModel : PaginationParameter
    {

        public DateTime? GameStartTimeBegin { get; set; }

        public DateTime? GameStartTimeEnd { get; set; }

        /// <summary>
        /// 僅提供完全
        /// </summary>
        public string PlayerName { get; set; }


        /// <summary>
        /// 單錢包模式的投注單號
        /// </summary>
        public string  TransactionNo { get; set; }


        public bool IsJackpotOnly { get; set; }




        [SwaggerIgnore]
        public List<int> PlayerIDs { get; set; }

        [SwaggerIgnore]
        public List<int> GameIDs { get; set; }

        [SwaggerIgnore]
        public List<int> ServerIDs { get; set; }

        [SwaggerIgnore]
        public List<int> MerchantIDs { get; set; }

        [SwaggerIgnore]
        public List<long> PerGameLogIDs { get; set; }


        public bool IsLockQuery { get; set; }


    }


}
