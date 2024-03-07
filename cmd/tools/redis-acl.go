package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rnd *rand.Rand

// Helper struct representing a redis user.
type User struct {
	// Username
	Name string

	// Password
	Password string

	// True if password was generated, false if not.
	Generated bool
}

func NewUser(name, password string) User {
	if len(password) < 1 {
		return User{
			Name:      name,
			Password:  randomString(32),
			Generated: true,
		}
	}
	return User{Name: name, Password: password}
}

func (u *User) Hash() {
	u.Password = "#" + hash(u.Password)
}

func (u User) Print() {
	fmt.Println(u.Name+":", u.Password)
}

func (u User) PrintIfGeneratedPW() {
	if u.Generated {
		u.Print()
	}
}

func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUWXYZ0123456789"
	out := ""
	for i := 0; i < length; i++ {
		idx := rnd.Intn(len(charset))
		out += string(charset[idx])
	}
	return out
}

func hash(str string) string {
	data := sha256.Sum256([]byte(str))
	return hex.EncodeToString(data[:])
}

func writeTemplate(w io.Writer, defUser, serverUser, clientUser User, prefix string) error {
	tmplStr := `# Created by thalos-tools on {{.timestamp}}
user default on {{.defaultpw}} ~* &* +@all
user {{.server}} on {{.serverpw}} resetchannels ~{{.prefix}}::* &{{.prefix}}::* -@all +get +publish +set
user {{.client}} on {{.clientpw}} resetchannels &{{.prefix}}::* -@all +subscribe
`

	tmpl, err := template.New("acl").Parse(tmplStr)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, map[string]string{
		"defaultpw": defUser.Password,
		"client":    clientUser.Name,
		"clientpw":  clientUser.Password,
		"server":    serverUser.Name,
		"serverpw":  serverUser.Password,
		"prefix":    prefix,
		"timestamp": time.Now().Format(time.UnixDate),
	})
}

var RedisACLCmd = &cobra.Command{
	Use:   "redis-acl",
	Short: "create a users.acl file",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		out := os.Stdout

		rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

		flagDefUserPw, _ := cmd.Flags().GetString("default-pw")
		flagServer, _ := cmd.Flags().GetString("server")
		flagServerPw, _ := cmd.Flags().GetString("server-pw")
		flagClient, _ := cmd.Flags().GetString("client")
		flagClientPw, _ := cmd.Flags().GetString("client-pw")
		flagPrefix, _ := cmd.Flags().GetString("prefix")

		defaultUser := NewUser("default", flagDefUserPw)
		serverUser := NewUser(flagServer, flagServerPw)
		clientUser := NewUser(flagClient, flagClientPw)

		atleastOneGeneratedPw := defaultUser.Generated || serverUser.Generated || clientUser.Generated

		cleartext, _ := cmd.Flags().GetBool("cleartext")
		if !cleartext {
			if atleastOneGeneratedPw {
				println("Passwords")
			}

			defaultUser.PrintIfGeneratedPW()
			serverUser.PrintIfGeneratedPW()
			clientUser.PrintIfGeneratedPW()

			defaultUser.Hash()
			serverUser.Hash()
			clientUser.Hash()
		}

		filename, _ := cmd.Flags().GetString("file")
		if len(filename) > 0 {
			out, err = os.Create(filename)
			if err != nil {
				log.WithError(err).Fatal("Failed to create output file")
				return
			}
			defer out.Close()
		} else if !cleartext && atleastOneGeneratedPw {
			fmt.Println()
		}

		err = writeTemplate(out, defaultUser, serverUser, clientUser, flagPrefix)
		if err != nil {
			log.WithError(err).Fatal("Failed to writte config")
			return
		}
	},
}
