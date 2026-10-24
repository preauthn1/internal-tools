using System;
using System.Collections.Generic;
using System.Linq;
using example.invalid;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Utility.Utility
{
    public class ApiRequestLogTemporaryInfoModel
    {

        public ApiRequestLogTemporaryInfoModel()
        {
            this.Response = new ResponseTemporaryInfoModel();
        }


        public RequestParameters RequestParameters { get; set; }
        
        public HttpWebRequest HttpWebRequest { get; set; }

        public string RequestContent { get; set; }

        public ResponseTemporaryInfoModel Response { get; set; }

        public System.Diagnostics.Stopwatch Timer { get; set; }
        
    }




    public class ResponseTemporaryInfoModel
    {

        public ResponseTemporaryInfoModel()
        {
            this.ResponseCookies = new Dictionary<string, string>();
        }


        public string ResponseContent { get; set; }

        public Dictionary<string, string> ResponseCookies { get; set; }

        public int ResponseStatusCode { get; set; }

        public long ResponseTime { get; set; }

    }







}
