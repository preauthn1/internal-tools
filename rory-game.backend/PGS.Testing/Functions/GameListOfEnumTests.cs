using Microsoft.VisualStudio.TestTools.UnitTesting;
using target.Models.Attributes;
using target.Models.Enums;
using target.Models.Models;
using target.Utility.Extensions;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class GameListOfEnumTests : BaseTest
    {

        /// <summary>
        /// 印出哪些遊戲設置為 IgnoreZeroBet
        /// </summary>
        [TestMethod]
        public void GetIgnoreBetZeroGames()
        {
            var allGames = ExtensionOfEnum.GetEnumAllValue<GameListEnum>();
            
            var enableIgnoreGames = allGames.Where(x => x.GetAttributeFromEnum<IgnoreZeroBet>()?.IsIgnore == true).ToList();

            //印出哪些遊戲屬於這種定義
            var gameIDs = string.Join("、", enableIgnoreGames.Select(z => z.ToInteger()));

            //印出筆數
            var totalCount = enableIgnoreGames.Count();

            return;
        }



    }
}
