using target.Models.ViewModels;
using target.Utility;
using target.Utility.Utility;
using target.Utility.Utility.ThirdParty;
using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;

namespace target.Console
{
    class Program
    {
        private static int _count = 1;


        static void Main(string[] args)
        {
            //相關初始化
            InitializeHelper.Initialize();

            System.Console.WriteLine("Main 1");

            ExecuteProcess();

            System.Console.WriteLine("Main 2");


            System.Console.WriteLine("What is your name? ");
            var name = System.Console.ReadLine();
        }

        private static async Task ExecuteProcess()
        {

            var class1 = new AsyncClass();
            var class2 = new SyncClass();


            System.Console.WriteLine("step 1");

            class1.DoMyTask("no.1");

            System.Console.WriteLine("step 2");

        }


        private static async Task GetUser(string no)
        {
           
        }


        public abstract class BaseClass
        {
            public abstract Task DoMyTask(string no);
        }


        public class SyncClass : BaseClass
        {
            public async override Task DoMyTask(string no)
            {
                Task.Delay(1000);
                System.Console.WriteLine("DoTask:" + no);
            }
        }

        public class AsyncClass : BaseClass
        {
            public async override Task DoMyTask(string no)
            {
                await Task.Delay(1000);
                System.Console.WriteLine("DoTask:" + no);
            }
        }




    }
}
