package grafana

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"

	gapi "github.com/grafana/grafana-api-golang-client"
)

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		Create: CreateUser,
		Read:   ReadUser,
		Update: UpdateUser,
		Delete: DeleteUser,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"login": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func CreateUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)
	user := gapi.User{
		Email:    d.Get("email").(string),
		Name:     d.Get("name").(string),
		Login:    d.Get("login").(string),
		Password: d.Get("password").(string),
	}
	id, err := client.CreateUser(user)
	if err != nil {
		return err
	}
	d.SetId(strconv.FormatInt(id, 10))
	return ReadUser(d, meta)
}

func ReadUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}
	user, err := client.User(id)
	if err != nil {
		return err
	}
	d.Set("email", user.Email)
	d.Set("name", user.Name)
	d.Set("login", user.Login)
	return nil
}

func UpdateUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}
	u := gapi.User{
		Id:    id,
		Email: d.Get("email").(string),
		Name:  d.Get("name").(string),
		Login: d.Get("login").(string),
	}
	err = client.UserUpdate(u)
	if err != nil {
		return err
	}
	if d.HasChange("password") {
		err = client.UpdateUserPassword(id, d.Get("password").(string))
		if err != nil {
			return err
		}
	}
	return ReadUser(d, meta)
}

func DeleteUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gapi.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}
	return client.DeleteUser(id)
}
