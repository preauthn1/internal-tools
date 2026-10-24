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
using System.Collections.Generic;
using System;

namespace target.Testing
{
    [TestClass]
    public class GameListServiceTests : BaseTest
    {

        public IGameListService ServiceInstance { get; set; }


        [TestMethod]
        public void GameList_GetByPagingFilter_Test()
        {
            var query = ServiceInstance.GetByPagingFilter(new QueryGameListApiRequestParameter());
            Assert_NotEmptyCollection(query.Items);
            Assert.IsNotNull(new CRUDResultModel());
            return;
        }



        [TestMethod]
        public void CheckNotDefineEnums()
        {
            var queryGames = ServiceInstance.GetByPagingFilter(new QueryGameListApiRequestParameter
            {
                PageSize = 99999
            });

            //大於1萬的視為測試垃圾資料
            queryGames.Items = queryGames.Items.Where(x => x.GameID <= 10000).ToList();

            var allGames = ExtensionOfEnum.GetEnumAllValue<GameListEnum>();

            List<GameListViewModel> notDefineEnumGames = new List<GameListViewModel>();
            
            foreach(var item in queryGames.Items)
            {
                bool isExist = allGames.Any(x => x.ToInteger() == item.GameID);

                if(!isExist)
                {
                    notDefineEnumGames.Add(item);
                }
            }

            if(notDefineEnumGames.IsNotEmptyNull())
            {
                string msg = string.Join("、", notDefineEnumGames.Select(x => $"{x.GameID}_{x.GameName}"));
                string test = "";
            }

            return;
        }



    }
}

