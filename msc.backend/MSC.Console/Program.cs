using MSC.BL.Services;
using MSC.Utility;
using MSC.Utility.IOC;
using target.Console;
using System;
using System.Linq;
using MSC.Utility.Extensions;
using System.Collections.Generic;
using MSC.Utility.Helper;

namespace MSC.Console
{
    class Program
    {
        static void Main(string[] args)
        {
            //相關初始化
            InitializeHelper.Initialize();




            List<string> cardNo = new List<string>
            {
                 "X4XZ8PVRZ66A8L7P",
                 "XC3DAZ958YXNW4MF",
                 "XG5UJXWHM2535QG4",
                 "XNV68B97BUZT3V6E",
                 "XNVUZL976GB3PR2X",
                 "XQ4X5G3AMNQPR5KS",
                 "Y5YJNR7D8LSYY7L5",
                 "Y6VT44PTMSV2L8CZ",
                 "Y77R32MX55KVH9YA"
            };

            foreach(var card in cardNo)
            {
              

            }




            CardNumberHelper.CheckCardIsValid("4ZA8KP87HJSKWTS6");


            //DateTime ts = DateTimeHelper.GetCurrentTime();

            //List<string> cards = new List<string>();

            //for (int i = 0; i < 10000; i++)
            //{
            //    var gen = service.GenerateCardNumber();
            //    var is_valid = service.CheckCardIsValid(gen);
            //    cards.Add(gen);
            //}

            //DateTime ts_end = DateTimeHelper.GetCurrentTime();



            //var qqqqq = cards.Select(x => x.LastOrDefault()).Distinct().ToList(); 





            //var qqqq = (ts_end - ts).TotalMilliseconds;



            //EB3HEHFQ24X7CWL


            //var is_valid_1 = service.CheckCardIsValid(gen_1);
            //var is_valid_2 = service.CheckCardIsValid(gen_2);


            System.Console.WriteLine("--------------------");
            System.Console.WriteLine("最終卡號為:");

        }





    }
}
