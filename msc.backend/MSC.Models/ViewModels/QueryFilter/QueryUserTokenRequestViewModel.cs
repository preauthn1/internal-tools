using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{
    public class QueryUserTokenRequestViewModel : PaginationParameter
    {

        public List<int> IDs { get; set; }

        public List<int> UserIDs { get; set; }

        public string Token { get; set; }

        public bool IsNotExpiredOnly { get; set; }
    }


}
