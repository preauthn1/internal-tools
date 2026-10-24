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
using target.Models.Enums;

namespace target.API.AdminSys.Controllers
{
    [AdminTokenrAuthorizeFilter]
    public class LanguageTextController : BaseController
    {

        public ILanguageTextService LanguageTextService { get; set; }


        [Route("language-text")]
        [HttpGet]
        public PagedList<LanguageTextViewModel> GetMerchantPagingByFilter([FromQuery] QueryLanguageTextApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = LanguageTextService.GetByPagingFilter(userID , filter);
            return query;
        }


        [Route("language-text/{languageCode}")]
        [HttpGet]
        public LanguageTextSingleInfoViewModel GetAllTextByLanguageCode(string languageCode)
        {
            int userID = GetCurrentUserID();
            var query = LanguageTextService.GetAllTextByLanguageCode(languageCode);
            return query;
        }


        [Route("language-text")]
        [HttpPost]
        public void Create(LanguageTextUpdateRequestModel model)
        {
            int userID = GetCurrentUserID();
            LanguageTextService.Create(userID , model);
        }


        [Route("language-text/batch-update")]
        [HttpPut]
        public void BatchUpdate(LanguageTextBatchUpdateRequestModel model)
        {
            int userID = GetCurrentUserID();
            LanguageTextService.BatchUpdate(userID, model);
        }


        [Route("language-text/initialize/game-text")]
        [HttpPost]
        public void CreateNotExistGameNameLanguageText()
        {
            int userID = GetCurrentUserID();
            LanguageTextService.CreateNotExistGameNameLanguageText();
        }


    }
}
