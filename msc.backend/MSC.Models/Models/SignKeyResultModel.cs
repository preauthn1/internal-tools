using MSC.Models.Enums;
using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Text;

namespace MSC.Models.Models
{
    public class SignKeyResultModel
    {
  
        public string BeforeMD5 { get; set; }

        public string AfterMD5 { get; set; }

        public string AfterMD5_UTF8 { get; set; }

    }
}
