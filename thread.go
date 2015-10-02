package gmuch

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/gmuch/gmuch/model"
	"github.com/gmuch/gmuch/server"
	"github.com/jordan-wright/email"
	"github.com/zenhack/go.notmuch"
)

const threadPrefix = "thread:"

var (
	// ErrIDInvalid is returned if the ID is invalid.
	ErrIDInvalid = errors.New("thread ID is invalid, it should only contain alphanumeric characters")

	// ErrIDMatchedSeveral is returned when the ID has matched several threads.
	ErrIDMatchedSeveral = errors.New("thread ID matched several threads")

	// ErrIDNotFound is returned when the ID did not match any threads.
	ErrIDNotFound = errors.New("thread ID matched no threads")

	idRegex = regexp.MustCompile(`^[a-z0-9]+$`)
)

// Thread returns a *ThreadResponse for a given thread ID.
func (g *Gmuch) Thread(id string) (*server.ThreadResponse, error) {
	if !idRegex.MatchString(id) {
		return nil, ErrIDInvalid
	}
	db, err := notmuch.Open(g.dbPath, notmuch.DBReadOnly)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	q := db.NewQuery(fmt.Sprintf("thread:%s", id))
	if q.CountThreads() > 1 {
		return nil, ErrIDMatchedSeveral
	}

	threads, err := q.Threads()
	if err != nil {
		return nil, err
	}

	thread := &notmuch.Thread{}
	if !threads.Next(thread) {
		return nil, ErrIDNotFound
	}
	m, um := thread.Authors()
	t := model.Thread{
		ID:       thread.ID(),
		Subject:  thread.Subject(),
		Authors:  append(m, um...),
		Messages: make([]*model.Message, 0, thread.Count()),
	}
	messages := thread.Messages()
	msg := &notmuch.Message{}
	for messages.Next(msg) {
		e, err := g.getMessage(msg)
		if err != nil {
			return nil, err
		}
		t.Messages = append(t.Messages, e)
	}

	return &server.ThreadResponse{
		Thread: t,
	}, nil
}

func (g *Gmuch) getMessage(m *notmuch.Message) (*model.Message, error) {
	msg := &model.Message{
		ID:       m.ID(),
		ThreadID: m.ThreadID(),
		Emails:   make([]*email.Email, 0, 1),
	}

	var fn string
	fns := m.Filenames()
	for fns.Next(&fn) {
		f, err := os.Open(fn)
		if err != nil {
			g.logger.Log(
				"method", "getMessage",
				"messageID", msg.ID,
				"operation", "open email file",
				"path", fn,
				"error", err,
			)
			return nil, err
		}
		defer f.Close()
		e, err := email.NewEmailFromReader(f)
		if err != nil {
			g.logger.Log(
				"method", "getMessage",
				"messageID", msg.ID,
				"operation", "parse email file",
				"path", fn,
				"error", err,
			)
			return nil, err
		}
		msg.Emails = append(msg.Emails, e)
	}

	tags := m.Tags()
	tag := &notmuch.Tag{}
	for tags.Next(tag) {
		msg.Tags = append(msg.Tags, tag.String())
	}

	return msg, nil
}
