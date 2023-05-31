package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"phrasmotica/bore-score-api/models"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type TableStorageDatabase struct {
	Client *aztables.ServiceClient
}

// TODO: put this in a more central place, or inject it as a dependency
var (
	Error *log.Logger = log.New(os.Stdout, "ERROR: ", log.LstdFlags|log.Lshortfile)
)

func CreateTableStorageClient(connStr string) *aztables.ServiceClient {
	client, err := aztables.NewServiceClientFromConnectionString(connStr, &aztables.ClientOptions{
		ClientOptions: policy.ClientOptions{},
	})

	if err != nil {
		Error.Fatal(err)
		return nil
	}

	return client
}

// AddApproval implements IDatabase
func (d *TableStorageDatabase) AddApproval(ctx context.Context, newApproval *models.Approval) bool {
	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: newApproval.ResultID,
			RowKey:       newApproval.ID,
		},
		Properties: map[string]interface{}{
			"ResultID":       newApproval.ResultID,
			"TimeCreated":    aztables.EDMInt64(newApproval.TimeCreated),
			"Username":       newApproval.Username,
			"ApprovalStatus": string(newApproval.ApprovalStatus),
		},
	}

	marshalled, err := json.Marshal(entity)
	if err != nil {
		Error.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("Approvals").AddEntity(ctx, marshalled, nil)
	if addErr != nil {
		Error.Println(addErr)
		return false
	}

	return true
}

// GetApprovals implements IDatabase
func (d *TableStorageDatabase) GetApprovals(ctx context.Context, resultId string) (bool, []models.Approval) {
	approvals := list(ctx, d.Client, "Approvals", createApproval, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("ResultID eq '%s'", resultId)),
	})
	return true, approvals
}

// GetAllGames implements IDatabase
func (d *TableStorageDatabase) GetAllGames(ctx context.Context) (bool, []models.Game) {
	games := list(ctx, d.Client, "Games", createGame, nil)
	return true, games
}

// GetGameByName implements IDatabase
func (d *TableStorageDatabase) GetGame(ctx context.Context, id string) (bool, *models.Game) {
	result := d.findGame(ctx, id)
	if result == nil {
		return false, nil
	}

	game := createGame(result)
	return true, &game
}

// GameExists implements IDatabase
func (d *TableStorageDatabase) GameExists(ctx context.Context, id string) bool {
	result := d.findGame(ctx, id)
	return result != nil
}

// AddGame implements IDatabase
func (d *TableStorageDatabase) AddGame(ctx context.Context, newGame *models.Game) bool {
	links, linksErr := json.Marshal(newGame.Links)
	if linksErr != nil {
		Error.Println(linksErr)
		return false
	}

	// TODO: add "CreatedBy" column and use it as partition key
	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "Games",
			RowKey:       newGame.ID,
		},
		Properties: map[string]interface{}{
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
		Error.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("Games").AddEntity(ctx, marshalled, nil)
	if addErr != nil {
		Error.Println(addErr)
		return false
	}

	return true
}

// DeleteGame implements IDatabase
func (d *TableStorageDatabase) DeleteGame(ctx context.Context, id string) bool {
	game := d.findGame(ctx, id)
	if game == nil {
		return false
	}

	_, err := d.Client.NewClient("Games").DeleteEntity(ctx, game.PartitionKey, game.RowKey, nil)
	if err != nil {
		Error.Println(err)
		return false
	}

	return true
}

// GetAllGroups implements IDatabase
func (d *TableStorageDatabase) GetAllGroups(ctx context.Context) (bool, []models.Group) {
	groups := list(ctx, d.Client, "Groups", createGroup, nil)
	return true, groups
}

// GetGroups implements IDatabase
func (d *TableStorageDatabase) GetGroups(ctx context.Context) (bool, []models.Group) {
	groups := list(ctx, d.Client, "Groups", createGroup, &aztables.ListEntitiesOptions{
		Filter: to.Ptr("Visibility ne 'global'"),
	})

	return true, groups
}

// GetGroup implements IDatabase
func (d *TableStorageDatabase) GetGroup(ctx context.Context, id string) (bool, *models.Group) {
	entity := d.findGroup(ctx, id)
	if entity == nil {
		return false, nil
	}

	group := createGroup(entity)
	return true, &group
}

// GetGroup implements IDatabase
func (d *TableStorageDatabase) GetGroupByName(ctx context.Context, name string) (bool, *models.Group) {
	entity := d.findGroupByName(ctx, name)
	if entity == nil {
		return false, nil
	}

	group := createGroup(entity)
	return true, &group
}

