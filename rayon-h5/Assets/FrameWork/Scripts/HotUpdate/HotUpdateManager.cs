using System;
using System.Collections;
using System.Collections.Generic;
using System.ComponentModel;
using example.invalid;
using example.invalid;
using System.Security;
using UnityEngine;

namespace HotUpdate
{
    public enum HotUpdateState
    {
        /// <summary>
        /// 无
        /// </summary>
        None = 0,
        /// <summary>
        /// 载入本地版本文件
        /// </summary>
        LoadLocalVersionFile,
        /// <summary>
        /// 载入本地资源清单文件
        /// </summary>
        LoadLocalFileRecords,
        /// <summary>
        /// 载入远程资源清单资源包
        /// </summary>
        LoadRemoteVersionFile,
        /// <summary>
        /// 载入远程资源清单文件
        /// </summary>
        LoadRemoteFileRecords,
        /// <summary>
        /// 下载资源状态
        /// </summary>
        DownloadData,
        /// <summary>
        /// 清理过期资源
        /// </summary>
        ClearObsoleteAssets,
        /// <summary>
        /// 校验本地数据
        /// </summary>
        ValidateLocalData,
        /// <summary>
        /// 完成
        /// </summary>
        Done,
        /// <summary>
        /// 失败
        /// </summary>
        Failed,
    }

    /// <summary>
    /// 进度更新事件的参数   
    /// 解压只用到ProgressPercentage，TipsContent字段
    /// 更新用到所有字段
    /// </summary>
    public class UpdateProgressInfo
    {
        public UpdateProgressInfo()
        {
            ProgressPercentage = 0;
            TipsContent = string.Empty;
            DownLoadSizeInfo = string.Empty;
            DownLoadSpeed = "0 KB/S";
            NeedTime = "正在计算...";
        }

        /// <summary>
        /// 进度百分比
        /// </summary>
        public float ProgressPercentage { get; set; }

        /// <summary>
        /// 提示的信息
        /// </summary>
        public string TipsContent { get; set; }

        /// <summary>
        /// 下载文件大小  0.00MB/50.00MB
        /// </summary>
        public string DownLoadSizeInfo { get; set; }

        /// <summary>
        /// 每秒下载速度  200 KB/S
        /// </summary>
        public string DownLoadSpeed { get; set; }

        /// <summary>
        /// 下载需要时间  0:02:59    
        /// </summary>
        public string NeedTime { get; set; }
    }

    public class HotUpdateManager : Singleton<HotUpdateManager>
    {
        private static bool s_ShowFullLog = true;

        private static string UpdateProjectName = "";      // 当前更新模块操作的项目  大厅 或者 游戏
        public void SetUpdateProjectName(string projectName)
        {
            m_UpdateFaildRetryTimes = 1;        // 重置重试次数
            UpdateProjectName = projectName;
        }

        public const string DEFAULT_VERSION = "0.0.0";
        /// <summary>
        /// 进度更新事件
        /// </summary>
        public static Action<UpdateProgressInfo> OnHotUpdateProgressChanged = null;
        public static Action OnHotUpdateComplated = () => { };
        public static Action OnHotUpdateFaild = null;

        public string GetProjectVersion(string projectName)
        {
            string versionNo;

            versionNo = FileUtility.ReadFileText(projectName + "_version.txt", PathType.Local);

            if (string.IsNullOrEmpty(versionNo))
                return DEFAULT_VERSION;

            return versionNo;
        }

        /// <summary>
        /// 是否工作中
        /// </summary>
        public bool IsWorking { get; private set; }

        /// <summary>
        /// 是否载入了新的资源文件
        /// </summary>
        public bool IsLoadedNewAsset { get; private set; }

        private HotUpdateState m_UpdaterState = HotUpdateState.None;
        /// <summary>
        /// 是否是验证过的资源
        /// </summary>
        public bool IsValidatedAssets { get; private set; }

        /// <summary>
        /// 能否重试的更新错误
        /// </summary>
        private bool m_CanRetryUpdateFaild = true;

        /// <summary>
        /// 更新器的状态
        /// </summary>
        public HotUpdateState UpdaterState
        {
            get { return m_UpdaterState; }
            private set
            {
                if (m_UpdaterState != value)
                {
                    Debug.LogWarningFormat("HotUpdater state:[{0}] ====>[{1}]", m_UpdaterState, value);
                }
                m_UpdaterState = value;
                Debug.LogWarning("HotUpdater State " + UpdaterState);
            }
        }

        /// <summary>
        /// 远程版本号
        /// </summary>
        private string m_LocalVersion = string.Empty;
        /// <summary>
        /// 本地版本号
        /// </summary>
        private string m_RemoteVersion = string.Empty;
        /// <summary>
        /// 本地安装包版本号
        /// </summary>
        private string m_InitVersion = string.Empty;

        /// <summary>
        /// 版本文件名称
        /// </summary>
        private static string VERSION_FILENAME
        {
            get
            {
                return UpdateProjectName + "_version.txt";
            }
        }

        public static string HallVersionFileName
        {
            get { return AppDefine.HallNumber + "_version.txt"; }
        }

