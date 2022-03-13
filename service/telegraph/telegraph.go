package telegraph

import (
	"github.com/faryne/api-server/config"
	t "gitlab.com/toby3d/telegraph"
)

func New() (*t.Account, error) {
	//("telegraph.default"),
	account := t.Account{
		AccessToken: config.Config.TelegraphToken,
	}
	if account.AccessToken == "" { // 如果沒有access token 的話，先呼叫 createAccount
		return t.CreateAccount(account)
	}
	return &account, nil
}
