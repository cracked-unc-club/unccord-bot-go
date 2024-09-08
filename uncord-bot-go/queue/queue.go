package queue

import (
	"math/rand"

	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type Queue struct {
	Tracks []lavalink.Track
}

func (q *Queue) Add(track lavalink.Track) {
	q.Tracks = append(q.Tracks, track)
}

func (q *Queue) Next() (lavalink.Track, bool) {
	if len(q.Tracks) == 0 {
		return lavalink.Track{}, false
	}
	track := q.Tracks[0]
	q.Tracks = q.Tracks[1:]
	return track, true
}

func (q *Queue) Shuffle() {
	rand.Shuffle(len(q.Tracks[1:]), func(i, j int) {
		q.Tracks[i+1], q.Tracks[j+1] = q.Tracks[j+1], q.Tracks[i+1]
	})
}

func (q *Queue) Skip(amount int) (lavalink.Track, bool) {
	if len(q.Tracks) == 0 {
		return lavalink.Track{}, false
	}
	if amount > len(q.Tracks) {
		amount = len(q.Tracks)
	}
	q.Tracks = q.Tracks[amount:]
	if len(q.Tracks) == 0 {
		return lavalink.Track{}, false
	}
	return q.Tracks[0], true
}

type QueueManager struct {
	Queues map[snowflake.ID]*Queue
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		Queues: make(map[snowflake.ID]*Queue),
	}
}

func (qm *QueueManager) Get(guildID snowflake.ID) *Queue {
	queue, ok := qm.Queues[guildID]
	if !ok {
		queue = &Queue{}
		qm.Queues[guildID] = queue
	}
	return queue
}
