using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{

    public class ExternalApiRequestLogEntity
    {
        [Key]
        public long Id { get; set; }
        public string ApiCategory { get; set; }
        public string Group { get; set; }
        public string RequestUrl { get; set; }
        public string RequestContent { get; set; }
        public string RequestUUID { get; set; }
        public string ResponseText { get; set; }
        public long ResponseTime { get; set; }
        public int ResponseStatusCode { get; set; }
        public string WebSiteName { get; set; }
        public string Memo { get; set; }
        public DateTime CreateTime { get; set; }

    }
}
