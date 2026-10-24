using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QueryVendorMaintenanceApiRequestParameterModel : QueryVendorMaintenanceRequestViewModel
    {

    }

    public class QueryVendorMaintenanceRequestViewModel : PaginationParameter
    {

        public List<int> VendorIDs { get; set; }

        public bool IsIncludeDelete { get; set; }

        public bool IsNotExpired { get; set; }

        public int? Status { get; set; }
    }


}
