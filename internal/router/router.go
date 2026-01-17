package router

import (
	"fmt"

	"github.com/gliderlabs/ssh"
	"github.com/rnzor/poor_man_exe/internal/caddy"
	"github.com/rnzor/poor_man_exe/internal/cli"
	"github.com/rnzor/poor_man_exe/internal/config"
	"github.com/rnzor/poor_man_exe/internal/db"
	"github.com/rnzor/poor_man_exe/internal/runner"
)

type Router struct {
	DB     *db.Database
	Runner *runner.DockerRunner
	Cfg    *config.Config
	Caddy  *caddy.Client
}

func NewRouter(d *db.Database, r *runner.DockerRunner, cfg *config.Config, c *caddy.Client) *Router {
	return &Router{DB: d, Runner: r, Cfg: cfg, Caddy: c}
}

func (r *Router) HandleSession(sess ssh.Session) {
	username := sess.User()
	command := sess.Command()

	// If username is one of these, it's management mode
	if username == "root" || username == "exedev" || username == "" || username == "poor-exe" {
		if len(command) > 0 {
			cli.ExecuteCommand(sess, command, r.DB, r.Runner, r.Caddy, r.Cfg)
		} else {
			cli.StartInteractiveCLI(sess, r.DB, r.Runner, r.Caddy, r.Cfg)
		}
		return
	}

	// Otherwise, username is the app name
	r.AttachToApp(sess, username)
}

func (r *Router) AttachToApp(sess ssh.Session, appName string) {
	userID := sess.Context().Value("user_id").(int)

	// Verify user has access to this app
	var exists bool
	err := r.DB.Conn.QueryRow("SELECT EXISTS(SELECT 1 FROM apps WHERE name = ? AND user_id = ?)", appName, userID).Scan(&exists)
	if err != nil || !exists {
		fmt.Fprintf(sess, "Error: App '%s' not found or access denied.\n", appName)
		sess.Exit(1)
		return
	}

	// Attach to the container
	err = r.Runner.Attach(sess.Context(), appName, sess, sess, sess.Stderr(), sess)
	if err != nil {
		fmt.Fprintf(sess, "Error attaching to app: %v\n", err)
		sess.Exit(1)
		return
	}
}
