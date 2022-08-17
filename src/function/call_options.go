type tokenOptions struct {
  account string
  verifier string
}
// tokenOption is an optional argument for ByToken().
// I've chosen to keep it private, but that is not required.
type tokenOption interface {
  token() // This must be kept private.
}
func (c *Client) ByToken(token string, options ...tokenOption) (string, error) {
  opts := tokenOptions{account: "default"} 
  if err := calloptions.ApplyOptions(&opts, options); err != nil {
    return "", err
  }
  // Put the rest of your code here.
}
type userPassOptions struct {
  account string
}
type userPassOption interface {
  userPass()
}
func (c *Client) ByUserPass(user, pass string, options ...userPassOption) (string, error) {
  opts := userPassOptions{account: "default"}
  if err := calloptions.ApplyOptions(&opts, options); err != nil {
    return "", err
  }
  // Put the rest of your code here.
}
