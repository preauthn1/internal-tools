using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class BaseVendorAPIRequestModel
    {
        public string RandomID { get; set; }

        public int MerchantID { get; set; }
        
        public long TS { get; set; }
    }


    public class BaseProviderAPIRequestModel
    {
        /// <summary>
        /// 經銷商ID
        /// </summary>
        public int ProviderID { get; set; }


        public long TS { get; set; }
    }






}
