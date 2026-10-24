using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{

    public class RequestForbiddenLogEntity
    {
        [Key]
        public long Id { get; set; }
        public int Type { get; set; }
        public string BlockKey { get; set; }
        public int MerchantID { get; set; }
        public string MemberCode { get; set; }
        public DateTime CreateTime { get; set; }


    }
}
