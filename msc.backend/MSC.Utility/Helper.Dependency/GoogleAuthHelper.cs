using System;
using System.Collections.Generic;
using System.Drawing;
using example.invalid;
using System.Linq;
using System.Threading.Tasks;
using Google.Authenticator;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;
using MSC.Models;
using MSC.Models.Enums;
using MSC.Models.Paging;
using MSC.Models.ViewModels;
using MSC.Models.ViewModels.API.Point;
using MSC.Models.ViewModels.RequestViewModel;
using MSC.Utility;

namespace MSC.Utility.Helper.Dependency
{
    public class GoogleAuthHelper
    {

        private static string SystemName = SettingHelper.AdminSysWebsiteSetting.GoogleAuth_SystemName_EN;

        private static string GoogleAuthSaltkey = SettingHelper.AdminSysWebsiteSetting.GoogleAuth_Saltkey;


        /// <summary>
        /// 檢查輸入的驗證碼是否正確
        /// </summary>
        public static bool ValidTextIsCorrect(int userID , string text)
        {
            string systemKey = GenerateGoogleKey(userID);

            TwoFactorAuthenticator tfA = new TwoFactorAuthenticator();
            var isRight = tfA.ValidateTwoFactorPIN(systemKey, text);
            return isRight;
        }


        /// <summary>
        /// 產生QR Code 與 加密金鑰
        /// </summary>
        public static GoogleAuthQRCodeViewModel CreateSecretKeyAndQrCode(string username , int userID)
        {
            string systemKey = GenerateGoogleKey(userID);

            TwoFactorAuthenticator tfA = new TwoFactorAuthenticator();
            var setupCode = tfA.GenerateSetupCode(SystemName, username, systemKey, false, 3);

            string base64 = "";

            //1. QRCode圖片從記憶體轉到畫面上
            using (MemoryStream ms = new MemoryStream(Convert.FromBase64String(setupCode.QrCodeSetupImageUrl.Replace("data:image/png;base64,", ""))))
            {
                var image = Image.FromStream(ms);
                base64 = ImageToBase64(image, System.Drawing.Imaging.ImageFormat.Jpeg);
            }

            return new GoogleAuthQRCodeViewModel
            {
                Base64 = base64,
                ManualEntryKey = setupCode.ManualEntryKey
            };
        }

        public static string GenerateGoogleKey(int userID)
        {
            return (GoogleAuthSaltkey + "_" + userID);
        }


        private static string ImageToBase64(Image image, System.Drawing.Imaging.ImageFormat format)
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


    }


    public class GoogleAuthRequestViewModel
    {
        public string UserName { get; set; }

    }


    public class GoogleAuthQRCodeViewModel
    {
        public string Base64 { get; set; }

        public string ManualEntryKey { get; set; }

    }

}
