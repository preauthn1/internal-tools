using UnityEngine;
using System.Collections;
using System.Collections.Generic;

public enum BaseNodeType
{
    /// <summary>
    /// 登陆，大厅，游戏等ui用此等级
    /// </summary>
    Main = 0,
    /// <summary>
    /// 跑马灯专用
    /// </summary>
    AboveMain = 1,
    /// <summary>
    /// 普通的界面，如Hall UI，Game UI，商城UI等等
    /// </summary>
    Normal = 2,
    /// <summary>
    /// Normal 界面层之上的一层
    /// </summary>
    AboveNormal = 3,
    /// <summary>
    /// 界面浮出提示信息
    /// </summary>
    Tips = 4,
    /// <summary>
    /// 进度条界面
    /// </summary>
    Loading = 5,
    /// <summary>
    /// 错误提示界面，，如断网提示等
    /// </summary>
    Propmt = 6,
}

public class WindowManager : MonoBehaviour
{
    [Tooltip("是否跨场景销毁")]
    public bool makePersistent = true;

    #region 常量

    public const string UIRootAsset = "[UIRoot]";

    public static void InitManager()
    {
        if (m_Instance == null)
        {
            GameObject go = Resources.Load<GameObject>(UIRootAsset);
            GameObject uiroot = Instantiate(go);
            uiroot.name = UIRootAsset;
            Utility.ReSetTransform(uiroot.transform, null);
        }
    }

    #endregion

    #region 定义根节点

    // 所有的界面挂的根节点
    [SerializeField]
    RectTransform m_UIParentTransform = null;

    // 隐藏界面的挂结点
    [SerializeField]
    Transform m_HideNodeTransform = null;

    // 根节点信息
    WindowNode[] m_RootWindowNodes = null;

    #endregion

    #region Manager Private Data

    /// <summary>
    /// 所有的窗体界面信息界面
    /// </summary>
    List<WindowNode> m_WindowNodeList = new List<WindowNode>();

    /// <summary>
    /// 缓存资源界面池
    /// </summary>
    WindowPool m_WindowPool = new WindowPool();

    #endregion

    private static WindowManager m_Instance = null;

    #region Engine Callbacks

    public void ClearWindowPool()
    {
        m_WindowPool.ClearPool();
    }

    void Awake()
    {
        CheckInstance();

        int max = (int)BaseNodeType.Propmt;
        m_RootWindowNodes = new WindowNode[max + 1];
        for (int i = 0; i <= max; i++)
        {
            WindowNode node = new WindowNode((BaseNodeType)i, m_UIParentTransform);
            m_WindowNodeList.Add(node);
            m_RootWindowNodes[i] = node;
        }
    }

    void Start()
    {
        if (makePersistent)
            DontDestroyOnLoad(this.gameObject);

        // Add By Rayon
        // 設置 Camera orthographicSize 會影響不同解析度下特效和角色動畫大小
        if (m_Instance != null)
            m_Instance.InitCamera();
    }

    // this is called after Awake() OR after the script is recompiled (Recompile > Disable > Enable)
    void OnEnable()
    {
        CheckInstance();
    }

    bool CheckInstance()
    {
        if (m_Instance == null)
        {
            m_Instance = this;
        }
        else if (m_Instance != this)
        {
            Debug.LogWarning("There is already an instance of WindowManager created (" + m_Instance.name + "). Destroying new one.");
            Destroy(this.gameObject);
            return false;
        }
        return true;
    }

    #endregion

    public static WindowManager Instance
    {
        get
        {
            if (m_Instance == null)
            {
                Debug.LogError("The instance was null of WindowManager, When you get it instance, please check it!");
            }
            return m_Instance;
        }
    }

    #region OpenWindow

    public float m_InitOrthographicSize;
    public Camera m_Camera;

    private void InitCamera()
    {
        m_Camera = this.gameObject.GetComponentInChildren<Camera>();
        if (m_Camera != null)
        {
            if (m_InitOrthographicSize == 0)
                m_InitOrthographicSize = m_Camera.orthographicSize;

            m_Camera.orthographicSize = m_InitOrthographicSize * Utility.GetScreenScale(0.0f);
        }
    }

    public void OpenGameWindow(string windowName, string projectName, Transform parent, object initParas)
    {
        var window = GameObject.Instantiate<GameObject>(AssetBundleManager.Instance.LoadAsset<GameObject>(projectName, windowName));
        window.transform.SetParent(parent);
        window.transform.localScale = Vector3.one;
        window.transform.localPosition = Vector3.zero;

        var luaWB = window.GetComponent<LuaWindowBase>();

        if (!luaWB)
            luaWB = window.AddComponent<LuaWindowBase>();

        luaWB.BindingLuaScript(projectName);
        luaWB.RefreshWindowData(initParas);
    }


