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

	// True if password should be hashed, false otherwise.
	Hash bool
}

func NewUser(name, password string, pass_len uint) User {
	if len(password) < 1 {
		return User{
			Name:      name,
			Password:  randomString(pass_len),
			Generated: true,
		}
	}
	return User{Name: name, Password: password}
}

func (u *User) GetPassword() string {
	if u.Hash {
		return "#" + hash(u.Password)
	}

	return ">" + u.Password
}

func (u User) Print() {
	fmt.Println(u.Name+":", u.Password)
}

func (u User) PrintIfGeneratedPW() {
	if u.Generated {
		u.Print()
	}
}

func randomString(length uint) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUWXYZ0123456789"
	out := ""
	for i := 0; i < int(length); i++ {
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
user {{.server}} on {{.serverpw}} resetchannels ~{{.prefix}}::* &{{.prefix}}::* -@all +ping +get +publish +set
user {{.client}} on {{.clientpw}} resetchannels &{{.prefix}}::* -@all +subscribe
`

	tmpl, err := template.New("acl").Parse(tmplStr)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, map[string]string{
		"defaultpw": defUser.GetPassword(),
		"client":    clientUser.Name,
		"clientpw":  clientUser.GetPassword(),
		"server":    serverUser.Name,
		"serverpw":  serverUser.GetPassword(),
		"prefix":    prefix,
		"timestamp": time.Now().Format(time.UnixDate),
	})
}

func CreateRedisACLCmd() *cobra.Command {
	cmd := &cobra.Command{
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
			flagPassLen, _ := cmd.Flags().GetUint("pass-len")

			defaultUser := NewUser("default", flagDefUserPw, flagPassLen)
			serverUser := NewUser(flagServer, flagServerPw, flagPassLen)
			clientUser := NewUser(flagClient, flagClientPw, flagPassLen)

			atleastOneGeneratedPw := defaultUser.Generated || serverUser.Generated || clientUser.Generated

			cleartext, _ := cmd.Flags().GetBool("cleartext")
			if !cleartext {
				if atleastOneGeneratedPw {
					println("Passwords")
				}

				defaultUser.PrintIfGeneratedPW()
				serverUser.PrintIfGeneratedPW()
				clientUser.PrintIfGeneratedPW()

				defaultUser.Hash = true
				serverUser.Hash = true
				clientUser.Hash = true
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

	cmd.Flags().String("default-pw", "", "Password to use for the default account, if not provided a random one will be generated")
	cmd.Flags().String("client", "thalos-client", "Thalos client account name")
	cmd.Flags().String("client-pw", "", "Password to use for the thalos client account, if not provided a random one will be generated")
	cmd.Flags().String("server", "thalos", "Thalos account name")
	cmd.Flags().String("server-pw", "", "Password to use for the thalos server account, if not provided a random one will be generated")
	cmd.Flags().String("prefix", "ship", "Redis key prefix")
	cmd.Flags().Bool("cleartext", false, "If passwords should be hashed or left in cleartext.")
	cmd.Flags().String("file", "", "Where the config should be written to (default: standard out)")
	cmd.Flags().Uint("pass-len", 32, "The length of generated passwords")

	return cmd
}
