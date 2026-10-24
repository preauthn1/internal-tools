using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class QueryUserApiRequestParameterModel : QueryUserRequestViewModel
    {

    }




    public class QueryUserRequestViewModel : PaginationParameter
    {

        public List<int> UserIDs { get; set; }


        public string UserName { get; set; }

        public string PasswordHash { get; set; }

    }


}
