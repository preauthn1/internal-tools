using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class LanguageTextViewModel
    {
        public string LanguageCode { get; set; }

        public int TotalCount { get; set; }

        public string AllTexts { get; set; }

        public List<string> AllTextList { get; set; }

    }


    public class LanguageTextSingleInfoViewModel
    {
        public string InsertSQL { get; set; }

        public List<LanguageTextItemViewModel> Items { get; set; }
    }



    public class LanguageTextItemViewModel
    {
        public int Id { get; set; }
        public string LanguageCode { get; set; }
        public string Language { get; set; }
        public string LanguageText { get; set; }
        public string Text { get; set; }
    }


    public class LanguageTextUpdateRequestModel
    {
        public int Id { get; set; }
        public string LanguageCode { get; set; }
        public string Language { get; set; }
        public string Text { get; set; }

    }

    public class LanguageTextBatchUpdateRequestModel
    {
        public string LanguageCode { get; set; }

        public List<LanguageTextBatchUpdateEachLanguageRequestModel> Items { get; set; }

    }


    public class LanguageTextBatchUpdateEachLanguageRequestModel
    {
        public string Language { get; set; }
        public string Text { get; set; }
    }






}
