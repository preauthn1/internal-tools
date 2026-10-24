using Microsoft.AspNetCore.Mvc;
using System.Collections.Generic;
using example.invalid;
using System;
using target.Web.Modules.Filter;
using System.Linq;
using Microsoft.AspNetCore.Http;
using System.Drawing;
using System.Drawing.Imaging;
using Azure.Storage.Blobs;
using Azure.Storage.Blobs.Models;
using example.invalid.Http;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.API.File.Helper;

namespace target.API.File.Controllers
{
    public class ImageUploadController : BaseController
	{
		
		[ApiConfig(IsLogApiEnable = false)]
		[Route("upload/images")]
		[HttpPost]
		public object UploadImageFile()
		{			
			List<string> resultFiles = new List<string>();

			foreach (var file in Request.Form.Files)
			{
				var maxSize = 2 * 1000 * 1000; //2MB
				if (file.Length > 0)
				{
					if(file.Length > maxSize)
                    {
						return new ImageUploadResultViewModel
						{
							Code = 9999,
							Message = "upload file size too large."
						};
					}

					if(!CheckIsValidFileExtension(file.FileName))
                    {
						return new ImageUploadResultViewModel
						{
							Code = 9999,
							Message = "invalid file type"
						};
                    }

					using (var fileStream = new MemoryStream())
					{
						file.CopyTo(fileStream);

						string filePath = UploadFileToRemote(fileStream);
						resultFiles.Add(filePath);
					}
				}
			}

			return new ImageUploadResultViewModel
			{
				File_Path = resultFiles.FirstOrDefault(),
				Files = resultFiles
			};
		}



		[ApiConfig(IsLogApiEnable = false)]
		[Route("upload/images/base64")]
		[HttpPost]
		public ImageUploadResultViewModel UploadImageFileByBase64([FromBody] ImageUploadRequestModel model)
		{
			model.RuleFor(x => x.Base64).Rule_NotEmptyText();

			if(model.Base64.StartsWith("data:image"))
			{
				model.Base64 = model.Base64.Replace("data:image/jpeg;base64,", "");
				model.Base64 = model.Base64.Replace("data:image/png;base64,", "");
			}

			var bytes = Convert.FromBase64String(model.Base64);

			var result = UploadFileToRemote(new MemoryStream(bytes));

			return new ImageUploadResultViewModel
			{
				File_Path = result
			};
		}


		private string UploadFileToRemote(MemoryStream fileStream)
        {
			return MinioStorageHelper.UploadFileToRemote(fileStream).Result;
		}


		private bool CheckIsValidFileExtension(string fileName)
		{
			string ext = Path.GetExtension(fileName);

			switch (ext.ToLower())
			{
				case ".gif":
					return true;
				case ".jpg":
					return true;
				case ".jpeg":
					return true;
				case ".png":
					return true;
				default:
					return false;
			}
		}


	}
}
