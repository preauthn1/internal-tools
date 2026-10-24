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

namespace target.API.AdminSys.Controllers
{
    [AdminTokenrAuthorizeFilter]
    public class GameListController : BaseController
    {
        public IGameListService GameListService { get; set; }


        [Route("games")]
        [HttpGet]
        public PagedList<GameListViewModel> GetGameListPagingByFilter([FromQuery] QueryGameListApiRequestParameter filter)
        {
            var query = GameListService.GetByPagingFilter(filter);
            return query;
        }


        [Route("games/{targetID}")]
        [HttpGet]
        public GameListViewModel GetByID(int targetID)
        {
            var query = GameListService.GetByID(targetID);
            return query;
        }


        [Route("games")]
        [HttpPost]
        public void Create([FromBody]GameListCreateRequestModel model)
        {
            GameListService.Create(model);
        }


        [Route("games")]
        [HttpPut]
        public void Update([FromBody]GameListUpdateRequestModel model)
        {
            GameListService.Update(model);
        }


        [Route("games/{targetID}")]
        [HttpDelete]
        public void Delete(int targetID)
        {
            GameListService.Delete(targetID);
        }

    }
}