        /// <summary>
        /// 文件列表文件名称
        /// </summary>
        private static string FILERECORD_FILENAME
        {
            get
            {
                return UpdateProjectName + "_FileList.xml";
            }
        }

        public static string GetFormatedVersionFileName(string projectName)
        {
            return string.Format("{0}_version.txt", projectName);
        }

        public static string GetFormatedSkinFileName(string projectName)
        {
            return string.Format("{0}_skin.txt", projectName);
        }

        public static string GetFormatedFileRecordFileName(string projectName)
        {
            return string.Format("{0}_FileList.xml", projectName);
        }

        /// <summary>
        /// 本地文件记录
        /// </summary>
        private Dictionary<string, FileRecord> m_LocalFileRecords = new Dictionary<string, FileRecord>();

        /// <summary>
        /// 远程文件记录
        /// </summary>
        private Dictionary<string, FileRecord> m_RemoteFileRecords = new Dictionary<string, FileRecord>();

        /// <summary>
        /// 需要更新的文件记录
        /// </summary>
        private List<FileRecord> m_NeedUpdateFileRecords = new List<FileRecord>();

        /// <summary>
        /// 下一个进入的状态
        /// </summary>
        HotUpdateState m_NextUpdateState = HotUpdateState.None;

        void HotUpdaterSwitchToState(HotUpdateState state)
        {
            m_NextUpdateState = state;
            Debug.LogWarning("HotUpdaterSwitchToState " + m_NextUpdateState);
        }

        void HotUpdaterSwitchToFaildState(bool canRetryUpdateFaild = true)
        {
            this.m_CanRetryUpdateFaild = canRetryUpdateFaild;
            HotUpdaterSwitchToState(HotUpdateState.Failed);
        }

        bool IsRefreshUpdateProgress;       // 刷新下载进度

        void Update()
        {
            if (IsRefreshUpdateProgress)
            {
                IsRefreshUpdateProgress = false;
                RefreshDownloadProgress();
            }

            if (!IsWorking || m_NextUpdateState == HotUpdateState.None || m_NextUpdateState == UpdaterState)
            {
                return;
            }

            UpdaterState = m_NextUpdateState;
            m_NextUpdateState = HotUpdateState.None;

            switch (UpdaterState)
            {
                case HotUpdateState.None:
                    break;
                case HotUpdateState.LoadLocalVersionFile:
                    updateProgressInfo.TipsContent = "校验版本信息";
                    ChangeHotUpdateProgress();
                    StartLoadLocalVersionFile();
                    break;
                case HotUpdateState.LoadLocalFileRecords:
                    updateProgressInfo.TipsContent = "校验版本信息";
                    ChangeHotUpdateProgress();
                    StartLoadLocalRecordFile();
                    break;
                case HotUpdateState.LoadRemoteVersionFile:
                    updateProgressInfo.TipsContent = "校验版本信息";
                    ChangeHotUpdateProgress();
                    StartLoadRemoteVersionFile();
                    break;
                case HotUpdateState.LoadRemoteFileRecords:
                    updateProgressInfo.TipsContent = "校验版本信息";
                    ChangeHotUpdateProgress();
                    StartLoadRemoteRecordFile();
                    break;
                case HotUpdateState.DownloadData:
                    StartDownLoadRemoteUpdatedAssets();
                    break;
                case HotUpdateState.ClearObsoleteAssets:
                    StartClearObsoleteAssetOfDownload();
                    break;
                case HotUpdateState.ValidateLocalData:
                    StartValidateLocalAssetFiles();
                    break;
                case HotUpdateState.Done:
                    updateProgressInfo.TipsContent = "校验版本信息";
                    ChangeHotUpdateProgress();
                    HotUpdateComplated();
                    break;
                case HotUpdateState.Failed:
                    HotUpdateFaild();
                    break;
                default:
                    break;
            }
        }

        void RestoreLocalFileRecords(Dictionary<string, FileRecord> fileRecords)
        {
            Debug.LogWarning("RestoreLocalFileRecords");
            FileUtility.StoreFileRecords(FileUtility.GetAssetFilePath(FILERECORD_FILENAME, PathType.Local), fileRecords);
        }

        void OnApplicationQuit()
        {
            StopAllCoroutines();
            WebClientDownloader.StopDownloader();
        }

        #region  First StartUp

        /// <summary>
        /// 解压大厅以及游戏资源
        /// </summary>
        public IEnumerator TryUpdateAllAssetFromStreamingAssets()
        {
            yield return FileUtility.CopyStreamingAssetsToFile(FileUtility.GetAssetFilePath(AppDefine.ProjectListFileName, PathType.InitData), FileUtility.GetAssetFilePath(AppDefine.ProjectListFileName, PathType.Cache));
            string projectList = FileUtility.ReadFileText(AppDefine.ProjectListFileName, PathType.Cache);
            string[] projectArray = string.IsNullOrEmpty(projectList) ? new string[] { AppDefine.HallNumber } : projectList.Split(',');

            for (int i = 0; i < projectArray.Length; i++)
            {
                UpdateProjectName = projectArray[i];
                yield return TryUpdateAssetFromStreamingAssets(1f / projectArray.Length, i * 1f / projectArray.Length);
            }
        }

