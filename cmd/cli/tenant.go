package main

import (
	"fmt"

	avisha "github.com/jackmordaunt/avisha-fn"
)

var tenant = Command{
	Name:   "tenant",
	Action: nil,
	Children: []Command{
		{
			Name: "register",
			Action: func(app *avisha.App, args []string) error {
				tenant := avisha.Tenant{}

				matcher := ArgMap{
					Handlers: map[string]func(string){
						"name":    Assigner(&tenant.Name),
						"contact": Assigner(&tenant.Contact),
					},
				}

				matcher.Match(args)

				fmt.Printf("tenant: %#v\n", tenant)

				if err := app.RegisterTenant(tenant); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
