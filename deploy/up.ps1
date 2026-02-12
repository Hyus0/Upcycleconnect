param(
  [switch]$Build
)

$cmd = "docker compose -f deploy/docker-compose/docker-compose.yml up -d"
if ($Build) {
  $cmd += " --build"
}
Write-Host "Running: $cmd"
Invoke-Expression $cmd