// GroupExists implements IDatabase
func (d *TableStorageDatabase) GroupExists(ctx context.Context, name string) bool {
	result := d.findGroupByName(ctx, name)
	return result != nil
}

// AddGroup implements IDatabase
func (d *TableStorageDatabase) AddGroup(ctx context.Context, newGroup *models.Group) bool {
	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: newGroup.CreatedBy,
			RowKey:       newGroup.ID,
		},
		Properties: map[string]interface{}{
			"TimeCreated":    aztables.EDMInt64(newGroup.TimeCreated),
			"DisplayName":    newGroup.DisplayName,
			"Description":    newGroup.Description,
			"ProfilePicture": newGroup.ProfilePicture,
			"Visibility":     string(newGroup.Visibility),
			"CreatedBy":      newGroup.CreatedBy,
		},
	}

	marshalled, err := json.Marshal(entity)
	if err != nil {
		Error.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("Groups").AddEntity(ctx, marshalled, nil)
	if addErr != nil {
		Error.Println(addErr)
		return false
	}

	return true
}

// DeleteGroup implements IDatabase
func (d *TableStorageDatabase) DeleteGroup(ctx context.Context, id string) bool {
	group := d.findGroup(ctx, id)
	if group == nil {
		return false
	}

	_, err := d.Client.NewClient("Groups").DeleteEntity(ctx, group.PartitionKey, group.RowKey, nil)
	if err != nil {
		Error.Println(err)
		return false
	}

	return true
}

// GetGroupInvitation implements IDatabase
func (d *TableStorageDatabase) GetGroupInvitation(ctx context.Context, invitationId string) (bool, *models.GroupInvitation) {
	result := d.findGroupInvitation(ctx, invitationId)
	if result == nil {
		return false, nil
	}

	invitation := createGroupInvitation(result)
	return true, &invitation
}

// GetGroupInvitations implements IDatabase
func (d *TableStorageDatabase) GetGroupInvitations(ctx context.Context, username string) (bool, []models.GroupInvitation) {
	if !d.UserExists(ctx, username) {
		return false, []models.GroupInvitation{}
	}

	invitations := list(ctx, d.Client, "GroupInvitations", createGroupInvitation, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("Username eq '%s'", username)),
	})

	return true, invitations
}

// GetGroupInvitationsForGroup implements IDatabase
func (d *TableStorageDatabase) GetGroupInvitationsForGroup(ctx context.Context, groupId string) (bool, []models.GroupInvitation) {
	success, group := d.GetGroup(ctx, groupId)
	if !success {
		return false, []models.GroupInvitation{}
	}

	invitations := list(ctx, d.Client, "GroupInvitations", createGroupInvitation, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("GroupID eq '%s'", group.ID)),
	})

	return true, invitations
}

// IsInvitedToGroup implements IDatabase
func (d *TableStorageDatabase) IsInvitedToGroup(ctx context.Context, groupId string, username string) bool {
	success, invitations := d.GetGroupInvitations(ctx, username)

	return success && slices.ContainsFunc(invitations, func(i models.GroupInvitation) bool {
		return i.GroupID == groupId && i.Username == username
	})
}

// AddGroupInvitation implements IDatabase
func (d *TableStorageDatabase) AddGroupInvitation(ctx context.Context, newGroupInvitation *models.GroupInvitation) bool {
	newGroupInvitation.ID = uuid.NewString()
	newGroupInvitation.TimeCreated = time.Now().UTC().Unix()
	newGroupInvitation.InvitationStatus = models.Sent

	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: newGroupInvitation.GroupID,
			RowKey:       newGroupInvitation.ID,
		},
		Properties: map[string]interface{}{
			"GroupID":          newGroupInvitation.GroupID,
			"TimeCreated":      aztables.EDMInt64(newGroupInvitation.TimeCreated),
			"Username":         newGroupInvitation.Username,
			"InviterUsername":  newGroupInvitation.InviterUsername,
			"InvitationStatus": string(newGroupInvitation.InvitationStatus),
		},
	}

	marshalled, err := json.Marshal(entity)
	if err != nil {
		Error.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("GroupInvitations").AddEntity(ctx, marshalled, nil)
	if addErr != nil {
		Error.Println(addErr)
		return false
	}

	return true
}

