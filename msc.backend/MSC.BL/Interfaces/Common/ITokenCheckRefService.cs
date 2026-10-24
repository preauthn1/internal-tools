using System;
using System.Collections.Generic;
using MSC.Models.CoreLibs;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Models;
using System.Linq;
using MSC.Utility.Helper;
using MSC.Models;
using MSC.DA.Repository;

namespace MSC.BL.Interface
{

    public interface ITokenCheckRefService
    {
        /// <summary>
        /// 僅檢查格式 (無檢查DB)
        /// </summary>
        LinkTokenMemberInfoModel CheckIsTokenFormatCorrect(string token);

        /// <summary>
        /// 僅檢查DB是否有資料
        /// </summary>
        void CheckIsDbExistToken(string token);

        /// <summary>
        /// 檢查cache，是否已失效
        /// </summary>
        void CheckTokenDisabledCache(string tokenHash);

        /// <summary>
        /// 檢查格式、檢查DB 的詳細檢查
        /// </summary>
        LinkTokenMemberInfoModel CheckIsValidToken(string token);
    }
}
