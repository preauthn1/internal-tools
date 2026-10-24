using MSC.Models.Enums;
using System;
using System.Collections.Generic;
using System.Text;

namespace MSC.Models.ViewModels.API.Point
{
    public class CardGenerateResultViewModel
    {
        public List<string> ListOfCardNo { get; set; }
    }

    public class CardInfoViewModel
    {
        public string CardNo { get; set; }
        public string CardPassword { get; set; }
        public decimal Point { get; set; }
        public string Memo { get; set; }
        public string Currency { get; set; }

        public DateTime CreateTime { get; set; }
        public string CreateBy { get; set; }
        public DateTime UpdateTime { get; set; }
        public string UpdateBy { get; set; }
        public CardInfoStatusEnum CardStatus { get; set; }
        public string CardStatusText { get; set; }
        public bool IsConsume { get; set; }

    }


}
