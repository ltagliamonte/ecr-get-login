package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"os"
	"strings"
)

func main() {
	var region string
	var helpShort, helpLong bool

	flags := flag.NewFlagSet("ecr-get-login", flag.ContinueOnError)
	flags.StringVar(&region, "region", "", "Connect to this AWS region.")
	flags.BoolVar(&helpShort, "h", false, "Display usage.")
	flags.BoolVar(&helpLong, "help", false, "Display usage.")
	err := flags.Parse(os.Args[1:])

	if helpShort || helpLong {
		usage()
		os.Exit(0)
	} else if err != nil || flags.NArg() == 0 {
		usage()
		os.Exit(1)
	}

	if region == "" {
		region = os.Getenv("AWS_REGION")
	}

	if err := login(region, flags.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Retrieve ECR Docker login credentials.")
	fmt.Fprintln(os.Stderr, "usage: ecr-get-login [-region=REGION] ACCOUNT...")
}

// Log Docker into the ECR accounts in the given region.
func login(region string, accounts []string) error {
	svc := ecr.New(session.New(), &aws.Config{
		Region: aws.String(region),
	})

	res, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{
		RegistryIds: aws.StringSlice(accounts),
	})
	if err != nil {
		return err
	}

	format := "docker login -u AWS -p %s %s\n"
	for _, data := range res.AuthorizationData {
		password, err := decodeAuth(aws.StringValue(data.AuthorizationToken))
		if err != nil {
			return nil
		}
		fmt.Printf(format, password, aws.StringValue(data.ProxyEndpoint))
	}
	return nil
}

// Decode the base64 auth string.
func decodeAuth(auth string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return "", err
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", errors.New("auth data contains invalid payload")
	}
	return parts[1], nil
}
