using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text;
using target.Models.CultureText;


namespace target.Models.Enums
{
    /// <summary>
    /// 遊戲商類別
    /// </summary>  
    public enum GameVendorTypeEnum
    {
       
        /// <summary>
        /// 預設:天馬
        /// </summary>
        [Description("target")]
        target = 0,
           
        /// <summary>
        /// slot88
        /// </summary>
        [Description("Slot88")]
        Slot88 = 2001,

        

        /* 重要** 每個遊戲商編號 請至少間隔10000 (不要連號) 因為遊戲編號是以遊戲商來區別 */
        /* 重要** 每個遊戲商編號 請至少間隔10000 (不要連號) 因為遊戲編號是以遊戲商來區別 */
        /* 重要** 每個遊戲商編號 請至少間隔10000 (不要連號) 因為遊戲編號是以遊戲商來區別 */
        /* 重要** 每個遊戲商編號 請至少間隔10000 (不要連號) 因為遊戲編號是以遊戲商來區別 */
        /* 重要** 每個遊戲商編號 請至少間隔10000 (不要連號) 因為遊戲編號是以遊戲商來區別 */
        /* 重要** 每個遊戲商編號 請至少間隔10000 (不要連號) 因為遊戲編號是以遊戲商來區別 */
        /* 重要** 每個遊戲商編號 請至少間隔10000 (不要連號) 因為遊戲編號是以遊戲商來區別 */
        /* 重要** 每個遊戲商編號 請至少間隔10000 (不要連號) 因為遊戲編號是以遊戲商來區別 */

        [Description("Funky")]
        Funky = 250000,

        [Description("AGG")]
        AGG = 290000,

        [Description("CrazyFishing")]
        CrazyFishing = 310000,

        [Description("AGGLive")]
        AGGLive = 330000,
    }
}
