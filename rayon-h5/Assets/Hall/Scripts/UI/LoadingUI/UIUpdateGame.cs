using DG.Tweening;
using HotUpdate;
using UnityEngine;
using UnityEngine.UI;

/// <summary>
/// 加载界面,显示解压以及大厅更新进度
/// </summary>
public class UIUpdateGame : WindowBase
{
    public static readonly string UIAssetName = "UIUpdateGame";

    static UpdateProgressInfo m_UIData = new UpdateProgressInfo();

    /// <summary>
    /// 提示标签
    /// </summary>
    [SerializeField]
    Text m_DownLoadSizeInfo = null;

    /// <summary>
    /// 进度
    /// </summary>
    [SerializeField]
    Slider m_ProgressValue = null;

    // Use this for initialization
    protected void Start()
    {
        m_ProgressValue.value = 0;

        if (WindowNode != openedWindowNode)
        {
            WindowManager.Instance.CloseWindow(WindowNode, false);
            return;
        }

        if (m_ProgressValue == null)
        {
            Debug.LogError("Some ui element was null in LoadingUI script, please check it!");
            return;
        }
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
            HotUpdateManager.OnHotUpdateProgressChanged += HotUpdate_ProgressChangedEvent;
            HotUpdateManager.OnHotUpdateComplated += Hide;
            HotUpdateManager.OnHotUpdateFaild += Hide;
        }
    }

    /// <summary>
    /// 隐藏数据加载界面
    /// </summary>
    public static void Hide()
    {
        if (openedWindowNode != null)
        {
            HotUpdateManager.OnHotUpdateProgressChanged -= HotUpdate_ProgressChangedEvent;
            HotUpdateManager.OnHotUpdateComplated -= Hide;
            HotUpdateManager.OnHotUpdateFaild -= Hide;
            WindowManager.Instance.CloseWindow(openedWindowNode, false);
            openedWindowNode = null;
        }
    }

    static void HotUpdate_ProgressChangedEvent(HotUpdate.UpdateProgressInfo e)
    {
        m_UIData = e;
    }

    void Update()
    {
        m_DownLoadSizeInfo.text = m_UIData.DownLoadSizeInfo;

        if (m_UIData.ProgressPercentage == 0)
            m_ProgressValue.value = 0;
        else
            m_ProgressValue.DOValue(m_UIData.ProgressPercentage, 0.3f);
    }
}