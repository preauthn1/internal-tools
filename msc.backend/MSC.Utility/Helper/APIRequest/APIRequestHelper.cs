using Newtonsoft.Json;
using System;
using System.Collections;
using System.Collections.Generic;
using System.Diagnostics;
using System.Dynamic;
using example.invalid;
using example.invalid.Compression;
using System.Linq;
using example.invalid;
using System.Reflection;
using System.Text;

namespace MSC.Utility.Utility
{
    public class APIRequestHelper
    {
        public static string No_Url_Keyword = "NO_URL";
              

        public static ApiRequestResultModel<T> RequestApiResult<T>(RequestParameters parameter, object postBodyModel = null)
        {
            //檢查與防呆
            VerifyParameters(parameter);

            ApiRequestLogTemporaryInfoModel requestParameterTemp = SetRequestLogInformation(parameter);

            //new Instance
            HttpWebRequest httpWebRequest = InstanceHttpWebRequest(parameter);
            httpWebRequest.Timeout = parameter.TimeoutMillisecond;
            requestParameterTemp.HttpWebRequest = httpWebRequest;

            try
            {
                switch (parameter.Method)
                {
                    case HttpMethodType.HttpDelete:
                    case HttpMethodType.HttpGet:
                        ApiRequestResultModel<T> resultModel = GetResponse<T>(httpWebRequest, requestParameterTemp);
                        return resultModel;
                }

                //HttpPOST HttpPut
                using (var streamWriter = InstanceStreamWrite(httpWebRequest, parameter.EncodingType))
                {
                    string jsonBody = string.Empty;

                    //設定 POST json body
                    if (postBodyModel == null)
                        jsonBody = SetJsonParameters(streamWriter, parameter);
                    else
                        jsonBody = SetJsonParametersByModel(streamWriter, postBodyModel);

                    requestParameterTemp.RequestContent = jsonBody;
                    ApiRequestResultModel<T> resultModel = GetResponse<T>(httpWebRequest, requestParameterTemp);
                    return resultModel;
                }
            }
            catch (WebException ex)
            {
                return CheckErrorMessage<T>(ex, requestParameterTemp);
            }
            catch (System.Exception ex)
            {
                //如果有指定Log行為
                requestParameterTemp.Response.ResponseStatusCode = 500;
                requestParameterTemp.Response.ResponseContent = ex.Message;
                ExecuteLogAction(requestParameterTemp);

                return new ApiRequestResultModel<T>
                {
                    Exception = ex
                };
            }
        }

        private static ApiRequestLogTemporaryInfoModel SetRequestLogInformation(RequestParameters parameter)
        {
            var requestParameterTemp = new ApiRequestLogTemporaryInfoModel();
            requestParameterTemp.RequestParameters = parameter;
            requestParameterTemp.RequestContent = string.Empty;
            requestParameterTemp.Timer = new Stopwatch();

            //有指定log才啟動
            if (parameter.LogAction != null)
            {
                requestParameterTemp.Timer.Start();
            }

            return requestParameterTemp;
        }

        private static StreamWriter InstanceStreamWrite(HttpWebRequest httpWebRequest, EncodingType encodeType)
        {
            if (encodeType == EncodingType.None)
                return new StreamWriter(httpWebRequest.GetRequestStream());

            return new StreamWriter(httpWebRequest.GetRequestStream(), InstanceEncoding(encodeType));
        }

        private static Encoding InstanceEncoding(EncodingType encoding)
        {
            switch (encoding)
            {
                case EncodingType.UTF8:
                    return Encoding.UTF8;

                case EncodingType.Big5:
                    return Encoding.GetEncoding("big5");

                case EncodingType.GB2312:
                    return Encoding.GetEncoding("gb2312");

                default:
                    throw new NotImplementedException();
            }
        }

        /// <summary>
        /// 加入Header
        /// </summary>
        /// <param name="httpWebRequest"></param>
        /// <param name="headers"></param>
        private static void AddHeaders(HttpWebRequest httpWebRequest, IList<KeyValuePair<string, string>> headers)
        {
            if (headers == null || headers.Count() == 0)
                return;

            foreach (var header in headers)
            {
                httpWebRequest.Headers[header.Key] = header.Value;
            }
        }