// UpdateGroupInvitation implements IDatabase
func (d *TableStorageDatabase) UpdateGroupInvitation(ctx context.Context, newGroupInvitation *models.GroupInvitation) bool {
	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: newGroupInvitation.GroupID,
			RowKey:       newGroupInvitation.ID,
		},
		Properties: map[string]interface{}{
			"GroupID":          newGroupInvitation.GroupID,
			"TimeCreated":      aztables.EDMInt64(newGroupInvitation.TimeCreated),
			"Username":         newGroupInvitation.Username,
			"InviterUsername":  newGroupInvitation.InviterUsername,
			"InvitationStatus": string(newGroupInvitation.InvitationStatus),
		},
	}

	marshalled, err := json.Marshal(entity)
	if err != nil {
		Error.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("GroupInvitations").UpdateEntity(ctx, marshalled, nil)
	if addErr != nil {
		Error.Println(addErr)
		return false
	}

	return true
}

// GetGroupMemberships implements IDatabase
func (d *TableStorageDatabase) GetGroupMemberships(ctx context.Context, username string) (bool, []models.GroupMembership) {
	if !d.UserExists(ctx, username) {
		return false, []models.GroupMembership{}
	}

	memberships := list(ctx, d.Client, "GroupMemberships", createGroupMembership, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("Username eq '%s'", username)),
	})

	return true, memberships
}

// GetGroupMembershipsForGroup implements IDatabase
func (d *TableStorageDatabase) GetGroupMembershipsForGroup(ctx context.Context, groupId string) (bool, []models.GroupMembership) {
	success, group := d.GetGroup(ctx, groupId)
	if !success {
		return false, []models.GroupMembership{}
	}

	memberships := list(ctx, d.Client, "GroupMemberships", createGroupMembership, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("GroupID eq '%s'", group.ID)),
	})

	return true, memberships
}

// IsInGroup implements IDatabase
func (d *TableStorageDatabase) IsInGroup(ctx context.Context, groupId string, username string) bool {
	success, memberships := d.GetGroupMemberships(ctx, username)
	if !success {
		return false
	}

	return success && slices.ContainsFunc(memberships, func(m models.GroupMembership) bool {
		return m.GroupID == groupId && m.Username == username
	})
}

// AddGroupMembership implements IDatabase
func (d *TableStorageDatabase) AddGroupMembership(ctx context.Context, newGroupMembership *models.GroupMembership) bool {
	newGroupMembership.ID = uuid.NewString()
	newGroupMembership.TimeCreated = time.Now().UTC().Unix()

	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: newGroupMembership.GroupID,
			RowKey:       newGroupMembership.ID,
		},
		Properties: map[string]interface{}{
			"GroupID":      newGroupMembership.GroupID,
			"TimeCreated":  aztables.EDMInt64(newGroupMembership.TimeCreated),
			"Username":     newGroupMembership.Username,
			"InvitationID": newGroupMembership.InvitationID,
		},
	}

	marshalled, err := json.Marshal(entity)
	if err != nil {
		Error.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("GroupMemberships").AddEntity(ctx, marshalled, nil)
	if addErr != nil {
		Error.Println(addErr)
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

// GetPlayersInGroup implements IDatabase
func (d *TableStorageDatabase) GetPlayersInGroup(ctx context.Context, groupId string) (bool, []models.Player) {
	success, players := d.GetAllPlayers(ctx)
	if !success {
		return false, []models.Player{}
	}

	success, memberships := d.GetGroupMembershipsForGroup(ctx, groupId)
	if !success {
		return false, []models.Player{}
	}

	playersInGroup := []models.Player{}

	for _, m := range memberships {
		// returns the index of the player with this membership
		playerIndex := slices.IndexFunc(players, func(p models.Player) bool {
			return p.Username == m.Username
		})

		if playerIndex >= 0 {
			playersInGroup = append(playersInGroup, players[playerIndex])
		}
	}

	return true, playersInGroup
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
			RowKey:       newPlayer.ID,
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
		Error.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("Players").AddEntity(ctx, marshalled, nil)
	if addErr != nil {
		Error.Println(addErr)
		return false
	}

	return true
}

// UpdatePlayer implements IDatabase
func (d *TableStorageDatabase) UpdatePlayer(ctx context.Context, player *models.Player) bool {
	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "Players",
			RowKey:       player.ID,
		},
		Properties: map[string]interface{}{
			"DisplayName":    player.DisplayName,
			"ProfilePicture": player.ProfilePicture,
		},
	}

	marshalled, err := json.Marshal(entity)
	if err != nil {
		Error.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("Players").UpdateEntity(ctx, marshalled, nil)
	if addErr != nil {
		Error.Println(addErr)
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
		Error.Println(err)
		return false
	}

	return true
}

// GetAllResults implements IDatabase
func (d *TableStorageDatabase) GetAllResults(ctx context.Context) (bool, []models.Result) {
	results := list(ctx, d.Client, "Results", createResult, nil)
	return true, results
}

// GetResultsForGroup implements IDatabase
func (d *TableStorageDatabase) GetResultsForGroup(ctx context.Context, groupId string) (bool, []models.Result) {
	results := list(ctx, d.Client, "Results", createResult, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("GroupID eq '%s'", groupId)),
	})
	return true, results
}