    public WindowNode OpenWindow(string windowName, string projectName, BaseNodeType nodeType = BaseNodeType.Normal)
    {
        return OpenWindow(new WindowNodeInitParam(windowName, projectName) { NodeType = nodeType });
    }

    public WindowNode OpenWindow(string windowName, string projectName, WindowNode parentNode)
    {
        return OpenWindow(new WindowNodeInitParam(windowName, projectName) { ParentNode = parentNode });
    }

    public WindowNode OpenWindow(WindowNodeInitParam initParam)
    {
        if (initParam == null)
            return null;
        if (initParam.ParentNode == null)
        {
            initParam.ParentNode = GetBaseWindowNodeByNodeType(initParam.NodeType);
        }
        WindowNode windowNode = new WindowNode(initParam);
        windowNode.ParentWindow = initParam.ParentNode;
        SetWindowNodeInfo(windowNode, initParam.NearNode, initParam.NearNodeIsPreNode);

        return windowNode;
    }

    public void CloseWindowByProjectName(string projectName)
    {
        if (projectName == null)
            return;

        // 统计SDK
        List<WindowNode> findNodes = m_WindowNodeList.FindAll(temp => temp.WindowProjectName == projectName);

        if (findNodes == null || findNodes.Count == 0)
            return;

        for (int i = 0; i < findNodes.Count; i++)
            RemoveWindowNode(findNodes[0], false);
    }

    void SetWindowNodeInfo(WindowNode windowNode, WindowNode nearNode, bool nearNodeIsPreNode)
    {
        WindowNode preNearNode = null;
        // 附近节点为空或者父节点不包括此附近节点
        if (nearNode == null || !windowNode.ParentWindow.ChildWindows.Contains(nearNode))
        {
            preNearNode = windowNode.ParentWindow.ChildWindows.Count == 0 ? windowNode.ParentWindow : windowNode.ParentWindow.ChildWindows[windowNode.ParentWindow.ChildWindows.Count - 1];
            while (preNearNode.ChildWindows.Count > 0)
            {
                preNearNode = preNearNode.ChildWindows[preNearNode.ChildWindows.Count - 1];
            }
            windowNode.ParentWindow.ChildWindows.Add(windowNode);
        }
        else
        {
            int index = windowNode.ParentWindow.ChildWindows.IndexOf(nearNode);
            if (!nearNodeIsPreNode)// 新节点在 靠近节点的后面
            {
                preNearNode = nearNode;
                index += 1;
            }
            else // 新的节点在靠近节点的前面
            {
                preNearNode = index == 0 ? windowNode.ParentWindow : windowNode.ParentWindow.ChildWindows[index - 1];
            }

            while (preNearNode.ChildWindows.Count > 0)
            {
                preNearNode = preNearNode.ChildWindows[preNearNode.ChildWindows.Count - 1];
            }
            windowNode.ParentWindow.ChildWindows.Insert(index, windowNode);
        }
        m_WindowNodeList.Insert(m_WindowNodeList.IndexOf(preNearNode) + 1, windowNode);
        SetWindowNodeGameObject(windowNode);
    }

    WindowNode GetBaseWindowNodeByNodeType(BaseNodeType nodeType)
    {
        int value = (int)nodeType;
        if (value < m_RootWindowNodes.Length && value > -1)
            return m_RootWindowNodes[value];
        else
            return null;
    }

    void SetWindowNodeGameObject(WindowNode windowNode)
    {
        GameObject poolGameObject = m_WindowPool.GetWindowGameObject(windowNode.WindowName);
        if (poolGameObject != null)
        {
            SetWindowNodeGameObject(windowNode, poolGameObject, false);
        }
        else
        {
            StartCoroutine(LoadWindowGameObjectByWindowNode(windowNode));
        }
    }

