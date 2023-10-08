$BASE_URL = "https://raw.githubusercontent.com/orangekame3/winget-pkgs/main/manifests/g/orangekame3/gitfetch"
$FILES = @("orangekame3.gitfetch.yaml", "orangekame3.gitfetch.installer.yaml", "orangekame3.gitfetch.locale.en-US.yaml")

if (-not (Test-Path ./tmp)) {
    New-Item -ItemType Directory -Path ./tmp
}

foreach ($file in $FILES) {
    Invoke-WebRequest -Uri "$BASE_URL/$file" -OutFile "./tmp/$file"
}

winget install -m ./tmp/orangekame3.gitfetch.yaml


foreach ($file in $FILES) {
    Remove-Item "./tmp/$file"
}
