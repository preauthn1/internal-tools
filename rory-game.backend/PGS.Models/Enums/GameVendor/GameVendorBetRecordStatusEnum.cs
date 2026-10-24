using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
   
    /// <summary>
    /// 如果有些遊戲商沒有定義明確的注單狀態 
    /// 使用我方自訂的
    /// </summary>
    public enum GameVendorBetRecordStatusEnum
    {        
        /// <summary>
        /// 等待中 未結算
        /// </summary>
        Pending = 0,

        /// <summary>
        /// 已結算
        /// </summary>
        SettleComplete = 1,
               
        /// <summary>
        /// 不明原因失敗(?)
        /// </summary>
        Fail = 9,

        /// <summary>
        /// 人工手動退款、處理完成
        /// </summary>
        ManualRefund = 39,

        /// <summary>
        /// 注單取消
        /// </summary>
        Cancel = 77,

    }
}
