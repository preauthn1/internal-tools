using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class MerchantAvailableAmountViewModel
    {
        public int MerchantID { get; set; }

        public List<decimal> Amounts { get; set; }
    }



}