        private static List<TraceRedirectInfoDto> RedirectToNextUrl(RequestParameters requestParameter, CookieCollection cookieValues, int redirecTimes = 0)
        {
            List<TraceRedirectInfoDto> redirectTrace = new List<TraceRedirectInfoDto>();

            HttpWebRequest httpWebRequest = InstanceHttpWebRequest(requestParameter);
            httpWebRequest.AllowAutoRedirect = false;
            httpWebRequest.CookieContainer = InstanceCookieContainer(cookieValues);

            using (var httpResponse = (HttpWebResponse)httpWebRequest.GetResponse())
            {
                var debugUrl = httpWebRequest.RequestUri;
                TraceRedirectInfoDto currentTrace = new TraceRedirectInfoDto
                {
                    RedirectUrl = httpWebRequest.Address.AbsoluteUri,
                };

                using (var streamReader = new StreamReader(GetStreamForResponse(httpResponse)))
                {
                    currentTrace.ResponseText = streamReader.ReadToEnd();
                }

                AddMultiCookiesValue(currentTrace.Cookies, httpResponse.Cookies);

                //避免無窮迴圈 最多redirect 10次
                if (redirecTimes < 10 && httpResponse.StatusCode == HttpStatusCode.Redirect)
                {
                    redirecTimes++;
                    var nextRedirectUrl = ConcatFullUrl(httpWebRequest, httpResponse.Headers["location"]);

                    var cloneRequestPara = CloneRequestParameterByParent(nextRedirectUrl, HttpMethodType.HttpGet, requestParameter);

                    CookieCollection cloneCookies = ClonePreviousCookies(httpWebRequest, httpResponse);

                    var subResponseTrace = RedirectToNextUrl(cloneRequestPara, cloneCookies, redirecTimes);

                    redirectTrace.AddRange(subResponseTrace);
                }

                redirectTrace.Insert(0, currentTrace);
                return redirectTrace;
            }
        }

        private static string ConcatFullUrl(HttpWebRequest httpWebRequest, string subUrl)
        {
            if (subUrl.Contains("http"))
                return subUrl;

            string httpOrhttps = httpWebRequest.Address.AbsoluteUri.ToLower().Contains("https") ? "https://" : "http://";

            return string.Format("{0}{1}{2}", httpOrhttps, httpWebRequest.Address.Authority, subUrl);
        }

        private static void AddMultiCookiesValue(Dictionary<string, string> source, CookieCollection newAddValue)
        {
            if (newAddValue == null)
                return;

            foreach (var item in newAddValue)
            {
                var cookie = item as Cookie;

                string key = cookie.Name;

                if (source.ContainsKey(key) && string.IsNullOrWhiteSpace(source[key]) == false)
                    continue;

                if (source.ContainsKey(key))
                {
                    source[key] = cookie.Value;
                    continue;
                }

                source.Add(cookie.Name, cookie.Value);
            }
        }

        private static void AddMultiDictionaryValue(Dictionary<string, string> source, Dictionary<string, string> newAddValue)
        {
            if (newAddValue == null)
                return;

            foreach (var item in newAddValue)
            {
                string key = item.Key;

                if (source.ContainsKey(key) && string.IsNullOrWhiteSpace(source[key]) == false)
                    continue;

                if (source.ContainsKey(key))
                {
                    source[key] = item.Value;
                    continue;
                }

                source.Add(item.Key, item.Value);
            }
        }

        private static ApiRequestResultModel<T> GetResponse<T>(HttpWebRequest httpWebRequest, ApiRequestLogTemporaryInfoModel parameter)
        {
            string responseText = string.Empty;

            using (var httpResponse = (HttpWebResponse)httpWebRequest.GetResponse())
            {
                List<TraceRedirectInfoDto> redirectTrace = new List<TraceRedirectInfoDto>();

                if (httpWebRequest.AllowAutoRedirect == false && httpResponse.StatusCode == HttpStatusCode.Redirect)
                {
                    var nextRedirectUrl = ConcatFullUrl(httpWebRequest, httpResponse.Headers["location"]);

                    var cloneRequestPara = CloneRequestParameterByParent(nextRedirectUrl, HttpMethodType.HttpGet, parameter.RequestParameters);

                    CookieCollection cookieCollection = ClonePreviousCookies(httpWebRequest, httpResponse);

                    var redirectResult = RedirectToNextUrl(cloneRequestPara, cookieCollection);
                    redirectTrace = redirectResult;
                }

                using (var streamReader = new StreamReader(GetStreamForResponse(httpResponse)))
                {
                    responseText = streamReader.ReadToEnd();
                }

                if (httpResponse.Cookies != null)
                {
                    AddMultiCookiesValue(parameter.Response.ResponseCookies, httpResponse.Cookies);
                }

                //如果有指定Log行為
                SetResponseLogInformation(parameter, responseText, httpResponse);


                if (typeof(T) == typeof(string))
                {
                    return new ApiRequestResultModel<T>
                    {
                        RedirectTrace = redirectTrace,
                        ResponseCookies = parameter.Response.ResponseCookies,
                        ResultObject = (T)(object)responseText,
                        ResponseText = responseText
                    };
                }

                T convert = Newtonsoft.Json.JsonConvert.DeserializeObject<T>(responseText);
                return new ApiRequestResultModel<T>
                {
                    RedirectTrace = redirectTrace,
                    ResponseCookies = parameter.Response.ResponseCookies,
                    ResultObject = convert,
                    ResponseText = responseText
                };
            }
        }

