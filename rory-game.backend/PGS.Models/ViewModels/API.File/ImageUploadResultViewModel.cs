using Newtonsoft.Json;
using System.Collections.Generic;

namespace target.Models.ViewModels
{
    public class ImageUploadResultViewModel
    {

        [JsonProperty("code")]
        public int Code { get; set; }

        [JsonProperty("message")]
        public string Message { get; set; }

        [JsonProperty("file_path")]
        public string File_Path { get; set; }

        [JsonProperty("files")]
        public List<string> Files { get; set; }

    }


}
