package auth

import (
	"net"
	"sync"

	"github.com/gliderlabs/ssh"
	"github.com/rnzor/poor_man_exe/internal/db"
	gossh "golang.org/x/crypto/ssh"
)

type Authenticator struct {
	DB          *db.Database
	Limiter     *RateLimiter
	Sessions    map[string]int // fingerprint -> count
	SessionMu   sync.Mutex
	MaxSessions int
}

func NewAuthenticator(d *db.Database) *Authenticator {
	return &Authenticator{
		DB:          d,
		Limiter:     NewRateLimiter(0.1, 5.0), // 1 conn every 10s, burst of 5
		Sessions:    make(map[string]int),
		MaxSessions: 10,
	}
}

func (a *Authenticator) PublicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	fingerprint := gossh.FingerprintSHA256(key)

	// Check session cap
	a.SessionMu.Lock()
	count := a.Sessions[fingerprint]
	if count >= a.MaxSessions {
		a.SessionMu.Unlock()
		return false
	}
	a.Sessions[fingerprint] = count + 1
	a.SessionMu.Unlock()

	var userID int
	err := a.DB.Conn.QueryRow("SELECT user_id FROM public_keys WHERE fingerprint = ?", fingerprint).Scan(&userID)
	if err != nil {
		return false
	}

	// Store data in context
	ctx.SetValue("user_id", userID)
	ctx.SetValue("fingerprint", fingerprint)
	ctx.SetValue("remote_ip", ctx.RemoteAddr().String())

	return true
}

func (a *Authenticator) ConnCallback(ctx ssh.Context, conn net.Conn) net.Conn {
	ip := conn.RemoteAddr().String()
	if !a.Limiter.Allow(ip) {
		conn.Close()
		return nil
	}
	return conn
}

func (a *Authenticator) OnSessionClose(sess ssh.Session) {
	fingerprint, ok := sess.Context().Value("fingerprint").(string)
	if ok {
		a.SessionMu.Lock()
		a.Sessions[fingerprint]--
		if a.Sessions[fingerprint] < 0 {
			a.Sessions[fingerprint] = 0
		}
		a.SessionMu.Unlock()
	}
}

func (a *Authenticator) BannerHandler(ctx ssh.Context) string {
	return "\n╔════════════════════════════════════════════════════════════╗\n║  Poor Man's exe.dev SSH Gateway                           ║\n║  Welcome to the machine.                                   ║\n╚════════════════════════════════════════════════════════════╝\n\n"
}