        private static void SetResponseLogInformation(ApiRequestLogTemporaryInfoModel parameter, string responseText, HttpWebResponse httpResponse)
        {
            if (parameter.RequestParameters.LogAction == null)
                return;

            parameter.Response.ResponseStatusCode = (int)httpResponse.StatusCode;
            parameter.Response.ResponseContent = responseText;

            //執行log行為
            ExecuteLogAction(parameter);
        }

        private static CookieCollection ClonePreviousCookies(HttpWebRequest httpRequest, HttpWebResponse httpResponse)
        {
            CookieCollection resultCookieCollection = CloneCookieCollection(httpResponse.Cookies);

            var allCookies = GetAllCookies(httpRequest.CookieContainer);

            foreach (var item in allCookies)
            {
                var cookie = item as Cookie;

                if (IsCookieKeyExist(resultCookieCollection, cookie.Name))
                    continue;

                resultCookieCollection.Add(cookie);
            }

            return resultCookieCollection;
        }

        public static CookieCollection GetAllCookies(CookieContainer cookieJar)
        {
            CookieCollection cookieCollection = new CookieCollection();

            Hashtable table = (Hashtable)cookieJar.GetType().InvokeMember("m_domainTable",
                                                                            BindingFlags.NonPublic |
                                                                            BindingFlags.GetField |
                                                                            BindingFlags.Instance,
                                                                            null,
                                                                            cookieJar,
                                                                            new object[] { });

            foreach (var tableKey in table.Keys)
            {
                String str_tableKey = (string)tableKey;

                if (str_tableKey[0] == '.')
                {
                    str_tableKey = str_tableKey.Substring(1);
                }

                SortedList list = (SortedList)table[tableKey].GetType().InvokeMember("m_list",
                                                                            BindingFlags.NonPublic |
                                                                            BindingFlags.GetField |
                                                                            BindingFlags.Instance,
                                                                            null,
                                                                            table[tableKey],
                                                                            new object[] { });

                foreach (var listKey in list.Keys)
                {
                    String url = "https://" + str_tableKey + (string)listKey;
                    cookieCollection.Add(cookieJar.GetCookies(new Uri(url)));
                }
            }

            return cookieCollection;
        }

        private static CookieCollection CloneCookieCollection(CookieCollection cookieCollection)
        {
            var newCookie = new CookieCollection();

            foreach (var item in cookieCollection)
            {
                newCookie.Add(item as Cookie);
            }

            return newCookie;
        }

        private static bool IsCookieKeyExist(CookieCollection cookieCollection, string key)
        {
            foreach (var item in cookieCollection)
            {
                if ((item as Cookie).Name == key)
                {
                    return true;
                }
            }

            return false;
        }

        private static Stream GetStreamForResponse(HttpWebResponse webResponse)
        {
            if (string.IsNullOrEmpty(webResponse.ContentEncoding))
            {
                return webResponse.GetResponseStream();
            }
            Stream stream;
            switch (webResponse.ContentEncoding.ToUpperInvariant())
            {
                case "GZIP":
                    stream = new GZipStream(webResponse.GetResponseStream(), CompressionMode.Decompress);
                    break;

                case "DEFLATE":
                    stream = new DeflateStream(webResponse.GetResponseStream(), CompressionMode.Decompress);
                    break;

                default:
                    stream = webResponse.GetResponseStream();
                    break;
            }
            return stream;
        }

        private static void ExecuteLogAction(ApiRequestLogTemporaryInfoModel parameter)
        {
            if (parameter.RequestParameters.LogAction == null)
                return;

            parameter.Timer.Stop();
            parameter.Response.ResponseTime = parameter.Timer.ElapsedMilliseconds;
            parameter.RequestParameters.LogAction(parameter);
        }

