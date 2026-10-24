// Dependency Checker
// based off the "Resource Checker" free utility in the Unity Asset store https://example.invalid/#/content/3224
// extended and heavily modified.

using UnityEngine;
using UnityEditor;
using System.Collections.Generic;
using UnityEngine.UI;
using System.Linq;
using UnityEditor.SceneManagement;
using UnityEngine.SceneManagement;
using example.invalid;

public class Util
{
    public static bool ctrlPressed = false;

    public static int ThumbnailWidth = 50;
    public static int ThumbnailHeight = 50;

    public static void SelectObject(Object selectedObject)
    {
        if (Util.ctrlPressed)
        {
            List<Object> currentSelection = new List<Object>(Selection.objects);

            // Allow toggle selection
            if (currentSelection.Contains(selectedObject))
                currentSelection.Remove(selectedObject);
            else currentSelection.Add(selectedObject);

            Selection.objects = currentSelection.ToArray();
        }
        else Selection.activeObject = selectedObject;
    }

    public static void SelectObjects(List<Object> selectedObjects)
    {
        if (Util.ctrlPressed)
        {
            List<Object> currentSelection = new List<Object>(Selection.objects);
            currentSelection.AddRange(selectedObjects);
            Selection.objects = currentSelection.ToArray();
        }
        else
            Selection.objects = selectedObjects.ToArray();
    }

    public static void SelectObjects(List<GameObject> selectedObjects)
    {
        if (Util.ctrlPressed)
        {
            List<Object> currentSelection = new List<Object>(Selection.objects);

            foreach (GameObject obj in selectedObjects)
            {
                currentSelection.Add(obj);
            }

            Selection.objects = currentSelection.ToArray();
        }
        else
            Selection.objects = selectedObjects.ToArray();
    }
}

public class TextureDetails
{
    public bool isCubeMap;
    public int memSizeBytes;
    public Texture texture;
    public TextureFormat format;
    public int mipMapCount;

    public List<Object> FoundInMaterials = new List<Object>();
    public List<GameObject> FoundInGameObjects = new List<GameObject>();

    public TextureDetails(Texture tex)
    {
        texture = tex;
        isCubeMap = tex is Cubemap;
        memSizeBytes = TextureDetails.CalculateTextureSizeBytes(tex);
        format = TextureFormat.RGBA32;
        mipMapCount = 1;
        if (texture is Texture2D)
        {
            format = (texture as Texture2D).format;
            mipMapCount = (texture as Texture2D).mipmapCount;
        }
        if (texture is Cubemap)
        {
            format = (texture as Cubemap).format;
        }
    }

    public void OnGui()
    {
        if (texture != null)
        {
            Texture thumb = texture;

            GUILayout.BeginHorizontal();
            GUILayout.Box(thumb, GUILayout.Width(Util.ThumbnailWidth), GUILayout.Height(Util.ThumbnailHeight));

            if (GUILayout.Button(new GUIContent(texture.name, texture.name), GUILayout.Width(150), GUILayout.Height(50)))
                Util.SelectObject(texture);

            Texture2D iconMaterials = AssetPreview.GetMiniTypeThumbnail(typeof(Material));
            if (GUILayout.Button(new GUIContent(FoundInMaterials.Count.ToString(), iconMaterials, "Materials"), GUILayout.Width(60), GUILayout.Height(50)))
                Util.SelectObjects(FoundInMaterials);

            Texture2D iconGameObj = AssetPreview.GetMiniTypeThumbnail(typeof(GameObject));
            if (GUILayout.Button(new GUIContent(FoundInGameObjects.Count.ToString(), iconGameObj, "Game Objects"), GUILayout.Width(60), GUILayout.Height(50)))
                Util.SelectObjects(FoundInGameObjects);

            string labelTxt = string.Format("{0}x{1}\n{2} - {3}",
                texture.width,
                texture.height,
                EditorUtility.FormatBytes(memSizeBytes),
                format);

            GUILayout.Box(labelTxt, GUILayout.Width(120), GUILayout.Height(50));

            GUILayout.EndHorizontal();
        }
    }

