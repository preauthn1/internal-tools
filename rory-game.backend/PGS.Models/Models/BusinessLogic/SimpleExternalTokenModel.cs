using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Models
{
    public class SimpleExternalTokenModel
    {
        public string RequestID { get; set; }

        public DateTime ExpiredTime { get; set; }

        public string TransactionNo { get; set; }

        public string Sign { get; set; }
    }
}
