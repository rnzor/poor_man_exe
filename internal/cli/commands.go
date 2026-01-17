package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/moby/moby/client"
	"github.com/rnzor/poor_man_exe/internal/caddy"
	"github.com/rnzor/poor_man_exe/internal/config"
	"github.com/rnzor/poor_man_exe/internal/db"
	"github.com/rnzor/poor_man_exe/internal/runner"
	gossh "golang.org/x/crypto/ssh"
)

func ExecuteCommand(sess ssh.Session, args []string, d *db.Database, r *runner.DockerRunner, c *caddy.Client, cfg *config.Config) {
	if len(args) == 0 {
		return
	}

	cmd := args[0]
	ctx := sess.Context()
	userID := ctx.Value("user_id").(int)

	isJSON := HasFlag(args, "--json")

	switch cmd {
	case "ls":
		handleLs(sess, d, r, userID, isJSON)
	case "new":
		handleNew(sess, args[1:], d, r, c, cfg, userID, isJSON)
	case "rm":
		handleRm(sess, args[1:], d, r, c, userID, isJSON)
	case "share":
		handleShare(sess, args[1:], d, c, cfg, userID, isJSON)
	case "keys":
		handleKeys(sess, args[1:], d, userID, isJSON)
	case "whoami":
		handleWhoami(sess, d, userID, isJSON)
	case "help":
		handleHelp(sess)
	default:
		if isJSON {
			WriteJSON(sess, false, "", nil, fmt.Errorf("unknown command: %s", cmd))
		} else {
			fmt.Fprintf(sess, "Unknown command: %s. Type 'help' for available commands.\n", cmd)
		}
	}
}

func StartInteractiveCLI(sess ssh.Session, d *db.Database, r *runner.DockerRunner, c *caddy.Client, cfg *config.Config) {
	fmt.Fprintf(sess, "Poor Man's exe.dev CLI\nType 'help' for commands.\n\n")

	for {
		fmt.Fprintf(sess, "poor-exe> ")
		line := make([]byte, 1024)
		n, err := sess.Read(line)
		if err != nil {
			break
		}

		input := strings.TrimSpace(string(line[:n]))
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			break
		}

		args := strings.Fields(input)
		ExecuteCommand(sess, args, d, r, c, cfg)
	}
}

func handleLs(sess ssh.Session, d *db.Database, r *runner.DockerRunner, userID int, isJSON bool) {
	rows, err := d.Conn.Query("SELECT name, image, status, created_at FROM apps WHERE user_id = ?", userID)
	if err != nil {
		if isJSON {
			WriteJSON(sess, false, "", nil, err)
		} else {
			fmt.Fprintf(sess, "Error listing apps: %v\n", err)
		}
		return
	}
	defer rows.Close()

	type appInfo struct {
		Name    string `json:"vm_name"`
		Image   string `json:"image"`
		Status  string `json:"status"`
		Created string `json:"created_at"`
	}
	var apps []appInfo

	if !isJSON {
		fmt.Fprintf(sess, "%-25s %-25s %-12s %-20s\n", "NAME", "IMAGE", "STATUS", "CREATED")
		fmt.Fprintf(sess, "%-25s %-25s %-12s %-20s\n", "----", "-----", "------", "-------")
	}

	for rows.Next() {
		var name, image, dbStatus, created string
		rows.Scan(&name, &image, &dbStatus, &created)

		// Sync with Docker
		status, err := r.GetAppStatus(sess.Context(), name)
		if err != nil {
			status = dbStatus // Fallback to DB
		}

		if isJSON {
			apps = append(apps, appInfo{name, image, status, created})
		} else {
			fmt.Fprintf(sess, "%-25s %-25s %-12s %-20s\n", name, image, status, created)
		}
	}

	if isJSON {
		WriteJSON(sess, true, "", map[string]interface{}{"vms": apps}, nil)
	}
}

