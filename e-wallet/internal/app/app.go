package app

//nolint:tagliatelle
type Conf struct {
	StartingBalance float64 `yaml:"starting_balance"`
	MinimalAmount   float64 `yaml:"minimal_amount"`
}
