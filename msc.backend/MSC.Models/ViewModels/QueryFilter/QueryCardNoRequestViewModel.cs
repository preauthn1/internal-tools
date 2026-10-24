using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class QueryCardNoRequestViewModel : PaginationParameter
    {

        public List<string> ListOfCardNo { get; set; }

    }


}
