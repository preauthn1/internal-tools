using System;
using System.Collections.Generic;
using System.Text;

namespace MSC.Models.Enums
{
    public enum NotifyQueueStatusEnum
    {
        /// <summary>
        /// 尚未執行
        /// </summary>
        Wait = 0,


        Success = 1,


        /// <summary>
        /// 需要retry
        /// </summary>
        Must_Retry = 3,


        /// <summary>
        /// 正在執行中
        /// </summary>
        On_Executed = 7,


        Fail = 9
    }
}
