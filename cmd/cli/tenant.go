package main

import avisha "github.com/jackmordaunt/avisha-fn"

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
						"name": func(value string) {
							tenant.Name = value
						},
						"contact": func(value string) {
							tenant.Contact = value
						},
					},
				}

				matcher.Match(args)

				if err := app.RegisterTenant(tenant); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
