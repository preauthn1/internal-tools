using System;
using System.Collections.Generic;
using System.Text;

namespace MSC.Models.ViewModels.RequestViewModel
{
  

    public class CardGenerateRequestModel
    {

        public int ExpectCount { get; set; }

        public decimal Point { get; set; }

        public string Currency { get; set; }

        public string Memo { get; set; }

    }
}
