using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using Microsoft.AspNetCore.Cors;

namespace target.API.Public.Controllers
{
    public class ErrorController : Controller
    {
        public IGameLaunchService GameLaunchService { get; set; }


        /// <summary>
        /// 錯誤頁
        /// </summary>
        [Route("error")]
        [ApiExplorerSettings(IgnoreApi = true)]
        [HttpGet]
        public IActionResult ErrorPage(int code, string msg)
        {
            ViewBag.Code = code;
            ViewBag.Msg = msg;
            return View();
        }

    }
}
