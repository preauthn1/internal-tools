using System;
using System.Collections.Generic;
using System.Text;
using System.Linq;

namespace MSC.Models.Models
{
    public class ProviderKeySettingModel
    {
        public int ProviderID { get; set; }

        public string ProviderSecret { get; set; }

        public string Memo { get; set; }
    }


}
