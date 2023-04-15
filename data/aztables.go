package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"phrasmotica/bore-score-api/models"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/google/uuid"
)

type TableStorageDatabase struct {
	Client *aztables.ServiceClient
}

func CreateTableStorageClient(connStr string) *aztables.ServiceClient {
	client, err := aztables.NewServiceClientFromConnectionString(connStr, &aztables.ClientOptions{
		ClientOptions: policy.ClientOptions{},
	})

	if err != nil {
		log.Fatal(err)
	}

	return client
}

// GetAllGames implements IDatabase
func (d *TableStorageDatabase) GetAllGames(ctx context.Context) (bool, []models.Game) {
	games := list(ctx, d.Client, "Games", createGame, nil)
	return true, games
}

func (d *TableStorageDatabase) GetGame(ctx context.Context, name string) (bool, *models.Game) {
	result := d.findGame(ctx, name)
	if result == nil {
		return false, nil
	}

	game := createGame(result)
	return true, &game
}

// GameExists implements IDatabase
func (d *TableStorageDatabase) GameExists(ctx context.Context, name string) bool {
	result := d.findGame(ctx, name)
	return result != nil
}

// AddGame implements IDatabase
func (d *TableStorageDatabase) AddGame(ctx context.Context, newGame *models.Game) bool {
	links, linksErr := json.Marshal(newGame.Links)
	if linksErr != nil {
		log.Println(linksErr)
		return false
	}

	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "Games",
			RowKey:       uuid.NewString(), // TODO: let the client generate this
		},
		Properties: map[string]interface{}{
			"Name":        newGame.Name,
			"TimeCreated": aztables.EDMInt64(newGame.TimeCreated),
			"DisplayName": newGame.DisplayName,
			"Synopsis":    newGame.Synopsis,
			"Description": newGame.Description,
			"MinPlayers":  newGame.MinPlayers,
			"MaxPlayers":  newGame.MaxPlayers,
			"WinMethod":   newGame.WinMethod,
			"ImageLink":   newGame.ImageLink,
			"Links":       string(links),
		},
	}

	marshalled, err := json.Marshal(entity)
	if err != nil {
		log.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("Games").AddEntity(ctx, marshalled, nil)
	if addErr != nil {
		log.Println(addErr)
		return false
	}

	return true
}

// DeleteGame implements IDatabase
func (d *TableStorageDatabase) DeleteGame(ctx context.Context, name string) bool {
	game := d.findGame(ctx, name)
	if game == nil {
		return false
	}

	_, err := d.Client.NewClient("Games").DeleteEntity(ctx, game.PartitionKey, game.RowKey, nil)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// GetAllGroups implements IDatabase
func (d *TableStorageDatabase) GetAllGroups(ctx context.Context) (bool, []models.Group) {
	groups := list(ctx, d.Client, "Groups", createGroup, &aztables.ListEntitiesOptions{
		Filter: to.Ptr("Visibility eq 'public' or Visibility eq 'global'"),
	})

	return true, groups
}

// GetGroups implements IDatabase
func (d *TableStorageDatabase) GetGroups(ctx context.Context) (bool, []models.Group) {
	groups := list(ctx, d.Client, "Groups", createGroup, &aztables.ListEntitiesOptions{
		Filter: to.Ptr("Visibility eq 'public'"),
	})

	return true, groups
}

// GetGroup implements IDatabase
func (d *TableStorageDatabase) GetGroup(ctx context.Context, name string) (RetrieveGroupResult, *models.Group) {
	result := d.findGroup(ctx, name)
	if result == nil {
		return Failure, nil
	}

	group := createGroup(result)

	if group.Visibility == models.Private {
		return Unauthorised, nil
	}

	return Success, &group
}

// GroupExists implements IDatabase
func (d *TableStorageDatabase) GroupExists(ctx context.Context, name string) bool {
	result := d.findGroup(ctx, name)
	return result != nil
}

// AddGroup implements IDatabase
func (d *TableStorageDatabase) AddGroup(ctx context.Context, newGroup *models.Group) bool {
	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "Groups",
			RowKey:       uuid.NewString(), // TODO: let the client generate this
		},
		Properties: map[string]interface{}{
			"Name":           newGroup.Name,
			"TimeCreated":    aztables.EDMInt64(newGroup.TimeCreated),
			"DisplayName":    newGroup.DisplayName,
			"Description":    newGroup.Description,
			"ProfilePicture": newGroup.ProfilePicture,
			"Visibility":     string(newGroup.Visibility),
		},
	}

	marshalled, err := json.Marshal(entity)
	if err != nil {
		log.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("Groups").AddEntity(ctx, marshalled, nil)
	if addErr != nil {
		log.Println(addErr)
		return false
	}

	return true
}

