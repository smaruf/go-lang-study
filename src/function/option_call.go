type SomeClient struct {
  http *http.Client
  endpoint string
}
// Option is an optional argument to New().
type Option func(c *SomeClient)
// WithEndpoint is an option that makes the client use the 
// endpoint passed.
func WithEndpoint(endpoint string) Option {
  return func(c *SomeClient) {
    c.endpoint = endpoint
  }
}
func New(*httpClient *http.Client, options ...Options) *SomeClient { 
  c := &SomeClient{
    http: httpClient,
    endpoint: "https://myclientendpoint.com",
  }

  for _, o := range options {
    o(c)
  }
  return c
}
