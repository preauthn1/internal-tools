using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;

namespace MSC.Models
{
    public class ExcelConvertToClass
    {
        public bool IsKeyObject { get; set; }

        public string ClassName { get; set; }

        public string Languages { get; set; }

        public List<ConvertClassPropertyValues> PropertyValues { get; set; }
    }



    public class ConvertClassPropertyValues
    {
        public List<string> ReferenceText { get; set; }

        public string PropertyValue { get; set; }

        public string PropertyText { get; set; }
    }


}