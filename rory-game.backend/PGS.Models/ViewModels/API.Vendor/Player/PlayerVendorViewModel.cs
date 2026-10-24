using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryPlayerAccountExistRequestModel: BaseVendorRequestParameter
    {

        public string PlayerName { get; set;}

    }


    public class CreatePlayerAccountExistRequestModel : BaseVendorRequestParameter
    {
        public string PlayerName { get; set; }

    }



    public class PlayerVendorViewModel
    {

        public bool IsExist { get; set; }

        public DateTime? CreateTime { get; set; }

    }


    public class PlayerCreateResultVenodrApiViewModel
    {

    }


}