// DeleteGroup implements IDatabase
func (d *TableStorageDatabase) DeleteGroup(ctx context.Context, name string) bool {
	group := d.findGroup(ctx, name)
	if group == nil {
		return false
	}

	_, err := d.Client.NewClient("Groups").DeleteEntity(ctx, group.PartitionKey, group.RowKey, nil)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (d *TableStorageDatabase) GetAllLinkTypes(ctx context.Context) (bool, []models.LinkType) {
	linkTypes := list(ctx, d.Client, "LinkTypes", createLinkType, nil)
	return true, linkTypes
}

// GetAllPlayers implements IDatabase
func (d *TableStorageDatabase) GetAllPlayers(ctx context.Context) (bool, []models.Player) {
	players := list(ctx, d.Client, "Players", createPlayer, nil)
	return true, players
}

// GetPlayer implements IDatabase
func (d *TableStorageDatabase) GetPlayer(ctx context.Context, username string) (bool, *models.Player) {
	result := d.findPlayer(ctx, username)
	if result == nil {
		return false, nil
	}

	player := createPlayer(result)
	return true, &player
}

// PlayerExists implements IDatabase
func (d *TableStorageDatabase) PlayerExists(ctx context.Context, username string) bool {
	result := d.findPlayer(ctx, username)
	return result != nil
}

// AddPlayer implements IDatabase
func (d *TableStorageDatabase) AddPlayer(ctx context.Context, newPlayer *models.Player) bool {
	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "Players",
			RowKey:       uuid.NewString(), // TODO: let the client generate this
		},
		Properties: map[string]interface{}{
			"Username":       newPlayer.Username,
			"TimeCreated":    aztables.EDMInt64(newPlayer.TimeCreated),
			"DisplayName":    newPlayer.DisplayName,
			"ProfilePicture": newPlayer.ProfilePicture,
		},
	}

	marshalled, err := json.Marshal(entity)
	if err != nil {
		log.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("Players").AddEntity(ctx, marshalled, nil)
	if addErr != nil {
		log.Println(addErr)
		return false
	}

	return true
}

// DeletePlayer implements IDatabase
func (d *TableStorageDatabase) DeletePlayer(ctx context.Context, username string) bool {
	player := d.findPlayer(ctx, username)
	if player == nil {
		return false
	}

	_, err := d.Client.NewClient("Players").DeleteEntity(ctx, player.PartitionKey, player.RowKey, nil)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// GetAllResults implements IDatabase
func (d *TableStorageDatabase) GetAllResults(ctx context.Context) (bool, []models.Result) {
	results := list(ctx, d.Client, "Results", createResult, nil)
	return true, results
}

// AddResult implements IDatabase
func (d *TableStorageDatabase) AddResult(ctx context.Context, newResult *models.Result) bool {
	scores, scoresErr := json.Marshal(newResult.Scores)
	if scoresErr != nil {
		log.Println(scoresErr)
		return false
	}

	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "Results",
			RowKey:       uuid.NewString(), // TODO: let the client generate this
		},
		Properties: map[string]interface{}{
			"GameName":         newResult.GameName,
			"GroupName":        newResult.GroupName,
			"TimeCreated":      aztables.EDMInt64(newResult.TimeCreated),
			"TimePlayed":       aztables.EDMInt64(newResult.TimePlayed),
			"Notes":            newResult.Notes,
			"CooperativeScore": newResult.CooperativeScore,
			"CooperativeWin":   newResult.CooperativeWin,
			"Scores":           string(scores),
		},
	}

	marshalled, err := json.Marshal(entity)
	if err != nil {
		log.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("Results").AddEntity(ctx, marshalled, nil)
	if addErr != nil {
		log.Println(addErr)
		return false
	}

	return true
}

// DeleteResultsWithGame implements IDatabase
func (d *TableStorageDatabase) DeleteResultsWithGame(ctx context.Context, gameName string) (bool, int64) {
	game := d.findGame(ctx, gameName)
	if game == nil {
		return false, 0
	}

	client := d.Client.NewClient("Results")
	entities := listEntities(ctx, client, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("GameName eq '%s'", gameName)),
	})

	if len(entities) <= 0 {
		return false, 0
	}

	deleteCount := 0
	for i := 0; i < len(entities); i++ {
		result := entities[i]

		_, err := client.DeleteEntity(ctx, result.PartitionKey, result.RowKey, nil)
		if err != nil {
			log.Println(err)
		} else {
			deleteCount++
		}
	}

	return true, int64(deleteCount)
}