        /// <summary>
        /// 尝试从初始化目录更新资源
        /// </summary>        
        /// <param name="percentage">当前解压对象占总解压进度的百分比</param>
        /// <returns></returns> 
        IEnumerator TryUpdateAssetFromStreamingAssets(float percentage, float processed)
        {
            // 拷贝包里的版本文件到本地
            yield return FileUtility.CopyStreamingAssetsToFile(FileUtility.GetAssetFilePath(HotUpdateManager.VERSION_FILENAME, PathType.InitData), FileUtility.GetAssetFilePath(HotUpdateManager.VERSION_FILENAME, PathType.Cache));
            m_InitVersion = FileUtility.ReadFileText(HotUpdateManager.VERSION_FILENAME, PathType.Cache);
            if (FileUtility.ExistsFile(FileUtility.GetAssetFilePath(VERSION_FILENAME, PathType.Local)))
            {
                m_LocalVersion = FileUtility.ReadFileText(VERSION_FILENAME, PathType.Local);
            }
            else
            {
                m_LocalVersion = DEFAULT_VERSION;
            }

            int checkReslut = CheckVersion(m_InitVersion, m_LocalVersion);
            // 比较本地版本号和本地安装包的版本号            
            if (checkReslut == 1 || checkReslut == 2)
            {
                // 拷贝包里的资源记录文件到本地
                yield return FileUtility.CopyStreamingAssetsToFile(FileUtility.GetAssetFilePath(HotUpdateManager.FILERECORD_FILENAME, PathType.InitData), FileUtility.GetAssetFilePath(HotUpdateManager.FILERECORD_FILENAME, PathType.Cache));
                string initRecordFileName = FileUtility.GetAssetFilePath(HotUpdateManager.FILERECORD_FILENAME, PathType.Cache);
                if (!FileUtility.ExistsFile(initRecordFileName))
                {
                    Application.Quit();
                }
                else
                {
                    string initRecordContent = FileUtility.ReadFileText(HotUpdateManager.FILERECORD_FILENAME, PathType.Cache);
                    var initFileRecords = ConvertToFileRecords(initRecordContent);
                    string localRecordContent = FileUtility.ReadFileText(FILERECORD_FILENAME, PathType.Local);
                    Dictionary<string, FileRecord> localFileRecords = new Dictionary<string, FileRecord>();
                    if (localRecordContent != null)
                    {
                        localFileRecords = ConvertToFileRecords(localRecordContent);
                    }
                    // 筛选出对应需要更新的文件记录信息
                    var needUpdateFileRecords = FilterNeedUpdateFileRecodes(initFileRecords, localFileRecords);
                    if (needUpdateFileRecords.Count > 0)
                    {
                        ChangeDecompressProgress(processed);
                        for (int index = 0; index < needUpdateFileRecords.Count; index++)
                        {
                            FileRecord fileRecord = needUpdateFileRecords[index];
                            string fromFileName = FileUtility.GetAssetFilePath(fileRecord.FileName, PathType.InitData);
                            string toFileName = FileUtility.GetAssetFilePath(fileRecord.FileName, PathType.Local);
                            yield return StartCoroutine(FileUtility.CopyStreamingAssetsToFile(fromFileName, toFileName));
                            if (FileUtility.ExistsFile(toFileName))
                            {
                                localFileRecords[needUpdateFileRecords[index].FileName] = needUpdateFileRecords[index];
                            }
                            else
                            {
                                localFileRecords[needUpdateFileRecords[index].FileName] = needUpdateFileRecords[index];
                                localFileRecords[needUpdateFileRecords[index].FileName].MD5 = "";
                            }

                            float progressVal = (index + 1) * 1f / needUpdateFileRecords.Count * percentage + processed;   // 进度值
                            ChangeDecompressProgress(progressVal);

                            if (Math.Abs(progressVal - 1f) <= 0.0001f)
                                yield return new WaitForSeconds(0.3f);
                        }
                    }

                    Debug.Log("TryUpdateAssetFromStreamingAssets");

                    RestoreLocalFileRecords(localFileRecords);
                    // 拷贝初始目录的版本文件到本地
                    yield return StartCoroutine(FileUtility.CopyStreamingAssetsToFile(FileUtility.GetAssetFilePath(VERSION_FILENAME, PathType.InitData), FileUtility.GetAssetFilePath(VERSION_FILENAME, PathType.Local)));
                    // 清理过期资源                    
                    yield return StartCoroutine(ClearObsoleteAssetData(initFileRecords, localFileRecords));
                }
            }
        }

        void ChangeDecompressProgress(float progressPercentage)
        {
            if (OnHotUpdateProgressChanged != null)
            {
                OnHotUpdateProgressChanged(new UpdateProgressInfo()
                {
                    ProgressPercentage = progressPercentage,
                    TipsContent = string.Format("首次加载耗时较长，请耐心等待({0}%)", (progressPercentage * 100).ToString("F2"))
                });
            }
        }

        #endregion

