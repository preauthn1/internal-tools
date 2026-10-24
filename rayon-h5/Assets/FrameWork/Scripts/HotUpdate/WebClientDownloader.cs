using System;
using System.Collections.Generic;
using System.ComponentModel;
using example.invalid;
using System.Threading;
using UnityEngine;

namespace HotUpdate
{
    public class Calculagraph
    {
        /// <summary>
        /// 时间到事件
        /// </summary>
        public event TimeoutCaller TimeOver;

        /// <summary>
        /// 开始时间
        /// </summary>
        private DateTime _startTime;
        private TimeSpan _timeout = new TimeSpan(0, 0, 10);
        private bool _hasStarted = false;
        object _userdata;

        /// <summary>
        /// 计时器构造方法
        /// </summary>
        /// <param name="userdata">计时结束时回调的用户数据</param>
        public Calculagraph(object userdata)
        {
            TimeOver += new TimeoutCaller(OnTimeOver);
            _userdata = userdata;
        }

        /// <summary>
        /// 超时退出
        /// </summary>
        /// <param name="userdata"></param>
        public virtual void OnTimeOver(object userdata)
        {
            Stop();
        }

        /// <summary>
        /// 过期时间(秒)
        /// </summary>
        public int Timeout
        {
            get
            {
                return _timeout.Seconds;
            }
            set
            {
                if (value <= 0)
                    return;
                _timeout = new TimeSpan(0, 0, value);
            }
        }

        /// <summary>
        /// 是否已经开始计时
        /// </summary>
        public bool HasStarted
        {
            get
            {
                return _hasStarted;
            }
        }

        /// <summary>
        /// 开始计时
        /// </summary>
        public void Start()
        {
            Reset();
            _hasStarted = true;
            Thread th = new Thread(WaitCall);
            th.IsBackground = true;
            th.Start();
        }

        /// <summary>
        /// 重置
        /// </summary>
        public void Reset()
        {
            _startTime = DateTime.Now;
        }

        /// <summary>
        /// 停止计时
        /// </summary>
        public void Stop()
        {
            _hasStarted = false;
        }

        /// <summary>
        /// 检查是否过期
        /// </summary>
        /// <returns></returns>
        private bool checkTimeout()
        {
            return (DateTime.Now - _startTime).Seconds >= Timeout;
        }

        private void WaitCall()
        {
            try
            {
                //循环检测是否过期
                while (_hasStarted && !checkTimeout())
                {
                    Thread.Sleep(1000);
                }
                if (TimeOver != null)
                    TimeOver(_userdata);
            }
            catch (Exception)
            {
                Stop();
            }
        }
    }

    /// <summary>
    /// 过期时回调委托
    /// </summary>
    /// <param name="userdata"></param>
    public delegate void TimeoutCaller(object userdata);
     
    /// <summary>
    /// WebClient 超时设置
    /// </summary>
    public class WebDownload : WebClient
    {
        private static bool s_ShowFullLog = true;

        private Calculagraph _timer;
        private int _timeOut = 10;

        /// <summary>
        /// 过期时间
        /// </summary>
        public int Timeout
        {
            get
            {
                return _timeOut;
            }
            set
            {
                if (value <= 0)
                    _timeOut = 10;
                _timeOut = value;
            }
        }

        /// <summary>
        /// 重写GetWebRequest,添加WebRequest对象超时时间
        /// </summary>
        /// <param name="address"></param>
        /// <returns></returns>
        protected override WebRequest GetWebRequest(Uri address)
        {
            HttpWebRequest request = (HttpWebRequest)base.GetWebRequest(address);
            request.Timeout = 1000 * Timeout;
            request.ReadWriteTimeout = 1000 * Timeout;
            return request;
        }

        /// <summary>
        /// 带过期计时的下载
        /// </summary>
        public void DownloadFileAsyncWithTimeout(Uri address, string fileName)
        {
            if (_timer == null)
            {
                _timer = new Calculagraph(this);
                _timer.Timeout = Timeout;
                _timer.TimeOver += new TimeoutCaller(_timer_TimeOver);
                this.DownloadProgressChanged += new DownloadProgressChangedEventHandler(WebClient_DownloadProgressChanged);
            }

            // Add s_ShowFullLog By Rayon
            if (s_ShowFullLog == true)
                Debug.Log(address.AbsoluteUri);

            DownloadFileAsync(address, fileName);
            _timer.Start();
        }

