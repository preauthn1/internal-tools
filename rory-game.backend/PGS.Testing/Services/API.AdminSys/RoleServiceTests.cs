using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.Models.Models;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Linq;
using target.Models.CoreLibs;
using target.Models.ViewModels;
using target.DA.Repository;
using target.Models.Enums;

namespace target.Testing
{
    [TestClass]
    public class RoleServiceTests : BaseTest
    {

        public IRoleService RoleService { get; set; }


        [TestMethod]
        public void Roles_GetByPagingFilter_Test()
        {
            var query = RoleService.GetByPagingFilter(new QueryRoleApiRequestParameterModel { PageSize = 10});
            Assert_NotEmptyCollection(query.Items);
            return;
        }



        [TestMethod]
        public void Roles_Create_Test()
        {
            string roleName = "unitTest";

            RoleService.Create(new RoleCreateRequestModel
            {
                RoleName = roleName,
            });

            var queryTop10 = RoleService.GetByPagingFilter(new QueryRoleApiRequestParameterModel 
            { 
                PageSize = 10 ,
                OrderByColumn = "RoleID",
                IsASC = false
            });

            var queryMatch = queryTop10.Items.FirstOrDefault(x => x.RoleName == roleName);

            if(queryMatch != null)
            {
                RoleService.Delete(queryMatch.RoleID);
            }

            return;
        }



    }
}