        #region Update Remote Assets

        /// <summary>
        /// 尝试从远程服务器更新资源
        /// </summary>
        /// <returns></returns>
        public void TryStartUpdateRemoteAssetsToLocal()
        {
            IsWorking = true;
            IsLoadedNewAsset = false;
            m_NextUpdateState = HotUpdateState.None;
            UpdaterState = HotUpdateState.None;

            m_LocalFileRecords.Clear();
            m_RemoteFileRecords.Clear();
            m_NeedUpdateFileRecords.Clear();

            m_NeedUpdateSize = 0;      // 需要下载数据大小
            m_UpdatedFileSize = 0;     // 已经下载文件数据大小
            m_LastDownLoadedFileSize = 0;
            updateProgressInfo = new UpdateProgressInfo();

            m_RemoteVersion = m_LocalVersion = DEFAULT_VERSION;
            IsValidatedAssets = false;
            HotUpdaterSwitchToState(HotUpdateState.LoadLocalVersionFile);
        }

        #region Load local version file

        /// <summary>
        /// 开始加载本地版本记录文件
        /// </summary>
        void StartLoadLocalVersionFile()
        {
            Debug.Log("资源更新：加载本地版本文件");
            m_LocalVersion = FileUtility.ReadFileText(VERSION_FILENAME, PathType.Local);
            HotUpdaterSwitchToState(HotUpdateState.LoadLocalFileRecords);
        }

        #endregion

        #region Load local record file

        /// <summary>
        /// 开始加载本地文件记录文件
        /// </summary>
        void StartLoadLocalRecordFile()
        {
            Debug.Log("资源更新：加载本地Record文件");
            m_LocalFileRecords = ConvertToFileRecords(FileUtility.ReadFileText(FILERECORD_FILENAME, PathType.Local));
            HotUpdaterSwitchToState(HotUpdateState.LoadRemoteVersionFile);
        }

        #endregion

        #region Download remote version file

        void StartLoadRemoteVersionFile()
        {
            Debug.Log("资源更新：下载远程版本文件");
            WebClientDownloader.AppendDownloadFile(VERSION_FILENAME, DownLoadVersionFile_Complated, null);
        }

        /// <summary>
        /// 远程版本文件下载完成
        /// </summary>
        /// <param name="fileName"></param>
        /// <param name="e"></param>
        void DownLoadVersionFile_Complated(string fileName, AsyncCompletedEventArgs e)
        {
            if (fileName != VERSION_FILENAME)
                return;

            if (e.Error != null)
            {
                Debug.LogErrorFormat("读取服务器版本文件失败:\n{0}\n{1}", fileName, e.Error.Message, fileName);

                // 下载失败直接当出错处理
                HotUpdaterSwitchToFaildState(true);

                // 进入到本地文件校验环节
                // HotUpdaterSwitchToState(HotUpdateState.ValidateLocalData);
            }
            else if (e.Cancelled)
            {
                HotUpdaterSwitchToFaildState(true);
            }
            else
            {
                Debug.Log("资源更新：远程版本文件下载完成");
                m_RemoteVersion = FileUtility.ReadFileText(VERSION_FILENAME, PathType.Cache);

                Debug.LogFormat("资源更新：对比本地和远程版本号 本地{0}  远程{1}", m_LocalVersion, m_RemoteVersion);
                int versionResult = CheckVersion(m_RemoteVersion, m_LocalVersion);

                if (versionResult == 2)
                {
                    if (UpdateProjectName == AppDefine.HallNumber)
                    {
                        // 弹出 更新失败，需要前往重新下载游戏                    
                        HotUpdaterSwitchToFaildState(false);
                    }
                    else
                    {
                        // 进入更新文件列表                    
                        HotUpdaterSwitchToState(HotUpdateState.LoadRemoteFileRecords);
                    }
                }
                else if (versionResult == 1)
                {
                    // 进入更新文件列表                    
                    HotUpdaterSwitchToState(HotUpdateState.LoadRemoteFileRecords);
                }
                else
                {
                    HotUpdaterSwitchToState(HotUpdateState.ValidateLocalData);
                }
            }
        }

        #endregion

        #region Download remote record file

        /// <summary>
        /// 开始加载远程的列表文件
        /// </summary>
        void StartLoadRemoteRecordFile()
        {
            Debug.Log("资源更新：下载远程Record文件");
            WebClientDownloader.AppendDownloadFile(FILERECORD_FILENAME, DownLoadFileRecords_Complated, null);
        }

        /// <summary>
        /// 远程文件列表下载完成
        /// </summary>
        /// <param name="fileName"></param>
        /// <param name="e"></param>
        void DownLoadFileRecords_Complated(string fileName, AsyncCompletedEventArgs e)
        {
            if (fileName != FILERECORD_FILENAME) return;

            if (e.Error != null)
            {
                Debug.LogErrorFormat("Load remote records error: {0} ", e.Error.Message);
                HotUpdaterSwitchToFaildState(true);
            }
            else if (e.Cancelled)
            {
                HotUpdaterSwitchToFaildState(true);
            }
            else
            {
                Debug.Log("资源更新：远程Record文件下载完成");
                m_RemoteFileRecords = ConvertToFileRecords(FileUtility.ReadFileText(FILERECORD_FILENAME, PathType.Cache));
                HotUpdaterSwitchToState(HotUpdateState.DownloadData);
            }
        }