    IEnumerator LoadWindowGameObjectByWindowNode(WindowNode windowNode)
    {
        AssetBundleManager assetBundleManager = AssetBundleManager.Instance;
        if (assetBundleManager == null)
        {
            Debug.LogErrorFormat("Load window asset[{0}] error[AssetBundleManager is null when window been loading]. windowNode[ID:{1}] be clear!", windowNode.WindowAssetName, windowNode.WindowName);
            CloseWindow(windowNode.WindowName, false, false);
            yield break;
        }
        LoadAssetAsyncOperation operation = assetBundleManager.LoadAssetAsync<GameObject>(windowNode.WindowProjectName, windowNode.WindowAssetName, false);
        if (operation != null && !operation.IsDone)
        {
            yield return operation;
        }
        if (m_WindowNodeList.Contains(windowNode))
        {
            if (operation != null && operation.IsDone)
            {
                UnityEngine.Object resource = operation.GetAsset<UnityEngine.Object>();
                if (resource != null)
                {
                    GameObject windowGameObject = Instantiate(resource) as GameObject;
                    if (windowGameObject != null)
                    {
                        SetWindowNodeGameObject(windowNode, windowGameObject, true);
                    }
                    else
                    {
                        Debug.LogErrorFormat("Load window asset[{0}] error[Asset was not a gameObject]. windowNode[ID:{1}] be clear!", windowNode.WindowAssetName, windowNode.WindowName);
                        CloseWindow(windowNode.WindowName, false, false);
                    }
                }
                else
                {
                    Debug.LogErrorFormat("Load window asset[{0}] error[Asset was null]. windowNode[ID:{1}] be clear!", windowNode.WindowAssetName, windowNode.WindowName);
                    CloseWindow(windowNode.WindowName, false, false);
                }
            }
            else
            {
                Debug.LogErrorFormat("Load window asset[{0}] error. windowNode[ID:{1}] be clear!", windowNode.WindowAssetName, windowNode.WindowName);
                CloseWindow(windowNode.WindowName, false, false);
            }
        }
        else
        {
            // 在加载的过程中，本窗体已经被删除掉了
        }
    }

    private void SetWindowNodeGameObject(WindowNode windowNode, GameObject windowGameObject, bool needBindingLua)
    {
        Transform windowTransform = windowGameObject.transform;
        windowTransform.SetParent(windowNode.ParentWindow.WindowGameObject.transform, false);
        windowTransform.name = windowNode.WindowName;
        windowNode.WindowGameObject = windowGameObject;
        //SetTransformSiblingIndex(windowNode, windowTransform);
        if (!windowTransform.GetComponent<WindowBase>())
            windowNode.WindowMonoBehaviour = windowTransform.gameObject.AddComponent<LuaWindowBase>();

        windowNode.WindowMonoBehaviour = windowTransform.GetComponent<WindowBase>();

        if (needBindingLua)
            windowNode.WindowMonoBehaviour.BindingLuaScript(windowNode.WindowProjectName);

        windowNode.OpenIndex = GetMaxIndexNode().OpenIndex + 1;
        windowNode.ShowWindow();
    }

    void SetTransformSiblingIndex(WindowNode windowNode, Transform transform)
    {
        int index = 0;
        for (int i = 0; i < m_WindowNodeList.Count; i++)
        {
            if (m_WindowNodeList[i].IsRootNode)
                continue;
            if (m_WindowNodeList[i] == windowNode)
            {
                break;
            }
            index++;
        }
        transform.SetSiblingIndex(index);
    }

    #endregion

    #region CloseWindow

    public void CloseWindow(string windowName, bool isReleaseToPool, bool isOpenParentWindow = false)
    {
        // 统计SDK
        WindowNode findNode = m_WindowNodeList.Find(temp => temp.WindowName == windowName);
        if (findNode == null)
            return;

        WindowNode parentNode = findNode.ParentWindow;

        RemoveWindowNode(findNode, isReleaseToPool);

        if (isOpenParentWindow && parentNode != null)
            OpenWindow(parentNode.NodeInitParam);
    }

    /// <summary>
    /// 关闭窗体
    /// </summary>
    /// <param name="windowNode">窗体节点信息</param>
    /// <param name="isRemoveParentNode">是否关闭父节点</param>
    public void CloseWindow(WindowNode windowNode, bool isReleaseToPool = true, bool isRemoveParentNode = false)
    {
        if (windowNode == null)
            return;

        CloseWindow(windowNode.WindowName, isReleaseToPool, isRemoveParentNode);
    }

    /// <summary>
    /// 清理Window Node 信息
    /// </summary>
    /// <param name="windowNode">窗体节点信息</param>
    void RemoveWindowNode(WindowNode windowNode, bool isReleaseToPool)
    {
        for (int index = windowNode.ChildWindows.Count - 1; index >= 0; index--)
        {
            RemoveWindowNode(windowNode.ChildWindows[index], isReleaseToPool);
        }

        windowNode.CloseWindow();
        if (windowNode.WindowGameObject != null)
        {
            if (isReleaseToPool)
            {
                windowNode.WindowGameObject.transform.SetParent(m_HideNodeTransform, false);
                m_WindowPool.ReleaseGameObject(windowNode.WindowName, windowNode.WindowGameObject);
            }
            else
            {
                Destroy(windowNode.WindowGameObject);
            }
        }

        windowNode.ParentWindow.ChildWindows.Remove(windowNode);// 从父节点删除当前结点
        windowNode.ParentWindow = null;
        m_WindowNodeList.Remove(windowNode);
        windowNode.ChildWindows.Clear();
        windowNode.WindowData = null;
        windowNode.WindowGameObject = null;
        windowNode = null;
    }