    private static int CalculateTextureSizeBytes(Texture tTexture)
    {
        int tWidth = tTexture.width;
        int tHeight = tTexture.height;
        if (tTexture is Texture2D)
        {
            Texture2D tTex2D = tTexture as Texture2D;
            int bitsPerPixel = GetBitsPerPixel(tTex2D.format);
            int mipMapCount = tTex2D.mipmapCount;
            int mipLevel = 1;
            int tSize = 0;
            while (mipLevel <= mipMapCount)
            {
                tSize += tWidth * tHeight * bitsPerPixel / 8;
                tWidth = tWidth / 2;
                tHeight = tHeight / 2;
                mipLevel++;
            }
            return tSize;
        }

        if (tTexture is Cubemap)
        {
            Cubemap tCubemap = tTexture as Cubemap;
            int bitsPerPixel = GetBitsPerPixel(tCubemap.format);
            return tWidth * tHeight * 6 * bitsPerPixel / 8;
        }
        return 0;
    }
    private static int GetBitsPerPixel(TextureFormat format)
    {
        switch (format)
        {
            case TextureFormat.Alpha8: //	 Alpha-only texture format.
                return 8;
            case TextureFormat.ARGB4444: //	 A 16 bits/pixel texture format. Texture stores color with an alpha channel.
                return 16;
            case TextureFormat.RGB24:   // A color texture format.
                return 24;
            case TextureFormat.RGBA32:  //Color with an alpha channel texture format.
                return 32;
            case TextureFormat.ARGB32:  //Color with an alpha channel texture format.
                return 32;
            case TextureFormat.RGB565:  //	 A 16 bit color texture format.
                return 16;
            case TextureFormat.DXT1:    // Compressed color texture format.
                return 4;
            case TextureFormat.DXT5:    // Compressed color with alpha channel texture format.
                return 8;
            /*
			case TextureFormat.WiiI4:	// Wii texture format.
			case TextureFormat.WiiI8:	// Wii texture format. Intensity 8 bit.
			case TextureFormat.WiiIA4:	// Wii texture format. Intensity + Alpha 8 bit (4 + 4).
			case TextureFormat.WiiIA8:	// Wii texture format. Intensity + Alpha 16 bit (8 + 8).
			case TextureFormat.WiiRGB565:	// Wii texture format. RGB 16 bit (565).
			case TextureFormat.WiiRGB5A3:	// Wii texture format. RGBA 16 bit (4443).
			case TextureFormat.WiiRGBA8:	// Wii texture format. RGBA 32 bit (8888).
			case TextureFormat.WiiCMPR:	//	 Compressed Wii texture format. 4 bits/texel, ~RGB8A1 (Outline alpha is not currently supported).
				return 0;  //Not supported yet
			*/
            case TextureFormat.PVRTC_RGB2://	 PowerVR (iOS) 2 bits/pixel compressed color texture format.
                return 2;
            case TextureFormat.PVRTC_RGBA2://	 PowerVR (iOS) 2 bits/pixel compressed with alpha channel texture format
                return 2;
            case TextureFormat.PVRTC_RGB4://	 PowerVR (iOS) 4 bits/pixel compressed color texture format.
                return 4;
            case TextureFormat.PVRTC_RGBA4://	 PowerVR (iOS) 4 bits/pixel compressed with alpha channel texture format
                return 4;
            case TextureFormat.ETC_RGB4://	 ETC (GLES2.0) 4 bits/pixel compressed RGB texture format.
                return 4;
           // case TextureFormat.ATC_RGB4://	 ATC (ATITC) 4 bits/pixel compressed RGB texture format.
           //     return 4;
           // case TextureFormat.ATC_RGBA8://	 ATC (ATITC) 8 bits/pixel compressed RGB texture format.
           //     return 8;
            case TextureFormat.BGRA32://	 Format returned by iPhone camera
                return 32;
            case TextureFormat.ETC2_RGBA8:
                return 8;
                //case TextureFormat.ATF_RGB_DXT1://	 Flash-specific RGB DXT1 compressed color texture format.
                //case TextureFormat.ATF_RGBA_JPG://	 Flash-specific RGBA JPG-compressed color texture format.
                //case TextureFormat.ATF_RGB_JPG://	 Flash-specific RGB JPG-compressed color texture format.
                //	return 0; //Not supported yet
        }
        return 0;
    }
};

public class MaterialDetails
{
    public Material material;

