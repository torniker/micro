package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/signup"
	"github.com/micro/micro/v3/client/cli/token"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/internal/report"
	"github.com/micro/micro/v3/service/auth"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

// login flow.
// For documentation of the flow please refer to https://github.com/micro/development/pull/223
func login(ctx *cli.Context) error {
	// assuming --otp go to platform.Signup
	if isOTP := ctx.Bool("otp"); isOTP {
		return signup.Run(ctx)
	}

	// otherwise assume username/password login

	// get the environment
	env := util.GetEnv(ctx)
	// get the username
	username := ctx.String("username")

	// username is blank
	if len(username) == 0 {
		fmt.Print("Enter username: ")
		// read out the username from prompt if blank
		reader := bufio.NewReader(os.Stdin)
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)
	}

	// clear tokens and try again
	if err := token.Remove(env.Name); err != nil {
		report.Errorf(ctx, "%v: Token remove: %v", username, err.Error())
		return err
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return err
	}

	password := ctx.String("password")
	if len(password) == 0 {
		pw, err := getPassword()
		if err != nil {
			return err
		}
		password = strings.TrimSpace(pw)
		fmt.Println()
	}
	tok, err := auth.Token(auth.WithCredentials(username, password), auth.WithTokenIssuer(ns))
	if err != nil {
		report.Errorf(ctx, "%v: Getting token: %v", username, err.Error())
		return err
	}
	token.Save(env.Name, tok)

	fmt.Println("Successfully logged in.")
	report.Success(ctx, username)
	return nil
}

// taken from https://stackoverflow.com/questions/2137357/getpasswd-functionality-in-go
func getPassword() (string, error) {
	fmt.Print("Enter password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	password := string(bytePassword)
	return strings.TrimSpace(password), nil
}

func logout(ctx *cli.Context) error {
	env := util.GetEnv(ctx)
	return token.Remove(env.Name)
}
