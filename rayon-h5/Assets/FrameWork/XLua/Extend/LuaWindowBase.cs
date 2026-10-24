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
public class LuaWindowBase : WindowBase
{
    private Action luaStart;
    private Action luaUpdate;
    private Action luaOnDestroy;
    private Action<object> luaRefreshWindowData;
    private Action luaEscapeWindow;

    private bool isStarted = false;
    private LuaTable scriptEnv;

    protected void Start()
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

    public override void BindingLuaScript(string projectName)
    {
        string scriptName = gameObject.name.Replace("(Clone)", string.Empty);

        byte[] luaContent = LuaManager.Instance.LoadCustomLuaFile(scriptName, projectName);
        LuaEnv luaEnv = LuaManager.LuaEnv;
        scriptEnv = luaEnv.NewTable();
        LuaTable meta = luaEnv.NewTable();
        meta.Set("__index", luaEnv.Global);
        scriptEnv.SetMetaTable(meta);
        meta.Dispose();

        scriptEnv.Set("this", this);

        luaEnv.DoString(Utility.BytesToUTF8String(luaContent), scriptName, scriptEnv);

        Action luaAwake = scriptEnv.Get<Action>("Awake");
        scriptEnv.Get("Start", out luaStart);
        scriptEnv.Get("Update", out luaUpdate);
        scriptEnv.Get("OnDestroy", out luaOnDestroy);
        scriptEnv.Get("RefreshWindowData", out luaRefreshWindowData);
        scriptEnv.Get("EscapeWindow", out luaEscapeWindow);


        if (luaAwake != null)
        {
            luaAwake();
        }

        if (isStarted)
        {
            CallLuaStart();
        }

        base.BindingLuaScript(projectName);
    }

    void Update()
    {
        if (luaUpdate != null)
        {
            luaUpdate();
        }
    }

    void OnDestroy()
    {
        if (luaOnDestroy != null)
        {
            luaOnDestroy();
        }

        luaOnDestroy = null;
        luaUpdate = null;
        luaStart = null;

        if (scriptEnv != null)
        {
            scriptEnv.Dispose();
        }
    }

    public override void RefreshWindowData(object windowData)
    {
        base.RefreshWindowData(windowData);
        if (luaRefreshWindowData != null)
        {
            luaRefreshWindowData(windowData);
        }
    }

    public override void EscapeWindow()
    {
        base.EscapeWindow();
        if (luaEscapeWindow != null)
        {
            luaEscapeWindow();
        }
    }
}
