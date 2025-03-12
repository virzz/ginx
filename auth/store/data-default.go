package store

type DefaultData struct {
	Token_   string          `json:"token" redis:"token"`
	ID_      string          `json:"id" redis:"id"`
	Account_ string          `json:"account" redis:"account"`
	Roles_   DataStringSlice `json:"roles" redis:"roles"`
}

var _ Data = (*DefaultData)(nil)

func (d *DefaultData) New()            { d.Token_ = generateRandomKey() }
func (d *DefaultData) Token() string   { return d.Token_ }
func (d *DefaultData) ID() string      { return d.ID_ }
func (d *DefaultData) Account() string { return d.Account_ }
func (d *DefaultData) Roles() []string { return []string(d.Roles_) }
func (d *DefaultData) SetToken(v string) Data {
	d.Token_ = v
	return d
}
func (d *DefaultData) SetID(v string) Data {
	d.ID_ = v
	return d
}
func (d *DefaultData) SetAccount(v string) Data {
	d.Account_ = v
	return d
}
func (d *DefaultData) SetRoles(v []string) Data {
	d.Roles_ = DataStringSlice(v)
	return d
}

func (d *DefaultData) Set(key string, val any) Data {
	switch key {
	case "id":
		if v, ok := val.(string); ok {
			d.ID_ = v
		}
	case "account":
		if v, ok := val.(string); ok {
			d.Account_ = v
		}
	case "roles":
		if v, ok := val.([]string); ok {
			d.Roles_ = DataStringSlice(v)
		}
	}
	return d
}

func (d *DefaultData) Get(key string) any {
	switch key {
	case "id":
		return d.ID_
	case "account":
		return d.Account_
	case "roles":
		return []string(d.Roles_)
	}
	return nil
}
