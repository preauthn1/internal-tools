using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Utility.Helper
{
    /// <summary>
    /// 若有標記此attr，呼叫SimpleExcelHelper 會轉千分位符號 (但不除以100)
    /// </summary>
    public class ToMoneyFormatAttribute : System.Attribute
    {

    }
}