func handleNew(sess ssh.Session, args []string, d *db.Database, r *runner.DockerRunner, c *caddy.Client, cfg *config.Config, userID int, isJSON bool) {
	name := ""
	image := "alpine:latest"

	for _, arg := range args {
		if strings.HasPrefix(arg, "--name=") {
			name = strings.TrimPrefix(arg, "--name=")
		} else if strings.HasPrefix(arg, "--image=") {
			image = strings.TrimPrefix(arg, "--image=")
		}
	}

	if name == "" {
		if isJSON {
			WriteJSON(sess, false, "", nil, fmt.Errorf("usage: new --name=<name> [--image=<image>]"))
		} else {
			fmt.Fprintf(sess, "Usage: new --name=<name> [--image=<image>]\n")
		}
		return
	}

	err := r.CreateApp(sess.Context(), name, image, userID)
	if err != nil {
		if isJSON {
			WriteJSON(sess, false, "", nil, err)
		} else {
			fmt.Fprintf(sess, "Error creating app: %v\n", err)
		}
		return
	}

	_, err = d.Conn.Exec("INSERT INTO apps (name, image, user_id, status) VALUES (?, ?, ?, 'running')", name, image, userID)
	if err != nil {
		if isJSON {
			WriteJSON(sess, false, "App created in Docker but failed to update registry", nil, err)
		} else {
			fmt.Fprintf(sess, "App created in Docker but failed to update registry: %v\n", err)
		}
		return
	}

	// Add to Caddy
	if err := c.UpsertRoute(name, cfg.Domain, 80); err != nil {
		if !isJSON {
			fmt.Fprintf(sess, "Warning: Failed to configure HTTP proxy: %v\n", err)
		}
	}

	// Audit log
	remoteIP, _ := sess.Context().Value("remote_ip").(string)
	d.LogAudit("app_create", userID, name, remoteIP, "image="+image)

	if isJSON {
		WriteJSON(sess, true, fmt.Sprintf("Successfully created app '%s'", name), map[string]string{
			"vm_name":  name,
			"image":    image,
			"status":   "running",
			"endpoint": fmt.Sprintf("https://%s.%s", name, cfg.Domain),
		}, nil)
	} else {
		fmt.Fprintf(sess, "Successfully created app '%s' using image '%s'\n", name, image)
		fmt.Fprintf(sess, "Endpoint: https://%s.%s\n", name, cfg.Domain)
	}
}

func handleRm(sess ssh.Session, args []string, d *db.Database, r *runner.DockerRunner, c *caddy.Client, userID int, isJSON bool) {
	if len(args) == 0 {
		if isJSON {
			WriteJSON(sess, false, "", nil, fmt.Errorf("usage: rm <app_name>"))
		} else {
			fmt.Fprintf(sess, "Usage: rm <app_name>\n")
		}
		return
	}

	name := args[0]
	// Verify ownership
	var exists bool
	d.Conn.QueryRow("SELECT EXISTS(SELECT 1 FROM apps WHERE name = ? AND user_id = ?)", name, userID).Scan(&exists)
	if !exists {
		if isJSON {
			WriteJSON(sess, false, "", nil, fmt.Errorf("app '%s' not found or access denied", name))
		} else {
			fmt.Fprintf(sess, "Error: App '%s' not found or access denied.\n", name)
		}
		return
	}

	// Remove from Caddy
	if err := c.DeleteRoute(name); err != nil {
		if !isJSON {
			fmt.Fprintf(sess, "Warning: Failed to remove HTTP proxy: %v\n", err)
		}
	}

	// Remove from Docker
	containerName := fmt.Sprintf("poor-exe-%s", name)
	_, err := r.Cli.ContainerRemove(context.Background(), containerName, client.ContainerRemoveOptions{Force: true})
	if err != nil {
		if !isJSON {
			fmt.Fprintf(sess, "Warning: Failed to remove container from Docker: %v\n", err)
		}
	}

	// Remove from DB
	_, err = d.Conn.Exec("DELETE FROM apps WHERE name = ? AND user_id = ?", name, userID)
	if err != nil {
		if isJSON {
			WriteJSON(sess, false, "", nil, err)
		} else {
			fmt.Fprintf(sess, "Error removing app from registry: %v\n", err)
		}
		return
	}

	// Audit log
	remoteIP, _ := sess.Context().Value("remote_ip").(string)
	d.LogAudit("app_delete", userID, name, remoteIP, "")

	if isJSON {
		WriteJSON(sess, true, fmt.Sprintf("Successfully removed app '%s'", name), nil, nil)
	} else {
		fmt.Fprintf(sess, "Successfully removed app '%s'\n", name)
	}
}

