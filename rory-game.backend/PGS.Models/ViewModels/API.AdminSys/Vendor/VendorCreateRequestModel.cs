using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class VendorCreateRequestModel
    {
        public int VendorID { get; set; }
        public string VendorCode { get; set; }
        public string Name { get; set; }

    }


    public class VendorUpdateRequestModel
    {
        public int VendorID { get; set; }
        public string VendorCode { get; set; }
        public string Name { get; set; }
        public int Status { get; set; }

    }



}
