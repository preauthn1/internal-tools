using System;

namespace MSC.Web.AdminSys.Models
{
    public class DataTableParameter
    {
        public DataTableParameter()
        {
            this.Name = "dataTable";
        }

        public DataTableParameter (string name)
        {
            this.Name = name;
        }

        public string Name { get; set; }
    }
}
