using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IGameVendorRoutineJobService
    {
        /// <summary>
        /// 取得遊戲商遊戲清單
        /// </summary>
        void GetVendorGames(GameVendorTypeEnum vendorID);


        /// <summary>
        /// 昨日 (根據廠商的時區決定時間區間)
        /// </summary>
        /// <param name="vendorType"></param>
        void CalculateYesterdayReportIsCorrect(GameVendorTypeEnum vendorType);


        /// <summary>
        /// 重算指定時區、指定日期
        /// </summary>
        void CalculateTargetDateReportIsCorrect(GameVendorTypeEnum vendorType, DateTime targetDate, double timeZone);


        void CalculateDateRangeReportIsCorrect(GameVendorTypeEnum vendorType, DateTime beginUTC, DateTime endUTC);
    }
}
