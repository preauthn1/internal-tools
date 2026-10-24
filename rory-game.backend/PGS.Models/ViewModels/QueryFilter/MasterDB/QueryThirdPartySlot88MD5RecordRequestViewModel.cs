using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryThirdPartySlot88MD5RecordRequestViewModel : PaginationParameter
    {

        public string MD5Text { get; set; }


        public List<int> PlayerIDs { get; set; }

    }


}
