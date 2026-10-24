using UnityEngine;
using System.Collections;
using System;
using System.Collections.Generic;

public class SoundManager : Singleton<SoundManager>
{
    private AudioSource audioSource;
    private AudioSource bgAudioSource;

    private AudioListener audioListener;

    private Dictionary<string, AudioClip> soundsDic = new Dictionary<string, AudioClip>();

    // ǷЧ   Ч
    public bool IsSoundOn
    {
        get
        {
            return PlayerPrefs.GetInt("IsSoundOn", 1) == 1;
        }
        set
        {
            PlayerPrefs.SetInt("IsSoundOn", value ? 1 : 0);
        }
    }

    // Ƿ   
    public bool IsMusicOn
    {
        get
        {
            return PlayerPrefs.GetInt("IsMusicOn", 1) == 1;
        }
        set
        {
            bgAudioSource.enabled = value;
            PlayerPrefs.SetInt("IsMusicOn", value ? 1 : 0);
        }
    }

    void Start()
    {
        audioSource = gameObject.AddComponent<AudioSource>();
        bgAudioSource = gameObject.AddComponent<AudioSource>();
        bgAudioSource.loop = true;

        audioListener = gameObject.AddComponent<AudioListener>();
    }

    /// <summary>
    /// ű
    /// </summary>
    /// <param name="projectName"></param>
    /// <param name="audioClipName"></param>
    public void PlayBGM(string projectName, string audioClipName)
    {
        bgAudioSource.enabled = IsMusicOn;

        LoadAudioClip(projectName, audioClipName, (audioClip) =>
        {
            bgAudioSource.clip = audioClip;
            bgAudioSource.Play();
        });
    }

    Dictionary<string, AudioSource> m_loopAudioDic = new Dictionary<string, AudioSource>();
    /// <summary>
    /// һ
    /// </summary>
    /// <param name="projectName"></param>
    /// <param name="audioClipName"></param>
    public void PlaySound(string projectName, string audioClipName, bool isLoop = false)
    {
        if (!IsSoundOn)
            return;

        LoadAudioClip(projectName, audioClipName, (audioClip) =>
        {
            if (isLoop && !m_loopAudioDic.ContainsKey(audioClipName))
            {
                AudioSource loopSource = gameObject.AddComponent<AudioSource>();
                loopSource.loop = true;

                loopSource.clip = audioClip;
                loopSource.Play();

                m_loopAudioDic.Add(audioClipName, loopSource);
            }
            else
            {
                audioSource.PlayOneShot(audioClip);
            }
        });
    }

    public void PlaySound(AudioClip audioClip)
    {
        if (!IsSoundOn)
            return;

        audioSource.PlayOneShot(audioClip);
    }

    public void CloseLoopSound(string audioClipName)
    {
        if (m_loopAudioDic.ContainsKey(audioClipName))
        {
            m_loopAudioDic[audioClipName].Stop();
            Destroy(m_loopAudioDic[audioClipName]);
            m_loopAudioDic.Remove(audioClipName);
        }
    }

    public void CloseSound()
    {
        audioSource.Stop();

        foreach (var item in m_loopAudioDic)
            Destroy(item.Value);

        soundsDic.Clear();
        m_loopAudioDic.Clear();
    }

    public void CloseBGMSound()
    {
        bgAudioSource.clip = null;
        bgAudioSource.Stop();
    }

    /// <summary>
    /// һƵ
    /// </summary>
    void LoadAudioClip(string projectName, string audioClipName, Action<AudioClip> callBack)
    {
        StartCoroutine(LoadAudioClipYid(projectName, audioClipName, callBack));
    }

    IEnumerator LoadAudioClipYid(string projectName, string audioClipName, Action<AudioClip> callBack)
    {
        AudioClip audioClip = null;

        string hashKey = projectName + "." + audioClipName; // ĿЧļΪHashKeyֵ

        if (soundsDic.ContainsKey(hashKey) && soundsDic[hashKey] != null)
            audioClip = soundsDic[hashKey] as AudioClip;
        else
        {
            LoadAssetAsyncOperation operation = AssetBundleManager.Instance.LoadAssetAsync<AudioClip>(projectName, audioClipName, false);

            if (operation != null && !operation.IsDone)
                yield return operation;

            if (operation != null)
            {
                audioClip = operation.GetAsset<AudioClip>();

                if (audioClip)
                {
                    if (!soundsDic.ContainsKey(hashKey))
                        soundsDic.Add(hashKey, audioClip);
                    else
                        soundsDic[hashKey] = audioClip;
                }
            }
        }

        callBack(audioClip);
    }

    private void ClearMemory()
    {
        GC.Collect();
        Resources.UnloadUnusedAssets();
    }
}