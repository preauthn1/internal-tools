using System;
using System.Collections.Generic;

namespace target.Models.ViewModels
{
    public class PerGameLogFullInformationViewModel
    {

        public List<PerGameDetailViewModel> PlayerDetails { get; set; }


        public long PerGameLogID { get; set; }

        public int ServerID { get; set; }

        public string ServerName { get; set; }

        public int GameID { get; set; }

        public string GameName { get; set; }

        public string GameData { get; set; }

        public decimal TaxAmount { get; set; }


        public DateTime BeginTime { get; set; }

        public DateTime EndTime { get; set; }

    }



}
