param(
    [switch] $Push
)

docker build -t phrasmotica/bore-score-api .

if ($Push) {
    docker push phrasmotica/bore-score-api
}
