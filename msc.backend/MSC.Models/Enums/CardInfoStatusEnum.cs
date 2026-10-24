using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;


namespace MSC.Models.Enums
{
    
    /// <summary>
    /// 此欄位DB還沒建
    /// </summary>
    public enum CardInfoStatusEnum
    {

        /// <summary>
        /// 未使用
        /// </summary>
        Pending = 1,        

        /// <summary>
        /// 已被使用
        /// </summary>
        Consumed = 7,


        /// <summary>
        /// 已銷毀
        /// </summary>
        Destroy = 999
    }
}
