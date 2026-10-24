using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QueryMerchantDailySummaryRequestViewModel
    {
        
        public DateTime Start { get; set; }
    
        public DateTime End { get; set; }
        
        public List<int> MerchantIDs { get; set; }

    }


}
