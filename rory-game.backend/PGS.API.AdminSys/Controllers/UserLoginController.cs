using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.API.AdminSys.Controllers
{
    public class UserLoginController : BaseController
    {
        public IUserLoginService UserLoginService { get; set; }


        [Route("users/login")]
        [HttpPost]
        public UserLoginViewModel Login(UserLoginRequestModel model)
        {
            var query = UserLoginService.Login(model);
            return query;
        }

    }
}
