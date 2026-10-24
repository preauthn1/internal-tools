using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QueryCronJobSettingRequestViewModel : PaginationParameter
    {
        public List<int> IDs { get; set; }

        public bool? IsEnable { get; set; }
    }


}
