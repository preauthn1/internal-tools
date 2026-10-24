using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class GamePayTableImageViewModel
    {
        public int Id { get; set; }
        public int GameID { get; set; }
        public string GameName { get; set; }
        public string Code { get; set; }

        public string Memo { get; set; }
        public string ImageUrl { get; set; }
        public DateTime CreateTime { get; set; }

    }


    public class GamePayTableImageUpdateRequestModel
    {
        public int Id { get; set; }
        public int GameID { get; set; }
        public string Code { get; set; }
        public string ImageUrl { get; set; }
        public string Memo { get; set; }
    }


}
