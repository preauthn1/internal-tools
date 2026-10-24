using game.Entity.Game;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using target.BL.Interface;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Web.Modules.Filter;
using System;
using System.Collections.Generic;
using example.invalid;

namespace target.API.AdminSys.Controllers
{
    [AdminTokenrAuthorizeFilter]
    public class GamePayTableImageController : BaseController
    {
        public IGamePayTableImageService GamePayTableImageService { get; set; }


        [Route("game/pay-table-images")]
        [HttpGet]
        public PagedList<GamePayTableImageViewModel> GetGamePayTableImagePagingByFilter([FromQuery] QueryGamePayTableImageApiRequestParameterModel filter)
        {
            int userID = GetCurrentUserID();
            var query = GamePayTableImageService.GetByPagingFilter(userID, filter);
            return query;
        }


        [Route("game/pay-table-images")]
        [HttpPost]
        public void CreateOrUpdate(GamePayTableImageUpdateRequestModel model)
        {
            int userID = GetCurrentUserID();
            GamePayTableImageService.CreateOrUpdate(userID, model);
        }


        [Route("game/pay-table-images/{targetID}")]
        [HttpDelete]
        public void Delete(int targetID)
        {
            int userID = GetCurrentUserID();
            GamePayTableImageService.Delete(userID, targetID);
        }

    }
}
