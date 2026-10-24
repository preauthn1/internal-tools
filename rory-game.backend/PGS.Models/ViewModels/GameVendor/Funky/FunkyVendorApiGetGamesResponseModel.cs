using target.Models.CultureText;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.ComponentModel;

namespace target.Models.ViewModels.Base
{


    public class FunkyVendorApiRootResponseModel : BaseFunkyVendorApiResponseModel
    {

        public List<FunkyVendorApiGetGamesItemModel> gameList { get; set; }

    }


    public class FunkyVendorApiGetGamesItemModel
    {
        public string gameCode { get; set; }
        public string gameName { get; set; }
        public string gameTypeDescription { get; set; }
        public List<string> regulations { get; set; }
        public bool marginAdjustable { get; set; }
        public int suggestedViewWidth { get; set; }
        public int suggestedViewHeight { get; set; }
        public List<string> supportedOrientation { get; set; }
        public bool responsive { get; set; }
        public List<string> supportedLanguage { get; set; }
        public string gameStatus { get; set; }
        public bool isNewGame { get; set; }
        public bool demoGameAvailable { get; set; }
        public string demoGameUrl { get; set; }
    }


}
