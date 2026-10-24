using target.Utility;
using target.Utility.IOC;
using System;
using System.Collections.Generic;
using System.Text;

namespace target.Console
{
    public class InitializeHelper
    {
        public static void Initialize()
        {
            //如果沒呼叫這段 讀config會出問題
            RuntimeHelper.Initialize("Dev");
            InitializeIOC();
        }


        private static void InitializeIOC()
        {
            List<string> projectList = new List<string>
            {
                "target.BL"
            };

            IOCInitializeHelper.InitializeIOC(projectList);
        }

    }
}
