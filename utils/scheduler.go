package utils

import(
    "kalaxia-game-api/exception"
    //"math/rand" // TODO remove
    "time"
)

type Scheduling struct {
    Ticker time.Ticker
    Queue map[string]Task
}
type Task struct {
    Id string
    Timer *time.Timer
}

var Scheduler Scheduling
var counter uint32

func init() {
    counter = 0;
    Scheduler = Scheduling{
        Queue: make(map[string]Task, 0),
    }
}

func (s *Scheduling) AddTask(seconds uint, callback func()) {
    id := getNewIdForTask();
    task := Task{
        Id: id,
        Timer: time.AfterFunc(time.Second * time.Duration(seconds), func() {
            callback()
            s.RemoveTask(id)
        }),
    }
    s.Queue[id] = task
}

func (s *Scheduling) AddHourlyTask(callback func()) {
    now := time.Now()
    nextHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour() + 1, 0, 0, 0, time.UTC)

    s.AddTask(uint(time.Until(nextHour).Seconds() + 1), func() {
        callback()
        s.AddHourlyTask(callback)
    })
}

func (s *Scheduling) RemoveTask(id string) {
    _, isset := s.Queue[id]
    if !isset {
        panic(exception.NewException("unknown scheduler id", nil))
    }
    delete(s.Queue, id)
}

func (s *Scheduling) CancelTask(id string) {
    task, isset := s.Queue[id]
    if !isset {
        panic(exception.NewHttpException(500, "unknown scheduler id", nil))
    }
    task.Timer.Stop()
    delete(s.Queue, id)
}


func getNewIdForTask () string{
    //  string(rand.Intn(100000)) +"-"+   // TODO remove
    id := string(counter)+"-"+ time.Now().Format(time.UnixDate); // this sould be more robust than a random number
    counter++;
    return id;
}
