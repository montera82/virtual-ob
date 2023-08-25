package mock

type (
	CircuitBreaker struct {
		ExecuteFunc        func(action any) (any, error)
		ExecuteFuncInvoked bool
	}
)

func (c *CircuitBreaker) Execute(action func() (any, error)) (any, error) {
	c.ExecuteFuncInvoked = true
	return c.ExecuteFunc(action)
}
