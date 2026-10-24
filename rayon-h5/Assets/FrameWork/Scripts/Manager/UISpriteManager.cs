using UnityEngine;
using System.Collections;
using System.Collections.Generic;

public class UISpriteManager : Singleton<UISpriteManager>
{
    /// <summary>
    /// 设置动态设置图片的 Sprite
    /// </summary>
    /// <param name="image">图片控件</param>
    /// <param name="spriteName">图片名称</param>
    /// <returns></returns>
    public void SetImageSprite(UnityEngine.UI.Image image, string projcetName, string spriteName, bool isNativeSize = false)
    {
        spriteName = "sprite_" + spriteName;
        GameObject spriteGo = AssetBundleManager.Instance.LoadAsset<GameObject>(projcetName, spriteName, unloadAssetBundle: false);
        if (spriteGo != null)
        {
            SpriteRenderer spriteRenderer = spriteGo.GetComponent<SpriteRenderer>();
            if (spriteRenderer != null)
            {
                if (spriteRenderer.sprite != null)
                {
                    image.sprite = spriteRenderer.sprite;
                    if (isNativeSize)
                    {
                        image.SetNativeSize();
                    }
                }
            }
        }
    }

    /// <summary>
    /// 异步设置图片
    /// </summary>
    /// <param name="image"></param>
    /// <param name="spriteName"></param>
    /// <param name="isNativeSize"></param>
    public void SetImageSpriteAsync(UnityEngine.UI.Image image, string projcetName, string spriteName, bool isNativeSize = false)
    {
        spriteName = "sprite_" + spriteName;
        StartCoroutine(SetImageSpriteCoroutine(image, projcetName, spriteName, isNativeSize));
    }

    IEnumerator SetImageSpriteCoroutine(UnityEngine.UI.Image image, string projcetName, string spriteName, bool isNativeSize)
    {
        LoadAssetAsyncOperation operation = AssetBundleManager.Instance.LoadAssetAsync<GameObject>(projcetName, spriteName, false);
        if (operation != null && !operation.IsDone)
        {
            yield return operation;
        }

        if (operation != null)
        {
            if (image != null)
            {
                GameObject spriteGo = operation.GetAsset<GameObject>();
                if (spriteGo != null)
                {
                    SpriteRenderer spriteRenderer = spriteGo.GetComponent<SpriteRenderer>();
                    if (spriteRenderer != null)
                    {
                        if (spriteRenderer.sprite != null)
                        {
                            image.sprite = spriteRenderer.sprite;
                        }
                    }
                }
            }
        }
    }

    public void SetImageSpriteByRemoteUrl(UnityEngine.UI.Image image, string remoteUrl, bool isNativeSize = false)
    {
        StartCoroutine(SetImageSpriteByWWWCoroutine(image, remoteUrl, isNativeSize));
    }

    IEnumerator SetImageSpriteByWWWCoroutine(UnityEngine.UI.Image image, string remoteUrl, bool isNativeSize, bool saveToCache = true)
    {
        if (string.IsNullOrEmpty(remoteUrl) || image == null)
        {
            yield break;
        }
        bool isLoadFromLocal = false;
        string localFileName = FileUtility.GetAssetFilePath(string.Format("HeadIcon/{0}.face", Utility.GetMD5Code(remoteUrl)), PathType.Local);

        if (FileUtility.ExistsFile(localFileName))
        {
            isLoadFromLocal = true;
        }

        using (WWW loader = new WWW(remoteUrl))
        {
            if (!loader.isDone)
            {
                yield return loader;
            }

            if (loader.isDone)
            {
                if (string.IsNullOrEmpty(loader.error))
                {
                    if (image != null)
                    {
                        if (!isLoadFromLocal)
                        {
                            FileUtility.WriteBytesToFile(localFileName, loader.bytes, loader.bytes.Length);
                        }
                        Texture2D tempImage = loader.texture;
                        if (tempImage != null)
                        {
                            Sprite sprite = Sprite.Create(tempImage, new Rect(0, 0, tempImage.width, tempImage.height), new Vector2(0, 0));
                            image.sprite = sprite;
                            if (isNativeSize)
                            {
                                image.SetNativeSize();
                            }
                        }
                    }
                }
            }
            else
            {
                Debug.LogErrorFormat("load sprite[{0}] error:{1}", remoteUrl, loader.error);
            }
        }
    }

    public void ResetSpriteByLocalDataPath(UnityEngine.UI.Image image, string projectName, string spriteName, bool isNativeSize = false)
    {
        byte[] contentBytes = FileUtility.ReadFileBytes(FileUtility.GetDataFileFullPath(projectName, spriteName));

        Texture2D texture = new Texture2D(4, 4, TextureFormat.ARGB32, false); //初始大小随便设置，只要大于1就好，false是不开启mipamaps
        texture.LoadImage(contentBytes);

        if (texture != null)
        {
            Sprite sprite = Sprite.Create(texture, new Rect(0, 0, texture.width, texture.height), new Vector2(0, 0));
            image.sprite = sprite;
            if (isNativeSize)
            {
                image.SetNativeSize();
            }
        }
    }

    public Sprite GetSprite(string projcetName, string spriteName)
    {
        spriteName = "sprite_" + spriteName;
        GameObject spriteGo = AssetBundleManager.Instance.LoadAsset<GameObject>(projcetName, spriteName, unloadAssetBundle: false);
        if (spriteGo != null)
        {
            SpriteRenderer spriteRenderer = spriteGo.GetComponent<SpriteRenderer>();
            if (spriteRenderer != null && spriteRenderer.sprite != null)
            {
                return spriteRenderer.sprite;
            }
        }

        return null;
    }
}