// ScrubResultsWithPlayer implements IDatabase
func (d *TableStorageDatabase) ScrubResultsWithPlayer(ctx context.Context, username string) (bool, int64) {
	// TODO: restructure data so that we can find the results containing this player more easily.
	// filter expressions don't support a string "contains" operator, so we have to fetch
	// all results and then filter them afterwards...

	success, results := d.GetAllResults(ctx)
	if !success {
		return false, 0
	}

	relevantResults := []models.Result{}

	// pick out the results that this player was involved in
	for i := range results {
		scores := results[i].Scores
		for j := range scores {
			if scores[j].Username == username {
				relevantResults = append(relevantResults, results[i])
			}
		}
	}

	if len(relevantResults) <= 0 {
		return false, 0
	}

	client := d.Client.NewClient("Results")

	updateCount := 0
	for i := 0; i < len(relevantResults); i++ {
		result := relevantResults[i]

		scores := result.Scores
		for j := range scores {
			if scores[j].Username == username {
				result.Scores[j].Username = ""
			}
		}

		marshalledScores, scoresErr := json.Marshal(result.Scores)
		if scoresErr != nil {
			log.Println(scoresErr)
			continue
		}

		// create new entity with scrubbed scores data for merging
		entity := aztables.EDMEntity{
			Entity: aztables.Entity{
				PartitionKey: "Results",
				RowKey:       result.ID,
			},
			Properties: map[string]interface{}{
				"Scores": string(marshalledScores),
			},
		}

		marshalled, err := json.Marshal(entity)
		if err != nil {
			log.Println(err)
			continue
		}

		_, updateErr := client.UpdateEntity(ctx, marshalled, nil)
		if updateErr != nil {
			log.Println(updateErr)
		} else {
			updateCount++
		}
	}

	return true, int64(updateCount)
}

func (d *TableStorageDatabase) GetAllWinMethods(ctx context.Context) (bool, []models.WinMethod) {
	winMethods := list(ctx, d.Client, "WinMethods", createWinMethod, nil)
	return true, winMethods
}

// GetSummary implements IDatabase
func (d *TableStorageDatabase) GetSummary(ctx context.Context) (bool, *Summary) {
	gamesSuccess, games := d.GetAllGames(ctx)
	if !gamesSuccess {
		log.Println("Failed to get all games for summary")
		return false, nil
	}

	groupsSuccess, groups := d.GetAllGroups(ctx)
	if !groupsSuccess {
		log.Println("Failed to get all groups for summary")
		return false, nil
	}

	playersSuccess, players := d.GetAllPlayers(ctx)
	if !playersSuccess {
		log.Println("Failed to get all players for summary")
		return false, nil
	}

	resultsSuccess, results := d.GetAllResults(ctx)
	if !resultsSuccess {
		log.Println("Failed to get all results for summary")
		return false, nil
	}

	return true, &Summary{
		GameCount:   int64(len(games)),
		GroupCount:  int64(len(groups)),
		PlayerCount: int64(len(players)),
		ResultCount: int64(len(results)),
	}
}

func (d *TableStorageDatabase) findGame(ctx context.Context, name string) *aztables.EDMEntity {
	client := d.Client.NewClient("Games")

	entities := listEntities(ctx, client, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("Name eq '%s'", name)),
	})

	if len(entities) == 1 {
		return &entities[0]
	}

	return nil
}

func (d *TableStorageDatabase) findGroup(ctx context.Context, name string) *aztables.EDMEntity {
	client := d.Client.NewClient("Groups")

	entities := listEntities(ctx, client, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("Name eq '%s'", name)),
	})

	if len(entities) == 1 {
		return &entities[0]
	}

	return nil
}

func (d *TableStorageDatabase) findPlayer(ctx context.Context, username string) *aztables.EDMEntity {
	client := d.Client.NewClient("Players")

	entities := listEntities(ctx, client, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("Username eq '%s'", username)),
	})

	if len(entities) == 1 {
		return &entities[0]
	}

	return nil
}

func list[T interface{}](ctx context.Context, client *aztables.ServiceClient, tableName string, convert func(*aztables.EDMEntity) T, options *aztables.ListEntitiesOptions) []T {
	entities := listEntities(ctx, client.NewClient(tableName), options)
	data := []T{}

	for i := range entities {
		data = append(data, convert(&entities[i]))
	}

	return data
}

