using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    /// <summary>
    /// 請求參數
    /// </summary>
    public class PlayerCreateParameterModel
    {


        public int MerchantID { get; set; }

        public string PlayerName { get; set; }

        public string Currency { get; set; }

    }


    /// <summary>
    /// 回應參數
    /// </summary>
    public class PlayerCreateResultModel
    {

        public int PlayerID { get; set; }

        public decimal DefaultBalance { get; set; }

    }


}
