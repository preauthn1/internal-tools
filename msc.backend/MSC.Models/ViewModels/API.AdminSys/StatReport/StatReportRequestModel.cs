using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class StatReportRequestModel
    {
        public int MerchantID { get; set; }

        public double TimeZone { get; set; }

        public DateTime BeginDate { get; set; }

        /// <summary>
        /// 如果查特定某一天會帶入 2020-09-01 - 2020-09-01
        /// </summary>
        public DateTime EndDate { get; set; }
    }


}
