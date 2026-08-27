package contract4chessraiders_test

import (
	"context"

	"github.com/sneat-co/sneat-ext-contracts/chessraiders/contract4chessraiders"
)

// The stubs below exist to prove, at compile time, that every interface in
// this package is actually implementable with the exact method set declared
// — the same guarantee contract4competios's application_contract_test.go
// gives its own interfaces. A contract whose interface cannot be satisfied
// by a plain struct is not a usable contract.

type matchLobbyApplicationStub struct{}

func (matchLobbyApplicationStub) CreateLobby(context.Context, contract4chessraiders.Player, contract4chessraiders.CaptureMode) (contract4chessraiders.LobbyView, error) {
	return contract4chessraiders.LobbyView{}, nil
}

func (matchLobbyApplicationStub) Join(context.Context, string, contract4chessraiders.Player) (contract4chessraiders.LobbyView, error) {
	return contract4chessraiders.LobbyView{}, nil
}

func (matchLobbyApplicationStub) Leave(context.Context, string, string) (contract4chessraiders.LobbyView, error) {
	return contract4chessraiders.LobbyView{}, nil
}

func (matchLobbyApplicationStub) SetCaptureMode(context.Context, string, string, contract4chessraiders.CaptureMode) (contract4chessraiders.LobbyView, error) {
	return contract4chessraiders.LobbyView{}, nil
}

func (matchLobbyApplicationStub) VoteReady(context.Context, string, string, bool) (contract4chessraiders.LobbyView, error) {
	return contract4chessraiders.LobbyView{}, nil
}

func (matchLobbyApplicationStub) Start(context.Context, string, string) (contract4chessraiders.LobbyView, error) {
	return contract4chessraiders.LobbyView{}, nil
}

func (matchLobbyApplicationStub) Get(context.Context, string, string) (contract4chessraiders.LobbyView, error) {
	return contract4chessraiders.LobbyView{}, nil
}

var _ contract4chessraiders.MatchLobbyApplication = matchLobbyApplicationStub{}

type externalIdentityVerifierStub struct {
	provider contract4chessraiders.IdentityProvider
}

func (s externalIdentityVerifierStub) Provider() contract4chessraiders.IdentityProvider {
	return s.provider
}

func (externalIdentityVerifierStub) Verify(context.Context, string) (contract4chessraiders.ExternalIdentity, error) {
	return contract4chessraiders.ExternalIdentity{}, nil
}

var _ contract4chessraiders.ExternalIdentityVerifier = externalIdentityVerifierStub{}

type identityLinkApplicationStub struct{}

func (identityLinkApplicationStub) LinkExternalIdentity(context.Context, string, contract4chessraiders.ExternalIdentity) error {
	return nil
}

func (identityLinkApplicationStub) ResolveExternalIdentity(context.Context, contract4chessraiders.IdentityProvider, string) (string, bool, error) {
	return "", false, nil
}

var _ contract4chessraiders.IdentityLinkApplication = identityLinkApplicationStub{}

type portalSessionApplicationStub struct{}

func (portalSessionApplicationStub) OpenPortalSession(context.Context, contract4chessraiders.PortalSessionRequest) (contract4chessraiders.PortalSession, error) {
	return contract4chessraiders.PortalSession{}, nil
}

func (portalSessionApplicationStub) GetPortalSession(context.Context, contract4chessraiders.IdentityProvider, string) (contract4chessraiders.PortalSession, error) {
	return contract4chessraiders.PortalSession{}, nil
}

var _ contract4chessraiders.PortalSessionApplication = portalSessionApplicationStub{}
