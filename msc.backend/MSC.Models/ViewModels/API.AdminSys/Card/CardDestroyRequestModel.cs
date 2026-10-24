using System.Collections.Generic;

namespace MSC.Models.ViewModels
{
    public class CardDestroyRequestModel 
    {
        public List<string> CardNoList { get; set; }

        public string Memo { get; set; }
    }


    public class CardDestroyResultModel
    {

        public string ResultMessage { get; set; }
    }



}
