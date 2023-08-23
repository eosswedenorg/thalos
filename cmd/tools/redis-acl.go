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

	"github.com/urfave/cli/v2"
)

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
		idx := rand.Intn(len(charset))
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

var RedisACLCmd = &cli.Command{
	Name:  "redis-acl",
	Usage: "create a users.acl file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "default-pw",
			Usage: "Password to use for the default account, if not provided a random one will be generated",
		},
		&cli.StringFlag{
			Name:  "client",
			Value: "thalos-client",
			Usage: "Thalos client account name",
		},
		&cli.StringFlag{
			Name:  "client-pw",
			Usage: "Password to use for the thalos client account, if not provided a random one will be generated",
		},
		&cli.StringFlag{
			Name:  "server",
			Value: "thalos",
			Usage: "Thalos account name",
		},
		&cli.StringFlag{
			Name:  "server-pw",
			Usage: "Password to use for the thalos server account, if not provided a random one will be generated",
		},
		&cli.StringFlag{
			Name:  "prefix",
			Value: "ship",
			Usage: "Redis key prefix",
		},
		&cli.BoolFlag{
			Name:  "cleartext",
			Usage: "If passwords should be hashed or left in cleartext.",
		},
		&cli.StringFlag{
			Name:        "file",
			DefaultText: "Standard out",
			Usage:       "Where the config should be written to",
		},
	},
	Action: func(ctx *cli.Context) error {
		var err error
		var out *os.File = os.Stdout

		rand.Seed(time.Now().Unix())

		defaultUser := NewUser("default", ctx.String("default-pw"))
		serverUser := NewUser(ctx.String("server"), ctx.String("server-pw"))
		clientUser := NewUser(ctx.String("client"), ctx.String("client-pw"))

		atleastOneGeneratedPw := defaultUser.Generated || serverUser.Generated || clientUser.Generated

		if !ctx.Bool("cleartext") {
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

		filename := ctx.String("file")
		if len(filename) > 0 {
			out, err = os.Create(filename)
			if err != nil {
				return err
			}
			defer out.Close()
		} else if !ctx.Bool("cleartext") && atleastOneGeneratedPw {
			fmt.Println()
		}

		return writeTemplate(out, defaultUser, serverUser, clientUser, ctx.String("prefix"))
	},
}
