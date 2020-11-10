package handler

import (
	"testing"

	mstore "github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/memory"

	memail "github.com/m3o/services/emails/proto/fakes"
	mt "github.com/m3o/services/internal/test"
	pb "github.com/m3o/services/invite/proto"

	. "github.com/onsi/gomega"
)

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
	userCtx := mt.ContextWithAccount("foo", mt.TestEmail())
	emails := inviteSvc.emailSvc.(*memail.FakeEmailsService)
	err := inviteSvc.User(userCtx, &pb.CreateRequest{
		Email:  "foo@bar.com",
		Resend: false,
	}, &pb.CreateResponse{})
	g.Expect(err).To(BeNil())
	g.Expect(emails.SendCallCount()).To(Equal(1))

	err = inviteSvc.User(userCtx, &pb.CreateRequest{
		Email:  "foo@bar.com",
		Resend: false,
	}, &pb.CreateResponse{})
	g.Expect(err).To(BeNil())
	g.Expect(emails.SendCallCount()).To(Equal(1))

	err = inviteSvc.User(userCtx, &pb.CreateRequest{
		Email:  "foo@bar.com",
		Resend: true,
	}, &pb.CreateResponse{})
	g.Expect(err).To(BeNil())
	g.Expect(emails.SendCallCount()).To(Equal(2))

}

func TestEmailValidation(t *testing.T) {
	g := NewWithT(t)
	inviteSvc := mockInvite()
	userCtx := mt.ContextWithAccount("foo", mt.TestEmail())
	err := inviteSvc.User(userCtx, &pb.CreateRequest{
		Email:  "notanemail.com",
		Resend: false,
	}, &pb.CreateResponse{})
	g.Expect(err).To(HaveOccurred())

}

func TestUserInviteLimit(t *testing.T) {
	g := NewWithT(t)
	inviteSvc := mockInvite()
	userCtx := mt.ContextWithAccount("foo", mt.TestEmail())

	for i := 0; i < 5; i++ {
		err := inviteSvc.User(userCtx, &pb.CreateRequest{
			Email:  mt.TestEmail(),
			Resend: false,
		}, &pb.CreateResponse{})
		g.Expect(err).To(BeNil())
	}
	err := inviteSvc.User(userCtx, &pb.CreateRequest{
		Email:  mt.TestEmail(),
		Resend: false,
	}, &pb.CreateResponse{})
	g.Expect(err).To(HaveOccurred())

}

func TestUserInviteToNotOwnedNamespace(t *testing.T) {
	g := NewWithT(t)
	inviteSvc := mockInvite()
	userCtx := mt.ContextWithAccount("foo", mt.TestEmail())
	err := inviteSvc.User(userCtx, &pb.CreateRequest{
		Email:     mt.TestEmail(),
		Namespace: "baz",
		Resend:    false,
	}, &pb.CreateResponse{})
	g.Expect(err).To(HaveOccurred())

}

func TestInviteUserToNamespace(t *testing.T) {
	g := NewWithT(t)
	inviteSvc := mockInvite()
	userCtx := mt.ContextWithAccount("foo", mt.TestEmail())
	invitee := mt.TestEmail()
	err := inviteSvc.User(userCtx, &pb.CreateRequest{
		Email:     invitee,
		Namespace: "foo",
		Resend:    false,
	}, &pb.CreateResponse{})
	g.Expect(err).To(BeNil())
	rsp := &pb.ValidateResponse{}
	err = inviteSvc.Validate(userCtx, &pb.ValidateRequest{
		Email: invitee,
	}, rsp)
	g.Expect(err).To(BeNil())
	g.Expect(rsp.Namespaces).To(HaveLen(1))
	g.Expect(rsp.Namespaces[0]).To(Equal("foo"))
}

func TestValidate(t *testing.T) {
	g := NewWithT(t)
	inviteSvc := mockInvite()
	userCtx := mt.ContextWithAccount("foo", mt.TestEmail())
	invitee := mt.TestEmail()
	err := inviteSvc.User(userCtx, &pb.CreateRequest{
		Email:  invitee,
		Resend: false,
	}, &pb.CreateResponse{})
	g.Expect(err).To(BeNil())
	rsp := &pb.ValidateResponse{}
	err = inviteSvc.Validate(userCtx, &pb.ValidateRequest{
		Email: invitee,
	}, rsp)
	g.Expect(err).To(BeNil())
	g.Expect(rsp.Namespaces).To(BeEmpty())
}
