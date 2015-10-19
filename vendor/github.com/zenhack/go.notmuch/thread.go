package notmuch

// Copyright © 2015 The go.notmuch Authors. Authors can be found in the AUTHORS file.
// Licensed under the GPLv3 or later.
// See COPYING at the root of the repository for details.

// #cgo LDFLAGS: -lnotmuch
// #include <stdlib.h>
// #include <notmuch.h>
import "C"
import (
	"runtime"
	"strings"
	"time"
)

// Thread represents a notmuch thread.
type Thread struct {
	cptr    *C.notmuch_thread_t
	threads *Threads
}

// Subject returns the subject of a thread.
func (t *Thread) Subject() string {
	cstr := C.notmuch_thread_get_subject(t.cptr)
	str := C.GoString(cstr)
	return str
}

// ID returns the ID of the thread.
func (t *Thread) ID() string {
	return C.GoString(C.notmuch_thread_get_thread_id(t.cptr))
}

// Count returns the total number of messages in the current thread.
func (t *Thread) Count() int {
	return int(C.notmuch_thread_get_total_messages(t.cptr))
}

// CountMatched returns the total number of messages in the current thread that
// matched the search.
func (t *Thread) CountMatched() int {
	return int(C.notmuch_thread_get_matched_messages(t.cptr))
}

// TopLevelMessages returns an iterator for the top-level messages in the
// current thread in oldest-first order.
func (t *Thread) TopLevelMessages() *Messages {
	return &Messages{
		thread: t,
		cptr:   C.notmuch_thread_get_toplevel_messages(t.cptr),
	}
}

// Messages returns an iterator for all messages in the current thread in
// oldest-first order.
func (t *Thread) Messages() *Messages {
	return &Messages{
		thread: t,
		cptr:   C.notmuch_thread_get_messages(t.cptr),
	}
}

// Authors returns the list of authors, the first are the authors that matched
// the query whilst the second return are the rest of the authors. All authors
// are ordered by date.
func (t *Thread) Authors() ([]string, []string) {
	var matched, unmatched []string

	as := C.GoString(C.notmuch_thread_get_authors(t.cptr))
	munm := strings.Split(as, "|")
	if len(munm) > 1 {
		matched = strings.Split(munm[0], ",")
		unmatched = strings.Split(munm[1], ",")
	} else {
		unmatched = strings.Split(munm[0], ",")
	}
	for i, s := range matched {
		matched[i] = strings.Trim(s, " ")
	}
	for i, s := range unmatched {
		unmatched[i] = strings.Trim(s, " ")
	}
	return matched, unmatched
}

// OldestDate returns the date of the oldest message in the thread.
func (t *Thread) OldestDate() time.Time {
	ctime := C.notmuch_thread_get_oldest_date(t.cptr)
	return time.Unix(int64(ctime), 0)
}

// NewestDate returns the date of the oldest message in the thread.
func (t *Thread) NewestDate() time.Time {
	ctime := C.notmuch_thread_get_newest_date(t.cptr)
	return time.Unix(int64(ctime), 0)
}

// Tags returns the tags for the current thread, returning a *Tags which can be
// used to iterate over all tags using `Tags.Next(Tag)`
//
// Note: In the Notmuch database, tags are stored on individual messages, not
// on threads. So the tags returned here will be all tags of the messages which
// matched the search and which belong to this thread.
func (t *Thread) Tags() *Tags {
	ts := &Tags{
		cptr:   C.notmuch_thread_get_tags(t.cptr),
		thread: t,
	}
	runtime.SetFinalizer(ts, func(ts *Tags) {
		C.notmuch_tags_destroy(ts.cptr)
	})
	return ts
}
