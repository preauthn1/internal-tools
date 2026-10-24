using MSC.Models.Enums;
using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{


    public class QueryCardInfoApiRequestParameterModel : QueryCardInfoRequestViewModel
    {
       
    }



    public class QueryCardInfoRequestViewModel : PaginationParameter
    {

        public bool IsLockQuery { get; set; }

        public string CardNo { get; set; }


        public string Currency { get; set; }

        public string Memo { get; set; }

        public decimal? Point { get; set; }

        public DateTime? CreateTimeBegin { get; set; }

        public DateTime? CreateTimeEnd { get; set; }

        /// <summary>
        /// DB暫無此欄位 也許未來會補上
        /// </summary>
        public CardInfoStatusEnum? CardStatus { get; set; }


        [SwaggerIgnore]
        public List<string> ListOfCardNo { get; set; }

    }


}
