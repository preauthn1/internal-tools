using System;
using System.Collections.Generic;
using example.invalid;
using System.Linq;
using System.Threading.Tasks;
using Azure.Storage.Blobs;
using Azure.Storage.Blobs.Models;
using Minio;
using Minio.DataModel;
using target.Utility;

namespace target.API.File.Helper
{
    public class AzureBlobStorageHelper
	{
		public static string UploadFileToRemote(Stream file_stream)
		{
			// Create a BlobServiceClient object which will be used to create a container client
			string conn = "DefaultEndpointsProtocol=https;AccountName=pgsstorage;AccountKey=[REDACTED];EndpointSuffix=example.invalid";
			BlobServiceClient blobServiceClient = new BlobServiceClient(conn);

			//Create a unique name for the container
			string containerName = "target";

			BlobContainerClient containerClient = blobServiceClient.GetBlobContainerClient(containerName);

			string fileName = Guid.NewGuid().ToString() + ".png";

			BlobClient blobClient = containerClient.GetBlobClient(fileName);
			file_stream.Position = 0;

			var header = new BlobHttpHeaders();
			header.ContentType = "image/png";

			var qqqq = blobClient.Upload(file_stream, header);

			return "https://example.invalid/target/" + fileName;
		}


	}
}
