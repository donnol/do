package do_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/donnol/do"
)

var (
	_ = func() struct{} {
		rand.Seed(time.Now().Unix())
		return struct{}{}
	}
)

func TestToken(t *testing.T) {
	var (
		issuer = "jd"
	)

	secret := []byte("Xadfdfoere2324212afasf34wraf090uadfafdIEJF039038")
	token := do.New[uint64](secret, time.Second*10, issuer)
	id := rand.Uint64()
	r, err := token.Sign(id)
	if err != nil {
		t.Fatal(err)
	}

	userID, err := token.Verify(r)
	if err != nil {
		t.Fatal(err)
	}

	do.Assert(t, userID, id, "Bad userID, got %v\n", userID)

	t.Run("VerifyEmpty", func(t *testing.T) {
		r, err := token.Verify("")
		if err != nil {
			t.Fatal(err)
		}
		if r != 0 {
			do.Assert(t, r, 0)
		}
	})
}

func TestTokenObject(t *testing.T) {
	type User struct {
		Id       uint64
		TenantId uint64
	}

	var (
		issuer = "jd"
	)

	secret := []byte("Xadfdfoere2324212afasf34wraf090uadfafdIEJF039038")
	token := do.New[User](secret, time.Second*10, issuer)
	id := rand.Uint64()
	tid := rand.Uint64()
	u := User{
		Id:       id,
		TenantId: tid,
	}
	r, err := token.Sign(u)
	if err != nil {
		t.Fatal(err)
	}

	wu, err := token.Verify(r)
	if err != nil {
		t.Fatal(err)
	}

	do.Assert(t, wu.Id, u.Id, "Bad user id, got %v\n", wu)
	do.Assert(t, wu.TenantId, u.TenantId, "Bad user tid, got %v\n", wu)
	do.Assert(t, wu, u, "Bad user, got %v\n", wu)

	t.Run("VerifyEmpty", func(t *testing.T) {
		r, err := token.Verify("")
		if err != nil {
			t.Fatal(err)
		}
		if r.Id != 0 {
			do.Assert(t, r.Id, 0)
		}
	})
}
