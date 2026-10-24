using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using Microsoft.AspNetCore.Cors;
using target.Web.Modules.Filter;
using target.BL.Services;
using target.Models.Enums;
using target.Models.ViewModels.GameVendor;
using target.Models;
using target.Utility.Helper;
using target.BL.Services;
using System.Threading;
using target.DA.Repository;
using System.Drawing;
using example.invalid;
using System.Drawing.Imaging;

namespace target.API.File.Controllers
{
    public class ImageConvertController : BaseController
    {


		[ApiConfig(ApiCategory = ApiConfigCategoryEnum.DebugApi)]
		[Route("image/convert-to-base64/remote-url")]
		[HttpGet, HttpPost]
		public object ConvertToBase64StringFromRemoteImageUrl(string url = "")
        {
			try
            {
				string base64Result = string.Empty;

				using (WebClient client = new WebClient())
				{
					byte[] data = client.DownloadData(url);

					using (MemoryStream mem = new MemoryStream(data))
					{
						using (var myImage = Image.FromStream(mem))
						{
							base64Result = ImageToBase64(myImage, ImageFormat.Png);
						}
					}
				}

				return base64Result;
			}
			catch(Exception ex)
            {
				return ex.ToString();
            }
        }



		private string ImageToBase64(Image image, ImageFormat format)
		{
			using (MemoryStream ms = new MemoryStream())
			{
				// Convert Image to byte[]
				image.Save(ms, format);
				byte[] imageBytes = ms.ToArray();

				// Convert byte[] to base 64 string
				string base64String = Convert.ToBase64String(imageBytes);
				return base64String;
			}
		}


		/// <summary>
		/// base 64字串格式的圖片轉成Image物件
		/// 傳入的字串 請勿包含 data:image/png;base64,
		/// </summary>
		private Image Base64StringToImage(string base64String)
		{
			// Convert base 64 string to byte[]
			byte[] Buffer = Convert.FromBase64String(base64String);

			byte[] data = null;
			Image oImage = null;
			MemoryStream oMemoryStream = null;
			Bitmap oBitmap = null;
			//建立副本
			data = (byte[])Buffer.Clone();
			try
			{
				oMemoryStream = new MemoryStream(data);
				//設定資料流位置
				oMemoryStream.Position = 0;
				oImage = System.Drawing.Image.FromStream(oMemoryStream);
				//建立副本
				oBitmap = new Bitmap(oImage);
			}
			catch
			{
				throw;
			}
			finally
			{
				oMemoryStream.Close();
				oMemoryStream.Dispose();
				oMemoryStream = null;
			}
			//return oImage;
			return oBitmap;
		}


	}
}
