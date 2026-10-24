using System;
using UnityEngine;

public class AppManager : Singleton<AppManager>
{
    public Action<bool> OnAppPause = (isPause) => { };
    public Action OnAppQuit = () => { };

    public void ReLaunch()
    {
        OnAppPause = (isPause) => { };
        OnAppQuit = () => { };
    }

    void OnApplicationPause(bool isPause)
    {
        if (OnAppPause != null)
            OnAppPause(isPause);
    }

    public void OnApplicationQuit()
    {
        if (OnAppQuit != null)
            OnAppQuit();
    }

    public void AddAppPause(Action<bool> action)
    {
        OnAppPause += action;
    }

    public void AddAppQuit(Action action)
    {
        OnAppQuit += action;
    }
}
