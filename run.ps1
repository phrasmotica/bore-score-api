$name = "bore-score-api"
$tag = "phrasmotica/bore-score-api"

$runningName = $(docker ps --filter name=$name --format "{{.Names}}")

if ($runningName) {
    docker stop $name | Out-Null
    Write-Host "Stopped running container $name"
}

$stoppedName = $(docker ps -a --filter name=$name --format "{{.Names}}")

if ($stoppedName) {
    docker rm $name | Out-Null
    Write-Host "Removed stopped container $name"
}

$createdId = $(docker run --name $name -e BORESCORE_ENV=production -d -p 8000:8000 $tag)
if ($createdId) {
    Write-Host "Started container $name ($createdId)"
}
