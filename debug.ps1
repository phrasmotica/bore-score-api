param(
    [switch] $Debugger = $false,
    [switch] $RunOnly = $false
)

$file = ".\docker-compose.yml"

if ($Debugger) {
    $file = ".\docker-compose.debug.yml"
}

if (!$RunOnly) {
    docker-compose -f .\docker-compose.dependencies.yml -f $file build
}

docker-compose -f .\docker-compose.dependencies.yml -f $file up
