using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;

namespace MSC.Models.Enums
{
    /// <summary>
    /// 後台用戶操作紀錄類型
    /// </summary>
    public enum UserOperateTypeEnum
    {
        /// <summary>
        /// 登入
        /// </summary>
        Login = 1001,



        ExportCardNumber = 3101,

        UpdateSelfPassword = [REDACTED],

        UpdateOtherPassword = [REDACTED],


        /// <summary>
        /// 產生卡號
        /// </summary>
        GenerateCardNumber = 5501,

        /// <summary>
        /// 銷毀卡號
        /// </summary>
        DestroyCardNumber = 5505,
    }
}
