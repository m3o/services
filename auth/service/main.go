package main

import (
	"context"
	"io/ioutil"

	"github.com/micro/go-micro/v2"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	log "github.com/micro/go-micro/v2/logger"
	"gopkg.in/yaml.v2"
)

func main() {
	srv := micro.NewService(
		micro.Name("go.micro.service.auth"),
		micro.Version("latest"),
	)
	srv.Init()

	// Setup the auth service
	auth := pb.NewRulesService("go.micro.auth", srv.Client())

	// Create the new rules
	newRules := loadRulesFromFile()
	for _, r := range newRules {
		_, err := auth.Create(context.TODO(), &pb.CreateRequest{
			Role:     r.Role,
			Access:   r.Access,
			Priority: r.Priority,
			Resource: r.Resource,
		})

		if err != nil {
			log.Fatalf("Error creating rule: %v", err)
		}
	}

	// Get all the rules
	aRsp, err := auth.List(context.TODO(), &pb.ListRequest{})
	if err != nil {
		log.Fatalf("Error retrieving existing rules: %v", err)
	}

	// Compare the rules and delete the ones which no longer exist
loop:
	for _, rule := range aRsp.GetRules() {
		for _, r := range newRules {
			if rulesMatch(rule, r) {
				continue loop
			}
		}

		_, err := auth.Delete(context.TODO(), &pb.DeleteRequest{
			Role:     rule.Role,
			Access:   rule.Access,
			Priority: rule.Priority,
			Resource: rule.Resource,
		})

		if err != nil {
			log.Fatalf("Error deleting rule: %v", err)
		}
	}
}

func rulesMatch(a, b *pb.Rule) bool {
	if a.Role != b.Role {
		return false
	}
	if a.Access != b.Access {
		return false
	}
	if a.Priority != b.Priority {
		return false
	}
	if a.Resource.Namespace != b.Resource.Namespace {
		return false
	}
	if a.Resource.Type != b.Resource.Type {
		return false
	}
	if a.Resource.Name != b.Resource.Name {
		return false
	}
	if a.Resource.Endpoint != b.Resource.Endpoint {
		return false
	}
	return true
}

func loadRulesFromFile() []*pb.Rule {
	// Parse the rules yaml file
	bytes, err := ioutil.ReadFile("rules.yaml")
	if err != nil {
		log.Fatalf("Error reading rules file: %v", err)
	}

	// Unmarshal into custom object because of the enum
	// used for access
	var rules []struct {
		Role     string
		Priority int32
		Access   string
		Resource *pb.Resource
	}

	if err := yaml.Unmarshal(bytes, &rules); err != nil {
		log.Fatalf("Error parsing rules file: %v", err)
	}

	pbRules := make([]*pb.Rule, 0, len(rules))
	for _, r := range rules {
		var access pb.Access
		switch r.Access {
		case "granted":
			access = pb.Access_GRANTED
		case "denied":
			access = pb.Access_DENIED
		}

		pbRules = append(pbRules, &pb.Rule{Role: r.Role, Priority: r.Priority, Access: access, Resource: r.Resource})
	}

	return pbRules
}
