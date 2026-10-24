using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.BL.Interface;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Models;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class GameServerSettingServiceTests : BaseTest
    {

        public IGameServerSettingService GameServerSettingService { get; set; }


        [TestMethod]
        public void ServerReRegister_Test()
        {
            GameServerSettingService.ServerReRegister(new SystemReregisterRequestModel
            {
                ServerID = _TEST_DATA.Default.ServerID,
            });

            return;
        }
        

        [TestMethod]
        public void RunServerStartProcess_Test()
        {
            ServerListRepository rep = new ServerListRepository();

            int serverID = _TEST_DATA.Default.ServerID;

            var queryServer = rep.GetById(serverID);

            Assert.IsNotNull(queryServer);

            var updateModel = queryServer.ReturnInjecteObject<ServerStartProcessRequestModel>();
            GameServerSettingService.RunServerStartProcess(updateModel);
            return;
        }


    }
}

