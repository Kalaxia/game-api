package utils

import(
    "errors"
    "time"
    "github.com/rs/xid"
)

type Scheduling struct {
    Queue map[string]Task
}
type Task struct {
    Id string
    Timer *time.Timer
}

var Scheduler Scheduling
var guid xid.ID

func init() {
    Scheduler = Scheduling{
        Queue: make(map[string]Task, 0),
    }
	guid = xid.New()
}

func (s *Scheduling) AddTask(seconds uint, callback func()) string {
    id := guid.String()
    task := Task{
        Id: id,
        Timer: time.AfterFunc(time.Second * time.Duration(seconds), func() {
            callback()
            s.RemoveTask(id)
        }),
    }
    s.Queue[id] = task
    return id
}

func (s *Scheduling) RemoveTask(id string) {
    _, isset := s.Queue[id]
    if !isset {
        panic(errors.New("unknown scheduler id"))
    }
    delete(s.Queue, id)
}

func (s *Scheduling) CancelTask(id string) {
    task, isset := s.Queue[id]
    if !isset {
        panic(errors.New("unknown scheduler id"))
    }
    task.Timer.Stop()
    delete(s.Queue, id)
}
