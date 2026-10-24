using System.Collections.Generic;
using UnityEngine;

/// <summary>
/// 程序固定基本信息定义
/// </summary>
public partial class AppDefine
{
    public static string m_CurrentProjectPath = string.Empty;

    public static string CurrentProjectPath
    {
        get
        {
            if (m_CurrentProjectPath == string.Empty)
                m_CurrentProjectPath = PlayerPrefs.GetString("CurrentProjectPath");

            return m_CurrentProjectPath;
        }
        set
        {
            if (m_CurrentProjectPath != value)
            {
                m_CurrentProjectPath = value;
                PlayerPrefs.SetString("CurrentProjectPath", value);
            }
        }
    }

    public static string CurrentProjectFolderPath
    {
        get { return CurrentProjectPath.Replace("Assets/", string.Empty); }
    }

    public static string CurrentProjectFolder
    {
        get { return CurrentProjectFolderPath.Replace("Game/", string.Empty); }
    }
     
    public static string ProjectListFileName = "ProjectList.txt";

    /// <summary>
    /// 是否加密
    /// </summary>
    public static bool IsEncrypted = false/* true*/;

    /// <summary>
    /// 资源密码
    /// </summary>
    public static string AssetSecretKey = "207102htialgnahc";

    /// <summary>
    /// 远程资源的根目录
    /// </summary>
    public static string UpdateFileUrl = "";     // 在Lua启动之后根据Lua配置项修改

    /// <summary>
    /// 整包更新地址
    /// </summary>
    public static string GameUrl = "";          // 完整更新包地址

    public static bool OpenHotUpdate = true;

    /// <summary>
    /// 资源记录文件
    /// </summary>
    public static string AssetRecordsFileName
    {
        get { return "AssetRecords.bytes"; }
    }

    public static string HallNumber
    {
        get { return "Hall"; }
    }

    /// <summary>
    /// 获取平台字符串
    /// </summary>
    /// <returns></returns>
    static string GetPlatformString()
    {
#if UNITY_EDITOR
        switch (UnityEditor.EditorUserBuildSettings.activeBuildTarget)
        {
            case UnityEditor.BuildTarget.Android:
                return "Android";
            case UnityEditor.BuildTarget.iOS:
                return "iOS";
            case UnityEditor.BuildTarget.WebGL:
                return "WebGL";
            case UnityEditor.BuildTarget.StandaloneWindows:
            case UnityEditor.BuildTarget.StandaloneWindows64:
            case UnityEditor.BuildTarget.StandaloneOSXIntel64:
            case UnityEditor.BuildTarget.StandaloneOSXIntel:
            case UnityEditor.BuildTarget.StandaloneOSX:
            case UnityEditor.BuildTarget.StandaloneLinux64:
            case UnityEditor.BuildTarget.StandaloneLinux:
            case UnityEditor.BuildTarget.StandaloneLinuxUniversal:
                return "Windows";
            default:
                return null;
        }
#else
        switch (Application.platform)
        {
            case RuntimePlatform.Android:
                return "Android";
            case RuntimePlatform.IPhonePlayer:
                return "iOS";
            case RuntimePlatform.WebGLPlayer:
                return "WebGL";
            case RuntimePlatform.WindowsPlayer:
            case RuntimePlatform.WindowsEditor:
            case RuntimePlatform.LinuxEditor:
            case RuntimePlatform.OSXEditor:
                return "Windows";
            default:
                return null;
        }
#endif
    }

    private static string m_PlatfromPath = GetPlatformString();
    /// <summary>
    /// 平台路径
    /// </summary>
    public static string PlatformPath
    {
        get { return m_PlatfromPath; }
    }

    /// <summary>
    /// 远程更新数据根目录
    /// </summary>
    public static string REMOTE_DATA_PATH
    {
        get { return string.Format("{0}/{1}", UpdateFileUrl, PlatformPath); }
    }

    private static string m_Local_Init_Data_Path = string.Format("{0}/{1}", Application.streamingAssetsPath, PlatformPath);
    /// <summary>
    /// 本地初始数据根目录(streamingAssetsPath + PlatformPath)
    /// </summary>
    public static string LOCAL_INIT_DATA_PATH
    {
        get { return m_Local_Init_Data_Path; }
    }

#if UNITY_EDITOR
    private static string m_Loacl_Data_Path = string.Format("C:/Test/{0}/{1}", PlatformPath, "Cache/Data");
#else    
    private static string m_Loacl_Data_Path = string.Format("{0}/{1}/{2}", Application.persistentDataPath, PlatformPath, "Cache/Data");
#endif

