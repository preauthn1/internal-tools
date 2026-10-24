using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Models
{
    public class DefineTypeCheckResultModel
    {
        public Type Type { get; set; }

        public string TypeName { get; set; }

        public bool IsEnums { get; set; }
    }
}