func listEntities(ctx context.Context, client *aztables.Client, options *aztables.ListEntitiesOptions) []aztables.EDMEntity {
	var entities = make([]aztables.EDMEntity, 0)

	client.CreateTable(ctx, nil)

	pager := client.NewListEntitiesPager(options)

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range response.Entities {
			entity := unmarshal(e)
			entities = append(entities, *entity)
		}
	}

	return entities
}

func unmarshal(bytes []byte) *aztables.EDMEntity {
	var entity aztables.EDMEntity

	err := json.Unmarshal(bytes, &entity)
	if err != nil {
		log.Fatal(err)
	}

	return &entity
}

func createGame(entity *aztables.EDMEntity) models.Game {
	return models.Game{
		ID:          entity.RowKey,
		Name:        propString(entity, "Name"),
		TimeCreated: propInt64(entity, "TimeCreated"),
		DisplayName: propString(entity, "DisplayName"),
		Synopsis:    propString(entity, "Synopsis"),
		Description: propString(entity, "Description"),
		MinPlayers:  propInt(entity, "MinPlayers"),
		MaxPlayers:  propInt(entity, "MaxPlayers"),
		WinMethod:   propString(entity, "WinMethod"),
		ImageLink:   propString(entity, "ImageLink"),
		Links:       createLinks(entity),
	}
}

// returns an array of Link objects by converting the JSON string in the
// table entity's "Links" column
func createLinks(entity *aztables.EDMEntity) []models.Link {
	linksStr := propString(entity, "Links")

	data := []models.Link{}
	json.Unmarshal([]byte(linksStr), &data)

	return data
}

func createGroup(entity *aztables.EDMEntity) models.Group {
	return models.Group{
		ID:             entity.RowKey,
		Name:           propString(entity, "Name"),
		TimeCreated:    propInt64(entity, "TimeCreated"),
		DisplayName:    propString(entity, "DisplayName"),
		Description:    propString(entity, "Description"),
		ProfilePicture: propString(entity, "ProfilePicture"),
		Visibility:     models.GroupVisibilityName(propString(entity, "Visibility")),
	}
}

func createLinkType(entity *aztables.EDMEntity) models.LinkType {
	return models.LinkType{
		ID:          entity.RowKey,
		Name:        models.LinkTypeName(propString(entity, "Name")),
		TimeCreated: propInt64(entity, "TimeCreated"),
		DisplayName: propString(entity, "DisplayName"),
	}
}

func createPlayer(entity *aztables.EDMEntity) models.Player {
	return models.Player{
		ID:             entity.RowKey,
		Username:       propString(entity, "Username"),
		TimeCreated:    propInt64(entity, "TimeCreated"),
		DisplayName:    propString(entity, "DisplayName"),
		ProfilePicture: propString(entity, "ProfilePicture"),
	}
}

func createResult(entity *aztables.EDMEntity) models.Result {
	return models.Result{
		ID:               entity.RowKey,
		GameName:         propString(entity, "GameName"),
		GroupName:        propString(entity, "GroupName"),
		TimeCreated:      propInt64(entity, "TimeCreated"),
		TimePlayed:       propInt64(entity, "TimePlayed"),
		Notes:            propString(entity, "Notes"),
		CooperativeScore: propInt(entity, "CooperativeScore"),
		CooperativeWin:   propBool(entity, "CooperativeWin"),
		Scores:           createScores(entity),
	}
}

// returns an array of PlayerScore objects by converting the JSON string in the
// table entity's "Scores" column
func createScores(entity *aztables.EDMEntity) []models.PlayerScore {
	scoresStr := propString(entity, "Scores")

	data := []models.PlayerScore{}
	json.Unmarshal([]byte(scoresStr), &data)

	return data
}

func createWinMethod(entity *aztables.EDMEntity) models.WinMethod {
	return models.WinMethod{
		ID:          entity.RowKey,
		Name:        models.WinMethodName(propString(entity, "Name")),
		TimeCreated: propInt64(entity, "TimeCreated"),
		DisplayName: propString(entity, "DisplayName"),
	}
}

func propString(entity *aztables.EDMEntity, name string) string {
	return entity.Properties[name].(string)
}

func propBool(entity *aztables.EDMEntity, name string) bool {
	return entity.Properties[name].(bool)
}

func propInt(entity *aztables.EDMEntity, name string) int {
	return int(entity.Properties[name].(int32))
}

func propInt64(entity *aztables.EDMEntity, name string) int64 {
	return int64(entity.Properties[name].(aztables.EDMInt64))
}