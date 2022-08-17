// tokenOpts holds all our optional token values.
type tokenOpts struct {
  account string
}
// TokenOption is an optional argument to ByToken().
type TokenOption func(t *tokenOpts)
// WithAccount sets the account to name. For use with ByToken().
func WithAccount(name string) TokenOption {
  return func(t *tokenOpts) {
    t.account = name
  }
}
// ByToken does something with tokens.
func (a *Auth) ByToken(token string, options ...TokenOption) (string, error) {
  opts := tokenOpts{account: "default"}
  for _, o := range options {
    o(&opts)
  }
  // Do whatever the method does below here
}
