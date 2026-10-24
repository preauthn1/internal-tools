using DG.Tweening;
using HotUpdate;
using UnityEngine;
using UnityEngine.UI;

/// <summary>
/// 加载界面,显示解压以及大厅更新进度
/// </summary>
public class UILoading : WindowBase
{
    public static readonly string UIAssetName = "UILoading";

    static UpdateProgressInfo m_UIData = new UpdateProgressInfo();

    /// <summary>
    /// 提示标签
    /// </summary>
    [SerializeField]
    Text m_TipsInfoText = null;

    /// <summary>
    /// 提示标签
    /// </summary>
    [SerializeField]
    Text m_DownLoadSizeInfo = null;

    /// <summary>
    /// 提示标签
    /// </summary>
    [SerializeField]
    Text m_DownLoadSpeed = null;

    /// <summary>
    /// 提示标签
    /// </summary>
    [SerializeField]
    Text m_NeedTime = null;

    /// <summary>
    /// 进度
    /// </summary>
    [SerializeField]
    Slider m_ProgressValue = null;

    /// <summary>
    /// 下载状态信息
    /// </summary>
    [SerializeField]
    RectTransform DownLoadInfo = null;


    /// <summary>
    /// 游戏版本号
    /// </summary>
    [SerializeField]
    Text m_VersionInfoText = null;

    private void Awake()
    {
        Transform _GO = transform.Find("LOGO_62");
        if (_GO != null)
            _GO.gameObject.AddComponent<LoadingStreamingAssetsTexture>();

        Transform _GO_1 = transform.Find("LOGO_62_1");
        if (_GO_1 != null)
            _GO_1.gameObject.AddComponent<LoadingStreamingAssetsTexture>();
    }

    // Use this for initialization
    protected void Start()
    {
        if (WindowNode != openedWindowNode)
        {
            WindowManager.Instance.CloseWindow(WindowNode, false);
            return;
        }

        if (m_TipsInfoText == null || m_ProgressValue == null)
        {
            Debug.LogError("Some ui element was null in LoadingUI script, please check it!");
            return;
        }

        if (m_VersionInfoText != null)
            m_VersionInfoText.text = "V." + HotUpdate.HotUpdateManager.Instance.GetProjectVersion(AppDefine.HallNumber);

        UpdateUIElement();
    }


    static WindowNode openedWindowNode = null;

    /// <summary>
    /// 显示数据加载界面
    /// </summary>    
    public static void Show(System.Action<WindowNode> complatedCallBack = null)
    {
        if (openedWindowNode == null)
        {
            WindowNodeInitParam initParam = new WindowNodeInitParam(UIAssetName, AppDefine.HallNumber);        // 大厅内置窗口
            initParam.NodeType = BaseNodeType.Loading;
            initParam.LoadComplatedCallBack = complatedCallBack;
            openedWindowNode = WindowManager.Instance.OpenWindow(initParam);
            HotUpdate.HotUpdateManager.OnHotUpdateProgressChanged += HotUpdate_ProgressChangedEvent;
        }
    }

    static void HotUpdate_ProgressChangedEvent(HotUpdate.UpdateProgressInfo e)
    {
        m_UIData = e;
    }

    /// <summary>
    /// 隐藏数据加载界面
    /// </summary>
    public static void Hide()
    {
        if (openedWindowNode != null)
        {
            HotUpdate.HotUpdateManager.OnHotUpdateProgressChanged -= HotUpdate_ProgressChangedEvent;
            WindowManager.Instance.CloseWindow(openedWindowNode, false);
            openedWindowNode = null;
            m_UIData.DownLoadSizeInfo = null;
        }
    }

    void UpdateUIElement()
    {
        // 有下载信息
        if (!string.IsNullOrEmpty(m_UIData.DownLoadSizeInfo) && !DownLoadInfo.gameObject.activeSelf)
            DownLoadInfo.gameObject.SetActive(true);

        m_TipsInfoText.text = m_UIData.TipsContent;
        m_DownLoadSizeInfo.text = m_UIData.DownLoadSizeInfo;
        m_DownLoadSpeed.text = m_UIData.DownLoadSpeed;
        m_NeedTime.text = m_UIData.NeedTime;

        if (m_ProgressValue.value > m_UIData.ProgressPercentage)
            m_ProgressValue.value = m_UIData.ProgressPercentage;
        else
            m_ProgressValue.DOValue(m_UIData.ProgressPercentage, 0.3f);
    }

    void Update()
    {
        UpdateUIElement();
    }

    /// <summary>
    /// 设置提示消息
    /// </summary>
    /// <param name="msg"></param>
    public static void SetTipMessage(string msg, float progressPercentage)
    {
        m_UIData.TipsContent = msg;                         // 提示信息
        m_UIData.ProgressPercentage = progressPercentage;   // 进度信息
    }
}