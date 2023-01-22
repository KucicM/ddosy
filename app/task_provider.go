package ddosy

import (
	"log"
)

type TaskProvider struct {
	repo *TaskRepository
}

func NewTaskProvider(repostiroy *TaskRepository) *TaskProvider {
	return &TaskProvider{
		repo: repostiroy,
	}
}

func (p *TaskProvider) ScheduleTask(req ScheduleRequestWeb) (uint64, error) {
	return p.repo.Save(req)
}

func (p *TaskProvider) Next() *LoadTask {
	id, req, err := p.repo.GetNext()
	if err != nil {
		log.Printf("error getting new task from db %s\n", err)
		return nil
	}

	if id == 0 { // no new tasks
		return nil
	}

	if err := p.repo.UpdateStatus(id, Running); err != nil {
		log.Panicf("falied to update task status on kill event %s\n", err)
	}

	task := NewLoadTask(req)
	task.id = id
	return task
}

func (p *TaskProvider) Kill(id uint64) {
	if err := p.repo.UpdateStatus(id, Killed); err != nil {
		log.Panicf("falied to update task status on kill event %s\n", err)
	}
}

func (p *TaskProvider) Done(id uint64) {
	if err := p.repo.UpdateStatus(id, Done); err != nil {
		log.Printf("falied to update task status on done event %s\n", err)
	}
}
