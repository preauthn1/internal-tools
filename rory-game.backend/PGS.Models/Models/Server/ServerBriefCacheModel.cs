using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class ServerBriefCacheModel
    {
        public int ServerID { get; set; }

        public string ServerName { get; set; }

        public int FlagLevel { get; set; }

        public int Status { get; set; }

        public int GameID { get; set; }

        public int GameType { get; set; }

        public string GameCode { get; set; }

    }
}
