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

    public interface ICardInputService
    {
        /// <summary>
        /// 一般正常輸入卡號程序
        /// </summary>
        CardInputViewModel InputCard(CardInputRequestModel model);

        /// <summary>
        /// 內部測試使用 (正式環境禁止使用)
        /// </summary>
        CardInputViewModel FakeTestRandomInputCard(CardInputFakeTestRequestModel model);

        /// <summary>
        /// 內部測試使用 (指定某一筆requestID、測試輸入卡號是否正常)
        /// </summary>
        CardInputViewModel FakeTestInputTargetCardNumber(CardInputFakeTestTargetCardNumberRequestModel model);
    }
}
