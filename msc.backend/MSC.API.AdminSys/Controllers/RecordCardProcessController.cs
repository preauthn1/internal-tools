using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using MSC.BL.Interface;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Web.Modules.Filter;
using MSC.Models.Paging;
using MSC.Models.Enums;

namespace MSC.API.AdminSys.Controllers
{
    [AdminSysTokenrAuthorizeFilter]
    public class RecordCardProcessController : BaseController
    {
        public ICardConsumeRecordService CardConsumeRecordService { get; set; }

        public ICardLinkRequestRecordService CardLinkRequestRecordService { get; set; }

        public ICardConsumeMerchantNotifyQueueService CardConsumeMerchantNotifyQueueService { get; set; }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.Record_MerchantLinkRequest_List)]
        [Route("records/card-link-request")]
        [HttpGet]
        public PagedList<CardLinkRequestRecordViewModel> GetByID([FromQuery] QueryCardLinkRequestRecordApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = CardLinkRequestRecordService.GetByPagingFilter(userID, filter);
            return query;
        }



        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.Record_CardConsume_List)]
        [Route("records/card-consume")]
        [HttpGet]
        public PagedList<CardConsumeRecordViewModel> GetCardConsumeRecordPagingByFilter([FromQuery] QueryCardConsumeRecordApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = CardConsumeRecordService.GetByPagingFilter(userID, filter);
            return query;
        }


        /// <summary>
        /// 取得三大步驟的時間圖
        /// </summary>
        [AdminPermissionAuthorizeFilter(
            AdminPermissionTokenEnum.Record_CardConsumeMerchantNotifyQueue_List , 
            AdminPermissionTokenEnum.Record_MerchantLinkRequest_List ,
            AdminPermissionTokenEnum.Record_CardConsume_List
            )]
        [Route("records/process-detail")]
        [HttpGet]
        public MerchantRequestProcessDetailInfoViewModel GetProcessDetail([FromQuery] MerchantRequestProcessDetailInfoRequestModel filter)
        {
            int userID = GetCurrentUserID();
            var query = CardLinkRequestRecordService.GetProcessDetail(userID, filter);
            return query;
        }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.Record_CardConsumeMerchantNotifyQueue_List)]
        [Route("records/merchant-notify-queue")]
        [HttpGet]
        public PagedList<CardConsumeMerchantNotifyQueueViewModel> GetNotifyQueueByPagingFilter([FromQuery]QueryCardConsumeMerchantNotifyQueueApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = CardConsumeMerchantNotifyQueueService.GetByPagingFilter(userID , filter);
            return query;
        }



    }
}