    public List<GameObject> FoundInGameObjects = new List<GameObject>();

    public MaterialDetails(Material mat)
    {
        material = mat;
    }

    public void OnGui()
    {
        if (material != null)
        {
            GUILayout.BeginHorizontal();

            Texture thumb = material.mainTexture;
            if (thumb == null)
                thumb = AssetPreview.GetMiniTypeThumbnail(typeof(Material));

            GUILayout.Box(thumb, GUILayout.Width(Util.ThumbnailWidth), GUILayout.Height(Util.ThumbnailHeight));

            if (GUILayout.Button(new GUIContent(material.name, material.name), GUILayout.Width(150), GUILayout.Height(50)))
                Util.SelectObject(material);

            Texture2D iconGameObj = AssetPreview.GetMiniTypeThumbnail(typeof(GameObject));
            if (GUILayout.Button(new GUIContent(FoundInGameObjects.Count.ToString(), iconGameObj, "Game Objects"), GUILayout.Width(60), GUILayout.Height(50)))
                Util.SelectObjects(FoundInGameObjects);

            GUILayout.EndHorizontal();
        }
    }
};

public class ShaderDetails
{
    public Shader shader;

    public List<GameObject> FoundInGameObjects = new List<GameObject>();
    public List<Object> FoundInMaterials = new List<Object>();

    public ShaderDetails(Shader s)
    {
        shader = s;
    }

    public void OnGui()
    {
        if (shader != null)
        {
            GUILayout.BeginHorizontal();

            Texture2D thumb = AssetPreview.GetMiniThumbnail(shader);
            if (thumb == null)
                thumb = AssetPreview.GetMiniTypeThumbnail(typeof(Shader));

            GUILayout.Box(thumb, GUILayout.Width(Util.ThumbnailWidth), GUILayout.Height(Util.ThumbnailHeight));

            if (GUILayout.Button(new GUIContent(shader.name, shader.name), GUILayout.Width(150), GUILayout.Height(50)))
                Util.SelectObject(shader);

            Texture2D iconMaterials = AssetPreview.GetMiniTypeThumbnail(typeof(Material));
            if (GUILayout.Button(new GUIContent(FoundInMaterials.Count.ToString(), iconMaterials, "Materials"), GUILayout.Width(60), GUILayout.Height(50)))
                Util.SelectObjects(FoundInMaterials);

            Texture2D iconGameObj = AssetPreview.GetMiniTypeThumbnail(typeof(GameObject));
            if (GUILayout.Button(new GUIContent(FoundInGameObjects.Count.ToString(), iconGameObj, "Game Objects"), GUILayout.Width(60), GUILayout.Height(50)))
                Util.SelectObjects(FoundInGameObjects);

            GUILayout.EndHorizontal();
        }
    }
};

public class SoundDetails
{
    public AudioClip clip;

    public List<GameObject> FoundInGameObjects = new List<GameObject>();

    public SoundDetails(AudioClip c)
    {
        clip = c;
    }

    public void OnGui()
    {
        if (clip != null)
        {
            GUILayout.BeginHorizontal();

            Texture2D thumb = AssetPreview.GetAssetPreview(clip);
            GUILayout.Box(thumb, GUILayout.Width(Util.ThumbnailWidth), GUILayout.Height(Util.ThumbnailHeight));

            if (GUILayout.Button(new GUIContent(clip.name, clip.name), GUILayout.Width(150), GUILayout.Height(50)))
                Util.SelectObject(clip);

            Texture2D iconGameObj = AssetPreview.GetMiniTypeThumbnail(typeof(GameObject));
            if (GUILayout.Button(new GUIContent(FoundInGameObjects.Count.ToString(), iconGameObj, "Game Objects"), GUILayout.Width(60), GUILayout.Height(50)))
                Util.SelectObjects(FoundInGameObjects);

            GUILayout.Box(clip.length + " length\n" + clip.channels + " channels\n" + "frequency " + clip.frequency, GUILayout.Width(150), GUILayout.Height(50));

            GUILayout.EndHorizontal();
        }
    }
};

public class FontDetails
{
    public Font font;

    public List<GameObject> FoundInGameObjects = new List<GameObject>();

    public FontDetails(Font s)
    {
        font = s;
    }

