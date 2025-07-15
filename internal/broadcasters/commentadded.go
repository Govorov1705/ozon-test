package broadcasters

import (
	"sync"

	"github.com/Govorov1705/ozon-test/graph/model"
	"github.com/google/uuid"
)

type CommentAddedBroadcaster struct {
	mu          sync.RWMutex
	subscribers map[uuid.UUID][]chan *model.Comment
}

func NewCommentAddedBroadcaster() *CommentAddedBroadcaster {
	return &CommentAddedBroadcaster{
		subscribers: make(map[uuid.UUID][]chan *model.Comment),
	}
}

func (b *CommentAddedBroadcaster) Subscribe(postID uuid.UUID) <-chan *model.Comment {
	ch := make(chan *model.Comment, 1)

	b.mu.Lock()
	b.subscribers[postID] = append(b.subscribers[postID], ch)
	b.mu.Unlock()

	return ch
}

func (b *CommentAddedBroadcaster) Unsubscribe(postID uuid.UUID, ch <-chan *model.Comment) {
	b.mu.Lock()
	defer b.mu.Unlock()

	channels := b.subscribers[postID]
	for i, c := range channels {
		if c == ch {
			b.subscribers[postID] = append(channels[:i], channels[i+1:]...)
			break
		}
	}
}

func (b *CommentAddedBroadcaster) Publish(postID uuid.UUID, comment *model.Comment) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.subscribers[postID] {
		select {
		case ch <- comment:
		default:
		}
	}
}