    /// <summary>
    /// 本地数据根目录
    /// </summary>
    public static string LOCAL_DATA_PATH
    {
        get { return m_Loacl_Data_Path; }
    }

#if UNITY_EDITOR
    private static string m_Local_Temp_Path = string.Format("C:/Test/{0}/{1}", PlatformPath, "Cache/Temp");
#else
    private static string m_Local_Temp_Path = string.Format("{0}/{1}/{2}", Application.persistentDataPath, PlatformPath, "Cache/Temp");
#endif


    /// <summary>
    /// 本地数据临时根目录
    /// </summary>
    public static string LOCAL_TEMP_PATH
    {
        get { return m_Local_Temp_Path; }
    }
    
    
    // 使用伺服器切換功能
    public static bool UseChangeServerIP = false;

    /// <summary> 登入流程類型
    /// 0.正式流程(網路抓serverlist.txt檔案)
    /// 1.指定伺服器列表位置
    /// 2.指定大廳伺服器位置
    /// </summary>
    public static int LoginProcessType
    {
        get { 
            if (UseChangeServerIP == false)
                return 0;

            return PlayerPrefs.GetInt("mLoginProcessType", 0); 
        }
        set { PlayerPrefs.SetInt("mLoginProcessType", value); }
    }

    // 
    public static string ServerListDomain
    {
        get { return PlayerPrefs.GetString("mServerListDomain"); }
        set { PlayerPrefs.SetString("mServerListDomain", value); }
    }
    // 大厅服务器IP
    public static string DefaultServerIP
    {
        get { return PlayerPrefs.GetString("mServerIP"); }
        set { PlayerPrefs.SetString("mServerIP", value); }
    }
    // 大厅服务器端口
    public static int DefaultServerPort
    {
        get { return PlayerPrefs.GetInt("mServerPort", 10000); }
        set { PlayerPrefs.SetInt("mServerPort", value); }
    }
    // 資源更新IP
    public static string UpdateDomainIP
    {
        get { return PlayerPrefs.GetString("mUpdateDomainIP"); }
        set { PlayerPrefs.SetString("mUpdateDomainIP", value); }
    }
    // 資源更新Port
    public static string UpdateDomainPort
    {
        get { return PlayerPrefs.GetString("mUpdateDomainPort"); }
        set { PlayerPrefs.SetString("mUpdateDomainPort", value); }
    }
    // 充值IP
    public static string ChargeDomainIP
    {
        get { return PlayerPrefs.GetString("mChargeDomainIP"); }
        set { PlayerPrefs.SetString("mChargeDomainIP", value); }
    }
    // 获取本机IP
    public static string IPReuestIP
    {
        get { return PlayerPrefs.GetString("mIPReuestIP"); }
        set { PlayerPrefs.SetString("mIPReuestIP", value); }
    }

    // 游戏渠道编号
    public static int GameSiteID
    {
        get
        {
            if (UseChangeServerIP)
                return PlayerPrefs.GetInt("mSiteID", 1);
            
            if (Main.Instance != null && Main.Instance.m_CustomDefine != null)
                return Main.Instance.m_CustomDefine.StaticSiteID;

            return 2;
        }
        set { PlayerPrefs.SetInt("mSiteID", value); }
    }
    // 游戏平台编号
    public static int GamePlatformID
    {
        get
        {
            if (UseChangeServerIP)
                return PlayerPrefs.GetInt("mPlatformID", 1);

            if (Main.Instance != null && Main.Instance.m_CustomDefine != null)
                return Main.Instance.m_CustomDefine.PlatformID;

            return 2;
        }
        set { PlayerPrefs.SetInt("mPlatformID", value); }
    }
    // 伺服器選擇
    public static int ServerOption
    {
        get
        {
            if (UseChangeServerIP)
                return PlayerPrefs.GetInt("mServerOption", 1);

            if (Main.Instance != null && Main.Instance.m_CustomDefine != null)
                return Main.Instance.m_CustomDefine.ServerOption;

            return 3;
        }
        set { PlayerPrefs.SetInt("mServerOption", value); }
    }

    // 編輯器模式下遊客登入是否使用固定帳號
    public static bool useDeviceUniqueIdentifier = false;

    public static void StartLog()
    {
#if UNITY_EDITOR
        LogConsole.StartLog();
#endif
    }

    public static void StopLog()
    {
#if UNITY_EDITOR
        LogConsole.StopLog();
#endif
    }
}