        #endregion

        #region Download remote asset files

        private long m_NeedUpdateSize = 0;                  // 需要下载数据大小
        private long m_UpdatedFileSize = 0;                 // 已经下载文件数据大小
        private long m_LastDownLoadedFileSize = 0;          // 最近一次记录下载文件的大小，用于比较两次大小计算下载速度 

        private float m_LastRefreshTime = 0;

        UpdateProgressInfo updateProgressInfo = new UpdateProgressInfo();

        void StartDownLoadRemoteUpdatedAssets()
        {
            m_NeedUpdateFileRecords = FilterNeedUpdateFileRecodes(m_RemoteFileRecords, m_LocalFileRecords);

            Debug.LogFormat("资源更新：开始更新远程资源文件，需要下载的文件数量为{0}", m_NeedUpdateFileRecords.Count);
            if (m_NeedUpdateFileRecords.Count > 0)
            {
                m_NeedUpdateFileRecords.ForEach(o => { m_NeedUpdateSize += o.Size; });      // 计算总共需要下载大小
                RefreshDownloadProgress();

                for (int i = m_NeedUpdateFileRecords.Count - 1; i >= 0; i--)
                {
                    FileRecord record = m_NeedUpdateFileRecords[i];
                    WebClientDownloader.AppendDownloadFile(record.FileName, DownloadAssetFile_Complated, null);
                }
                IsLoadedNewAsset = true;
            }
            else
            {
                HotUpdaterSwitchToState(HotUpdateState.ValidateLocalData);
            }
        }

        private void RefreshDownloadProgress()
        {
            if (m_UpdatedFileSize > m_NeedUpdateSize)     // 最后一个文件下载完后 m_CurrentDownLoadBytesReceived 会被重复计算
                m_UpdatedFileSize = m_NeedUpdateSize;

            updateProgressInfo.ProgressPercentage = m_UpdatedFileSize * 1f / m_NeedUpdateSize;
            updateProgressInfo.TipsContent = "正在下载文件";
            updateProgressInfo.DownLoadSizeInfo = string.Format("{0}MB/{1}MB", (m_UpdatedFileSize / (1024f * 1024f)).ToString("F2"), (m_NeedUpdateSize / (1024f * 1024f)).ToString("F2"));

            if (m_LastDownLoadedFileSize == 0)
            {
                updateProgressInfo.DownLoadSpeed = "0 KB/S";
                updateProgressInfo.NeedTime = "正在计算...";
            }
            else
            {
                float speed = (m_UpdatedFileSize - m_LastDownLoadedFileSize) / (Time.timeSinceLevelLoad - m_LastRefreshTime);

                if (speed != 0)
                {
                    float needTime = (m_NeedUpdateSize - m_UpdatedFileSize) / speed;

                    if (speed > 1024)
                        updateProgressInfo.DownLoadSpeed = string.Format("{0} MB/S", (speed / 1024f / 1024f).ToString("F2"));
                    else
                        updateProgressInfo.DownLoadSpeed = string.Format("{0} KB/S", (speed / 1024f).ToString("F0"));

                    updateProgressInfo.NeedTime = new DateTime(1970, 01, 01, 00, 00, 00).AddSeconds(needTime).ToString("HH:mm:ss");
                }
                else
                {
                    updateProgressInfo.DownLoadSpeed = "0 KB/S";
                    updateProgressInfo.NeedTime = "正在计算...";
                }
            }

            m_LastDownLoadedFileSize = m_UpdatedFileSize;
            m_LastRefreshTime = Time.timeSinceLevelLoad;

            ChangeHotUpdateProgress();
        }

        Queue<FileRecord> complatedFileRecords = new Queue<FileRecord>();
        /// <summary>
        /// 完成多少个之后回填本地记录信息
        /// </summary>
        int ComplatedDownloadSaveCount = 5;

        void DownloadAssetFile_Complated(string fileName, AsyncCompletedEventArgs e)
        {
            FileRecord record = m_NeedUpdateFileRecords.Find(temp => temp.FileName == fileName);
            if (record == null) return;

            if (e.Error != null)
            {
                Debug.LogErrorFormat("download remote file:[{0}] error[{1}]", fileName, e.Error.Message);
                WebClientDownloader.StopDownloader();
                HotUpdaterSwitchToFaildState(true);
            }
            else if (e.Cancelled)
            {
                HotUpdaterSwitchToFaildState(true);
            }
            else
            {
                m_NeedUpdateFileRecords.Remove(record);

                if (FileUtility.CopyFileByPathType(record.FileName, PathType.Cache, PathType.Local))
                {
                    complatedFileRecords.Enqueue(record);
                    m_UpdatedFileSize += record.Size;
                    IsRefreshUpdateProgress = true;
                }

                if (m_NeedUpdateFileRecords.Count == 0)
                {
                    // 下载完成了
                    DownloadComplatedRestore();
                    HotUpdaterSwitchToState(HotUpdateState.ClearObsoleteAssets);
                }
                // 大廳熱更確保最後三個檔案一次下載
                // Hall/AssetRecords.bytes
                // Hall/Hall
                // Hall/hall_main
                else if (UpdateProjectName == AppDefine.HallNumber)
                { 
                    if (m_NeedUpdateFileRecords.Count < 3)
                    {

                    }
                    else if (complatedFileRecords.Count >= ComplatedDownloadSaveCount)
                    {
                        DownloadComplatedRestore();
                    }
                }
                else if (complatedFileRecords.Count >= ComplatedDownloadSaveCount)
                {
                    DownloadComplatedRestore();
                }
            }
        }

