package api

import(
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
var counter uint64

func InitScheduler() {
    counter = 0
    Scheduler = Scheduling{
        Queue: make(map[string]Task, 0),
    }
}

func (s *Scheduling) AddTask(date time.Time, callback func()) {
    id := getNewIdForTask()
    task := Task{
        Id: id,
        Timer: time.AfterFunc(time.Until(date), func() {
            defer CatchException(nil)
            callback()
            s.RemoveTask(id)
        }),
    }
    s.Queue[id] = task
}

func (s *Scheduling) AddDailyTask(callback func()) {
    now := time.Now()
    nextDay := time.Date(now.Year(), now.Month(), now.Day() + 1, 0, 0, 0, 0, time.UTC)

    s.AddTask(nextDay, func() {
        callback()
        s.AddDailyTask(callback)
    })
}

func (s *Scheduling) AddHourlyTask(callback func()) {
    now := time.Now()
    nextHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour() + 1, 0, 0, 0, time.UTC)

    s.AddTask(nextHour, func() {
        callback()
        s.AddHourlyTask(callback)
    })
}

func (s *Scheduling) RemoveTask(id string) {
    _, isset := s.Queue[id]
    if !isset {
        panic(NewException("unknown scheduler id", nil))
    }
    delete(s.Queue, id)
}

func (s *Scheduling) CancelTask(id string) {
    task, isset := s.Queue[id]
    if !isset {
        panic(NewHttpException(500, "unknown scheduler id", nil))
    }
    task.Timer.Stop()
    delete(s.Queue, id)
}


func getNewIdForTask () string {
    id := string(counter)+"-"+ time.Now().Format(time.UnixDate) // this sould be more robust than a random number
    if counter == 18446744073709551615 { // maxvalue for uint64
        counter = 0
    } else {
        counter++
    }
    
    return id
}
