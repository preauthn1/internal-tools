using MSC.Models.Enums;
using MSC.Models.Models;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{

    public class QueryApiRequestLogEachApiTotalCountModel
    {
        public string RelativeUrl { get; set; }

        public int TotalCount { get; set; }
    }
}
