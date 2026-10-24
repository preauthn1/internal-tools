using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Collections.Specialized;
using target.Utility.Extensions;
using example.invalid;
using example.invalid;

namespace target.Utility
{
    public static class ServerToServerUploadHelper
    {

        /// <summary>
        /// 上傳檔案至 Server 透過HttpWebRequest
        /// </summary>
        /// <param name="url">上傳網址</param>
        /// <param name="fileName">檔名</param>
        /// <param name="contentType">image type </param>
        /// <param name="nvc">其他參數</param>
        /// <returns></returns>
        public static string UploadFileToRemoteServer(string url, string fileName, string contentType, Stream myFileStram)
        {
            string paramName = "files";
            string boundary = "---------------------------" + DateTime.Now.Ticks.ToString("x");
            byte[] boundarybytes = System.Text.Encoding.ASCII.GetBytes("\r\n--" + boundary + "\r\n");

            HttpWebRequest wr = (HttpWebRequest)WebRequest.Create(url);
            wr.ContentType = "multipart/form-data; boundary=" + boundary;
            wr.Method = "POST";
            wr.KeepAlive = true;
            wr.Credentials = example.invalid.CredentialCache.DefaultCredentials;

            Stream rs = wr.GetRequestStream();
            rs.Write(boundarybytes, 0, boundarybytes.Length);

            string headerTemplate = "Content-Disposition: form-data; name=\"{0}\"; filename=\"{1}\"\r\nContent-Type: {2}\r\n\r\n";
            string header = string.Format(headerTemplate, paramName, fileName, contentType);
            byte[] headerbytes = System.Text.Encoding.UTF8.GetBytes(header);
            rs.Write(headerbytes, 0, headerbytes.Length);


            byte[] buffer = new byte[4096];
            int bytesRead = 0;
            while ((bytesRead = myFileStram.Read(buffer, 0, buffer.Length)) != 0)
            {
                rs.Write(buffer, 0, bytesRead);
            }

            byte[] trailer = System.Text.Encoding.ASCII.GetBytes("\r\n--" + boundary + "--\r\n");
            rs.Write(trailer, 0, trailer.Length);
            rs.Close();

            WebResponse wresp = null;
            var responseText = string.Empty;
            try
            {
                wresp = wr.GetResponse();
                Stream stream2 = wresp.GetResponseStream();
                StreamReader reader2 = new StreamReader(stream2);

                //成功回傳結果
                responseText = reader2.ReadToEnd();
            }
            catch (Exception ex)
            {
                if (wresp != null)
                {
                    wresp.Close();
                }

                responseText = ex.Message;
            }

            wr = null;
            return responseText;
        }







        public static string UploadFile(string URL, byte[] VideoFileData)
        {
            string Response = null;
            HttpWebRequest WebReq = null;
            HttpWebResponse WebRes = null;
            StreamReader StreamResponseReader = null;
            Stream requestStream = null;
            try
            {
                WebReq = (HttpWebRequest)WebRequest.Create(URL);
                WebReq.Method = "POST";
                WebReq.Accept = "*/*";
                WebReq.Timeout = 50000;
                WebReq.KeepAlive = false;
                WebReq.AllowAutoRedirect = false;
                WebReq.AllowWriteStreamBuffering = true;
                WebReq.ContentType = "image/jpeg";
                WebReq.ContentLength = VideoFileData.Length;


                requestStream = WebReq.GetRequestStream();
                requestStream.Write(VideoFileData, 0, VideoFileData.Length);

                requestStream.Close();

                WebRes = (HttpWebResponse)WebReq.GetResponse();
                StreamResponseReader = new StreamReader(WebRes.GetResponseStream(), Encoding.UTF8);
                Response = StreamResponseReader.ReadToEnd();
            }
            catch
            {
                throw;
            }
            finally
            {
                if (WebReq != null)
                {
                    WebReq.Abort();
                    WebReq = null;
                }
                if (WebRes != null)
                {
                    WebRes.Close();
                    WebRes = null;
                }
                if (StreamResponseReader != null)
                {
                    StreamResponseReader.Close();
                    StreamResponseReader = null;
                }
                if (requestStream != null)
                {
                    requestStream = null;
                }
            }


            return Response;
        }

    }

}