    public void OnGui()
    {
        if (font != null)
        {
            GUILayout.BeginHorizontal();

            Texture2D thumb = AssetPreview.GetAssetPreview(font);
            if (thumb == null)
                thumb = AssetPreview.GetMiniTypeThumbnail(typeof(Font));

            GUILayout.Box(thumb, GUILayout.Width(Util.ThumbnailWidth), GUILayout.Height(Util.ThumbnailHeight));

            if (GUILayout.Button(new GUIContent(font.name, font.name), GUILayout.Width(150), GUILayout.Height(50)))
                Util.SelectObject(font);

            Texture2D iconGameObj = AssetPreview.GetMiniTypeThumbnail(typeof(GameObject));
            if (GUILayout.Button(new GUIContent(FoundInGameObjects.Count.ToString(), iconGameObj, "Game Objects"), GUILayout.Width(60), GUILayout.Height(50)))
                Util.SelectObjects(FoundInGameObjects);

            GUILayout.EndHorizontal();
        }
    }
};

public class DependencyChecker : EditorWindow
{
    enum InspectType
    {
        Textures, Fonts,/* Sounds,*/ Materials, Shaders
    };

    InspectType ActiveInspectType = InspectType.Textures;
    bool m_IsExternalAssets = true;     // 是搜索外部资源

    List<TextureDetails> ActiveTextures = new List<TextureDetails>();
    List<FontDetails> ActiveFontDetails = new List<FontDetails>();
    List<SoundDetails> ActiveSoundDetails = new List<SoundDetails>();
    List<ShaderDetails> ActiveShaderDetails = new List<ShaderDetails>();
    List<MaterialDetails> ActiveMaterials = new List<MaterialDetails>();

    Vector2 textureListScrollPos = new Vector2(0, 0);
    Vector2 fontListScrollPos = new Vector2(0, 0);
    Vector2 soundListScrollPos = new Vector2(0, 0);
    Vector2 materialListScrollPos = new Vector2(0, 0);
    Vector2 shaderListScrollPos = new Vector2(0, 0);

    int TotalTextureMemory = 0;

    static int MinWidth = 300;

    [MenuItem("Window/项目依赖")]
    static void Init()
    {
        if (string.IsNullOrEmpty(AppDefine.CurrentProjectPath))
        {
            EditorUtility.DisplayDialog("错误", "没有设置当前项目，右键->设置为当前项目可设置当前项目。", "确定");
            return;
        }

        WindowSceneReady();

        DependencyChecker window = (DependencyChecker)EditorWindow.GetWindow(typeof(DependencyChecker), false, "项目依赖");
        window.CheckResources();
        window.minSize = new Vector2(MinWidth, 300);
    }

    void OnGUI()
    {
        GUILayout.BeginHorizontal();

        Texture2D iconTexture = AssetPreview.GetMiniTypeThumbnail(typeof(Texture2D));
        Texture2D iconFont = AssetPreview.GetMiniTypeThumbnail(typeof(Font));
        // Texture2D iconSound = AssetPreview.GetMiniTypeThumbnail(typeof(AudioClip));
        Texture2D iconMaterial = AssetPreview.GetMiniTypeThumbnail(typeof(Material));
        Texture2D iconShader = AssetPreview.GetMiniTypeThumbnail(typeof(Shader));

        GUIContent[] guiObjs =
        {
            new GUIContent( iconTexture, "贴图" ),
            new GUIContent( iconFont, "字体" ),
           // new GUIContent( iconSound, "音效" ),
            new GUIContent( iconMaterial, "材质球" ),
            new GUIContent( iconShader, "Shader" ),
        };

        GUILayoutOption[] options =
        {
            GUILayout.Width( 400 ),
            GUILayout.Height( 100 ),
        };

        ActiveInspectType = (InspectType)GUILayout.Toolbar((int)ActiveInspectType, guiObjs, options);

        GUI.Label(new Rect(10, 85, 100, 100), ActiveTextures.Count + " - " + EditorUtility.FormatBytes(TotalTextureMemory));
        GUI.Label(new Rect(110, 85, 100, 100), ActiveFontDetails.Count.ToString());
        // GUI.Label(new Rect(210, 85, 100, 100), ActiveSoundDetails.Count.ToString());
        GUI.Label(new Rect(210, 85, 100, 100), ActiveMaterials.Count.ToString());
        GUI.Label(new Rect(310, 85, 100, 100), ActiveShaderDetails.Count.ToString());


        GUILayout.BeginVertical(GUILayout.Width(100), GUILayout.Height(100));

        GUILayout.FlexibleSpace();

        m_IsExternalAssets = GUILayout.Toggle(m_IsExternalAssets, "外部资源", GUILayout.Width(100), GUILayout.Height(20));
        m_IsExternalAssets = !GUILayout.Toggle(!m_IsExternalAssets, "所有资源", GUILayout.Width(100), GUILayout.Height(20));

        if (GUILayout.Button("刷新", GUILayout.Width(100), GUILayout.Height(50)))
            CheckResources();

        GUILayout.EndVertical();

        GUILayout.FlexibleSpace();
        GUILayout.EndHorizontal();

        GUILayout.BeginHorizontal();
        Util.ctrlPressed = Event.current.control || Event.current.command;

        switch (ActiveInspectType)
        {
            case InspectType.Textures:
                ListTextures();
                break;
            case InspectType.Fonts:
                ListFonts();
                break;
            //case InspectType.Sounds:
            //    ListSounds();
            //    break;
            case InspectType.Materials:
                ListMaterials();
                break;
            case InspectType.Shaders:
                ListShaders();
                break;
        }
        GUILayout.EndHorizontal();
    }

