package scheduler

import (
	"context"
	"github.com/Songmu/go-httpdate"
	"github.com/procyon-projects/chrono"
	"json-objects-vault/models"
	"log"
	"time"
)

var (
	scheduler = chrono.NewDefaultTaskScheduler()
)

// AddTask adds task for obj deletion
func AddTask(l *log.Logger, obj *models.Object) {
	if obj.DeleteOn != "" {
		t1, err := httpdate.Str2Time(obj.DeleteOn, time.UTC)
		if err != nil {
			l.Printf(
				"Passed Expires date for object %s has unappropriate style. Object can't be deleted by schedule.",
				obj.Id,
			)
		} else {
			l.Println(t1)
			now := time.Now()
			startTime := now.Add(t1.Sub(now))
			l.Printf("Deletion job will start at #%v", startTime)

			_, err = scheduler.Schedule(func(ctx context.Context) {
				//delete json object by schedule task
				models.DeleteObject(*obj.Id)
			}, chrono.WithTime(startTime))

			if err == nil {
				l.Printf("Task for id= %s has been scheduled", obj.Id)
			}
		}
	}
}