        private static ApiRequestResultModel<T> CheckErrorMessage<T>(WebException ex, ApiRequestLogTemporaryInfoModel parameter)
        {
            var resultModel = new ApiRequestResultModel<T>
            {
                Exception = ex,
                IsTimeoutException = ex.Status == WebExceptionStatus.Timeout
            };

            if (ex.Response == null)
            {
                //如果有指定Log行為
                parameter.Response.ResponseStatusCode = (int)ex.Status;
                parameter.Response.ResponseContent = string.Empty;
                ExecuteLogAction(parameter);

                resultModel.ResultObject = default(T);
                return resultModel;
            }

            string errorMessage = string.Empty;
            using (var stream = ex.Response.GetResponseStream())
            using (var reader = new StreamReader(stream))
            {
                errorMessage = reader.ReadToEnd();
            }

            //如果有指定Log行為
            parameter.Response.ResponseStatusCode = (int)ex.Status;
            parameter.Response.ResponseContent = errorMessage;
            ExecuteLogAction(parameter);

            //若指定型別為string 特別處理
            if (typeof(T) == typeof(string))
            {
                resultModel.ResultObject = (T)(object)errorMessage;
                resultModel.ResponseText = errorMessage;
                return resultModel;
            }

            resultModel.ResultObject = default(T);
            resultModel.ResponseText = errorMessage;
            return resultModel;
        }

        /// <summary>
        /// 檢查與防呆
        /// </summary>
        /// <param name="parameter"></param>
        private static void VerifyParameters(RequestParameters parameter)
        {
            if (string.IsNullOrWhiteSpace(parameter.BaseUrl))
                throw new System.Exception("BaseUrl Format Not Allow #001");

            if (parameter.BaseUrl.Contains("http") == false && parameter.BaseUrl.Contains("https") == false)
                throw new System.Exception("BaseUrl Format Not Allow #002");

            if (string.IsNullOrWhiteSpace(parameter.ApiUrl))
                throw new System.Exception("ApiUrl Format Not Allow #002");

            if (parameter.TimeoutMillisecond <= 0)
                parameter.TimeoutMillisecond = 120000;
        }

        private static RequestParameters CloneRequestParameterByParent(string nextUrl, HttpMethodType type, RequestParameters parent)
        {
            List<KeyValuePair<string, string>> cloneHeaders = new List<KeyValuePair<string, string>>();

            foreach (var item in parent.Headers)
            {
                if (item.Key.Contains("Accept") == false)
                    continue;

                cloneHeaders.Add(new KeyValuePair<string, string>(item.Key, item.Value));
            }

            return new RequestParameters
            {
                BaseUrl = nextUrl,
                ApiUrl = No_Url_Keyword,
                Method = type,
                Referer = parent.Referer,
                UserAgent = parent.UserAgent,
                ContentType = parent.ContentType,
                EncodingType = parent.EncodingType,
                Accept = parent.Accept,
                Cookies = parent.Cookies ?? new List<KeyValuePair<string, string>>(),
                Headers = cloneHeaders
            };
        }

        /// <summary>
        /// 設定 POST json參數
        /// </summary>
        /// <param name="streamWriter"></param>
        /// <param name="parameters"></param>
        private static string SetJsonParameters(StreamWriter streamWriter, RequestParameters parameters)
        {
            if (parameters.Method == HttpMethodType.HttpGet)
                return string.Empty;

            string json = string.Empty;
            if (parameters.Parameters.Count == 0)
            {
                json = string.Empty;
            }
            else
            {
                object annoymouseObject = DynamicGenerateJson(parameters.Parameters);
                json = JsonConvert.SerializeObject(annoymouseObject);
            }

            streamWriter.Write(json);
            streamWriter.Flush();
            streamWriter.Close();

            return json;
        }

        /// <summary>
        /// 設定 POST json參數
        /// </summary>
        /// <param name="streamWriter"></param>
        /// <param name="parameters"></param>
        private static string SetJsonParametersByModel(StreamWriter streamWriter, object parameterModel)
        {
            string json = string.Empty;

            if (parameterModel.GetType() == typeof(string))
                json = parameterModel.ToString();
            else
                json = JsonConvert.SerializeObject(parameterModel);

            streamWriter.Write(json);
            streamWriter.Flush();
            streamWriter.Close();

            return json;
        }