        /// <summary>
        /// 下载完成回存资源列表信息
        /// </summary>
        void DownloadComplatedRestore()
        {
            Debug.Log("DownloadComplatedRestore");

            while (complatedFileRecords.Count > 0)
            {
                FileRecord record = complatedFileRecords.Dequeue();
                m_LocalFileRecords[record.FileName] = record;
            }

            RestoreLocalFileRecords(m_LocalFileRecords);
        }

        #endregion

        #region Clear obsolete asset files

        /// <summary>
        /// 清理下载时残留的资源信息
        /// </summary>
        void StartClearObsoleteAssetOfDownload()
        {
            StartCoroutine(ClearObsoleteAssetData(m_RemoteFileRecords, m_LocalFileRecords));
        }

        int PreFrameClearAssetCount = 20;
        /// <summary>
        /// 清理过期的资源信息
        /// </summary>
        /// <param name="newFileRecords"></param>
        /// <param name="oldFileRecords"></param>
        /// <returns></returns>
        IEnumerator ClearObsoleteAssetData(Dictionary<string, FileRecord> newFileRecords, Dictionary<string, FileRecord> oldFileRecords)
        {
            List<FileRecord> obsoleteList = FilterObsoleteFileRecodes(newFileRecords, oldFileRecords);
            int handleCount = 0;
            if (obsoleteList.Count > 0)
            {
                for (int i = 0; i < obsoleteList.Count; i++)
                {
                    if (handleCount >= PreFrameClearAssetCount)
                    {
                        yield return 1;
                        handleCount = 0;

                        updateProgressInfo.TipsContent = "清理资源...";
                        ChangeHotUpdateProgress();
                    }
                    handleCount++;

                    //Debug.Log("ClearObsoleteAssetData");

                    if (FileUtility.DeleteFileOfType(obsoleteList[i].FileName, PathType.Local))
                    {
                        if (oldFileRecords.ContainsKey(obsoleteList[i].FileName))
                        {
                            oldFileRecords.Remove(obsoleteList[i].FileName);
                        }
                    }
                }
                obsoleteList.Clear();
                RestoreLocalFileRecords(oldFileRecords);
            }

            HotUpdaterSwitchToState(HotUpdateState.ValidateLocalData);
        }

        #endregion

        #region Validate local asset files

        /// <summary>
        /// 校验本地资源文件
        /// </summary>
        void StartValidateLocalAssetFiles()
        {
            StartCoroutine(ValidateLoadAssetsCoroutine());
        }

        /// <summary>
        /// 每帧校验文件的最大数量
        /// </summary>
        const int PerFrameValidateMax = 5;

