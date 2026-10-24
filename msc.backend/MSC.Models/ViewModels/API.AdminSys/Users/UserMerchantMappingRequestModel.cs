using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class UserMerchantMappingRequestModel
    {

        public int UserID { get; set; }


        public List<int> MerchantIDs { get; set; }

    }


}
