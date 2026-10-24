using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class PlayerLaunchGameServerChooseResultModel
    {
        public bool IsRedirectToHall { get; set; }

        public int ServerID { get; set; }

        public string ServerIPAddress { get; set; }

    }


    

}
