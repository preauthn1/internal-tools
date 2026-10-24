using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class ServerMerchantMappingEntity
    {
        [Key]
        public int Id { get; set; }
        public int ServerID { get; set; }
        public int MerchantID { get; set; }

    }
}