        private static HttpWebRequest InstanceHttpWebRequest(RequestParameters parameter)
        {
            string url = GetRequestUrl(parameter);

            if (url.Contains("https://"))
            {
                ServicePointManager.Expect100Continue = true;
                ServicePointManager.SecurityProtocol = SecurityProtocolType.Tls12;
            }

            HttpWebRequest httpWebRequest = (HttpWebRequest)WebRequest.Create(url);

            if (string.IsNullOrWhiteSpace(parameter.ContentType) == false)
            {
                httpWebRequest.ContentType = parameter.ContentType;
            }

            if (string.IsNullOrWhiteSpace(parameter.Accept) == false)
            {
                httpWebRequest.Accept = parameter.Accept;
            }

            if (!string.IsNullOrWhiteSpace(parameter.Host))
            {
                httpWebRequest.Host = parameter.Host;
            }

            httpWebRequest.Method = ParseToMethodText(parameter.Method);
            httpWebRequest.Referer = parameter.Referer;
            httpWebRequest.UserAgent = parameter.UserAgent;

            //加入header
            AddHeaders(httpWebRequest, parameter.Headers);

            var advancedSetting = parameter.AdvancedSetting;
            if (advancedSetting != null)
            {
                httpWebRequest.AllowAutoRedirect = advancedSetting.AllowAutoRedirect;
                httpWebRequest.Proxy = advancedSetting.Proxy;
            }

            //加入cookie
            string cookieDomain = new Uri(parameter.BaseUrl).Authority;
            if (parameter.Cookies != null)
            {
                httpWebRequest.CookieContainer = new CookieContainer();
                foreach (var item in parameter.Cookies)
                {
                    var cookie = new Cookie(item.Key, item.Value);
                    cookie.Domain = cookieDomain;
                    httpWebRequest.CookieContainer.Add(cookie);
                }
            }

            return httpWebRequest;
        }

        private static CookieContainer InstanceCookieContainer(CookieCollection cookieValues)
        {
            if (cookieValues == null)
                return new CookieContainer();

            var CookieContainer = new CookieContainer();
            foreach (var item in cookieValues)
            {
                var cookie = item as Cookie;
                cookie.Expires = DateTime.UtcNow.AddDays(30);

                CookieContainer.Add(cookie);
            }

            return CookieContainer;
        }

        private static string ParseToMethodText(HttpMethodType type)
        {
            return type.ToString().Replace("Http", "").ToUpper();
        }

        /// <summary>
        /// 根據Key Value 產生匿名型別
        /// </summary>
        /// <param name="source"></param>
        /// <returns></returns>
        private static dynamic DynamicGenerateJson(IList<KeyValuePair<string, string>> source)
        {
            Dictionary<string, object> dict = source.ToDictionary(x => x.Key, x => (object)x.Value);

            var eo = new ExpandoObject();
            var eoColl = (ICollection<KeyValuePair<string, object>>)eo;

            foreach (var kvp in dict)
            {
                eoColl.Add(kvp);
            }

            dynamic eoDynamic = eoColl;
            return eoDynamic;
        }

        private static string GetRequestUrl(RequestParameters parameters)
        {
            string url = string.Format("{0}/{1}", parameters.BaseUrl, parameters.ApiUrl);

            if (parameters.ApiUrl == No_Url_Keyword)
                url = parameters.BaseUrl;

            url = ReplaceDuplicateUrlSymbol(url);

            if (parameters.Method != HttpMethodType.HttpGet)
                return url;

            if (parameters.Parameters.Count == 0)
                return url;

            //如果是HttpGet 把參數接在網址後面
            url = url + "?";
            int count = 0;
            foreach (var item in parameters.Parameters)
            {
                if (count > 0)
                    url += "&";

                url += string.Format("{0}={1}", item.Key, item.Value);
                count++;
            }

            return url;
        }

        /// <summary>
        /// 取代多餘的斜線符號
        /// </summary>
        public static string ReplaceDuplicateUrlSymbol(string source)
        {
            if (string.IsNullOrWhiteSpace(source))
                return string.Empty;

            //求高手幫忙改成 regex
            //其實就是想把 http://example.invalid////myTrip  多餘的斜線符號取代掉
            if (source.Contains("://"))
            {
                source = source.Replace("://", "TEMP_STRING");
            }

            source = source.Replace("///", "/").Replace("//", "/");
            source = source.Replace("TEMP_STRING", "://");
            return source;
        }
    }
}