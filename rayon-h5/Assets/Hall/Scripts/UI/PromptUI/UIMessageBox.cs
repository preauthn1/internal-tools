using System;
using UnityEngine;
using UnityEngine.UI;

public class UIMessageBox : WindowBase
{
    public static readonly string UIAssetName = "UIMessageBox";
    static WindowNode openedWindowNode = null;

    [SerializeField]
    Button m_CloseButton = null;

    [SerializeField]
    Button m_CloseAreaButton = null;

    [SerializeField]
    Button m_OKButton = null;

    [SerializeField]
    Button m_CancelButton = null;

    [SerializeField]
    Text m_OKButtonText = null;

    [SerializeField]
    Text m_CancelButtonText = null;

    [SerializeField]
    Text m_TitleText = null;

    [SerializeField]
    Text m_ContentText = null;

    [SerializeField]
    RectTransform m_ContentArea = null;

    /// <summary>
    /// 使用倒计时
    /// </summary>
    bool m_UseCountDown = false;

    float m_CountDown = 0;

    System.Action<MessageBoxResult> CallBack = null;

    string m_OkButtonName = "OK";

    bool m_EnableCloseArea = false;

    const int ContentMinHeight = 350;

    protected void Start()
    {
        if (m_OKButton == null || m_CloseButton == null || m_CancelButton == null || m_CloseAreaButton == null || m_TitleText == null || m_ContentText == null)
        {
            Debug.LogError("Some ui element was null in MessageBox script, please check it!");
            return;
        }

        m_OKButton.onClick.AddListener(() => OnHandleButtonClick(MessageBoxResult.OK));
        m_CancelButton.onClick.AddListener(() => OnHandleButtonClick(MessageBoxResult.Cancel));
        m_CloseButton.onClick.AddListener(() => OnHandleButtonClick(MessageBoxResult.Close));
        m_CloseAreaButton.onClick.AddListener(() =>
        {
            if (m_EnableCloseArea)
                OnHandleButtonClick(MessageBoxResult.CloseArea);
        });
    }

    public override void RefreshWindowData(object windowData)
    {
        base.RefreshWindowData(windowData);
        MessageBoxData uiData = windowData as MessageBoxData;
        if (uiData == null)
        {
            return;
        }
        currentCountDownSecond = -1;
        m_CountDown = uiData.LastTime;
        m_UseCountDown = uiData.LastTime > 0;

        if (m_TitleText != null)
        {
            m_TitleText.text = uiData.Title;
        }

        if (m_ContentText != null)
        {
            m_ContentText.text = uiData.Content;
            m_ContentText.rectTransform.sizeDelta = new Vector2(m_ContentText.rectTransform.sizeDelta.x, Mathf.Max(ContentMinHeight - 60, m_ContentText.preferredHeight));
            m_ContentArea.sizeDelta = new Vector2(m_ContentArea.sizeDelta.x, Mathf.Max(ContentMinHeight, m_ContentText.preferredHeight + 60));
        }

        if (m_CancelButtonText != null)
        {
            m_CancelButtonText.text = uiData.CancelButtonName;
        }

        m_OkButtonName = uiData.OKButtonName;
        if (m_OKButtonText != null)
        {
            if (m_UseCountDown)
            {
                UpdateCountDownDisplay(m_CountDown);
            }
            else
            {
                m_OKButtonText.text = m_OkButtonName;
            }
        }
        this.CallBack = uiData.CallBack;
        this.m_EnableCloseArea = uiData.EnableCloseArea;

        SetUIButtonActiveByStyle(uiData.Style);
    }


    void SetUIButtonActiveByStyle(MessageBoxStyle style)
    {
        switch (style)
        {
            case MessageBoxStyle.OK:
                m_OKButton.gameObject.SetActive(true);
                m_CancelButton.gameObject.SetActive(false);
                m_CloseButton.gameObject.SetActive(false);
                break;
            case MessageBoxStyle.OKCancel:
                m_OKButton.gameObject.SetActive(true);
                m_CancelButton.gameObject.SetActive(true);
                m_CloseButton.gameObject.SetActive(false);
                break;
            case MessageBoxStyle.OKCancelClose:
                m_OKButton.gameObject.SetActive(true);
                m_CancelButton.gameObject.SetActive(true);
                m_CloseButton.gameObject.SetActive(true);
                break;
            case MessageBoxStyle.OKClose:
                m_OKButton.gameObject.SetActive(true);
                m_CancelButton.gameObject.SetActive(false);
                m_CloseButton.gameObject.SetActive(true);
                break;
        }
    }


    void OnHandleButtonClick(MessageBoxResult result)
    {
        openedWindowNode = null;

        if (CallBack != null)
        {
            CallBack(result);
        }

        if (WindowManager.Instance != null)
        {
            WindowManager.Instance.CloseWindow(WindowNode, true);
        }
        else
        {
            Destroy(gameObject);
        }
    }

    private void Update()
    {
        if (m_UseCountDown)
        {
            m_CountDown -= Time.deltaTime;
            UpdateCountDownDisplay(m_CountDown);
            if (m_CountDown <= 0)
            {
                m_UseCountDown = false;
                OnHandleButtonClick(MessageBoxResult.Close);
            }
        }
    }

    int currentCountDownSecond = 0;

