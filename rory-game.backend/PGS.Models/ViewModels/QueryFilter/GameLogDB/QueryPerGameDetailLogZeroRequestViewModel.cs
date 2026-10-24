using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryPerGameDetailLogZeroRequestViewModel : PaginationParameter
    {
        public List<long> IDs { get; set; }

        public List<int> ServerIDs { get; set; }

        public List<int> GameIDs { get; set; }

    }


}
