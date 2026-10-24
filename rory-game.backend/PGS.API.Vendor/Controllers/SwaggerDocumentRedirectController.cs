using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using target.Models.Paging;
using target.Web.Modules.Filter;
using target.Models.Enums;
using target.Models;

namespace target.API.Vendor.Controllers
{

    [ApiExplorerSettings(IgnoreApi = true)]
    public class SwaggerDocumentRedirectController : BaseController
    {


        private static string urlOfTransferMode = "/swagger/index.html";
        private static string urlOfSingleWalletMode = "/swagger/index.html?urls.primaryName=target.API.SeamlessMode";

        // api-document/transfer/EN
        // api-document/transfer/CN
        // api-document/singleWallet/EN
        // api-document/singleWallet/CN

        /// <summary>
        /// 透過中間一層網址導頁 (英文)
        /// </summary>
        [Route("api-document/transfer/EN")]
        [ApiConfig(IsLogApiEnable = false)]
        [HttpGet]
        public IActionResult RedirectTo_TransferDocument_EN()
        {
            Response.Cookies.Append(SystemConstant.HeaderKey_Lang, LanguageEnum.en_us.ToString());     
            return Redirect(urlOfTransferMode);
        }

        /// <summary>
        /// 透過中間一層網址導頁 (英文)
        /// </summary>
        [Route("api-document/transfer/CN")]
        [ApiConfig(IsLogApiEnable = false)]
        [HttpGet]
        public IActionResult RedirectTo_TransferDocument_CN()
        {
            Response.Cookies.Append(SystemConstant.HeaderKey_Lang, LanguageEnum.zh_cn.ToString());
            return Redirect(urlOfTransferMode);
        }

        /// <summary>
        /// 透過中間一層網址導頁 (英文)
        /// </summary>
        [Route("api-document/singleWallet/EN")]
        [ApiConfig(IsLogApiEnable = false)]
        [HttpGet]
        public IActionResult RedirectTo_SingleWalletDocument_EN()
        {
            Response.Cookies.Append(SystemConstant.HeaderKey_Lang, LanguageEnum.en_us.ToString());
            return Redirect(urlOfSingleWalletMode);
        }


        /// <summary>
        /// 透過中間一層網址導頁 (中文)
        /// </summary>
        [Route("api-document/singleWallet/CN")]
        [ApiConfig(IsLogApiEnable = false)]
        [HttpGet]
        public IActionResult RedirectTo_SingleWalletDocument_CN()
        {
            Response.Cookies.Append(SystemConstant.HeaderKey_Lang, LanguageEnum.zh_cn.ToString());
            return Redirect(urlOfSingleWalletMode);
        }
    }
}
