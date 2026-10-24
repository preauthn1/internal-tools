using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Models.StoreProcedures
{
    public class Stat_ChartDataSPResultModel
    {
        public int? MerchantID { get; set; }

        public string MerchantName { get; set; }

        public int? GameID { get; set; }

        public string GameName { get; set; }

        public DateTime StatTime { get; set; }

		public int DataCount { get; set; }
	}
}
