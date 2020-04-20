package handler

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/store"

	pb "github.com/micro/services/teams/service/proto/teams"
)

// Teams implements the teams service interface
type Teams struct {
	name  string
	store store.Store
}

// New returns an initialized teams handler
func New(service micro.Service) *Teams {
	return &Teams{
		name:  service.Name(),
		store: store.DefaultStore,
	}
}

const (
	// teamsPrefix is the store prefix for teams. Teams are stored with
	// keys in the following format "teams/{namespace}/{id}". This allows
	// us to lookup teams using both namespace and ID. Namespace is used
	// more commonly than ID so we'll use this as the first component of
	// the key.
	teamsPrefix = "teams/"
	// membersPrefix is the stroe prefix for memberships. Memberships are
	// stored with key in the following format "memberships/{teamID}/{userID}".
	// The value is the user ID (string, stored as bytes).
	membersPrefix = "members/"
)

var (
	// reservedNamespaces cannot be used by teams
	reservedNamespaces = []string{"default", "go.micro", "runtime"}
)

// Read looks up a team using ID or namespace
func (t *Teams) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// lookup the team
	var err error
	if len(req.Id) > 0 {
		rsp.Team, err = t.findTeamByID(req.Id)
	}
	if len(req.Namespace) > 0 && rsp.Team == nil {
		rsp.Team, err = t.findTeamByNamespace(req.Namespace)
	}
	if err != nil {
		return err
	}

	// lookup the team members
	recs, err := t.store.Read(membersPrefix+rsp.Team.Id+"/", store.ReadPrefix())
	if err != nil {
		return nil
	}
	rsp.Team.Members = make([]*pb.Member, 0, len(recs))
	for _, r := range recs {
		rsp.Team.Members = append(rsp.Team.Members, &pb.Member{
			Id: string(r.Value),
		})
	}

	return nil
}

// Create a new team
func (t *Teams) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if req.Team == nil {
		return errors.BadRequest(t.name, "Missing team")
	}
	if len(req.Team.Name) == 0 {
		return errors.BadRequest(t.name, "Missing team name")
	}
	if err := t.validateNamespace(req.Team.Namespace); err != nil {
		return err
	}

	// add the default fields
	req.Team.Id = uuid.New().String()

	// write to the store
	if err := t.writeTeamToStore(req.Team); err != nil {
		return err
	}

	// return the team in the response
	rsp.Team = req.Team
	return nil
}

// Update a team
func (t *Teams) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// validate the request
	if req.Team == nil {
		return errors.BadRequest(t.name, "Missing team")
	}
	if len(req.Team.Id) == 0 {
		return errors.BadRequest(t.name, "Missing team id")
	}
	if len(req.Team.Name) == 0 {
		return errors.BadRequest(t.name, "Missing team name")
	}

	// lookup the team
	team, err := t.findTeamByID(req.Team.Id)
	if err != nil {
		return errors.BadRequest(t.name, "Error finding team: %v", err)
	}

	// assign the update params
	team.Name = req.Team.Name
	team.WebDomain = req.Team.WebDomain
	team.ApiDomain = req.Team.ApiDomain

	// write to the store
	if err := t.writeTeamToStore(req.Team); err != nil {
		return errors.InternalServerError(t.name, "Error writing team: %v", err)
	}

	return nil
}

// List all the teams (does not return membership)
func (t *Teams) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	// get the records with the team prefix
	recs, err := t.store.Read(teamsPrefix, store.ReadPrefix())
	if err != nil {
		return err
	}

	// unmarshal and return in the response
	rsp.Teams = make([]*pb.Team, len(recs))
	for i, r := range recs {
		if err := json.Unmarshal(r.Value, &rsp.Teams[i]); err != nil {
			return errors.InternalServerError(t.name, "Error unmarsaling json: %v", err)
		}
	}

	return nil
}

// AddMember to a team
func (t *Teams) AddMember(ctx context.Context, req *pb.AddMemberRequest, rsp *pb.AddMemberResponse) error {
	// validate the request
	if _, err := t.findTeamByID(req.TeamId); err != nil {
		return err
	}
	if len(req.MemberId) == 0 {
		return errors.BadRequest(t.name, "Missing member id")
	}

	// write the membership to the store
	return t.store.Write(&store.Record{
		Key:   membersPrefix + req.TeamId + "/" + req.MemberId,
		Value: []byte(req.MemberId),
	})
}

