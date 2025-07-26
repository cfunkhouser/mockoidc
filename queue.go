package mockoidc

import "sync"

// UserQueue manages the queue of Users returned for each
// call to the authorize endpoint
type UserQueue struct {
	sync.Mutex
	queue       []User
	defaultUser User
}

// CodeQueue manages the queue of codes returned for each
// call to the authorize endpoint
type CodeQueue struct {
	sync.Mutex
	queue []string
}

// ErrorQueue manages the queue of errors for handlers to return
type ErrorQueue struct {
	sync.Mutex
	queue []*ServerError
}

// ServerError is a tester-defined error for a handler to return
type ServerError struct {
	Code        int
	Error       string
	Description string
}

// Push adds a User to the Queue to be set in subsequent calls to the
// `authorization_endpoint`
func (q *UserQueue) Push(user User) {
	q.Lock()
	defer q.Unlock()
	q.queue = append(q.queue, user)
}

// Pop a User from the Queue. If empty, return `DefaultUser()`
func (q *UserQueue) Pop() User {
	q.Lock()
	defer q.Unlock()

	if len(q.queue) == 0 {
		if q.defaultUser != nil {
			return q.defaultUser
		}

		return DefaultUser()
	}

	var user User
	user, q.queue = q.queue[0], q.queue[1:]
	return user
}

func (q *UserQueue) SetDefaultUser(user User) {
	q.Lock()
	defer q.Unlock()

	q.defaultUser = user
}

// Push adds a code to the Queue to be returned by subsequent
// `authorization_endpoint` calls as the code
func (q *CodeQueue) Push(code string) {
	q.Lock()
	defer q.Unlock()
	q.queue = append(q.queue, code)
}

// Pop a `code` from the Queue. If empty, return a random code
func (q *CodeQueue) Pop() (string, error) {
	q.Lock()
	defer q.Unlock()

	if len(q.queue) == 0 {
		code, err := randomNonce(24)
		if err != nil {
			return "", err
		}
		return code, nil
	}

	var code string
	code, q.queue = q.queue[0], q.queue[1:]
	return code, nil
}

// Push adds a ServerError to the Queue to be returned in subsequent
// handler calls
func (q *ErrorQueue) Push(se *ServerError) {
	q.Lock()
	defer q.Unlock()
	q.queue = append(q.queue, se)
}

// Pop a ServerError from the Queue. If empty, return nil
func (q *ErrorQueue) Pop() *ServerError {
	q.Lock()
	defer q.Unlock()

	if len(q.queue) == 0 {
		return nil
	}

	var se *ServerError
	se, q.queue = q.queue[0], q.queue[1:]
	return se
}
