using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Diagnostics;
using System.Linq;

namespace target.Utility.Extensions
{
    public static class ExtensionOfException
    {
        /// <summary>
        ///取得 最底層的 inner exception
        /// </summary>
        public static Exception GetMostInnerException(this Exception ex)
        {
            int times = 0;

            Exception resultEx = ex;
            while (true)
            {
                times++;

                if (times >= 25)
                    break;

                if(resultEx.InnerException == null)
                {
                    break;
                }      
                else
                {
                    resultEx = resultEx.InnerException;
                }
            }

            return resultEx;
        }


        public static string GetMessageAndErrorLine(this Exception ex)
        {
            return $"{ex.Message}，{GetErrorOccurLine(ex)}";
        }


        /// <summary>
        /// 取得程式碼錯在第幾行的重點資訊
        /// </summary>
        public static string GetErrorOccurLine(this Exception ex)
        {
            var result = GetErrorOccurLineFromStackTrace(ex.StackTrace);
            return result;
        }


        /// <summary>
        /// 取得程式碼錯在第幾行的重點資訊
        /// </summary>
        public static string GetErrorOccurLineFromStackTrace(string stackTrace)
        {
            List<string> listResult = new List<string>();

            try
            {
                if (string.IsNullOrWhiteSpace(stackTrace))
                    return string.Empty;

                var split = stackTrace.Split(new string[] { "at " }, StringSplitOptions.None);

                foreach (var eachText in split)
                {
                    var splitAgain = eachText.Split(new string[] { @"\" }, StringSplitOptions.None);

                    foreach (var text in splitAgain)
                    {
                        if (text.Contains(".cs"))
                        {
                            listResult.Add(text.Truncate(150));
                        }
                    }
                }

            }
            catch
            {

            }
            return " " + string.Join("", listResult);
        }

    }
}
