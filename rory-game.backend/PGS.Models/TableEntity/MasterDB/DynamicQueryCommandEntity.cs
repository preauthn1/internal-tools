using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{
    public class DynamicQueryCommandEntity
    {
        [Key]
        public int Id { get; set; }
        public string DisplayName { get; set; }
        public string Memo { get; set; }
        public string Sql { get; set; }
        public string FilterFormItems { get; set; }

        public string RequiredColumns { get; set; }


    }
}
