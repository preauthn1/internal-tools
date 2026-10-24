using Newtonsoft.Json;
using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class OrganizationViewModel
    {
        public int OrganizationID { get; set; }

        public string OrganizationName { get; set; }

        public string OrganizationCode { get; set; }

        public string Memo { get; set; }

        public int Status { get; set; }

        public string Lang { get; set; }

        public int MerchantTotalCount { get; set; }

        public int UsersTotalCount { get; set; }

        public int OrganizationType { get; set; }

        public string CreateBy { get; set; }

        public DateTime CreateTime { get; set; }

        public string UpdateBy { get; set; }

        public DateTime UpdateTime { get; set; }


        [SwaggerIgnoreOutputWhenPagedList]
        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        public List<int> Roles { get; set; }


        [SwaggerIgnoreOutputWhenPagedList]
        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        public List<int> Merchants { get; set; }

    }


}
