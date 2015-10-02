package grpc

import (
	"github.com/gmuch/gmuch/server"
	"golang.org/x/net/context"
)

// Binding provides gRPC bindings for the server.
type Binding struct {
	server.GmuchService
}

// Query implements GmuchService.Query.
func (g Binding) Query(ctx context.Context, req *QueryRequest) (*QueryResponse, error) {
	ts, err := g.GmuchService.Query(req.Query, int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, err
	}

	qr := &QueryResponse{
		Threads: make([]*Thread, len(ts)),
	}
	for i, t := range ts {
		qr.Threads[i] = &Thread{
			Id:      t.ID,
			Subject: t.Subject,
			Authors: t.Authors,
		}
	}
	return qr, nil
}

// Thread implements GmuchService.Thread.
func (g Binding) Thread(ctx context.Context, req *ThreadRequest) (*ThreadResponse, error) {
	t, err := g.GmuchService.Thread(req.Id)
	if err != nil {
		return nil, err
	}

	tr := &ThreadResponse{
		Thread: &Thread{
			Id:       t.ID,
			Subject:  t.Subject,
			Authors:  t.Authors,
			Messages: make([]*Message, len(t.Messages)),
		},
	}

	for i, m := range t.Messages {
		ems := make([]*Email, len(m.Emails))

		for j, em := range m.Emails {
			ems[j] = &Email{
				From:        em.From,
				To:          em.To,
				Cc:          em.Cc,
				Bcc:         em.Bcc,
				Subject:     em.Subject,
				Text:        string(em.Text),
				Html:        string(em.HTML),
				Headers:     make([]*Pair, 0, len(em.Headers)),
				Attachments: make([]*Attachment, len(em.Attachments)),
			}

			for k, v := range em.Headers {
				ems[j].Headers = append(ems[j].Headers, &Pair{Key: k, Value: v})
			}

			for l, a := range em.Attachments {
				pa := &Attachment{
					Filename: a.Filename,
					Content:  a.Content,
				}

				for k, v := range a.Header {
					pa.Headers = append(pa.Headers, &Pair{Key: k, Value: v})
				}

				ems[j].Attachments[l] = pa
			}
		}

		tr.Thread.Messages[i] = &Message{
			Id:       m.ID,
			ThreadId: m.ThreadID,
			Tags:     m.Tags,
			Emails:   ems,
		}
	}

	return tr, nil
}
