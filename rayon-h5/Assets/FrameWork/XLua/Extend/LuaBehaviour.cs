/*
 * Tencent is pleased to support the open source community by making xLua available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://example.invalid/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

using UnityEngine;
using System.Collections;
using System.Collections.Generic;
using XLua;
using System;

[LuaCallCSharp]
public class LuaBehaviour : MonoBehaviour
{
    private Action luaStart;
    private Action luaUpdate;
    private Action<Collider2D> luaOnTriggerEnter2D;
    private Action<Collider2D> luaOnTriggerStay2D;
    private Action<Collider2D> luaOnTriggerExit2D;
    private Action<Collision2D> luaOnCollisionEnter2D;
    private Action<Collision2D> luaOnCollisionStay2D;
    private Action<Collision2D> luaOnCollisionExit2D;
    private Action luaOnDestroy;

    private bool isStarted = false;
    private LuaTable scriptEnv;

    public LuaTable LuaScript
    {
        get { return scriptEnv; }
    }

    bool isLoadedLuaScript = false;

    [HideInInspector]
    public string ProjectName;

    // [HideInInspector]
    public string ScriptName;

    public void BindingLuaScript(string projectName, string scriptName)
    {
        ProjectName = projectName;
        ScriptName = scriptName;

        if (string.IsNullOrEmpty(scriptName))
            scriptName = gameObject.name.Replace("(Clone)", string.Empty);

        if (!isLoadedLuaScript)
        {
            if (string.IsNullOrEmpty(scriptName))
                return;

            LuaEnv luaEnv = LuaManager.LuaEnv;
            scriptEnv = luaEnv.NewTable();
            LuaTable meta = luaEnv.NewTable();
            meta.Set("__index", luaEnv.Global);
            scriptEnv.SetMetaTable(meta);
            meta.Dispose();

            scriptEnv.Set("this", this);

            string luaString = LuaManager.Instance.LoadCustomLuaFileString(scriptName, projectName);
            luaEnv.DoString(luaString, scriptName, scriptEnv);
            isLoadedLuaScript = true;
        }

        Action luaAwake = scriptEnv.Get<Action>("Awake");

        scriptEnv.Get("Start", out luaStart);
        scriptEnv.Get("Update", out luaUpdate);
        scriptEnv.Get("OnTriggerEnter2D", out luaOnTriggerEnter2D);
        scriptEnv.Get("OnTriggerStay2D", out luaOnTriggerStay2D);
        scriptEnv.Get("OnTriggerExit2D", out luaOnTriggerExit2D);
        scriptEnv.Get("OnCollisionEnter2D", out luaOnCollisionEnter2D);
        scriptEnv.Get("OnCollisionStay2D", out luaOnCollisionStay2D);
        scriptEnv.Get("OnCollisionExit2D", out luaOnCollisionExit2D);
        scriptEnv.Get("OnDestroy", out luaOnDestroy);

        if (luaAwake != null)
        {
            luaAwake();
        }

        if (isStarted)
        {
            CallLuaStart();
        }
    }

    void Start()
    {
        isStarted = true;
        CallLuaStart();
    }

    void CallLuaStart()
    {
        if (luaStart != null)
        {
            luaStart();
        }
    }

    // Update is called once per frame
    void Update()
    {
        if (luaUpdate != null)
        {
            luaUpdate();
        }
    }

    public void OnTriggerEnter2D(Collider2D collision)
    {
        if (luaOnTriggerEnter2D != null)
        {
            luaOnTriggerEnter2D(collision);
        }
    }

    public void OnTriggerStay2D(Collider2D collision)
    {
        if (luaOnTriggerStay2D != null)
        {
            luaOnTriggerStay2D(collision);
        }
    }

    public void OnTriggerExit2D(Collider2D collision)
    {
        if (luaOnTriggerExit2D != null)
        {
            luaOnTriggerExit2D(collision);
        }
    }

    public void OnCollisionEnter2D(Collision2D collision)
    {
        if (luaOnCollisionEnter2D != null)
        {
            luaOnCollisionEnter2D(collision);
        }
    }

    public void OnCollisionStay2D(Collision2D collision)
    {
        if (luaOnCollisionStay2D != null)
        {
            luaOnCollisionStay2D(collision);
        }
    }

    public void OnCollisionExit2D(Collision2D collision)
    {
        if (luaOnCollisionExit2D != null)
        {
            luaOnCollisionExit2D(collision);
        }
    }

    void OnDestroy()
    {
        if (luaOnDestroy != null)
        {
            luaOnDestroy();
        }

        luaStart = null;
        luaUpdate = null;
        luaOnTriggerEnter2D = null;
        luaOnTriggerStay2D = null;
        luaOnTriggerExit2D = null;
        luaOnCollisionEnter2D = null;
        luaOnCollisionStay2D = null;
        luaOnCollisionExit2D = null;
        luaOnDestroy = null;

        if (scriptEnv != null)
        {
            scriptEnv.Dispose();
        }
        scriptEnv = null;
    }
}
