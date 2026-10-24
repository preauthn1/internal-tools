using target.Utility.Helper;
using System;

namespace target.Utility
{
    public class RandomHelper
    {

        /// <summary>
        /// 种子码
        /// </summary>
        private static char[] character = { '0','1','2','3','4','5','6','7','8','9',
                                            'a','b','c','d','e','f','g','h','i','j','k','l','m','n','o','p','q','r','s','t','u','v','w','x','y','z'
                                          };
        
        /// <summary>
        /// 种子码
        /// </summary>
        private static char[] character1 = { '0', '1', '2', '3', '4', '5', '6', '7', '8', '9' };
        

        public static long GetTradeNo()
        {
            var now = DateTimeHelper.GetCurrentTime();

            string prefix = now.ToString("yyyyMM");
            string suffix = now.ToString("HHmmss");

            return Int64.Parse(prefix + "00000" + suffix);
        }


        /// <summary>
        /// 取得随机数(含英文)
        /// </summary>
        public static string GetRand(int count)
        {
            string strNo = "";

            //设置种子,不重复
            Random rand = new Random(Guid.NewGuid().GetHashCode());
            for (int i = 0; i < count; i++)
            {
                strNo += character[rand.Next(36)].ToString();
            }

            return strNo;
        }

        /// <summary>
        /// 生成附言，月份+四位字数
        /// </summary>
        /// <returns></returns>
        public static string GetRemarkNo()
        {
            //设置种子
            Random rand = new Random(Guid.NewGuid().GetHashCode());

            string strNo = "";
            for (int i = 0; i < 4; i++)
            {
                strNo += character1[rand.Next(10)].ToString();
            }

            strNo = DateTime.Now.ToString("mm") + strNo;

            return strNo;
        }


        /// <summary>
        /// 取得隨機數字
        /// </summary>
        public static int GetRandIntegerNumbers(int min , int max)
        {
            Random rand = new Random(Guid.NewGuid().GetHashCode());
            int rndNumber = rand.Next(min, max);
            return rndNumber;
        }


        /// <summary>
        /// 取得隨機數字
        /// </summary>
        /// <returns></returns>
        public static string GetRandNumbers(int length = 5)
        {
            //设置种子
            Random rand = new Random(Guid.NewGuid().GetHashCode());

            int maxNumber = (int)(Math.Pow(10, length) - 1);

            string strNo = rand.Next(1, maxNumber).ToString().PadLeft(length, '0');

            return strNo;
        }

        /// <summary>
        /// 取得隨機數字
        /// </summary>
        /// <returns></returns>
        public static int GetRandIntegerNumbers(int length = 5)
        {
            //设置种子
            Random rand = new Random(Guid.NewGuid().GetHashCode());
            int maxNumber = (int)(Math.Pow(10, length) - 1);
            var result = rand.Next(1, maxNumber);
            return result;
        }


        /// <summary>
        /// 生成订单编号
        /// </summary>
        /// <param name="pre">前缀</param>
        public static string GetOrderNo(string pre)
        {

            string strNo = GetRand(5);

            strNo = DateTimeHelper.GetCurrentTime().ToString("yyMMddHHmmssfff") + strNo;

            if (pre != string.Empty)
                strNo = pre + strNo;

            return strNo;
        }


        /// <summary>
        /// 不包含dash符號的guid
        /// </summary>
        /// <param name="pre">前缀</param>
        public static string GetGuidWithoutDash()
        {
            string str = Guid.NewGuid().ToString().Replace("-", "");
            return str;
        }



    }
}
