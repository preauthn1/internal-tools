using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class QueryLobbyTokenLogRequestViewModel : PaginationParameter
    {

        public string IdentityHash { get; set; }

        public Guid? IdentityGuid { get; set; }

    }


}
