
public static class ExtendDefine
{
    /// <summary>
    /// 扩展 Image 方法，设置Sprite名称 刷新 Sprite
    /// </summary>
    /// <param name="image">扩展的对象</param>
    /// <param name="spriteName">图片资源名称</param>
    /// <param name="isNativeSize"></param>
    public static void ResetSpriteByName(this UnityEngine.UI.Image image, string projcetName, string spriteName, bool isNativeSize = false)
    {
        if (UISpriteManager.Instance != null)
        {
            UISpriteManager.Instance.SetImageSprite(image, projcetName, spriteName, isNativeSize);
        }
    }

    public static void ResetSpriteByNameAsync(this UnityEngine.UI.Image image, string projcetName, string spriteName, bool isNativeSize = false)
    {
        if (UISpriteManager.Instance != null)
        {
            UISpriteManager.Instance.SetImageSpriteAsync(image, projcetName, spriteName, isNativeSize);
        }
    }

    /// <summary>
    /// 重置图片的Sprite，通过Url 或者 名称
    /// (先加载默认资源)
    /// </summary>
    /// <param name="image"></param>
    /// <param name="spriteName"></param>
    /// <param name="remoteUrl"></param>
    /// <param name="isNativeSize"></param>
    public static void ResetSpriteByRemoteUrl(this UnityEngine.UI.Image image, string remoteUrl, bool isNativeSize = false)
    {
        if (UISpriteManager.Instance != null)
        {
            UISpriteManager.Instance.SetImageSpriteByRemoteUrl(image, remoteUrl, isNativeSize);
        }
    }
     
    /// <summary>
    /// 读取本地目录的一个图片显示到Image上 
    /// </summary>
    /// <param name="image"></param>
    /// <param name="spriteName"></param>
    /// <param name="remoteUrl"></param>
    /// <param name="isNativeSize"></param>
    public static void ResetSpriteByLocalDataPath(UnityEngine.UI.Image image, string projectName, string spriteName)
    {
        if (UISpriteManager.Instance != null)
        {
            UISpriteManager.Instance.ResetSpriteByLocalDataPath(image, projectName, spriteName);
        }
    }
}
