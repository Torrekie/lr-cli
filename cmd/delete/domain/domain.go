package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/loginradius/lr-cli/api"
	"github.com/loginradius/lr-cli/cmdutil"

	"github.com/spf13/cobra"
)

var fileName string

type domain struct {
	Domain string `json:"domain"`
}

type Result struct {
	CallbackUrl string `json:"CallbackUrl"`
}

func NewdomainCmd() *cobra.Command {
	opts := &domain{}

	cmd := &cobra.Command{
		Use:   "domain",
		Short: "delete domain",
		Long:  `This commmand deletes domain`,
		Example: heredoc.Doc(`$ lr delete domain --domain <domain>
		<doamin> is now removed from whitelisted domain."
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Domain == "" {
				return &cmdutil.FlagError{Err: errors.New("`domain` is required argument")}
			}
			p, err := api.GetSites()
			if err != nil {
				return err
			}
			urls := strings.Split(p.Callbackurl, ";")
			if !strings.Contains(p.Callbackurl, opts.Domain) {
				return &cmdutil.FlagError{Err: errors.New("Entered Domain not found")}
			}
			if len(urls) == 1 {
				return &cmdutil.FlagError{Err: errors.New("Cannot delete the last domain")}
			}
			for index, url := range urls {
				if url == opts.Domain {
					urls = append(urls[:index], urls[index+1:]...)
					break
				}
			}
			return delete(opts.Domain, urls)

		},
	}

	fl := cmd.Flags()
	fl.StringVarP(&opts.Domain, "domain", "d", "", "Enter Domain Value")

	return cmd
}

func delete(remVal string, allDomain []string) error {
	err := api.UpdateDomain(allDomain)
	if err != nil {
		return err
	}
	fmt.Println(remVal + " is now removed from whitelisted domain.")
	return nil
}
