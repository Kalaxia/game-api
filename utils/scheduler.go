package utils

import(
    "errors"
    "math/rand"
    "time"
)

type Scheduling struct {
    Queue map[string]Task
}
type Task struct {
    Id string
    Timer *time.Timer
}

var Scheduler Scheduling

func init() {
    Scheduler = Scheduling{
        Queue: make(map[string]Task, 0),
    }
}

func (s *Scheduling) AddTask(seconds uint, callback func()) {
    id := string(rand.Intn(100000)) + time.Now().Format(time.UnixDate)
    task := Task{
        Id: id,
        Timer: time.AfterFunc(time.Second * time.Duration(seconds), func() {
            callback()
            s.RemoveTask(id)
        }),
    }
    s.Queue[id] = task
}

func (s *Scheduling) RemoveTask(id string) {
    _, isset := s.Queue[id]
    if !isset {
        panic(errors.New("Unknown scheduler id"))
    }
    delete(s.Queue, id)
}

func (s *Scheduling) CancelTask(id string) {
    task, isset := s.Queue[id]
    if !isset {
        panic(errors.New("Unknown scheduler id"))
    }
    task.Timer.Stop()
    delete(s.Queue, id)
}
