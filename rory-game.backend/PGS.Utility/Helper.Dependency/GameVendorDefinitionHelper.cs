using Microsoft.AspNetCore.Mvc;
using target.Models.ViewModels.Base;
using System;
using System.Collections.Generic;
using System.Text;
using System.Linq;
using target.Utility.Extensions;
using target.Utility.Helper.Dependency;
using target.Models.Table;
using target.Utility;
using target.Models.Models;
using target.Models;
using target.Utility.Helper;
using System.Reflection;
using target.Web.Modules.Filter;
using target.Utility.Utility.ThirdParty;
using target.Models.Enums;

namespace target.Utility.Extensions.Web
{
    /// <summary>
    /// 懶得用interface去定義的一些 遊戲商差異性
    /// </summary>
    public class GameVendorDefinitionHelper
    {
        /// <summary>
        /// 哪些遊戲商有支援外部連結
        /// </summary>
        private static List<GameVendorTypeEnum> SupportExternalLinkScreenShotVendors = new List<GameVendorTypeEnum>
        {
            GameVendorTypeEnum.AGG,
        };

        public static bool IsSupportExternalLinkUrl(int vendorID)
        {
            return IsSupportExternalLinkUrl((GameVendorTypeEnum)vendorID);
        }

        public static bool IsSupportExternalLinkUrl(GameVendorTypeEnum vendorID)
        {
            return SupportExternalLinkScreenShotVendors.Contains(vendorID);
        }

    }
}