func handleWhoami(sess ssh.Session, d *db.Database, userID int, isJSON bool) {
	var email string
	d.Conn.QueryRow("SELECT email FROM users WHERE id = ?", userID).Scan(&email)
	if isJSON {
		WriteJSON(sess, true, "", map[string]interface{}{
			"user_id": userID,
			"email":   email,
		}, nil)
	} else {
		fmt.Fprintf(sess, "User ID: %d\nEmail: %s\n", userID, email)
	}
}

func handleShare(sess ssh.Session, args []string, d *db.Database, c *caddy.Client, cfg *config.Config, userID int, isJSON bool) {
	if len(args) < 2 {
		usage := "Usage: share <cmd> <vm> [args]\nCmds: set-public, set-private, port, add, remove"
		if isJSON {
			WriteJSON(sess, false, "", nil, errors.New(usage))
		} else {
			fmt.Fprintln(sess, usage)
		}
		return
	}

	cmd := args[0]
	vmName := args[1]

	// Verify ownership
	var appID int
	var currentPort int
	err := d.Conn.QueryRow("SELECT id, http_port FROM apps WHERE name = ? AND user_id = ?", vmName, userID).Scan(&appID, &currentPort)
	if err != nil {
		if isJSON {
			WriteJSON(sess, false, "", nil, fmt.Errorf("app '%s' not found or access denied", vmName))
		} else {
			fmt.Fprintf(sess, "Error: App '%s' not found.\n", vmName)
		}
		return
	}

	switch cmd {
	case "set-public":
		_, err = d.Conn.Exec("UPDATE apps SET is_public = 1 WHERE id = ?", appID)
	case "set-private":
		_, err = d.Conn.Exec("UPDATE apps SET is_public = 0 WHERE id = ?", appID)
	case "port":
		if len(args) < 3 {
			err = fmt.Errorf("usage: share port <vm> <port>")
		} else {
			_, err = d.Conn.Exec("UPDATE apps SET http_port = ? WHERE id = ?", args[2], appID)
			if err == nil {
				// Update Caddy too
				port := 80
				fmt.Sscanf(args[2], "%d", &port)
				c.UpsertRoute(vmName, cfg.Domain, port)
			}
		}
	case "add":
		if len(args) < 3 {
			err = fmt.Errorf("usage: share add <vm> <email>")
		} else {
			_, err = d.Conn.Exec("INSERT OR IGNORE INTO app_shares (app_id, email) VALUES (?, ?)", appID, args[2])
		}
	case "remove":
		if len(args) < 3 {
			err = fmt.Errorf("usage: share remove <vm> <email>")
		} else {
			_, err = d.Conn.Exec("DELETE FROM app_shares WHERE app_id = ? AND email = ?", appID, args[2])
		}
	default:
		err = fmt.Errorf("unknown share command: %s", cmd)
	}

	if err == nil {
		// Audit log
		remoteIP, _ := sess.Context().Value("remote_ip").(string)
		details := cmd
		if len(args) >= 3 {
			details = cmd + " " + args[2]
		}
		d.LogAudit("share_change", userID, vmName, remoteIP, details)
	}

	if isJSON {
		WriteJSON(sess, err == nil, "", nil, err)
	} else if err != nil {
		fmt.Fprintf(sess, "Error: %v\n", err)
	} else {
		fmt.Fprintf(sess, "Successfully updated sharing for '%s'\n", vmName)
	}
}

func handleHelp(sess ssh.Session) {
	help := `
Available commands:
  ls                     List your apps
  new --name=X           Create a new app
  rm <app>               Delete an app
  share <cmd> <vm>       Update sharing settings
  keys [add|rm]          Manage SSH keys
  whoami                 Show user info
  help                   Show this help
  exit                   Disconnect

Options:
  --json                 Output in JSON format
`
	fmt.Fprint(sess, help)
}

