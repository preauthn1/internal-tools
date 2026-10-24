using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Utility.Utility
{
    public class ApiRequestResultModel<T>
    {
       
        public T ResultObject { get; set; }

        public Dictionary<string,string> ResponseCookies { get; set; }

        public List<TraceRedirectInfoDto> RedirectTrace { get; set; }

        public Exception Exception { get; set; }

        public bool IsTimeoutException { get; set; }

        public string ResponseText { get; set; }

    }
}
