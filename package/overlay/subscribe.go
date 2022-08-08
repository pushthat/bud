package overlay

import "github.com/pushthat/bud/internal/pubsub"

func (f *FileSystem) Subscribe(name string) pubsub.Subscription {
	return f.ps.Subscribe(name)
}
