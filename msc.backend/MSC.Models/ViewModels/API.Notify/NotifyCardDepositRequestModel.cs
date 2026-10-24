using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class NotifyCardDepositRequestModel
    {

        public long CardConsumeRecordID { get; set; }

        public string NotifyApiUrl { get; set; }

    }


}
