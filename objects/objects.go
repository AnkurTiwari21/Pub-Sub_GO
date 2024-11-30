package objects

type HTTPRequestFormat struct {
	Task         string
	BindingKey   string
	QueueName    string
	ExchangeName string
	Durability   string
}

type CommunicationMessage struct {
	Task           string `json:"task,omitempty"`
	BindingKey     string `json:"binding_key,omitempty"`
	QueueName      string `json:"queue_name,omitempty"`
	ExchangeType   string `json:"exchange_type,omitempty"`
	Durability     bool   `json:"durability,omitempty"`
	Status         string `json:"status,omitempty"`
	Message        string `json:"message,omitempty"`
	UserConnString string `json:"user_uuid,omitempty"`
	RoutingKey     string `json:"routing_key,omitempty"`
}

type Response struct {
	Message        string `json:"message,omitempty"`
	Data           string `json:"data,omitempty"`
	UserWSConnUUID string `json:"conn,omitempty"`
	ErrorMessage   string `json:"error_message,omitempty"`
}

type CompleteResponse struct {
	QueueName string   `json:"queue_name"`
	Messages  []string `json:"message"`
}