func handleKeys(sess ssh.Session, args []string, d *db.Database, userID int, isJSON bool) {
	if len(args) == 0 {
		// List keys
		rows, err := d.Conn.Query("SELECT fingerprint, comment, created_at FROM public_keys WHERE user_id = ?", userID)
		if err != nil {
			if isJSON {
				WriteJSON(sess, false, "", nil, err)
			} else {
				fmt.Fprintf(sess, "Error listing keys: %v\n", err)
			}
			return
		}
		defer rows.Close()

		type keyInfo struct {
			Fingerprint string `json:"fingerprint"`
			Comment     string `json:"comment"`
			Created     string `json:"created_at"`
		}
		var keys []keyInfo

		if !isJSON {
			fmt.Fprintf(sess, "%-50s %-20s %-20s\n", "FINGERPRINT", "COMMENT", "CREATED")
			fmt.Fprintf(sess, "%-50s %-20s %-20s\n", "-----------", "-------", "-------")
		}

		for rows.Next() {
			var fp, comment, created string
			rows.Scan(&fp, &comment, &created)
			if isJSON {
				keys = append(keys, keyInfo{fp, comment, created})
			} else {
				fmt.Fprintf(sess, "%-50s %-20s %-20s\n", fp, comment, created)
			}
		}

		if isJSON {
			WriteJSON(sess, true, "", map[string]interface{}{"keys": keys}, nil)
		}
		return
	}

	cmd := args[0]
	switch cmd {
	case "add":
		if len(args) < 2 {
			msg := "Usage: keys add <public_key_data>"
			if isJSON {
				WriteJSON(sess, false, "", nil, errors.New(msg))
			} else {
				fmt.Fprintln(sess, msg)
			}
			return
		}
		// Join all remaining args as key data (may contain spaces)
		keyData := strings.Join(args[1:], " ")
		// Parse the public key and compute proper fingerprint
		pubKey, comment, _, _, err := gossh.ParseAuthorizedKey([]byte(keyData))
		if err != nil {
			if isJSON {
				WriteJSON(sess, false, "", nil, fmt.Errorf("invalid SSH key format: %v", err))
			} else {
				fmt.Fprintf(sess, "Error: invalid SSH key format: %v\n", err)
			}
			return
		}
		fingerprint := gossh.FingerprintSHA256(pubKey)

		_, err = d.Conn.Exec("INSERT INTO public_keys (user_id, fingerprint, key_data, comment) VALUES (?, ?, ?, ?)",
			userID, fingerprint, keyData, comment)
		if err != nil {
			if isJSON {
				WriteJSON(sess, false, "", nil, err)
			} else {
				fmt.Fprintf(sess, "Error adding key: %v\n", err)
			}
			return
		}
		if isJSON {
			WriteJSON(sess, true, "Key added", map[string]string{"fingerprint": fingerprint}, nil)
		} else {
			fmt.Fprintf(sess, "Key added: %s\n", fingerprint)
		}

	case "rm":
		if len(args) < 2 {
			msg := "Usage: keys rm <fingerprint>"
			if isJSON {
				WriteJSON(sess, false, "", nil, errors.New(msg))
			} else {
				fmt.Fprintln(sess, msg)
			}
			return
		}
		fp := args[1]
		result, err := d.Conn.Exec("DELETE FROM public_keys WHERE user_id = ? AND fingerprint = ?", userID, fp)
		if err != nil {
			if isJSON {
				WriteJSON(sess, false, "", nil, err)
			} else {
				fmt.Fprintf(sess, "Error removing key: %v\n", err)
			}
			return
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			msg := "Key not found"
			if isJSON {
				WriteJSON(sess, false, "", nil, errors.New(msg))
			} else {
				fmt.Fprintln(sess, msg)
			}
			return
		}
		if isJSON {
			WriteJSON(sess, true, "Key removed", nil, nil)
		} else {
			fmt.Fprintln(sess, "Key removed")
		}

	default:
		msg := "Usage: keys [add <key>|rm <fingerprint>]"
		if isJSON {
			WriteJSON(sess, false, "", nil, errors.New(msg))
		} else {
			fmt.Fprintln(sess, msg)
		}
	}
}
