using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QueryVendorConfigApiRequestParameterModel : QueryVendorConfigRequestViewModel
    {

    }


    public class QueryVendorConfigRequestViewModel : PaginationParameter
    {

        public List<int> VendorConfigID { get; set; }

        public List<int> VendorIDs { get; set; }

    }


}
