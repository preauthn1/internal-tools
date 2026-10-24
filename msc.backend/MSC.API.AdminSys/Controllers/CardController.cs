using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.BL.Interface;
using MSC.Models;
using MSC.Models.Enums;
using MSC.Models.Paging;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.API.Point;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Web.Modules.Filter;

namespace MSC.API.AdminSys.Controllers
{
    [AdminSysTokenrAuthorizeFilter]
    public class CardController : BaseController
    {       
        public ICardGenerateService CardGenerateService { get; set; }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.Card_Card_List)]
        [Route("cards")]
        [HttpGet]
        public PagedList<CardInfoViewModel> GetCardsByPagingFilter([FromQuery]QueryCardInfoApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var card_no = CardGenerateService.GetByPagingFilter(userID, filter);
            return card_no;
        }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.Card_Card_List)]
        [Route("cards/destroy")]
        [HttpPost]
        public CardDestroyResultModel DestoryCard([FromBody] CardDestroyRequestModel model)
        {
            int userID = GetCurrentUserID();
            var result = CardGenerateService.DestroyCards(userID, model);
            return result;
        }


        [AdminPermissionAuthorizeFilter(AdminPermissionTokenEnum.Card_Card_GenerateCardNo)]
        [Route("cards/batch-gen")]
        [HttpPost]
        public object GenearateLargeCardNo([FromBody]CardGenerateRequestModel model)
        {
            int userID = GetCurrentUserID();
            var card_no = CardGenerateService.BatchGenerateCardNumber(userID , model);
            return card_no;
        }

    }
}
