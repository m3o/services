package handler

import (
	"context"
	"fmt"
	"math/rand"
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

func TestMain(m *testing.M) {
	mstore.DefaultStore = memory.NewStore()
	m.Run()
}

func TestDuplicateInvites(t *testing.T) {
	g := NewWithT(t)
	inviteSvc := mockInvite()
	userCtx := contextWithAccount("foo", testEmail())
	emails := inviteSvc.emailSvc.(*memail.FakeEmailsService)
	err := inviteSvc.User(userCtx, &pb.CreateRequest{
		Email:     "foo@bar.com",
		Namespace: "foo",
		Resend:    false,
	}, &pb.CreateResponse{})
	g.Expect(err).To(BeNil())
	g.Expect(emails.SendCallCount()).To(Equal(1))

	err = inviteSvc.User(userCtx, &pb.CreateRequest{
		Email:     "foo@bar.com",
		Namespace: "foo",
		Resend:    false,
	}, &pb.CreateResponse{})
	g.Expect(err).To(BeNil())
	g.Expect(emails.SendCallCount()).To(Equal(1))

	err = inviteSvc.User(userCtx, &pb.CreateRequest{
		Email:     "foo@bar.com",
		Namespace: "foo",
		Resend:    true,
	}, &pb.CreateResponse{})
	g.Expect(err).To(BeNil())
	g.Expect(emails.SendCallCount()).To(Equal(2))

}

func TestEmailValidation(t *testing.T) {
	g := NewWithT(t)
	inviteSvc := mockInvite()
	userCtx := contextWithAccount("foo", testEmail())
	err := inviteSvc.User(userCtx, &pb.CreateRequest{
		Email:     "notanemail.com",
		Namespace: "foo",
		Resend:    false,
	}, &pb.CreateResponse{})
	g.Expect(err).To(HaveOccurred())

}

func TestUserInviteLimit(t *testing.T) {
	g := NewWithT(t)
	inviteSvc := mockInvite()
	userCtx := contextWithAccount("foo", testEmail())

	for i := 0; i < 5; i++ {
		err := inviteSvc.User(userCtx, &pb.CreateRequest{
			Email:     testEmail(),
			Namespace: "foo",
			Resend:    false,
		}, &pb.CreateResponse{})
		g.Expect(err).To(BeNil())
	}
	err := inviteSvc.User(userCtx, &pb.CreateRequest{
		Email:     testEmail(),
		Namespace: "foo",
		Resend:    false,
	}, &pb.CreateResponse{})
	g.Expect(err).To(HaveOccurred())

}

func TestUserInviteToNotOwnedNamespace(t *testing.T) {
	g := NewWithT(t)
	inviteSvc := mockInvite()
	userCtx := contextWithAccount("foo", testEmail())
	err := inviteSvc.User(userCtx, &pb.CreateRequest{
		Email:     testEmail(),
		Namespace: "baz",
		Resend:    false,
	}, &pb.CreateResponse{})
	g.Expect(err).To(HaveOccurred())
}

func testEmail() string {
	return fmt.Sprintf("foo%d@bar.com", rand.Int())
}
