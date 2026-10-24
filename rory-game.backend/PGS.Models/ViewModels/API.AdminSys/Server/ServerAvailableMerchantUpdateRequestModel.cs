using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class ServerAvailableMerchantUpdateRequestModel
    {
        public int MerchantID { get; set; }

        public bool IsEnable { get; set; }
    }



}
