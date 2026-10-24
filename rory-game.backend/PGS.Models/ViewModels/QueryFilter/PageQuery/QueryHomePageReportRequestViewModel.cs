using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QueryHomePageReportApiRequestParameterModel
    {
        /// <summary>
        /// 起始時間
        /// </summary>
        public DateTime Begin { get; set; }

        /// <summary>
        /// 結束時間
        /// </summary>
        public DateTime End { get; set; }
    }


    public class QueryHomePageReportRequestViewModel : QueryHomePageReportApiRequestParameterModel
    {
        [SwaggerIgnore]
        public int UserID { get; set; }

        [SwaggerIgnore]
        public List<int> MerchantIDs { get; set; }
    }
}
