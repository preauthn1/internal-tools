using System;
using System.Collections.Generic;
using UnityEditor;
using UnityEngine;
using example.invalid;
public partial class ResourceExporter
{
    static List<TextureFormatInfo> GetAllTex(bool bFilter, CheckPlatformTexFormat checkAndroid, CheckPlatformTexFormat checkIphone, ReplacePlatformTexFormat replaceAction = null)
    {
        List<string> texPathLs = new List<string>();
        GetAssetsTexInSubFolderForAtlas("assets", ref texPathLs);
        List<TextureFormatInfo> infoLs = new List<TextureFormatInfo>();
        int totalCount = texPathLs.Count;
        if (texPathLs.Count > 0)
        {
            int iCurCount = 0;
            foreach (var path in texPathLs)
            {
                TextureImporter texImporter = TextureImporter.GetAtPath(path) as TextureImporter;
                if (null != texImporter)
                {
                    TextureFormatInfo texFormatInfo = new TextureFormatInfo(texImporter.assetPath);
                    texImporter.GetPlatformTextureSettings("Standalone", out texFormatInfo.windowsMaxSize, out texFormatInfo.texFormatWindows);
                    texImporter.GetPlatformTextureSettings("iPhone", out texFormatInfo.iphoneMaxSize, out texFormatInfo.texFormatIphone);
                    texImporter.GetPlatformTextureSettings("Android", out texFormatInfo.androidMaxSize, out texFormatInfo.texFormatAndroid);
                    texFormatInfo.bMipMapEnabled = texImporter.mipmapEnabled;
                    bool bAdd = true;
                    texFormatInfo.bNeedAndroidFormat = checkAndroid(texFormatInfo.texFormatAndroid);
                    texFormatInfo.bNeedIphoneFormat = checkIphone(texFormatInfo.texFormatIphone);
                    if (bFilter)
                    {
                        bAdd = texFormatInfo.bMipMapEnabled || (!texFormatInfo.bNeedAndroidFormat) || (!texFormatInfo.bNeedIphoneFormat);
                    }
                    if (bAdd)
                    {
                        //Debug.Log(texFormatInfo.DataToString());
                        if (null != replaceAction)
                        {
                            replaceAction(texImporter, texFormatInfo);
                        }
                        infoLs.Add(texFormatInfo);
                    }
                }
                EditorUtility.DisplayCancelableProgressBar("Check TexFormat", "Wait......", (++iCurCount) * 1f / totalCount);
            }
        }
        totalCount = infoLs.Count;
        Debug.Log("PathCount:" + texPathLs.Count);
        Debug.Log("InfoCount:" + infoLs.Count);
        EditorUtility.ClearProgressBar();
        return infoLs;
    }
    static void SaveTextureFormatInfoPath(List<TextureFormatInfo> infoLs)
    {
        string fileSavePath = Application.dataPath + "\\FileLog.html";
        if (File.Exists(fileSavePath))
        {
            File.Delete(fileSavePath);
        }

        StreamWriter t_sStreamWriter; // 文件流信息
        FileInfo t_fFileInfo = new FileInfo(fileSavePath);
        if (!t_fFileInfo.Exists)
        {
            t_sStreamWriter = t_fFileInfo.CreateText();  // 如果此文件不存在则创建
        }
        else
        {
            t_sStreamWriter = t_fFileInfo.AppendText(); // 如果此文件存在则打开
        }

        t_sStreamWriter.WriteLine("size = " + infoLs.Count + "</br>");
        foreach (var info in infoLs)
        {
            t_sStreamWriter.WriteLine(info.ToHtmlString());
        }

        t_sStreamWriter.Close(); //关闭流
        t_sStreamWriter.Dispose(); // 销毁流
    }
    static bool CheckAndroidFormat(TextureImporterFormat format)
    {
        return format == TextureImporterFormat.ETC_RGB4;
    }
    static bool CheckIphoneFormat(TextureImporterFormat format)
    {
        return format.ToString().Contains("PVRTC");
    }
    delegate bool CheckPlatformTexFormat(TextureImporterFormat format);
    delegate void ReplacePlatformTexFormat(TextureImporter textureImporter, TextureFormatInfo texInfo);
    static void ReplacePlatformFormat(TextureImporter textureImporter, TextureFormatInfo texInfo)
    {
        textureImporter.mipmapEnabled = false;
        textureImporter.SetPlatformTextureSettings("Android", texInfo.androidMaxSize, TextureImporterFormat.ETC_RGB4);
        AssetDatabase.ImportAsset(textureImporter.assetPath);
    }
    class TextureFormatInfo
    {
        public string path;
        public bool bMipMapEnabled;
        public bool bNeedAndroidFormat;
        public bool bNeedIphoneFormat;
        public TextureImporterFormat texFormatIphone;
        public TextureImporterFormat texFormatAndroid;
        public TextureImporterFormat texFormatWindows;
        public int androidMaxSize;
        public int iphoneMaxSize;
        public int windowsMaxSize;
        public TextureFormatInfo(string filePath)
        {
            path = filePath;
        }
        public string DataToString()
        {
            string strFormat = "{0}: MipMapEnabled:{1},\nAndroid:{2},\nIphone:{3}";
            return string.Format(strFormat, path, bMipMapEnabled, texFormatAndroid, texFormatIphone, texFormatWindows);
        }
        public string ToHtmlString()
        {
            //<font color="red"><b> size = 2357 </b></font></br>
            string strFormat = "{0}: {1}, {2}, {3}</br>";
            string bMipStr = string.Format("MipMapEnabled:{0}", bMipMapEnabled);
            string parm1 = bMipMapEnabled ? GetSpecialTip(bMipStr, "red") : bMipStr;
            string androidStr = string.Format("Android:{0}", texFormatAndroid);
            string parm2 = !bNeedAndroidFormat ? GetSpecialTip(androidStr, "blue") : androidStr;
            string iphoneStr = string.Format("Iphone:{0}", texFormatIphone);
            string parm3 = !bNeedIphoneFormat ? GetSpecialTip(iphoneStr, "YellowGreen") : iphoneStr;
            return string.Format(strFormat, path, parm1, parm2, parm3);
        }
        private string GetSpecialTip(string name, string fontColor)
        {
            return string.Format("<font color=\"{0}\"><b>{1}</b></font>", fontColor, name);
        }
    }
    static void GetAssetsTexInSubFolderForAtlas(string srcFolder, ref List<string> atlas)
    {
        string searchPattern0 = "*.png";
        string searchPattern1 = "*.tga";
        string searchPattern2 = "*.psd";
        string searchFolder = srcFolder.Replace(@"\", @"/");
        string searchDir0 = searchFolder;
        if (Directory.Exists(searchDir0))
        {
            //string[] files = Directory.GetFiles(searchDir0, searchPattern);
            AddFilePathToList(Directory.GetFiles(searchDir0, searchPattern0), ref atlas);
            AddFilePathToList(Directory.GetFiles(searchDir0, searchPattern1), ref atlas);
            AddFilePathToList(Directory.GetFiles(searchDir0, searchPattern2), ref atlas);
        }
        string[] dirs = Directory.GetDirectories(searchFolder);
        foreach (string oneDir in dirs)
        {
            GetAssetsTexInSubFolderForAtlas(oneDir, ref atlas);
        }
    }
    static void AddFilePathToList(string[] files, ref List<string> ls)
    {
        foreach (string oneFile in files)
        {
            string srcFile = oneFile.Replace(@"\", @"/");
            string lowerFile = srcFile.ToLower();
            ls.Add(lowerFile);
        }
    }
    [MenuItem("Texture/顯示所有Texture格式")]
    public static void CheckTextureFormat()
    {
        List<TextureFormatInfo> infoLs = GetAllTex(true, CheckAndroidFormat, CheckIphoneFormat);
        SaveTextureFormatInfoPath(infoLs);
    }
    [MenuItem("Texture/轉換所有Texture格式")]
    public static void ReplaceTextureFormat()
    {
        if (EditorUtility.DisplayDialog("提示", "確定轉換所有Texture格式?", "确定", "取消") == false)
            return;

        GetAllTex(true, CheckAndroidFormat, CheckIphoneFormat, ReplacePlatformFormat);
        //SaveTextureFormatInfoPath(infoLs);
    }
}