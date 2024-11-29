package objects

type HTTPRequestFormat struct {
	Task         string
	BindingKey   string
	QueueName    string
	ExchangeName string
	Durability   string
}
