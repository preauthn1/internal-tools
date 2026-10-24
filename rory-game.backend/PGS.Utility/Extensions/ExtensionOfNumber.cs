using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Utility.Extensions
{

    public static class ExtensionOfNumber
    {
        public static decimal GetValueOrAssignValue(this decimal? amount , decimal defaultValue)
        {
            if (amount.HasValue)
                return amount.Value;

            return default;
        }

        public static string To_Cash(this long amount)
        {
            var divideValue = (((decimal)amount) / 100.00M);
            return divideValue.ToString("N2");
        }

        public static string ToMoneyFormat(this decimal amount)
        {
            return amount.ToString("N2");
        }


        /// <summary>
        /// 轉絕對值
        /// </summary>
        public static decimal ToAbs(this decimal amount)
        {
            return amount > 0 ? amount : amount * -1;
        }


        /// <summary>
        /// 轉絕對值
        /// </summary>
        public static long ToAbs(this long amount)
        {
            return amount > 0 ? amount : amount * -1;
        }

        /// <summary>
        /// 轉絕對值
        /// </summary>
        public static int ToAbs(this int amount)
        {
            return amount > 0 ? amount : amount * -1;
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
