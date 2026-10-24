using target.Models.CultureText;
using target.Utility.Helper;
using System;
using System.ComponentModel;

namespace target.Models.ViewModels.Base
{




    public class FunkyVendorApiLaunchResponseModel
    {

        public string gameUrl { get; set; }

        public string token { get; set; }

    }



    public class FunkyVendorApiResponseModel<T>
    {
        public int errorCode { get; set; }

        public string errorMessage { get; set; }

        public T data { get; set; }

    }


}
