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
    public class ServerListServiceTests : BaseTest
    {

        public IServerListService ServerListService { get; set; }


        [TestMethod]
        public void ServerList_GetByPagingFilter_Test()
        {
            var query = ServerListService.GetByPagingFilter(new QueryServerListApiRequestParameterModel { PageSize = 10});
            Assert_NotEmptyCollection(query.Items);
            return;
        }


    }
}

