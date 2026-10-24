using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.Table
{
    public class ApiRequestLogEntity
    {
        [Key]
        public long Id { get; set; }

        public string UserIdentity { get; set; }
        public string RequestUrl { get; set; }
        public string Doamin { get; set; }
        public string RelativeUrl { get; set; }
        public string HttpMethod { get; set; }
        public string RequestBody { get; set; }
        public string Headers { get; set; }
        public string ResponseText { get; set; }
        public int ResponseStatusCode { get; set; }
        public string IPAddress { get; set; }
        public string IPLocation { get; set; }
        public string Memo { get; set; }
        public string IdentityKey { get; set; }
        public DateTime CreateTime { get; set; }
        public int CreateTimeCategory { get; set; }
        public string StackTrace { get; set; }
        public string DebugInfo { get; set; }
        public long ResponseLength { get; set; }
        public string DebugInfo2 { get; set; }
        public string DebugInfo3 { get; set; }
        public string WebSiteName { get; set; }


    }
}
