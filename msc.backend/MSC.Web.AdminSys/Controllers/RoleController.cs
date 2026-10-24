using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.Models.Enums;
using MSC.Web.AdminSys.Models;

namespace MSC.Web.AdminSys.Controllers
{
    public class RoleController : Controller
    {
        public IActionResult Permissions()
        {
            PermissionHelper.CheckPermissionAvailable(AdminPermissionTokenEnum.System_User_List);
            return View();
        }

    }
}
