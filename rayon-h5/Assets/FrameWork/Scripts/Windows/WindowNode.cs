using UnityEngine;
using System.Collections.Generic;
using UnityEngine.UI;

public class WindowNode
{
    public WindowNodeInitParam NodeInitParam
    {
        get; private set;
    }

    private WindowBase m_WindowMonoBehaviour = null;
    /// <summary>
    /// 窗体脚本--设置窗体脚本的时候刷新脚本的 WindowNode
    /// </summary>
    public WindowBase WindowMonoBehaviour
    {
        get { return m_WindowMonoBehaviour; }
        set
        {
            if (value != null)
            {
                value.SetWindowNode(this);
            }
            m_WindowMonoBehaviour = value;
        }
    }

    /// <summary>
    /// 父窗体节点的ID
    /// </summary>
    public WindowNode ParentWindow { get; set; }

    /// <summary>
    /// 是否根节点(区分虚拟根节点)
    /// </summary>
    public bool IsRootNode { get; private set; }

    /// <summary>
    /// 根节点类型，用于挂在那个节点下
    /// </summary>
    public BaseNodeType RootNodeType { get; set; }

    public WindowNode(BaseNodeType nodeType, RectTransform rootTrans)
    {
        WindowGameObject = new GameObject();
        WindowGameObject.name = "WindowRoot_" + nodeType.ToString();
        WindowGameObject.transform.SetParent(rootTrans);
        WindowGameObject.transform.localScale = Vector3.one;
        WindowGameObject.transform.localPosition = Vector3.zero;

        WindowGameObject.layer = LayerMask.NameToLayer("UI");     // UI 层

        RectTransform rectTran = WindowGameObject.AddComponent<RectTransform>();
        rectTran.anchorMin = Vector2.zero;
        rectTran.anchorMax = Vector2.one;
        rectTran.anchoredPosition = Vector2.zero;
        rectTran.sizeDelta = Vector2.zero;

        Canvas canvas = WindowGameObject.AddComponent<Canvas>();
        canvas.overrideSorting = true;
        canvas.sortingOrder = (int)nodeType * 10 - 5;

        GraphicRaycaster gr = WindowGameObject.AddComponent<GraphicRaycaster>();
        gr.ignoreReversedGraphics = true;

        ChildWindows = new List<WindowNode>();
        WindowName = nodeType.ToString();
        RootNodeType = nodeType;
        IsRootNode = true;
    }

    public WindowNode(WindowNodeInitParam initParam)
    {
        NodeInitParam = initParam;

        ChildWindows = new List<WindowNode>();
        WindowName = initParam.WindowName;
        WindowProjectName = initParam.WindowProjectName;
        WindowAssetName = initParam.WindowAssetName;
        IsRootNode = false;
        m_WindowData = initParam.WindowData;
        ParentWindow = initParam.ParentNode;
        RootNodeType = initParam.ParentNode.RootNodeType;
        this.m_WindowLoadComplatedCallBack = initParam.LoadComplatedCallBack;
        this.m_WindowCloseCallBack = initParam.WindowCloseCallBack;
    }

    /// <summary>
    /// 窗体名称
    /// </summary>
    public string WindowName { get; private set; }
    /// <summary>
    /// 资源名称
    /// </summary>
    public string WindowAssetName { get; private set; }

    // 窗口被打开的层级顺序，用于判断接收Esc事件
    public int OpenIndex = -1;

    // 窗体所在项目的名称
    public string WindowProjectName { get; private set; }

    private object m_WindowData = null;
    /// <summary>
    /// 窗体的数据
    /// </summary>
    public object WindowData
    {
        get { return m_WindowData; }
        set
        {
            m_WindowData = value;
            RefreshWindowData(value);
        }
    }

    /// <summary>
    /// 窗体的对象
    /// </summary>
    public GameObject WindowGameObject { get; set; }

    /// <summary>
    /// 子节点
    /// </summary>
    public List<WindowNode> ChildWindows { get; private set; }

    /// <summary>
    /// 界面加载完成回调
    /// </summary>
    private System.Action<WindowNode> m_WindowLoadComplatedCallBack = null;

    /// <summary>
    /// 界面加载完成回调
    /// </summary>
    private System.Action m_WindowCloseCallBack = null;

    /// <summary>
    /// 显示窗体
    /// </summary>
    public void ShowWindow()
    {
        if (WindowGameObject != null)
        {
            WindowGameObject.SetActive(true);

            RefreshWindowData(WindowData);
            if (m_WindowLoadComplatedCallBack != null)
            {
                m_WindowLoadComplatedCallBack(this);
                m_WindowLoadComplatedCallBack = null;
            }
        }
    }

    /// <summary>
    /// 删除窗体
    /// </summary>
    public void CloseWindow()
    {
        if (WindowGameObject != null)
        {
            WindowGameObject.SetActive(false);
        }
        //清空事件列表
        m_WindowLoadComplatedCallBack = null;

        if (m_WindowCloseCallBack != null)
        {
            m_WindowCloseCallBack();
            m_WindowCloseCallBack = null;
        }
    }

    /// <summary>
    /// 设置界面的数据
    /// 界面被创建出来后和被现实出时，会调用一次此函数
    /// </summary>
    /// <param name="data">窗体的值</param>
    protected void RefreshWindowData(object data)
    {
        if (WindowMonoBehaviour != null && WindowMonoBehaviour.isActiveAndEnabled)
        {
            WindowMonoBehaviour.RefreshWindowData(data);
        }
    }

    public void EscapeWindow()
    {
        if (WindowMonoBehaviour != null && WindowMonoBehaviour.isActiveAndEnabled)
        {
            WindowMonoBehaviour.EscapeWindow();
        }
    }
}