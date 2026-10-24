using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class GameListCreateRequestModel
    {

        public int GameID { get; set; }


        public string GameName { get; set; }


        public string GameCode { get; set; }

        /// <summary>
        /// 對應 GameTypeEnum
        /// </summary>
        public int GameType { get; set; }

        /// <summary>
        /// 對應 GameVendorTypeEnum
        /// </summary>
        public int GameVendorType { get; set; }


        public int Status { get; set; }


        public bool IsSkipInsertZeroBet { get; set; }


    }


    public class GameListUpdateRequestModel
    {

        public int GameID { get; set; }


        public string GameName { get; set; }


        public string GameCode { get; set; }

        /// <summary>
        /// 對應 GameTypeEnum
        /// </summary>
        public int GameType { get; set; }


        /// <summary>
        /// 對應 GameVendorTypeEnum
        /// </summary>
        public int GameVendorType { get; set; }


        public int Status { get; set; }


        public bool IsSkipInsertZeroBet { get; set; }



    }



}
