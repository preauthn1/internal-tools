using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class RequestForbiddenLogViewModel
    {
        public long Id { get; set; }


        /// <summary>
        /// 請參考 RequestForbiddenTypeEnum
        /// </summary>
        public int Type { get; set; }

        public string TypeName { get; set; }        
        public string BlockKey { get; set; }
        public int MerchantID { get; set; }
        public string MerchantName { get; set; }
        public string MemberCode { get; set; }
        public DateTime CreateTime { get; set; }

    }



}
