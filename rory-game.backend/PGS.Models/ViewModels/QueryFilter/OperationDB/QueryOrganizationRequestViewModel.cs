using target.Models.Paging;
using target.Models.ViewModels.Base;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    /// <summary>
    /// API接口接到的
    /// </summary>
    public class QueryOrganizationApiRequestParameterModel : QueryOrganizationRequestViewModel
    {


    }
    
    /// <summary>
    /// 程式碼DA端 實際用到的
    /// </summary>
    public class QueryOrganizationRequestViewModel : PaginationParameter
    {

        public string OrganizationCode { get; set; }

        [SwaggerIgnore]
        public bool IsIncludeDeleteData { get; set; }

        [SwaggerIgnore]
        public List<int> OrganizationIDs { get; set; }

    }


}
