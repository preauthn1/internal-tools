using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{
    public class GamePayTableImageEntity
    {
        [Key]
        public int Id { get; set; }
        public int GameID { get; set; }
        public string Code { get; set; }
        public string ImageUrl { get; set; }
        public string Memo { get; set; }
        public DateTime CreateTime { get; set; }


    }
}
