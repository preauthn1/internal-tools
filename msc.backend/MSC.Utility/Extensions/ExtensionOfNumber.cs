using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Utility.Extensions
{

    public static class ExtensionOfNumber
    {
        /// <summary>
        /// 其實 AdminSys也有一個相同的擴充方法(但只能在admin專案上使用)
        /// </summary>
        public static string To_Cash(this long amount)
        {
            var divideValue = (((decimal)amount) / 100.00M);
            return divideValue.ToString("N2");
        }


        /// <summary>
        /// 是否為整數
        /// </summary>
        public static bool IsIntegerValue(this decimal number)
        {
            return (number % 1.0M) == 0;
        }


    }
}
