using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class AGGLiveBaseGame
    {
        public string brandgameid { get; set; }
        public string keyname { get; set; }
    }


    public class AGGLiveAuth
    {
        public string username { get; set; }
        public string passkey { get; set; }
        public string machinename { get; set; }
        public string locale { get; set; }
        public string brandid { get; set; }
    }


}
