using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class ExternalApiRequestLogViewModel
    {

        public long Id { get; set; }
        public string Category { get; set; }
        public string SubCategory { get; set; }
        public string RequestID { get; set; }
        public string RequestUrl { get; set; }
        public string RequestContent { get; set; }
        public string RequestHeader { get; set; }
        public string RequestBody { get; set; }
        public string ResponseText { get; set; }
        public long ResponseTime { get; set; }
        public int ResponseStatusCode { get; set; }
        public string WebSiteName { get; set; }
        public string Memo { get; set; }
        public DateTime CreateTime { get; set; }

    }




}
