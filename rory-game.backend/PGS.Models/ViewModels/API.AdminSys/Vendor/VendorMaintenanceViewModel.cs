using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class VendorMaintenanceViewModel
    {

        public int Id { get; set; }
        public int Type { get; set; }
        public int VendorID { get; set; }
        public string VendorName { get; set; }
        public int Status { get; set; }
        public string StatusText { get; set; }
        public DateTime StartTime { get; set; }
        public DateTime EndTime { get; set; }
        public DateTime CreateTime { get; set; }
        public string CreateBy { get; set; }
        public DateTime UpdateTime { get; set; }
        public string UpdateBy { get; set; }
        public bool IsExpired { get; set; }

    }


    public class VendorMaintenanceCreateViewModel
    {

        public int VendorID { get; set; }
        public DateTime StartTime { get; set; }
        public DateTime EndTime { get; set; }

    }



}
