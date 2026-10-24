
using UnityEngine;
using UnityEngine.UI;
using System.Collections;
using System.Collections.Generic;
using XLua;
using System;
using UnityEngine.SceneManagement;
using DG.Tweening;

[LuaCallCSharp]
public class LuaGameBehaviour : LuaCoroutine
{
    private Action luaStart;
    private Action luaUpdate;
    private Action luaOnDestroy;

    private Action luaEventStart;
    private Action luaEventDestory;


    private LuaTable scriptEnv;
    private LuaTable eventEnv;

    public string GameName;
    public LuaTable LuaScript
    {
        get { return scriptEnv; }
    }

    bool isLoadedLuaScript = false;

    // Add By Rayon
    // 不同解析度下要操作的攝影機，若沒指定，則抓取Canvas下的Camera
    public Camera m_Camera;
    // 攝影機在不同解析度下的清除方式
    // 若只有一層則使用CameraClearFlags.SolidColor
    // 若有多層如捕魚機，則使用CameraClearFlags.Depth
    public CameraClearFlags m_CameraClearFlags = CameraClearFlags.SolidColor;
    public float m_InitOrthographicSize = 0;

    void Awake()
    {
    }
    void Start()
    {
        // Add By Rayon
        // 不同解析度下畫面的處理
        InitSceneResolution();
    }

    public void StartGame()
    {
        BindingLuaScript();
        luaEventStart();
        CallLuaStart();
    }

    // Add By Rayon
    public void InitSceneResolution()
    {
        Canvas canvas = this.gameObject.GetComponent<Canvas>();
        if (canvas == null)
            canvas = this.gameObject.GetComponentInParent<Canvas>();

        if (canvas != null)
        {
            CanvasScaler canvasScaler = canvas.GetComponent<CanvasScaler>();
            if (canvasScaler != null)
            {
                Debug.Log("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! find canvasScaler");
                canvasScaler.screenMatchMode = CanvasScaler.ScreenMatchMode.Expand;
            }
            else
            {
                Debug.Log("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! canvasScaler is nil");
            }

            if (m_Camera == null)
                m_Camera = canvas.GetComponentInChildren<Camera>();
        }
        else
        {
            Debug.LogWarning("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! canvas is nil");
        }

        if (m_Camera != null)
        {
            Debug.Log("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! find camera");

            if (m_InitOrthographicSize == 0)
                m_InitOrthographicSize = m_Camera.orthographicSize;

            // 修改Camera ClearFlags為SolidColor，避免殘留之前的影像
            m_Camera.clearFlags = m_CameraClearFlags;
            // 修改背景顏色
            m_Camera.backgroundColor = new Color(0, 0, 0, 0);
            // 縮放特效和角色動畫
            m_Camera.orthographicSize = Camera.main.orthographicSize; //m_InitOrthographicSize * Utility.GetScreenScale(0.0f);
        }
        else
        {
            Debug.LogWarning("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! camera is nil");
        }
    }

    void CallLuaStart()
    {
        if (luaStart != null)
        {
            luaStart();
        }
    }
    public void BindingLuaScript()
    {
        if (!isLoadedLuaScript)
        {
            byte[] gameEventLuaContent = LuaManager.Instance.LoadCustomLuaFile("GameEvent", AppDefine.HallNumber);
            byte[] luaContent = LuaManager.Instance.LoadCustomLuaFile("GameStart", GameName);
            LuaEnv luaEnv = LuaManager.LuaEnv;
            scriptEnv = luaEnv.NewTable();
            eventEnv = luaEnv.NewTable();

            LuaTable meta = luaEnv.NewTable();
            meta.Set("__index", luaEnv.Global);
            scriptEnv.SetMetaTable(meta);
            eventEnv.SetMetaTable(meta);
            meta.Dispose();

            scriptEnv.Set("this", this);

            luaEnv.DoString(Utility.BytesToUTF8String(gameEventLuaContent), "GameEvent", eventEnv);
            luaEnv.DoString(Utility.BytesToUTF8String(luaContent), "GameStart", scriptEnv);

            isLoadedLuaScript = true;
        }

        eventEnv.Get("Start", out luaEventStart);
        eventEnv.Get("Destory", out luaEventDestory);

        Action luaAwake = scriptEnv.Get<Action>("Awake");
        scriptEnv.Get("Start", out luaStart);
        scriptEnv.Get("Update", out luaUpdate);
        scriptEnv.Get("OnDestroy", out luaOnDestroy);

        if (luaAwake != null)
        {
            luaAwake();
        }
    }

    // Update is called once per frame
    new void Update()
    {
        base.Update();

        if (luaUpdate != null)
        {
            luaUpdate();
        }
    }

    void OnDestroy()
    {
        if (luaEventDestory != null)
        {
            luaEventDestory();
        }
        if (luaOnDestroy != null)
        {
            luaOnDestroy();
        }
        luaOnDestroy = null;
        luaUpdate = null;
        luaEventStart = null;

        luaEventStart = null;
        luaEventDestory = null;

        if (eventEnv != null)
        {
            eventEnv.Dispose();
        }
        eventEnv = null;

        if (scriptEnv != null)
        {
            scriptEnv.Dispose();
        }
        scriptEnv = null;
    }
}
