```Go
type OneLogin struct {
	// Allow custom base URL to override the generated URL
	CustomURL string

	// OneLogin service shard (eu, us, etc)
	Shard string

	// Token struct for managing the OAuth token
	Token *OneLogin_Token

	SubDomain     string
	Client_id     string
	Client_secret string
}
```

```Go
func (o *OneLogin) Authenticate(username string, password string, token string) error
```

```Go
func (o *OneLogin) GetUrl(uri string, args ...string) string
```

```Go
func (o *OneLogin) Get_AppsForUser(userId int) ([]OneLoginApp, error)
```

```Go
func (o *OneLogin) Get_Role_Id(role string) (int, error)
```

```Go
func (o *OneLogin) Get_Roles(name string) ([]OneLoginRole, error)
```

```Go
func (o *OneLogin) Get_RolesForUser(userId int) ([]int, error)
```

```Go
func (o *OneLogin) Get_Token() (*OneLogin_Token, error)
```

```Go
func (o *OneLogin) Get_Users(filter map[string]string) (*[]OneLoginUser, error)
```

```Go
func (o *OneLogin) Get_UsersWithRole(role string) ([]OneLoginUser, error)
```

```Go
func (o *OneLogin) Set_CustomAttribute(userid int, attributeName string, attributeValue string) error
```

```Go
func (o *OneLogin) VerifyToken(device_id string, state_token string, token string) error
```
