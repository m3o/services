package handler

import (
	"context"
	"testing"

	"github.com/micro/micro/v3/service/store/memory"

	mstore "github.com/micro/micro/v3/service/store"

	"github.com/micro/micro/v3/service/auth"

	pb "github.com/m3o/services/invite/proto"
	. "github.com/onsi/gomega"

	memail "github.com/m3o/services/emails/proto/protofakes"
)

func contextWithAccount(issuer, id string) context.Context {
	return auth.ContextWithAccount(context.TODO(), &auth.Account{
		ID:       id,
		Type:     "",
		Issuer:   issuer,
		Metadata: nil,
		Scopes:   nil,
		Secret:   "",
		Name:     "",
	})
}

func mockInvite() *Invite {
	return &Invite{
		config:   inviteConfig{},
		name:     "",
		emailSvc: &memail.FakeEmailsService{},
	}

}

func TestDuplicateInvites(t *testing.T) {
	g := NewWithT(t)
	mstore.DefaultStore = memory.NewStore()
	inviteSvc := mockInvite()
	emails := inviteSvc.emailSvc.(*memail.FakeEmailsService)
	err := inviteSvc.User(contextWithAccount("foo", "foo@bar.com"), &pb.CreateRequest{
		Email:     "foo@bar.com",
		Namespace: "foo",
		Resend:    false,
	}, &pb.CreateResponse{})
	g.Expect(err).To(BeNil())
	g.Expect(emails.SendCallCount()).To(Equal(1))

	err = inviteSvc.User(contextWithAccount("foo", "foo@bar.com"), &pb.CreateRequest{
		Email:     "foo@bar.com",
		Namespace: "foo",
		Resend:    false,
	}, &pb.CreateResponse{})
	g.Expect(err).To(BeNil())
	g.Expect(emails.SendCallCount()).To(Equal(1))

	err = inviteSvc.User(contextWithAccount("foo", "foo@bar.com"), &pb.CreateRequest{
		Email:     "foo@bar.com",
		Namespace: "foo",
		Resend:    true,
	}, &pb.CreateResponse{})
	g.Expect(err).To(BeNil())
	g.Expect(emails.SendCallCount()).To(Equal(2))

}
