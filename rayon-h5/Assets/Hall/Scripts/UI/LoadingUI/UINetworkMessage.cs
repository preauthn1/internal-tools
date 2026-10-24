using UnityEngine;
using UnityEngine.UI;

public class UINetworkMessage : WindowBase
{
    public readonly static string UIAssetName = "UINetworkMessage";

    static WindowNode openedWindowNode = null;
     
    /// <summary>
    /// 内容文本
    /// </summary>
    [SerializeField]
    Text m_ContentText = null;

    public static void Show(string tipsContent)
    {
        if (openedWindowNode == null)
        {
            WindowNodeInitParam initParam = new WindowNodeInitParam(UIAssetName, AppDefine.HallNumber);       // 大厅内置窗口
            initParam.NodeType = BaseNodeType.Propmt;
            initParam.WindowData = tipsContent;
            openedWindowNode = WindowManager.Instance.OpenWindow(initParam);
        }
        else
        {
            openedWindowNode.WindowData = tipsContent;
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
            WindowManager.Instance.CloseWindow(openedWindowNode, true);
            openedWindowNode = null;
        }
    }

    public override void RefreshWindowData(object windowData)
    {
        base.RefreshWindowData(windowData); 
        this.m_ContentText.text = (string)windowData;
    }
}