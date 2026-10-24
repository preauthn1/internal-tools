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
    public class QueryMerchantApiRequestParameterModel : QueryMerchantRequestViewModel , IMerchantFilterModel
    {
       

        public int? OrganizationID { get; set; }

    }



    public class QueryMerchantRequestViewModel : PaginationParameter
    {
        [SwaggerIgnore]
        public List<int> MerchantIDs { get; set; }


        [SwaggerIgnore]
        public List<int> OrganizationIDs { get; set; }


        [SwaggerIgnore]
        public string ApiSecretKey { get; set; }        


        [SwaggerIgnore]
        public bool IsIncludeDeleteData { get; set; }


        public int? Status { get; set; }


        public string Currency { get; set; }


        public int? WalletMode { get; set; }


        /// <summary>
        /// 對應 MerchantTypeEnum
        /// </summary>
        public int? Type { get; set; }


    }


}
