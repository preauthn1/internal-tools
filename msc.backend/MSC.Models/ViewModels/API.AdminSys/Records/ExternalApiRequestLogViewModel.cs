using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class ExternalApiRequestLogViewModel
    {
        public long Id { get; set; }
        public string ApiCategory { get; set; }
        public string Group { get; set; }
        public string RequestUrl { get; set; }
        public string RequestContent { get; set; }    
        public string ResponseText { get; set; }
        public long ResponseTime { get; set; }
        public string Memo { get; set; }
        public DateTime CreateTime { get; set; }

    }



}
