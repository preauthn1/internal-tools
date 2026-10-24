using target.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class QueryGameListApiRequestParameter : QueryGameListRequestViewModel
    {
        public int? GameVendorType { get; set; }

    }



    public class QueryGameListRequestViewModel : PaginationParameter
    {
        public string GameCode { get; set; }

        public string ThirdPartyGameCode { get; set; }

        public List<int> GameVendorTypes { get; set; }

        public List<int> GameIDs { get; set; }

        public int? Status { get; set; }

        /// <summary>
        /// true = 過濾掉status = 999的資料
        /// </summary>
        public bool? IsNotDelete { get; set; }

    }


}