// GetResultsForGroupAndGame implements IDatabase
func (d *TableStorageDatabase) GetResultsForGroupAndGame(ctx context.Context, groupId string, gameId string) (bool, []models.Result) {
	results := list(ctx, d.Client, "Results", createResult, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("GroupID eq '%s' and GameID eq '%s'", groupId, gameId)),
	})
	return true, results
}

// GetResultsWithPlayer implements IDatabase
func (d *TableStorageDatabase) GetResultsWithPlayer(ctx context.Context, username string) (bool, []models.Result) {
	// TODO: restructure data so that we can find the results containing this player more easily.
	// filter expressions don't support a string "contains" operator, so we have to fetch
	// all results and then filter them afterwards...

	success, results := d.GetAllResults(ctx)
	if !success {
		return false, []models.Result{}
	}

	relevantResults := []models.Result{}

	// pick out the results that this player was involved in
	for i := range results {
		// TODO: use slices.ContainsFunc(...) to check
		scores := results[i].Scores
		for j := range scores {
			if scores[j].Username == username {
				relevantResults = append(relevantResults, results[i])
			}
		}
	}

	return true, relevantResults
}

// GetResult implements IDatabase
func (d *TableStorageDatabase) GetResult(ctx context.Context, id string) (bool, *models.Result) {
	entity := d.findResult(ctx, id)
	if entity == nil {
		return false, nil
	}

	result := createResult(entity)
	return true, &result
}

// ResultExists implements IDatabase
func (d *TableStorageDatabase) ResultExists(ctx context.Context, id string) bool {
	entity := d.findResult(ctx, id)
	return entity != nil
}

// AddResult implements IDatabase
func (d *TableStorageDatabase) AddResult(ctx context.Context, newResult *models.Result) bool {
	scores, scoresErr := json.Marshal(newResult.Scores)
	if scoresErr != nil {
		Error.Println(scoresErr)
		return false
	}

	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: newResult.GameID,
			RowKey:       newResult.ID,
		},
		Properties: map[string]interface{}{
			"GameID":           newResult.GameID,
			"GroupID":          newResult.GroupID,
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
		Error.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("Results").AddEntity(ctx, marshalled, nil)
	if addErr != nil {
		Error.Println(addErr)
		return false
	}

	return true
}

// DeleteResultsWithGame implements IDatabase
func (d *TableStorageDatabase) DeleteResultsWithGame(ctx context.Context, gameId string) (bool, int64) {
	game := d.findGame(ctx, gameId)
	if game == nil {
		return false, 0
	}

	client := d.Client.NewClient("Results")
	entities := listEntities(ctx, client, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("GameID eq '%s'", game.ID)),
	})

	deleteCount := 0
	for i := 0; i < len(entities); i++ {
		result := entities[i]

		_, err := client.DeleteEntity(ctx, result.PartitionKey, result.RowKey, nil)
		if err != nil {
			Error.Println(err)
		} else {
			deleteCount++
		}
	}

	return true, int64(deleteCount)
}

