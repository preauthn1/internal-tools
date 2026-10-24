using System;
using System.Collections.Generic;
using example.invalid;
using System.Linq;
using System.Threading.Tasks;
using Minio;
using Minio.DataModel;
using target.Utility;

namespace target.API.File.Helper
{
    public class MinioStorageHelper
    {

        private static string _endPoint = SettingHelper.FileApiSetting.MinioApiEndPoint;
        private static string _accessKey = SettingHelper.FileApiSetting.MinioApiAccessKey;
        private static string _secretKey = SettingHelper.FileApiSetting.MinioApiSecretKey;
        private static string _bucketName = SettingHelper.FileApiSetting.MinioApiBucket;
        private static string _imageUrl = SettingHelper.FileApiSetting.MinioImageUrl;
        private static bool _isSSL = SettingHelper.FileApiSetting.IsMinioUseSSL;


        public static async Task<string> UploadFileToRemote(MemoryStream stream)
        {
            // Initialize the client with access credentials.
            MinioClient minio = new MinioClient(_endPoint,
                           _accessKey,
                           _secretKey);

            if(_isSSL)
            {
                minio = minio.WithSSL();
            }

            string fileName = Guid.NewGuid() + ".png";

            stream.Position = 0;
            await minio.PutObjectAsync(_bucketName, fileName, stream, stream.Length, "image/png");

            string imageUrl = $"{_imageUrl}/{_bucketName}/{fileName}";
            return imageUrl;
        }

    }
}