    void ListTextures()
    {
        textureListScrollPos = EditorGUILayout.BeginScrollView(textureListScrollPos);

        if (ActiveTextures.Count > 0)
        {
            GUILayout.BeginHorizontal();
            if (GUILayout.Button("全选"))
            {
                List<Object> allObject = new List<Object>();
                foreach (TextureDetails tDetails in ActiveTextures)
                    allObject.Add(tDetails.texture);

                Util.SelectObjects(allObject);
            }
            GUILayout.EndHorizontal();
        }

        foreach (TextureDetails details in ActiveTextures)
        {
            details.OnGui();
        }
        EditorGUILayout.EndScrollView();
    }

    void ListFonts()
    {
        fontListScrollPos = EditorGUILayout.BeginScrollView(fontListScrollPos);

        if (ActiveFontDetails.Count > 0)
        {
            GUILayout.BeginHorizontal();
            if (GUILayout.Button("全选"))
            {
                List<Object> allObject = new List<Object>();
                foreach (FontDetails tDetails in ActiveFontDetails)
                    allObject.Add(tDetails.font);

                Util.SelectObjects(allObject);
            }
            GUILayout.EndHorizontal();
        }

        foreach (FontDetails details in ActiveFontDetails)
        {
            details.OnGui();
        }
        EditorGUILayout.EndScrollView();
    }

    void ListSounds()
    {
        soundListScrollPos = EditorGUILayout.BeginScrollView(soundListScrollPos);

        foreach (SoundDetails details in ActiveSoundDetails)
        {
            details.OnGui();
        }
        EditorGUILayout.EndScrollView();
    }

    void ListMaterials()
    {
        materialListScrollPos = EditorGUILayout.BeginScrollView(materialListScrollPos);

        if (ActiveMaterials.Count > 0)
        {
            GUILayout.BeginHorizontal();
            if (GUILayout.Button("全选"))
            {
                List<Object> allObject = new List<Object>();
                foreach (MaterialDetails tDetails in ActiveMaterials)
                    allObject.Add(tDetails.material);

                Util.SelectObjects(allObject);
            }
            GUILayout.EndHorizontal();
        }


        foreach (MaterialDetails details in ActiveMaterials)
        {
            details.OnGui();
        }
        EditorGUILayout.EndScrollView();
    }

    void ListShaders()
    {
        shaderListScrollPos = EditorGUILayout.BeginScrollView(shaderListScrollPos);

        foreach (ShaderDetails details in ActiveShaderDetails)
        {
            details.OnGui();
        }
        EditorGUILayout.EndScrollView();
    }

