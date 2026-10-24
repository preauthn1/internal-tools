using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace MSC.Models.Table
{
    public class UserTokenEntity
    {
        [Key]
        public int Id { get; set; }
        public string Token { get; set; }
        public int UserID { get; set; }
        public DateTime ExpiredTime { get; set; }
        public DateTime CreateTime { get; set; }


    }
}
