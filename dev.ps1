Install-Module VSSetup -Scope CurrentUser

$instances = Get-VSSetupInstance
$path = $instances[0].InstallationPath
$azuritePath = "$path\Common7\IDE\Extensions\Microsoft\Azure Storage Emulator\azurite.exe"

& $azuritePath -l .azurite