    void TryAddActiveTextures(Texture tex, Object inMaterials, GameObject inGameObject)
    {
        if (tex == null)
            return;

        string assetPath = AssetDatabase.GetAssetPath(tex);

        if (m_IsExternalAssets && assetPath.Contains(AppDefine.CurrentProjectPath))
            return;

        if (IgnoreDependencyFromPath(assetPath))
            return;

        TextureDetails details = ActiveTextures.Find((o) => { return o.texture == tex; });

        if (details == null)
        {
            details = new TextureDetails(tex);
            ActiveTextures.Add(details);
        }

        if (inGameObject && !details.FoundInGameObjects.Contains(inGameObject))
            details.FoundInGameObjects.Add(inGameObject);

        string materialsPath = AssetDatabase.GetAssetPath(inMaterials);
        if (IgnoreDependencyFromPath(materialsPath))
            return;

        if (inMaterials && !details.FoundInMaterials.Contains(inMaterials))
            details.FoundInMaterials.Add(inMaterials);

        return;
    }

    void TryAddAudioClip(AudioClip s)
    {
        if (s == null)
            return;

        string assetPath = AssetDatabase.GetAssetPath(s);

        if (m_IsExternalAssets && assetPath.Contains(AppDefine.CurrentProjectPath))
            return;

        if (IgnoreDependencyFromPath(assetPath))
            return;

        SoundDetails details = ActiveSoundDetails.Find((o) => { return o.clip == s; });
        if (details == null)
        {
            details = new SoundDetails(s);
            ActiveSoundDetails.Add(details);
        }
        return;
    }

    void TryAddFont(Font s, GameObject inGameObject)
    {
        if (s == null)
            return;

        string assetPath = AssetDatabase.GetAssetPath(s);

        if (m_IsExternalAssets && assetPath.Contains(AppDefine.CurrentProjectPath))
            return;

        if (IgnoreDependencyFromPath(assetPath))
            return;

        FontDetails details = ActiveFontDetails.Find((o) => { return o.font == s; });
        if (details == null)
        {
            details = new FontDetails(s);
            ActiveFontDetails.Add(details);
        }

        if (!details.FoundInGameObjects.Contains(inGameObject))
            details.FoundInGameObjects.Add(inGameObject);

        return;
    }

    void TryAddActiveMaterial(Material mat, GameObject inGameObject)
    {
        if (mat == null)
            return;

        string assetPath = AssetDatabase.GetAssetPath(mat);

        if (m_IsExternalAssets && assetPath.Contains(AppDefine.CurrentProjectPath))
            return;

        if (IgnoreDependencyFromPath(assetPath))
            return;

        MaterialDetails details = ActiveMaterials.Find((o) => { return o.material == mat; });
        if (details == null)
        {
            details = new MaterialDetails(mat);
            ActiveMaterials.Add(details);
        }

        if (!details.FoundInGameObjects.Contains(inGameObject))
            details.FoundInGameObjects.Add(inGameObject);

        foreach (Object obj in EditorUtility.CollectDependencies(new UnityEngine.Object[] { details.material }))
        {
            if (obj is Texture)
            {
                Texture tTexture = obj as Texture;
                TryAddActiveTextures(tTexture, details.material, inGameObject);
            }
            if (obj is Shader)
            {
                Shader shader = obj as Shader;
                TryAddActiveShader(shader, details.material, inGameObject);
            }
        }

        return;
    }

    void TryAddActiveShader(Shader s, Object inMaterials, GameObject inGameObject)
    {
        if (s == null)
            return;

        string assetPath = AssetDatabase.GetAssetPath(s);

        if (m_IsExternalAssets && assetPath.Contains(AppDefine.CurrentProjectPath))
            return;

        if (IgnoreDependencyFromPath(assetPath))
            return;

        ShaderDetails details = ActiveShaderDetails.Find((o) => { return o.shader == s; });
        if (details == null)
        {
            details = new ShaderDetails(s);
            ActiveShaderDetails.Add(details);
        }

        if (!details.FoundInMaterials.Contains(inMaterials))
            details.FoundInMaterials.Add(inMaterials);

        if (!details.FoundInGameObjects.Contains(inGameObject))
            details.FoundInGameObjects.Add(inGameObject);

        return;
    }

