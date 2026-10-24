using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QueryLanguageTextApiRequestParameterModel : QueryLanguageTextRequestViewModel
    {

    }


    public class QueryLanguageTextRequestViewModel : PaginationParameter
    {

        public string LanguageCodeLikeQuery { get; set; }

        public string LanguageCode { get; set; }

        public string Language { get; set; }

        public List<string> LanguageCodeList { get; set; }

    }


}