        IEnumerator ValidateLoadAssetsCoroutine()
        {
            Debug.LogFormat("资源更新：开始验证本地资源，需要验证文件数量:{0}", m_LocalFileRecords.Count);
            bool isContainsErrorFile = false;

            if (m_LocalFileRecords.Count > 0)
            {
                int handleCount = 0;
                int index = 0;

                foreach (var item in m_LocalFileRecords)
                {
                    index++;
                    handleCount++;
                    if (handleCount < PerFrameValidateMax)
                    {
                        FileRecord record = item.Value;
                        // 文件不存在，MD5为空，或者不相等
                        // Add s_ShowFullLog By Rayon
                        if (s_ShowFullLog == false)
                        {
                            if (!FileUtility.ExistsFileOfType(record.FileName, PathType.Local)
                                 || string.IsNullOrEmpty(record.MD5)
                                 || record.MD5 != FileUtility.GetMD5HashFromFile(record.FileName, PathType.Local))
                            {
                                Debug.LogErrorFormat("File error:{0}", record.FileName);
                                isContainsErrorFile = true;
                                record.MD5 = "";// 将MD5码置空，便于下次重新来过时可以更新到该资源
                            }
                        }
                        else
                        {

                            if (!FileUtility.ExistsFileOfType(record.FileName, PathType.Local))
                            {
                                Debug.LogErrorFormat("File error:(No file)\n{0}", record.FileName);
                                isContainsErrorFile = true;
                                record.MD5 = "";// 将MD5码置空，便于下次重新来过时可以更新到该资源
                            }

                            if (string.IsNullOrEmpty(record.MD5))
                            {
                                Debug.LogErrorFormat("File error:(No MD5)\n{0}", record.FileName);
                                isContainsErrorFile = true;
                                record.MD5 = "";// 将MD5码置空，便于下次重新来过时可以更新到该资源
                            }

                            string _MD5 = FileUtility.GetMD5HashFromFile(record.FileName, PathType.Local);
                            if (record.MD5 != _MD5)
                            {
                                Debug.LogErrorFormat("File error:(MD5 error)\n{0}\n{1}\n{2}",
                                                    record.FileName,
                                                    record.MD5,
                                                    _MD5);
                                isContainsErrorFile = true;
                                record.MD5 = "";// 将MD5码置空，便于下次重新来过时可以更新到该资源
                            }
                        }
                    }
                    else
                    {
                        yield return 1;

                        updateProgressInfo.TipsContent = "验证本地资源...";
                        ChangeHotUpdateProgress();
                        handleCount = 0;
                    }
                }
                updateProgressInfo.TipsContent = "验证本地资源...";
                ChangeHotUpdateProgress();
            }

            // 如果本地没有文件，验证不允许通过，否则后续流程取资源会出问题
            if (m_LocalFileRecords.Count == 0)
            {
                Debug.Log("如果本地没有文件，验证不允许通过");

                HotUpdaterSwitchToFaildState(true);
            }
            else
            {
                if (isContainsErrorFile)
                {
                    Debug.Log("资源更新：Record记录和实际文件MD5不相同");
                    RestoreLocalFileRecords(m_LocalFileRecords);
                    // 验证资源失败
                    HotUpdaterSwitchToFaildState(true);
                }
                else
                {
                    // 拷贝远程版本号到本地
                    FileUtility.CopyFileByPathType(VERSION_FILENAME, PathType.Cache, PathType.Local);

                    Debug.Log("资源更新：资源验证通过");
                    HotUpdaterSwitchToState(HotUpdateState.Done);
                }
                
                // 清理本地临时文件路径
                FileUtility.DeleteDirectory(AppDefine.LOCAL_TEMP_PATH, true);
            }
        }

        #endregion

        #region 更新失败

        int m_UpdateFaildRetryTimes = 1;

        public static bool s_HaveHotUpdateFaild = true;

        private System.Action m_CallbackHotUpdateFaild;
        /// <summary> 支援版本[CallbackHotUpdateFail]，下載失敗的事件 </summary>
        public void SetCallbackHotUpdateFaild(System.Action callback)
        {
            m_CallbackHotUpdateFaild = callback;
        }

        private System.Action m_CallbackHotUpdateFaildCancel;
        /// <summary> 支援版本[CallbackHotUpdateFail]，下載失敗，詢問視窗點取消按鈕的事件 </summary>
        public void SetCallbackHotUpdateFailCancel(System.Action callback)
        {
            m_CallbackHotUpdateFaildCancel = callback;
        }

        private bool m_HotUpdateFaildPassMessageBox = false;
        /// <summary> 支援版本[CallbackHotUpdateFail]，下載失敗，跳過詢問視窗 </summary>
        public void SetHotUpdateFaildPassMessageBox(bool yes_or_no)
        {
            m_HotUpdateFaildPassMessageBox = yes_or_no;
        }

        void HotUpdateFaild()
        {
            Debug.LogError("HotUpdateFaild");
            IsValidatedAssets = false;

            if (OnHotUpdateFaild != null)
                OnHotUpdateFaild();
            
            if (m_CallbackHotUpdateFaild != null)
                m_CallbackHotUpdateFaild();
            
            if (m_HotUpdateFaildPassMessageBox == true)
                return;

            if (m_CanRetryUpdateFaild)
            {
                ShowRetryUpdaterFaildMessageBox();
            }
            else
            {
                ShowCannotUpdateFaildMessageBox();
            }
        }

        void ShowRetryUpdaterFaildMessageBox()
        {
            MessageBoxData boxData = new MessageBoxData();
            boxData.Title = "错误";
            boxData.Content = "更新资源错误，是否重试?";
            boxData.OKButtonName = "重试";
            boxData.CancelButtonName = "关闭";
            boxData.Style = MessageBoxStyle.OKCancel;
            boxData.CallBack = (result) =>
            {
                switch (result)
                {
                    case MessageBoxResult.OK:
                        // 移除熱更詢問超過3次關掉遊戲 By Rayon
                        /*
                        if (this.m_UpdateFaildRetryTimes >= 3)
                        {
                            ShowCannotUpdateFaildMessageBox();
                        }
                        else
                        */
                        {
                            m_UpdateFaildRetryTimes++;
                            TryStartUpdateRemoteAssetsToLocal();
                        }
                        break;
                    case MessageBoxResult.Cancel:
                    default:
                        //Application.Quit();
                        if (m_CallbackHotUpdateFaildCancel != null)
                            m_CallbackHotUpdateFaildCancel();

                        break;
                }
            };
            UIMessageBox.Show(boxData, null);
        }

