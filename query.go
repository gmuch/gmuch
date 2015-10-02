package gmuch

import (
	"github.com/gmuch/gmuch/model"
	"github.com/zenhack/go.notmuch"
)

// Query returns a []*model.Thread for a given search query.
func (g *Gmuch) Query(qs string, offset, limit int) ([]*model.Thread, error) {
	db, err := notmuch.Open(g.dbPath, notmuch.DBReadOnly)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	q := db.NewQuery(qs)
	ts := make([]*model.Thread, 0, limit-offset)
	threads, err := q.Threads()
	if err != nil {
		return nil, err
	}

	thread := &notmuch.Thread{}
	for i := 0; i < offset; i++ {
		if !threads.Next(thread) {
			break
		}
	}

	for i := 0; threads.Next(thread); i++ {
		m, um := thread.Authors()

		t := &model.Thread{
			ID:      thread.ID(),
			Subject: thread.Subject(),
			Authors: append(m, um...),
		}
		ts = append(ts, t)

		if i == limit-offset-1 {
			break
		}
	}

	return ts, nil
}
