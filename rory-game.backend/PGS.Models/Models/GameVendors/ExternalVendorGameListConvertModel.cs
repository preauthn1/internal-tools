using target.Models.Enums;
using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{


    public class ExternalVendorGameListConvertModel
    {

        public string GameName { get; set; }

        public string GameCode { get; set; }

        public string ThirdPartyGameCode { get; set; }

        /// <summary>
        /// 對應 GameVendorTypeEnum
        /// </summary>
        public int GameVendorType { get; set; }


        public Dictionary<LanguageEnum , string> LanguageText { get; set; }

    }


}
