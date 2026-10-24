using System.Collections.Generic;

namespace target.Models.ViewModels
{
    public class DynamicQueryCommandViewModel
    {

        public int Id { get; set; }
        public string DisplayName { get; set; }
        public string Memo { get; set; }
        public string FilterFormItems { get; set; }
        public List<string> AvailableFilterFormItems { get; set; }

    }


}
