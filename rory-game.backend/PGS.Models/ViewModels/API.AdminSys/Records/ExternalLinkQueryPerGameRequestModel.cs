using target.Models.Table;
using target.Models.ViewModels.GameVendor;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class ExternalLinkQueryPerGameRequestModel
    {
        public string TransactionNo { get; set; }


        /// <summary>
        /// 會存成 SimpleExternalTokenModel json後 轉成 base64
        /// </summary>
        public string Token { get; set; }

    }





}