// ScrubResultsWithPlayer implements IDatabase
func (d *TableStorageDatabase) ScrubResultsWithPlayer(ctx context.Context, username string) (bool, int64) {
	success, relevantResults := d.GetResultsWithPlayer(ctx, username)

	if !success {
		return false, 0
	}

	client := d.Client.NewClient("Results")

	updateCount := 0
	for i := 0; i < len(relevantResults); i++ {
		result := relevantResults[i]

		scores := result.Scores
		for j := range scores {
			// TODO: use slices.ContainsFunc(...) to check
			if scores[j].Username == username {
				result.Scores[j].Username = ""
			}
		}

		marshalledScores, scoresErr := json.Marshal(result.Scores)
		if scoresErr != nil {
			Error.Println(scoresErr)
			continue
		}

		// create new entity with scrubbed scores data for merging
		entity := aztables.EDMEntity{
			Entity: aztables.Entity{
				PartitionKey: result.GameID,
				RowKey:       result.ID,
			},
			Properties: map[string]interface{}{
				"Scores": string(marshalledScores),
			},
		}

		marshalled, err := json.Marshal(entity)
		if err != nil {
			Error.Println(err)
			continue
		}

		_, updateErr := client.UpdateEntity(ctx, marshalled, nil)
		if updateErr != nil {
			Error.Println(updateErr)
		} else {
			updateCount++
		}
	}

	return true, int64(updateCount)
}

// GetUser implements IDatabase
func (d *TableStorageDatabase) GetUser(ctx context.Context, username string) (bool, *models.User) {
	result := d.findUser(ctx, username)
	if result == nil {
		return false, nil
	}

	user := createUser(result)
	return true, &user
}

// GetUserByEmail implements IDatabase
func (d *TableStorageDatabase) GetUserByEmail(ctx context.Context, email string) (bool, *models.User) {
	result := d.findUserByEmail(ctx, email)
	if result == nil {
		return false, nil
	}

	user := createUser(result)
	return true, &user
}

func (d *TableStorageDatabase) AddUser(ctx context.Context, newUser *models.User) bool {
	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "Users",
			RowKey:       newUser.ID,
		},
		Properties: map[string]interface{}{
			"Username":    newUser.Username,
			"TimeCreated": aztables.EDMInt64(newUser.TimeCreated),
			"Email":       newUser.Email,
			"Password":    newUser.Password,
			"Permissions": strings.Join(newUser.Permissions, ";"),
		},
	}

	marshalled, err := json.Marshal(entity)
	if err != nil {
		Error.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("Users").AddEntity(ctx, marshalled, nil)
	if addErr != nil {
		Error.Println(addErr)
		return false
	}

	return true
}

// UserExists implements IDatabase
func (d *TableStorageDatabase) UserExists(ctx context.Context, username string) bool {
	result := d.findUser(ctx, username)
	return result != nil
}

// UserExistsByEmail implements IDatabase
func (d *TableStorageDatabase) UserExistsByEmail(ctx context.Context, email string) bool {
	result := d.findUserByEmail(ctx, email)
	return result != nil
}

// UpdateUser implements IDatabase
func (d *TableStorageDatabase) UpdateUser(ctx context.Context, user *models.User) bool {
	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "Users",
			RowKey:       user.ID,
		},
		Properties: map[string]interface{}{
			"Password": user.Password,
		},
	}

	marshalled, err := json.Marshal(entity)
	if err != nil {
		Error.Println(err)
		return false
	}

	_, addErr := d.Client.NewClient("Users").UpdateEntity(ctx, marshalled, nil)
	if addErr != nil {
		Error.Println(addErr)
		return false
	}

	return true
}

func (d *TableStorageDatabase) GetAllWinMethods(ctx context.Context) (bool, []models.WinMethod) {
	winMethods := list(ctx, d.Client, "WinMethods", createWinMethod, nil)
	return true, winMethods
}

// GetSummary implements IDatabase
func (d *TableStorageDatabase) GetSummary(ctx context.Context) (bool, *Summary) {
	gamesSuccess, games := d.GetAllGames(ctx)
	if !gamesSuccess {
		Error.Println("Failed to get all games for summary")
		return false, nil
	}

	groupsSuccess, groups := d.GetAllGroups(ctx)
	if !groupsSuccess {
		Error.Println("Failed to get all groups for summary")
		return false, nil
	}

	playersSuccess, players := d.GetAllPlayers(ctx)
	if !playersSuccess {
		Error.Println("Failed to get all players for summary")
		return false, nil
	}

	resultsSuccess, results := d.GetAllResults(ctx)
	if !resultsSuccess {
		Error.Println("Failed to get all results for summary")
		return false, nil
	}

	return true, &Summary{
		GameCount:   int64(len(games)),
		GroupCount:  int64(len(groups)),
		PlayerCount: int64(len(players)),
		ResultCount: int64(len(results)),
	}
}

func (d *TableStorageDatabase) findGame(ctx context.Context, id string) *aztables.EDMEntity {
	client := d.Client.NewClient("Games")

	entities := listEntities(ctx, client, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("RowKey eq '%s'", id)),
	})

	if len(entities) == 1 {
		return &entities[0]
	}

	return nil
}

