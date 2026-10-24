using target.Models.CoreLibs;
using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QueryUserApiRequestParameterModel : QueryUserRequestViewModel
    {
        public int? OrganizationID { get; set; }

    }





    public class QueryUserRequestViewModel : PaginationParameter
    {

        [SwaggerIgnore]
        public List<int> UserIDs { get; set; }

        [SwaggerIgnore]
        public List<int> OrganizationIDs { get; set; }


        /// <summary>
        /// *允許不同廠商之間 帳號重複
        /// </summary>
        public string Account { get; set; }

        [SwaggerIgnore]
        public string Password { get; set; }


        /// <summary>
        /// 是否包含已刪除的資料
        /// </summary>
        [SwaggerIgnore]
        public bool IsIncludeDelete { get; set; }



    }


}
