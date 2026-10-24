using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QueryVendorApiRequestParameterModel : QueryVendorRequestViewModel
    {


    }


    public class QueryVendorRequestViewModel : PaginationParameter
    {

        public List<int> VendorIDs { get; set; }

    }


}
