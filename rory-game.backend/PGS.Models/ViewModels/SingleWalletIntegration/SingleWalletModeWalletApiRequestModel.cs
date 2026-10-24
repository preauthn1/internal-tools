using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Text;
using System.Text.Json.Serialization;
using System.Threading.Tasks;

namespace game.Entity.Game
{

    public class SingleWalletModeWalletApiRequestModel
    {

        public string SessionToken { get; set; }
    }


    public class SingleWalletModeWalletApiRequestParameter
    {
        public string SessionTokenPGS { get; set; }

        public int PlayerID { get; set; }
    }


    public class SingleWalletModeWalletApiResponseModel
    {
        public decimal Balance { get; set; }

    }


}
