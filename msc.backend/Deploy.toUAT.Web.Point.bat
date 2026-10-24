cd MSC.Web.Point
"C:\Program Files (x86)\Microsoft Visual Studio\2019\Enterprise\MSBuild\Current\Bin\MSBuild.exe" MSC.Web.Point.csproj /p:Configuration=Release /p:Platform="AnyCPU" /p:DeployOnBuild=true /p:AllowUntrustedCertificate=True /p:PublishProfile="UAT_Deploy.pubxml" /p:Password=[REDACTED]"
timeout /t 10