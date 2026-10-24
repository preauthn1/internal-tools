using UnityEngine;
using UnityEngine.UI;
using DG.Tweening;

public class UILoadingData : WindowBase
{
    public readonly static string UIAssetName = "UILoadingData";

    static WindowNode openedWindowNode = null;
    private bool IsAutoHide = true;

    [SerializeField]
    Image m_Bg;

    /// <summary>
    /// 内容文本
    /// </summary>
    [SerializeField]
    Text m_ContentText = null;

    /// <summary>
    /// 显示数据加载界面
    /// </summary>
    /// <param name="tipsContent">提示文字内容</param>
    /// <param name="timeOut">超时时间</param>
    /// <param name="timeOutCallBack">超时回调</param>
    public static void Show(string tipsContent, int timeOut = 5, System.Action timeOutCallBack = null)
    {
        LoadingDataUIData data = new LoadingDataUIData()
        {
            MaxLastTime = timeOut,
            TipsContent = tipsContent,
            MaxTimeCloseWindowCallBack = timeOutCallBack
        };

        Show(data);
    }

    public static void Show(LoadingDataUIData uiData)
    {
        if (openedWindowNode == null)
        {
            WindowNodeInitParam initParam = new WindowNodeInitParam(UIAssetName, AppDefine.HallNumber);       // 大厅内置窗口
            initParam.NodeType = BaseNodeType.AboveNormal;

            initParam.WindowData = uiData;

            openedWindowNode = WindowManager.Instance.OpenWindow(initParam);
        }
        else
        {
            openedWindowNode.WindowData = uiData;
            openedWindowNode.ShowWindow();
        }
    }


    /// <summary>
    /// 隐藏数据加载界面
    /// </summary>
    public static void Hide()
    {
        if (openedWindowNode != null)
        {
            StartCountDown = false;
            WindowManager.Instance.CloseWindow(openedWindowNode, true);
            openedWindowNode = null;
        }
    }

    private float MaxLastTime = 10f;

    private System.Action MaxTimeCloseWindowCallBack = null;

    public static bool StartCountDown = false;

    public override void RefreshWindowData(object windowData)
    {
        base.RefreshWindowData(windowData);

        if (windowData is LoadingDataUIData)
        {
            LoadingDataUIData loadingDataUIData = (LoadingDataUIData)windowData;

            this.m_ContentText.text = loadingDataUIData.TipsContent;
            m_Bg.gameObject.SetActive(loadingDataUIData.IsFrozenUI);

            MaxLastTime = loadingDataUIData.MaxLastTime;
            this.MaxTimeCloseWindowCallBack = loadingDataUIData.MaxTimeCloseWindowCallBack;
            StartCountDown = true;

            IsAutoHide = MaxLastTime != 0;
        }
    }

    void Update()
    {
        if (StartCountDown && IsAutoHide)
        {
            MaxLastTime -= Time.unscaledDeltaTime;
            if (MaxLastTime <= 0)
            {
                if (MaxTimeCloseWindowCallBack != null)
                {
                    MaxTimeCloseWindowCallBack();
                }
                WindowManager.Instance.CloseWindow(this.WindowNode, true);
                openedWindowNode = null;
            }
        }
    }

}

/// <summary>
/// 加载数据界面的数据
/// </summary>
public class LoadingDataUIData
{
    public LoadingDataUIData()
    {
        MaxLastTime = 10;
    }

    /// <summary>
    /// 提示的内容
    /// </summary>
    public string TipsContent { get; set; }

    private int m_MaxLastTime = 10;

    /// <summary>
    /// 最大持续时间 多少秒(范围[1-20s],默认 10s)
    /// </summary>
    public int MaxLastTime
    {
        get { return m_MaxLastTime; }
        set
        {
            m_MaxLastTime = value;
        }
    }

    public bool IsFrozenUI = true;

    /// <summary>
    /// 到达最大持续时间时关闭界面的回调
    /// </summary>
    public System.Action MaxTimeCloseWindowCallBack { get; set; }
}