func (d *TableStorageDatabase) findGroup(ctx context.Context, id string) *aztables.EDMEntity {
	client := d.Client.NewClient("Groups")

	entities := listEntities(ctx, client, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("RowKey eq '%s'", id)),
	})

	if len(entities) == 1 {
		return &entities[0]
	}

	return nil
}

func (d *TableStorageDatabase) findGroupByName(ctx context.Context, name string) *aztables.EDMEntity {
	client := d.Client.NewClient("Groups")

	entities := listEntities(ctx, client, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("Name eq '%s'", name)),
	})

	if len(entities) == 1 {
		return &entities[0]
	}

	return nil
}

func (d *TableStorageDatabase) findGroupInvitation(ctx context.Context, id string) *aztables.EDMEntity {
	client := d.Client.NewClient("GroupInvitations")

	entities := listEntities(ctx, client, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("RowKey eq '%s'", id)),
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

func (d *TableStorageDatabase) findResult(ctx context.Context, id string) *aztables.EDMEntity {
	client := d.Client.NewClient("Results")

	entities := listEntities(ctx, client, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("RowKey eq '%s'", id)),
	})

	if len(entities) == 1 {
		return &entities[0]
	}

	return nil
}

func (d *TableStorageDatabase) findUser(ctx context.Context, username string) *aztables.EDMEntity {
	client := d.Client.NewClient("Users")

	entities := listEntities(ctx, client, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("Username eq '%s'", username)),
	})

	if len(entities) == 1 {
		return &entities[0]
	}

	return nil
}

func (d *TableStorageDatabase) findUserByEmail(ctx context.Context, email string) *aztables.EDMEntity {
	client := d.Client.NewClient("Users")

	entities := listEntities(ctx, client, &aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("Email eq '%s'", email)),
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

	// TODO: don't do this if it already exists
	client.CreateTable(ctx, nil)

	pager := client.NewListEntitiesPager(options)

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			Error.Fatal(err)
			return entities
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
		Error.Fatal(err)
		return nil
	}

	return &entity
}

func createApproval(entity *aztables.EDMEntity) models.Approval {
	return models.Approval{
		ID:             entity.RowKey,
		ResultID:       propString(entity, "ResultID"),
		TimeCreated:    propInt64(entity, "TimeCreated"),
		Username:       propString(entity, "Username"),
		ApprovalStatus: models.ApprovalStatus(propString(entity, "ApprovalStatus")),
	}
}

func createGame(entity *aztables.EDMEntity) models.Game {
	return models.Game{
		ID:          entity.RowKey,
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
		TimeCreated:    propInt64(entity, "TimeCreated"),
		DisplayName:    propString(entity, "DisplayName"),
		Description:    propString(entity, "Description"),
		ProfilePicture: propString(entity, "ProfilePicture"),
		Visibility:     models.GroupVisibilityName(propString(entity, "Visibility")),
		CreatedBy:      propString(entity, "CreatedBy"),
	}
}

func createGroupInvitation(entity *aztables.EDMEntity) models.GroupInvitation {
	return models.GroupInvitation{
		ID:               entity.RowKey,
		GroupID:          propString(entity, "GroupID"),
		TimeCreated:      propInt64(entity, "TimeCreated"),
		Username:         propString(entity, "Username"),
		InviterUsername:  propString(entity, "InviterUsername"),
		InvitationStatus: models.InvitationStatus(propString(entity, "InvitationStatus")),
	}
}

func createGroupMembership(entity *aztables.EDMEntity) models.GroupMembership {
	return models.GroupMembership{
		ID:           entity.RowKey,
		GroupID:      propString(entity, "GroupID"),
		TimeCreated:  propInt64(entity, "TimeCreated"),
		Username:     propString(entity, "Username"),
		InvitationID: propString(entity, "InvitationID"),
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
		GameID:           propString(entity, "GameID"),
		GroupID:          propString(entity, "GroupID"),
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

func createUser(entity *aztables.EDMEntity) models.User {
	return models.User{
		ID:          entity.RowKey,
		Username:    propString(entity, "Username"),
		TimeCreated: propInt64(entity, "TimeCreated"),
		Email:       propString(entity, "Email"),
		Password:    propString(entity, "Password"),
		Permissions: strings.Split(propString(entity, "Permissions"), ";"),
	}
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