        /// <summary>
        /// WebClient下载过程事件，接收到数据时引发
        /// </summary>
        /// <param name="sender"></param>
        /// <param name="e"></param>
        void WebClient_DownloadProgressChanged(object sender, DownloadProgressChangedEventArgs e)
        {
            _timer.Reset();//重置计时器
        }

        /// <summary>
        /// 计时器过期
        /// </summary>
        /// <param name="userdata"></param>
        void _timer_TimeOver(object userdata)
        {
            this.CancelAsync();//取消下载
        }
    }

    /// <summary>
    /// 下载文件缓存信息
    /// </summary>
    class DownloadCacheInfo
    {
        /// <summary>
        /// 初始化下载信息
        /// </summary>
        /// <param name="key"></param>
        public DownloadCacheInfo(string key, string fileName)
        {
            this.Key = key;
            this.FileName = fileName;
        }

        public string Key { get; private set; }

        /// <summary>
        /// 文件名称
        /// </summary>
        public string FileName { get; set; }

        /// <summary>
        /// 远程路径字符串
        /// </summary>
        public string RemoteFileFullName
        {
            get
            {
                return FileUtility.GetAssetFilePath(FileName, PathType.Remote);
            }
        }

        /// <summary>
        /// 本地Cache 路径字符串
        /// </summary>
        public string CacheFileFullName
        {
            get
            {
                string cacheFullName = FileUtility.GetAssetFilePath(FileName, PathType.Cache);
                FileUtility.CreateDirectory(cacheFullName);
                return cacheFullName;
            }
        }

        /// <summary>
        /// 完成通知
        /// </summary>
        public Action<string, AsyncCompletedEventArgs> ComplatedCallBack { get; set; }

        /// <summary>
        /// 进度更新事件
        /// </summary>
        public Action<string, DownloadProgressChangedEventArgs> ProgressChanged { get; set; }
    }


    public class WebClientDownloader
    {
        /// <summary>
        /// 等待下载的队列
        /// </summary>
        static List<DownloadCacheInfo> m_WaitingDownloadList = new List<DownloadCacheInfo>();

        /// <summary>
        /// 锁的对象
        /// </summary>
        static object m_LockObject = new object();

        /// <summary>
        /// 当前下载的文件信息
        /// </summary>
        static DownloadCacheInfo m_CurrentDownload = null;

        /// <summary>
        /// WebClient 下载器
        /// </summary>
        static WebDownload m_Downloader = null;

        /// <summary>
        /// 下载文件完成通知
        /// </summary>
        /// <param name="sender"></param>
        /// <param name="e"></param>
        static void DownloadFile_Complated(object sender, AsyncCompletedEventArgs e)
        {
            try
            {
                if (m_CurrentDownload != null && m_CurrentDownload.ComplatedCallBack != null)
                {
                    m_CurrentDownload.ComplatedCallBack(m_CurrentDownload.Key, e);
                }
            }
            catch (Exception ex)
            {
                Debug.LogException(ex);
            }
            finally
            {
                m_CurrentDownload = null;
                TryStartDownloadOneFile();
            }
        }

        /// <summary>
        /// 下载文件进度更新
        /// </summary>
        /// <param name="sender"></param>
        /// <param name="e"></param>
        static void DownloadFile_ProgressChanged(object sender, DownloadProgressChangedEventArgs e)
        {
            if (m_CurrentDownload != null && m_CurrentDownload.ProgressChanged != null)
            {
                try
                {
                    m_CurrentDownload.ProgressChanged(m_CurrentDownload.Key, e);
                }
                catch (Exception ex)
                {
                    Debug.LogException(ex);
                }
            }
        }

