using HotUpdate;
using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.SceneManagement;

/// <summary>
/// 游戏管理器
/// </summary>
public class GameManager : Singleton<GameManager>
{
    public const string SCENENAME = "Game";

    bool m_isCutOut = false;

    public Action<string, Action> EnterGame;
    public Action<Action> QuitGame;
    public Action<Action> ReLoadGame;

    private Scene hallScene;

    /// <summary>
    /// 去掉游戏休眠
    /// </summary>
    /// <param name="focus"></param>
    void OnApplicationFocus(bool focus)
    {
        Screen.sleepTimeout = focus ? SleepTimeout.NeverSleep : SleepTimeout.SystemSetting;
    }

    void Start()
    {
        OnApplicationFocus(true);       // 开始游戏运行到此时，游戏已调用OnApplicationFocus，此处补充调用一次
        EnterGame = EnterGameFun;
        QuitGame = QuitGameFun;
        ReLoadGame = ReloadGameScene;
        hallScene = SceneManager.GetActiveScene();
    }

    private string m_CurrentEnterGameName;  // 当前载入的游戏项目名
    public string GetCurrentEnterGameName()
    {
        return m_CurrentEnterGameName;
    }

    // Add By Rayon
    private LuaGameBehaviour m_behaviour;   // 目前場景的LuaGameBehaviour

    public void StartGame()
    {
        LuaGameBehaviour behaviour = GameObject.FindObjectOfType<LuaGameBehaviour>();

        // Add By Rayon
        m_behaviour = behaviour;

        if (!behaviour)
            Debug.Log("找不到 LuaGameBehaviour");
        else
            behaviour.StartGame();
    }

    // Add By Rayon
    public Camera GetCamera()
    {
        if (m_behaviour != null)
            return m_behaviour.m_Camera;
        else
            return Camera.main;
    }

    private Queue<IEnumerator> m_OprationQueue = new Queue<IEnumerator>();
    private bool m_IsOprating = false;


    void Update()
    {
        if (!m_IsOprating)
        {
            if (m_OprationQueue.Count != 0)
            {
                IEnumerator ienum = m_OprationQueue.Dequeue();
                StartCoroutine(MainOpration(ienum));
            }
        }
    }

    IEnumerator MainOpration(IEnumerator ienum)
    {
        m_IsOprating = true;
        yield return ienum;
        m_IsOprating = false;
    }

    void EnterGameFun(string projectName, Action callBack = null)
    {
        m_OprationQueue.Enqueue(EnterGameYid(projectName, callBack));
    }

    void ReloadGameScene(Action callBack)
    {
        m_OprationQueue.Enqueue(ReLoadGameScene(callBack));
    }

    void QuitGameFun(Action callBack)
    {
        m_OprationQueue.Enqueue(QuitGameYid(callBack));
    }

    private string m_projectName;
    IEnumerator EnterGameYid(string projectName, Action callBack = null)
    {
        if (!string.IsNullOrEmpty(m_CurrentEnterGameName))
        {
            AssetBundleManager.Instance.unRegistProjectAssetBundleManager(m_CurrentEnterGameName);
            AsyncOperation unloadAsync = SceneManager.UnloadSceneAsync(SCENENAME);
            while (!unloadAsync.isDone)
                yield return 1;

            m_CurrentEnterGameName = string.Empty;
        }

        m_projectName = projectName;

        Debug.Log("框架：注册项目资源" + projectName);

        // 要使用游戏的资源必须注册
        AssetBundleManager.Instance.RegistProjectAssetBundleManager(projectName);

        Debug.Log("框架：加载场景");
        DateTime Time = DateTime.Now;

        AsyncOperation async = AssetBundleManager.Instance.LoadSceneAsync(projectName, SCENENAME);

#if UNITY_EDITOR

        if (AssetBundleOption.SimulateAssetBundleInEditor && async == null)
            Debug.LogError("在模拟AssetBundle时，需要将场景添加到File->Build Settings...中");
#endif

        while (!async.isDone)
            yield return 1;

        Debug.Log("加载场景花費時間" + (DateTime.Now - Time).TotalMilliseconds);

        m_CurrentEnterGameName = projectName;

        Debug.Log("框架：设置当前Acative场景");
        SceneManager.SetActiveScene(SceneManager.GetSceneByName(SCENENAME));

        Debug.Log("框架：加载场景回调");
        //只执行最后一个加载场景的回调
        if (callBack != null && m_OprationQueue.Count <= 0)
            callBack();
    }

    IEnumerator ReLoadGameScene(Action callBack = null)
    {
        Debug.Log("开始重新加载场景");

        AsyncOperation async = null;

        // By Rayon
        // 以完整方式載入場景
        // 有多个名称相同但是路径不同的场景，你应该使用全部路径
#if UNITY_EDITOR
        if (AssetBundleOption.SimulateAssetBundleInEditor)
        {
            if (string.IsNullOrEmpty(m_projectName) == false)
                async = SceneManager.LoadSceneAsync("Game/" + m_projectName + "/" + SCENENAME, LoadSceneMode.Additive);
        }
#endif

        if (async == null)
            async = SceneManager.LoadSceneAsync(SCENENAME, LoadSceneMode.Additive);

        while (!async.isDone)
            yield return 1;

        Scene newLoadedScene = SceneManager.GetSceneAt(SceneManager.sceneCount - 1);        

        // 先隐藏新场景中内容，等待卸载旧场景完成再显示
        List<GameObject> newSceneObj = new List<GameObject>();
        newLoadedScene.GetRootGameObjects(newSceneObj);
        newSceneObj = newSceneObj.FindAll(o => { return o.gameObject.activeSelf; });
        newSceneObj.ForEach(o => { o.SetActive(false); });

        Scene activeScene = SceneManager.GetActiveScene();

        Debug.Log("开始卸载旧场景");
        async = SceneManager.UnloadSceneAsync(activeScene);
        while (!async.isDone)
            yield return 1;

        Debug.Log("场景切换完成");
        SceneManager.SetActiveScene(newLoadedScene);

        newSceneObj.ForEach(o => { o.SetActive(true); });
        //只执行最后一个加载场景的回调
        if (callBack != null && m_OprationQueue.Count <= 0)
            callBack();
    }

    IEnumerator QuitGameYid(Action callBack)
    {
        if (!string.IsNullOrEmpty(m_CurrentEnterGameName))
        {
            LuaManager.Instance.ClearProjectCacheLuaScript(m_CurrentEnterGameName);
            AssetBundleManager.Instance.unRegistProjectAssetBundleManager(m_CurrentEnterGameName);

            AsyncOperation async = SceneManager.UnloadSceneAsync(SCENENAME);
            while (!async.isDone)
                yield return 1;

            m_CurrentEnterGameName = string.Empty;
        }

        SceneManager.SetActiveScene(hallScene);

        if (callBack != null)
            callBack();
    }
}