using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Collections.Specialized;
using target.Utility.Extensions;
using example.invalid;
using example.invalid;
using Microsoft.AspNetCore.Http;
using System.Reflection;
using System.Runtime.Versioning;
using System.Collections.Concurrent;

namespace target.Utility.Helper
{

    public static class ThreadLockHelper
    {

        private static readonly ConcurrentDictionary<string, string> concSafe = new ConcurrentDictionary<string, string>();

        /// <summary>
        /// 若需要以 字串當作執行緒 lock分組依據 
        /// </summary>
        public static string InternConcurrentSafe(string input_lock_key)
        {
            /*
             * 程式端用法 
                lock (ThreadLockHelper.InternConcurrentSafe("{Prefix_你的_lock_key}")) 
                {  
                    .... 
                }
            */
            return concSafe.GetOrAdd(input_lock_key, String.Copy);
        }

    }

}
