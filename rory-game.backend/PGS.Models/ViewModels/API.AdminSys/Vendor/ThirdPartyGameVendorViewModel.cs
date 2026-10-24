using target.Models.Enums;
using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class ThirdPartyGameVendorBetRecordDetailUrlViewModel
    {
        public string Url { get; set; }

        public string RawJson { get; set; }

        public ThirdPartyVendorInfoTypeEnum InfoType { get; set; }

    }


}
