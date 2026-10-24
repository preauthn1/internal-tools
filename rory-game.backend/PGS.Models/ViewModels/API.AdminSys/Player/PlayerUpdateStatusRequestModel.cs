using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class PlayerUpdateStatusRequestModel
    {
        [SwaggerIgnore]
        public int PlayerID { get; set; }

        public string Memo { get; set; }

        public int Status { get; set; }               

    }


    public class PlayerMultipleUpdateStatusRequestModel
    {

        public List<int> PlayerIDs { get; set; }

        public string Memo { get; set; }


        [SwaggerIgnore]
        public int Status { get; set; }

    }



}