        /// <summary>
        /// 显示不能更新的错误提示
        /// </summary>
        void ShowCannotUpdateFaildMessageBox()
        {
            MessageBoxData boxData = new MessageBoxData();
            boxData.Title = "更新失败";
            boxData.Content = "资源更新失败!";
            // boxData.OKButtonName = "重新下载";
            boxData.CancelButtonName = "退出游戏";
            boxData.Style = MessageBoxStyle.OK;
            boxData.CallBack = (result) =>
            {
                //switch (result)
                //{
                //    case MessageBoxResult.OK:
                //        Application.OpenURL(AppDefine.GameUrl);
                //        Application.Quit();
                //        break;
                //    case MessageBoxResult.Cancel:
                //    default:
                Application.Quit();
                //        break;
                //}
            };
            UIMessageBox.Show(boxData, null);
        }

        #endregion

        #region 更新完成

        void HotUpdateComplated()
        {
            IsValidatedAssets = true;
            IsWorking = false;

            if (OnHotUpdateComplated != null)
                OnHotUpdateComplated();
        }

        #endregion

        /// <summary>
        /// 校验版本号
        /// </summary>
        /// <param name="newVersionStr">新版本号</param>
        /// <param name="oldVersionStr">旧版本号</param>
        /// <returns>-1：不需要更新 0: 版本一致，1：更新，2：强制更新</returns>
        public static int CheckVersion(string newVersionStr, string oldVersionStr)
        {
            if (!AppDefine.OpenHotUpdate)
                return 0;

            if (string.IsNullOrEmpty(newVersionStr))
                newVersionStr = DEFAULT_VERSION;
            if (string.IsNullOrEmpty(oldVersionStr))
                oldVersionStr = DEFAULT_VERSION;
            Version newVersion = new Version(newVersionStr);
            Version oldVersion = new Version(oldVersionStr);

            if (oldVersion < newVersion) // 旧的版本小于新的版本，需要更新
            {
                //1, 2 位版本号变化，需要强制重新下载 3, 4 位版本号变化，只需要更新
                if (oldVersion.Major != newVersion.Major)
                {
                    return 2;
                }
                else
                {
                    if (oldVersion.Minor != newVersion.Minor)
                    {
                        return 2;
                    }
                }
                return 1;
            }
            else if (oldVersion == newVersion)
            {
                return 0; // 旧版本和新版本一致，不需要更新
            }
            else
            {
                return -1; // 旧版本大于新版本，不需要更新
            }
        }

        /// <summary>
        /// 转换文件列表内容为文件记录信息
        /// </summary>
        /// <param name="content"></param>
        /// <returns></returns>
        public static Dictionary<string, FileRecord> ConvertToFileRecords(string content)
        {
            Dictionary<string, FileRecord> fileRecords = new Dictionary<string, FileRecord>();
            Mono.Xml.SecurityParser sp = new Mono.Xml.SecurityParser();
            if (string.IsNullOrEmpty(content))
                return fileRecords;
            sp.LoadXml(content);
            SecurityElement root = sp.ToXml();
            if (root.Children != null)
            {
                foreach (SecurityElement child in root.Children)
                {
                    SecurityElement fileName = child.SearchForChildByTag("NAME");
                    SecurityElement md5 = child.SearchForChildByTag("MD5");
                    SecurityElement size = child.SearchForChildByTag("SIZE");

                    if (fileName == null || md5 == null || size == null)
                    {
                        Debug.LogError("'FileList.xml' has error!");
                        continue;
                    }

                    FileRecord record = new FileRecord();
                    record.FileName = fileName.Text;
                    record.MD5 = md5.Text;
                    record.Size = long.Parse(size.Text);
                    fileRecords.Add(record.FileName, record);
                }
            }
            return fileRecords;
        }

        /// <summary>
        /// 筛选出过期的资源文件
        /// </summary>
        /// <returns></returns>
        List<FileRecord> FilterObsoleteFileRecodes(Dictionary<string, FileRecord> newFileRecords, Dictionary<string, FileRecord> oldFileRecords)
        {
            List<FileRecord> fileRecordList = new List<FileRecord>();
            // 筛选出需要更新的文件
            foreach (var item in oldFileRecords)
            {
                if (!newFileRecords.ContainsKey(item.Key))
                {
                    fileRecordList.Add(item.Value);
                }
            }

            return fileRecordList;
        }

        /// <summary>
        /// 筛选出需要更新的文件
        /// </summary>
        /// <returns></returns>
        List<FileRecord> FilterNeedUpdateFileRecodes(Dictionary<string, FileRecord> newFileRecords, Dictionary<string, FileRecord> oldFileRecords)
        {
            List<FileRecord> fileRecordList = new List<FileRecord>();
            // 筛选出需要更新的文件
            foreach (var item in newFileRecords)
            {
                if (oldFileRecords.ContainsKey(item.Key))
                {
                    if (oldFileRecords[item.Key].MD5 == item.Value.MD5)
                    {
                        continue;
                    }
                }
                fileRecordList.Add(item.Value);
            }

            return fileRecordList;
        }

        #endregion

        #region Notify HotUpdate Progess Changed

        void ChangeHotUpdateProgress()
        {
            if (OnHotUpdateProgressChanged != null)
                OnHotUpdateProgressChanged(updateProgressInfo);
        }

        #endregion 
    }
}