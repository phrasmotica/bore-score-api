param(
    [switch] $Debugger = $false,
    [switch] $BackendOnly = $false,
    [switch] $DependenciesOnly = $false,
    [switch] $BuildOnly = $false,
    [switch] $RunOnly = $false
)

$file = ".\docker-compose.yml"
$serviceName = "backend"

if ($Debugger) {
    $file = ".\docker-compose.debug.yml"
    $serviceName = "backend-debug"
}

if (!$RunOnly) {
    if ($BackendOnly) {
        docker-compose -f .\docker-compose.dependencies.yml -f $file --env-file .\.env.docker build $serviceName
    }

    if ($DependenciesOnly) {
        docker-compose -f .\docker-compose.dependencies.yml .\.env.docker build
    }
}

if (!$BuildOnly) {
    if ($BackendOnly) {
        docker-compose -f .\docker-compose.dependencies.yml -f $file --env-file .\.env.docker up $serviceName
    }

    if ($DependenciesOnly) {
        docker-compose -f .\docker-compose.dependencies.yml .\.env.docker up
    }
}