// RemoveMember from a team
func (t *Teams) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest, rsp *pb.RemoveMemberResponse) error {
	return t.store.Delete(membersPrefix + req.TeamId + "/" + req.MemberId)
}

// ListMemberships returns all the teams a member belongs to
func (t *Teams) ListMemberships(ctx context.Context, req *pb.ListMembershipsRequest, rsp *pb.ListMembershipsResponse) error {
	// member id is the last component of the key, so list all
	// the keys in the store which relate to memberships
	keys, err := t.store.List(store.ListPrefix(membersPrefix))
	if err != nil {
		return err
	}

	// filter to get the team ids which the member belongs to
	var teamIDs []string
	for _, k := range keys {
		if strings.HasSuffix(k, "/"+req.MemberId) {
			teamIDs = append(teamIDs, strings.Split(k, "/")[1])
		}
	}

	// get each of the teams
	rsp.Teams = make([]*pb.Team, 0, len(teamIDs))
	for _, id := range teamIDs {
		team, err := t.findTeamByID(id)
		if err != nil {
			return err
		}
		rsp.Teams = append(rsp.Teams, team)
	}

	return nil
}

func (t *Teams) findTeamByID(id string) (*pb.Team, error) {
	// ID is stored as the last component of the key, so list
	// all the keys in the store which relate to teams
	keys, err := t.store.List(store.ListPrefix(teamsPrefix))
	if err != nil {
		return nil, err
	}

	// Check each key to see if it ends in the ID. If the key
	// is not found, return an error.
	var teamKey string
	for _, k := range keys {
		if strings.HasSuffix(k, "/"+id) {
			teamKey = k
			break
		}
	}
	if len(teamKey) == 0 {
		return nil, store.ErrNotFound
	}

	// Lookup the record and then decode the value
	recs, err := t.store.Read(teamKey)
	if err != nil {
		return nil, err
	}
	var team *pb.Team
	err = json.Unmarshal(recs[0].Value, &team)
	return team, err
}

// writeTeamToStore marshals a team and writes it to the store under
// the corresponding key (prefix + namespace + id)
func (t *Teams) writeTeamToStore(team *pb.Team) error {
	bytes, err := json.Marshal(team)
	if err != nil {
		return errors.InternalServerError(t.name, "Error marsaling json: %v", err)
	}

	key := teamsPrefix + team.Namespace + "/" + team.Id
	if err := t.store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.InternalServerError(t.name, "Error writing to the store: %v", err)
	}

	return nil
}

func (t *Teams) findTeamByNamespace(ns string) (*pb.Team, error) {
	// Namespace is the first component of the key, so lookup records which
	// have this as a prefix. Read does't return an error when using the
	// ReadPrefix option, so we also need to check for an empty slice.
	recs, err := t.store.Read(teamsPrefix+ns+"/", store.ReadPrefix())
	if err != nil {
		return nil, err
	} else if len(recs) != 1 {
		return nil, store.ErrNotFound
	}

	// Unmarshal and return the result
	var team *pb.Team
	err = json.Unmarshal(recs[0].Value, &team)
	return team, err
}

// validateNamespace returns an error if the namespace provided is invalid.
func (t *Teams) validateNamespace(ns string) error {
	// compare namespaces in lowercase
	ns = strings.ToLower(ns)

	// validate the length of the namespace
	if len(ns) < 3 {
		return errors.BadRequest(t.name, "Namespaces must be at least 3 characters long")
	}
	if len(ns) > 20 {
		return errors.BadRequest(t.name, "Namespaces must be at no more than 20 characters long")
	}

	// check against reserved namespaces
	for _, v := range reservedNamespaces {
		if v == ns {
			return errors.BadRequest(t.name, "%v is a reserved namespace", ns)
		}
	}

	// check against existing namespaces. The namespace is used
	// as the key in the store
	recs, err := t.store.Read(teamsPrefix+ns+"/", store.ReadPrefix())
	if err != nil {
		return err
	} else if len(recs) > 0 {
		return errors.BadRequest(t.name, "The namespace %v has already been taken", ns)
	}

	return nil
}
