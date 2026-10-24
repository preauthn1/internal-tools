using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{
    public class LanguageTextEntity
    {
        [Key]
        public int Id { get; set; }
        public string LanguageCode { get; set; }
        public string Language { get; set; }
        public string Text { get; set; }


    }
}