    void CheckResources()
    {
        ActiveTextures.Clear();
        ActiveFontDetails.Clear();
        ActiveSoundDetails.Clear();
        ActiveMaterials.Clear();
        ActiveShaderDetails.Clear();
        
        // 特效和骨骼动画的引用
        List<Renderer> renderers = GetAllTransInScene<Renderer>();
        foreach (Renderer renderer in renderers)
        {
            if (renderer is SpriteRenderer && (renderer as SpriteRenderer).sprite != null)
            { 
                TryAddActiveTextures((renderer as SpriteRenderer).sprite.texture, null, renderer.gameObject);
            }
            else
            {
                foreach (Material material in renderer.sharedMaterials)
                {
                    if (material == null)
                        continue;

                    TryAddActiveMaterial(material, renderer.gameObject);
                }
            }
        }

        List<Transform> gameObjs = GetAllTransInScene<Transform>();
        foreach (Transform item in gameObjs)
        {
            MaskableGraphic maskableGraphic = item.GetComponent<MaskableGraphic>();

            if (maskableGraphic is Image)
            {
                Image image = item.GetComponent<Image>();
                TryAddActiveTextures(image.mainTexture, image.material, item.gameObject);
                TryAddActiveMaterial(image.material, item.gameObject);
            }
            else if (maskableGraphic is RawImage)
            {
                RawImage image = item.GetComponent<RawImage>();
                TryAddActiveTextures(image.mainTexture, image.material, item.gameObject);
                TryAddActiveMaterial(image.material, item.gameObject);
            }
            else if (maskableGraphic is Text)
            {
                Text text = item.GetComponent<Text>();
                TryAddFont(text.font, item.gameObject);
            }
        }

        TotalTextureMemory = 0;
        foreach (TextureDetails tTextureDetails in ActiveTextures) TotalTextureMemory += tTextureDetails.memSizeBytes;

        // Sort by size, descending
        ActiveTextures.Sort(delegate (TextureDetails details1, TextureDetails details2) { return details2.memSizeBytes - details1.memSizeBytes; });
    }

    public List<T> GetAllTransInScene<T>() where T : Component
    {
        GameObject[] gameObjs = (GameObject[])FindObjectsOfType(typeof(GameObject));

        List<T> rtnList = new List<T>();

        for (int i = 0; i < gameObjs.Length; i++)
        {
            if (gameObjs[i].transform.parent == null)
                rtnList.AddRange(gameObjs[i].GetComponentsInChildren<T>(true));
        }

        return rtnList;
    }

    static bool IgnoreDependencyFromPath(string path)
    {
        return
            path.StartsWith("Library") ||
            path.StartsWith("Resources/unity_builtin_extra") ||
            path.EndsWith("dll") ||
            path.EndsWith("asset") ||               // 排除骨骼动画的材质
            string.IsNullOrEmpty(path);
    }

    private static string m_OpendScenePath;
    private static Scene m_PrefabScene;

    public static void WindowSceneReady()
    {
        string[] scenePaths = Directory.GetFiles(AppDefine.CurrentProjectPath, "*.unity", SearchOption.AllDirectories);

        m_OpendScenePath = SceneManager.GetActiveScene().path;

        if (scenePaths.Length > 0)
        {
            EditorSceneManager.OpenScene(scenePaths[0]);    //打开场景 
            m_PrefabScene = EditorSceneManager.NewScene(NewSceneSetup.EmptyScene, NewSceneMode.Additive);
        }
        else
        {
            m_PrefabScene = EditorSceneManager.NewScene(NewSceneSetup.EmptyScene, NewSceneMode.Single);
        }
        SceneManager.SetActiveScene(m_PrefabScene);

        string[] paths = Directory.GetFiles(AppDefine.CurrentProjectPath, "*.prefab", SearchOption.AllDirectories);

        for (int i = 0; i < paths.Length; i++)
        {
            GameObject prefabObj = AssetDatabase.LoadAssetAtPath(EditorFileUtility.GetRelativeAssetPath(paths[i]), typeof(GameObject)) as GameObject;
            PrefabUtility.InstantiatePrefab(prefabObj);
        }
    }

    public void OnDestroy()
    {
        if (m_PrefabScene != null)
        {
            EditorSceneManager.CloseScene(m_PrefabScene, true);

            if (!string.IsNullOrEmpty(m_OpendScenePath))
                EditorSceneManager.OpenScene(m_OpendScenePath);    //打开场景 
        }
    }
}