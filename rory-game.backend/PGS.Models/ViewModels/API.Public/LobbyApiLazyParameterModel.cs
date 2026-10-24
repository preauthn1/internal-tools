using System.Collections.Generic;

namespace target.Models.ViewModels
{
    public class LobbyApiGenerateLinkParameterModel
    {
        public string SignKey { get; set; }

        public long Ts { get; set; }

        public string PlayerName { get; set; }

        public int GameID { get; set; }

    }

    public class LobbyApiGetGamesRequestModel
    {
        public string SignKey { get; set; }

        public long Ts { get; set; }

        public string PlayerName { get; set; }

    }

    public class LobbyApiGetGamesInfoViewModel
    {
        public List<LobbyGamesViewModel> Games { get; set; }

        public string PlayerName { get; set; }

        public decimal Balance { get; set; }

    }




    public class LobbyGenerateLinkResultModel
    {
        public string Url { get; set; }

    }



}
