using System;
using System.Collections.Generic;
using System.Collections.Specialized;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Utility.Utility
{
    public class ApiRequestResultModel<T>
    {
       
        public T ResultObject { get; set; }

        public Exception Exception { get; set; }

        public bool IsTimeoutException { get; set; }

        public string ResponseText { get; set; }

        public int ResponseStatusCode { get; set; }

        public NameValueCollection ResponseHeaders { get; set; }

        public string RelativeUrl { get; set; }

        public long ExecutedTime { get; set; }


        public string ErrorLogMessage
        {
            get
            {
                return Exception?.Message + "_" + ResponseText;
            }
        }
    }


   

}