    void UpdateCountDownDisplay(float m_CountDown)
    {
        int number = Mathf.CeilToInt(m_CountDown);
        if (currentCountDownSecond != number)
        {
            currentCountDownSecond = number;
            if (m_OKButtonText != null)
            {
                m_OKButtonText.text = string.Format("{0}({1})", m_OkButtonName, currentCountDownSecond);
            }
        }
    }

    /// <summary>
    /// 显示提示框
    /// </summary>
    /// <param name="uiData"></param>
    /// <param name="parent"></param>
    public static void Show(MessageBoxData uiData, WindowNode parent = null)
    {
        WindowNodeInitParam initParam = new WindowNodeInitParam(UIAssetName, AppDefine.HallNumber);        // 大厅内置窗口

        if (openedWindowNode == null)
        {
            if (parent != null)
            {
                initParam.ParentNode = parent;
            }
            else
            {
                initParam.NodeType = BaseNodeType.Propmt;
            }

            initParam.WindowData = uiData;
            openedWindowNode = WindowManager.Instance.OpenWindow(initParam);
        }
        else
        {
            openedWindowNode.WindowData = uiData;

            openedWindowNode.OpenIndex = WindowManager.Instance.GetMaxIndexNode().OpenIndex + 1;
            openedWindowNode.ShowWindow();
        }
    }

    public static void Show(MessageBoxStyle style, string title, string content, Action<MessageBoxResult> callBack = null)
    {
        Show(style, title, content, "", "", 0, callBack);
    }

    public static void Show(MessageBoxStyle style, string title, string content, string okButtonText, Action<MessageBoxResult> callBack = null)
    {
        Show(style, title, content, okButtonText, "", 0, callBack);
    }

    public static void Show(MessageBoxStyle style, string title, string content, string okButtonText, string cancelButtonText, Action<MessageBoxResult> callBack = null)
    {
        Show(style, title, content, okButtonText, cancelButtonText, 0, callBack);
    }

    public static void Show(MessageBoxStyle style, string title, string content, int lastTime, Action<MessageBoxResult> callBack = null)
    {
        Show(style, title, content, "", "", lastTime, callBack);
    }

    public static void Show(MessageBoxStyle style, string title, string content, string okButtonText, int lastTime, Action<MessageBoxResult> callBack = null)
    {
        Show(style, title, content, okButtonText, "", lastTime, callBack);
    }

    public static void Show(MessageBoxStyle style, string title, string content, string okButtonText, string cancelButtonText, int lastTime, Action<MessageBoxResult> callBack = null)
    {
        MessageBoxData boxData = new MessageBoxData()
        {
            Style = style,
            CallBack = callBack,
            Title = title,
            Content = content,
            OKButtonName = okButtonText,
            CancelButtonName = cancelButtonText,
            EnableCloseArea = false,
            LastTime = lastTime
        };

        Show(boxData, null);
    }
    public static void Hide()
    {
        if (openedWindowNode != null)
        {
            WindowManager.Instance.CloseWindow(openedWindowNode, true);
            openedWindowNode = null;
        }
    }
    public override void EscapeWindow()
    {
        OnHandleButtonClick(MessageBoxResult.Close);
    }
}

public enum MessageBoxResult : int
{
    /// <summary>
    /// 返回结果
    /// </summary>
    OK = 1,
    /// <summary>
    /// 取消
    /// </summary>
    Cancel = 2,
    /// <summary>
    /// 关闭
    /// </summary>
    Close = 3,
    /// <summary>
    /// 关闭区域关闭
    /// </summary>
    CloseArea = 3,
}

public enum MessageBoxStyle : int
{
    /// <summary>
    /// 确定
    /// </summary>
    OK = 1,
    /// <summary>
    /// 确定，取消
    /// </summary>
    OKCancel = 2,
    /// <summary>
    /// 确定，关闭
    /// </summary>
    OKClose = 3,
    /// <summary>
    /// 确定，取消，关闭
    /// </summary>
    OKCancelClose = 4,
}

public class MessageBoxData
{
    public string Title { get; set; }

    public string Content { get; set; }

    /// <summary>
    /// 持续时间（0：不使用持续时间, 大于0 在确认按钮上显示持续时间）
    /// ex:OK(2) 每秒更新
    /// </summary>
    public int LastTime { get; set; }

    /// <summary>
    /// 对话框样式
    /// </summary>
    public MessageBoxStyle Style { get; set; }

    private string m_OKButtonName = "确定";

    /// <summary>
    /// 确定按钮名称
    /// </summary>
    public string OKButtonName
    {
        get { return m_OKButtonName; }
        set
        {
            if (!string.IsNullOrEmpty(value))
                m_OKButtonName = value;
        }
    }

    private string m_CancelButtonName = "取消";
    /// <summary>
    /// 取消按钮名称
    /// </summary>
    public string CancelButtonName
    {
        get { return m_CancelButtonName; }
        set
        {
            if (!string.IsNullOrEmpty(value))
                m_CancelButtonName = value;
        }
    }

    /// <summary>
    /// 关闭区域是否可用
    /// </summary>
    public bool EnableCloseArea { get; set; }

    /// <summary>
    /// 回掉处理
    /// </summary>
    public Action<MessageBoxResult> CallBack;
}