    /// <summary>
    /// 清除掉所有打开的界面
    /// </summary>
    /// <param name="unCloseWindowNames">不关闭的窗体</param>
    public void CloseAllWindows(List<string> unCloseWindowNames = null)
    {
        for (int i = 0; i < m_RootWindowNodes.Length; i++)
        {
            for (int j = m_RootWindowNodes[i].ChildWindows.Count - 1; j >= 0; j--)
            {
                if (unCloseWindowNames != null && unCloseWindowNames.Contains(m_RootWindowNodes[i].ChildWindows[j].WindowAssetName))
                {
                    continue;
                }
                RemoveWindowNode(m_RootWindowNodes[i].ChildWindows[j], false);
            }
        }
    }

    #endregion

    #region FindWindowNode

    /// <summary>
    /// 通过名称查找WindowNode
    /// </summary>
    /// <param name="windowName">窗口名称</param>
    /// <returns>返回同名窗体列表中，最后一个窗体ID的WindowNode，如果未找到则返回Null</returns>
    public WindowNode FindWindowNodeByName(string windowName)
    {
        return m_WindowNodeList.Find(temp => temp.WindowName == windowName);
    }

    public T FindWindowBaseByType<T>() where T : WindowBase
    {
        WindowNode node = m_WindowNodeList.Find(temp => temp.WindowMonoBehaviour is T);

        if (node != null)
            return node.WindowMonoBehaviour as T;
        else
            return null;
    }

    #endregion

    #region WindowPool

    public class WindowPool
    {
        public WindowPool()
        {
        }

        private Dictionary<string, GameObject> resourceDict = new Dictionary<string, GameObject>();

        public GameObject GetWindowGameObject(string windowName)
        {
            GameObject windowGameObject = null;
            if (resourceDict.ContainsKey(windowName))
            {
                windowGameObject = resourceDict[windowName];
                resourceDict.Remove(windowName);
            }
            return windowGameObject;
        }

        public bool ReleaseGameObject(string windowName, GameObject saveGameObject)
        {
            if (saveGameObject == null || string.IsNullOrEmpty(windowName)) return false;

            if (resourceDict.ContainsKey(windowName))
            {
                GameObject.Destroy(saveGameObject);
                return false;
            }
            else
            {
                resourceDict.Add(windowName, saveGameObject);
                return true;
            }
        }

        public void ClearPool()
        {
            foreach (var item in resourceDict)
                GameObject.Destroy(item.Value);

            resourceDict.Clear();
        }
    }

    #endregion

    public void Update()
    {
        if (Input.GetKeyDown(KeyCode.Escape))
            GetMaxIndexNode().EscapeWindow();
    }

    public WindowNode GetMaxIndexNode()
    {
        WindowNode rtnVal = null;
        m_WindowNodeList.ForEach((o) =>
        {
            if (rtnVal == null || (o.NodeInitParam != null && o.NodeInitParam.NodeType != BaseNodeType.Tips && o.OpenIndex > rtnVal.OpenIndex))
                rtnVal = o;
        });

        return rtnVal;
    }
}

/// <summary>
/// 窗体初始化参数
/// </summary>
public class WindowNodeInitParam
{
    /// <summary>
    /// 初始化
    /// </summary>
    /// <param name="windowName">窗体名称</param>
    public WindowNodeInitParam(string windowName, string projectName)
    {
        this.WindowAssetName = this.WindowName = windowName;
        this.WindowProjectName = projectName;
        NodeType = BaseNodeType.Normal;
        WindowData = null;
        LoadComplatedCallBack = null;
        ParentNode = null;
        NearNode = null;
        NearNodeIsPreNode = true;
    }
    /// <summary>
    /// 窗体名称
    /// </summary>
    public string WindowName { get; private set; }

    public string WindowProjectName { get; private set; }

    /// <summary>
    /// 窗体资源名称
    /// </summary>
    public string WindowAssetName { get; private set; }
    /// <summary>
    /// 所在的根节点类型
    /// </summary>
    public BaseNodeType NodeType { get; set; }
    /// <summary>
    /// 窗体数据
    /// </summary>
    public object WindowData { get; set; }
    /// <summary>
    /// 加载完成回调
    /// </summary>
    public System.Action<WindowNode> LoadComplatedCallBack { get; set; }

    /// <summary>
    /// 窗口关闭回调
    /// </summary>
    public System.Action WindowCloseCallBack { get; set; }

    /// <summary>
    /// 父节点
    /// </summary>
    public WindowNode ParentNode { get; set; }
    /// <summary>
    /// 靠近的节点
    /// </summary>
    public WindowNode NearNode { get; set; }
    /// <summary>
    /// 靠近的节点是否在当前结点之前
    /// </summary>
    public bool NearNodeIsPreNode { get; set; }
}
