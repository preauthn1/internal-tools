using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{




    public class FakeDataGenerateRequestModel
    {

        /// <summary>
        /// 非必填參數
        /// </summary>
        public int? MerchantID { get; set; }


        public string PlayerName { get; set; }

    }

    public class FakeDataAddPlayerWalletRequestModel
    {
               
        public int PlayerID { get; set; }


    }




}
