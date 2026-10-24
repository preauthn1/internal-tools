using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using target.Web.Modules.Filter;
using target.Models.Paging;
using target.Models;
using target.Utility.Helper;
using target.Utility.Extensions;
using target.Models.Models;
using target.Models.Enums;
using target.Utility.Helper.HttpContext;

namespace target.API.AdminSys.Controllers
{
    /// <summary>
    /// (已廢棄) 匯出excel
    /// </summary>
    [AdminTokenrAuthorizeFilter]
    public class ExportExcelController : BaseController
    {


        // 2020-09-25 因為docker上linux環境跑 會有 The type initializer for 'Gdip' threw an exception  錯誤，因此改成前端處理excel
        // 2020-09-25 因為docker上linux環境跑 會有 The type initializer for 'Gdip' threw an exception  錯誤，因此改成前端處理excel
        // 2020-09-25 因為docker上linux環境跑 會有 The type initializer for 'Gdip' threw an exception  錯誤，因此改成前端處理excel
        // 2020-09-25 因為docker上linux環境跑 會有 The type initializer for 'Gdip' threw an exception  錯誤，因此改成前端處理excel
        // 2020-09-25 因為docker上linux環境跑 會有 The type initializer for 'Gdip' threw an exception  錯誤，因此改成前端處理excel
        // 2020-09-25 因為docker上linux環境跑 會有 The type initializer for 'Gdip' threw an exception  錯誤，因此改成前端處理excel
        // 2020-09-25 因為docker上linux環境跑 會有 The type initializer for 'Gdip' threw an exception  錯誤，因此改成前端處理excel
        // 2020-09-25 因為docker上linux環境跑 會有 The type initializer for 'Gdip' threw an exception  錯誤，因此改成前端處理excel
        // 2020-09-25 因為docker上linux環境跑 會有 The type initializer for 'Gdip' threw an exception  錯誤，因此改成前端處理excel
        // 2020-09-25 因為docker上linux環境跑 會有 The type initializer for 'Gdip' threw an exception  錯誤，因此改成前端處理excel
        // 2020-09-25 因為docker上linux環境跑 會有 The type initializer for 'Gdip' threw an exception  錯誤，因此改成前端處理excel



        //public ITransactionRecordService TransactionRecordService { get; set; }

        //public ITransferRecordService TransferRecordService { get; set; }

        //public IPerGameDetailService PerGameDetailService { get; set; }

        //public IStatReportService IStatReportService { get; set; }



        ///// <summary>
        ///// 帳變紀錄匯出
        ///// </summary>
        //[Route("export/excel/transaction")]
        //[HttpPost, HttpGet]
        //public FileResult ExportTransactionExcelByPagingFilter([FromQuery]QueryTransactionRecordRequestParameterModel filter)
        //{
        //    int userID = GetCurrentUserID();
        //    Stream fileStream = TransactionRecordService.ExportExcelByPagingFilter(userID, filter);
        //    return InstanceFileResult("Transaction", fileStream);
        //}

        ///// <summary>
        ///// 轉帳紀錄匯出
        ///// </summary>
        //[Route("export/excel/transfer")]
        //[HttpPost, HttpGet]
        //public FileResult ExportTransferExcelByPagingFilter([FromQuery]QueryTransferRecordApiRequestParameterModel filter)
        //{
        //    int userID = GetCurrentUserID();
        //    Stream fileStream = TransferRecordService.ExportExcelByPagingFilter(userID, filter);
        //    return InstanceFileResult("Transfer", fileStream);
        //}

        ///// <summary>
        ///// 遊戲明細匯出
        ///// </summary>
        //[Route("export/excel/perGame")]
        //[HttpPost, HttpGet]
        //public FileResult ExportPerGameExcelByPagingFilter([FromQuery]QueryPerGameDetailLogApiRequestParameterModel filter)
        //{
        //    int userID = GetCurrentUserID();
        //    Stream fileStream = PerGameDetailService.ExportExcelByPagingFilter(userID, filter);
        //    return InstanceFileResult("PerGame", fileStream);
        //}


        ///// <summary>
        ///// 營運日報表匯出
        ///// </summary>
        //[Route("export/excel/dailyReport")]
        //[HttpPost , HttpGet]
        //public FileResult ExportDailySummaryExcelByFilter([FromQuery]QueryMerchantDailySummaryStatReportReportApiRequestParameterModel filter)
        //{
        //    int userID = GetCurrentUserID();
        //    Stream fileStream = IStatReportService.ExportDailySummaryExcelByFilter(userID, filter);
        //    return InstanceFileResult("DailyReport", fileStream);
        //}



        //private FileStreamResult InstanceFileResult(string dataName , Stream fileStream)
        //{
        //    string fileName = $"Export_{dataName}.xlsx";
        //    string contentType = SystemConstant.Excel_File_ContentType;
        //    return File(fileStream, contentType, fileName);
        //}


    }
}