        static void TryStartDownloadOneFile()
        {
            if (m_CurrentDownload == null)
            {
                if (m_WaitingDownloadList.Count > 0)
                {
                    lock (m_LockObject)
                    {
                        m_CurrentDownload = m_WaitingDownloadList[0];
                        m_WaitingDownloadList.RemoveAt(0);
                    }
                    if (m_Downloader == null)
                    {
                        m_Downloader = new WebDownload();

                        m_Downloader.Timeout = 10;
                        m_Downloader.DownloadFileCompleted += DownloadFile_Complated;
                        m_Downloader.DownloadProgressChanged += DownloadFile_ProgressChanged;
                    }

                    m_Downloader.DownloadFileAsyncWithTimeout(new Uri(m_CurrentDownload.RemoteFileFullName), m_CurrentDownload.CacheFileFullName);
                }
                else
                {
                    ReleaseDownLoader();
                }
            }
        }

        /// <summary>
        /// 添加到下载队列中
        /// </summary>
        /// <param name="fileName">文件名称</param>
        /// <param name="complated">完成通知</param>
        /// <param name="progressChanged">进度更新</param>
        /// <returns></returns>
        public static void AppendDownloadFile(string fileName, Action<string, AsyncCompletedEventArgs> complated = null, Action<string, DownloadProgressChangedEventArgs> progressChanged = null)
        {
            DownloadCacheInfo cacheInfo = null;
            if (m_CurrentDownload != null && m_CurrentDownload.Key == fileName)
            {
                cacheInfo = m_CurrentDownload;
            }
            else
            {
                m_WaitingDownloadList.Find(temp => temp.Key == fileName);
            }

            if (cacheInfo == null)
            {
                cacheInfo = new DownloadCacheInfo(fileName, fileName);
                m_WaitingDownloadList.Add(cacheInfo);
            }
            if (complated != null)
            {
                if (cacheInfo.ComplatedCallBack != null)
                {
                    cacheInfo.ComplatedCallBack += complated;
                }
                else
                {
                    cacheInfo.ComplatedCallBack = complated;
                }
            }

            if (progressChanged != null)
            {
                if (cacheInfo.ProgressChanged != null)
                {
                    cacheInfo.ProgressChanged += progressChanged;
                }
                else
                {
                    cacheInfo.ProgressChanged = progressChanged;
                }
            }

            TryStartDownloadOneFile();
        }

        /// <summary>
        /// 释放下载器
        /// </summary>
        static void ReleaseDownLoader()
        {
            if (m_Downloader != null)
            {
                m_Downloader.Dispose();
                m_Downloader = null;
            }
        }

        /// <summary>
        /// 重置下载器
        /// </summary>
        public static void ResetDownLoader()
        {
            m_WaitingDownloadList.Clear();
            ReleaseDownLoader();
        }

        /// <summary>
        /// 停止下载
        /// </summary>
        public static void StopDownloader()
        {
            lock (m_LockObject)
            {
                for (int i = 0; i < m_WaitingDownloadList.Count; i++)
                {
                    DownloadCacheInfo cacheInfo = m_WaitingDownloadList[i];
                    if (cacheInfo.ComplatedCallBack != null)
                    {
                        try
                        {
                            m_CurrentDownload.ComplatedCallBack(m_CurrentDownload.Key, new AsyncCompletedEventArgs(null, true, null));
                        }
                        catch (Exception ex)
                        {
                            Debug.LogException(ex);
                        }
                    }
                }
                m_WaitingDownloadList.Clear();
            }
            if (m_Downloader != null)
            {
                m_Downloader.CancelAsync();
            }
        }

        public static void StopDownloadFile(string fileName)
        {
            Debug.LogError(m_CurrentDownload.FileName);
            if (m_CurrentDownload != null && m_CurrentDownload.Key == fileName)
            {
                m_Downloader.CancelAsync();
            }
            else
            {
                lock (m_LockObject)
                {
                    DownloadCacheInfo cacheInfo = m_WaitingDownloadList.Find((temp) => fileName == temp.Key);
                    m_WaitingDownloadList.Remove(cacheInfo);
                    if (cacheInfo.ComplatedCallBack != null)
                    {
                        try
                        {
                            cacheInfo.ComplatedCallBack(cacheInfo.Key, new AsyncCompletedEventArgs(null, true, null));
                        }
                        catch (Exception ex)
                        {
                            Debug.LogException(ex);
                        }
                    }
                }
            }
        }
    }
}