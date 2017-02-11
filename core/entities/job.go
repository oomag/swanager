package entities

import (
	"time"

	"github.com/da4nik/swanager/config"
	"github.com/da4nik/swanager/core/db"
	"github.com/da4nik/swanager/lib"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// JobStateWorking - job state working
const JobStateWorking = "working"

// JobStateSuccess - job state success
const JobStateSuccess = "success"

// JobStateError - job state error
const JobStateError = "error"

const jobsCollectionName = "jobs"

// Job struct
type Job struct {
	ID         string      `json:"id" bson:"_id"`
	UserID     string      `json:"user_id" bson:"user_id"`
	State      string      `json:"state" bson:"state"`
	Result     interface{} `json:"result" bson:"result"`
	StartedAt  time.Time   `json:"started_at" bson:"started_at"`
	FinishedAt time.Time   `json:"finished_at" bson:"finished_at"`
}

// CreateJob creates new job with "working" state
func CreateJob(user *User) (job *Job, err error) {
	session := db.GetSession()
	defer session.Close()
	c := getJobsCollection(session)

	job = &Job{
		ID:        lib.GenerateUUID(),
		UserID:    user.ID,
		State:     JobStateWorking,
		StartedAt: time.Now(),
	}

	err = c.Insert(job)
	return
}

// GetJob return job ny id
func GetJob(id string) (job *Job, err error) {
	session := db.GetSession()
	defer session.Close()
	c := getJobsCollection(session)

	err = c.Find(bson.M{"id": id}).One(job)
	return
}

// SetState sets state and result for current job
func (j *Job) SetState(state string, result interface{}) error {
	session := db.GetSession()
	defer session.Close()
	c := getJobsCollection(session)

	err := c.UpdateId(j.ID, bson.M{"$set": bson.M{"result": result, "state": state, "finished_at": time.Now()}})
	if err != nil {
		return err
	}
	return nil
}

func getJobsCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(config.DatabaseName).C(jobsCollectionName)
}
