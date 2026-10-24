using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryUnionBetRecordRequestViewModel : PaginationParameter
    {
        public string PlayerName { get; set; }

        public int? PlayerID { get; set; }

        public int? MerchantID { get; set; }

        public List<int> MerchantIDs { get; set; }

        public int? GameID { get; set; }

        public string TransactionNoPGS { get; set; }

        public string TransactionNo { get; set; }

        public int? Status { get; set; }

        public DateTime? CreateTimeBegin { get; set; }

        public DateTime? CreateTimeEnd { get; set; }

    }


}
