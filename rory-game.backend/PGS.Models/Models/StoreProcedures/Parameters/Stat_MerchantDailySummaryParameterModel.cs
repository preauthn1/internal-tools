using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Models.StoreProcedures
{
    public class Stat_MerchantDailySummaryParameterModel
    {
        public DateTime DateStart { get; set; }

        public DateTime DateEnd { get; set; }

        public double Timezone { get; set; }

        public string MerchantList { get; set; }

    }
}
