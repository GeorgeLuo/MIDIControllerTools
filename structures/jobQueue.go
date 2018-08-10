package structures

/*
	threadsafe queue structure that contains the jobs to be proccessed.
	deviceJob is the underlying Node that contains an array of midi event with a source (port).
*/

type jobInterface interface {
	// use this value to check for rules mapping
    commandSource() int
    events() []portmidi.Event
}

type deviceJob struct {
	source int
	inputEvents []portmidi.Event
}

type queueJob struct {
	data jobInterface
	next *queueJob
}

type JobQueue struct {
	head  *queueJob
	tail  *queueJob
	count int
	lock  *sync.Mutex
}

func NewJobQueue() *JobQueue {
	q := &JobQueue{}
	q.lock = &sync.Mutex{}
	return q
}

func (q *JobQueue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.count
}

func (q *JobQueue) Push(event jobInterface) {
	q.lock.Lock()
	defer q.lock.Unlock()

	n := &queueJob{data: event}

	if q.tail == nil {
		q.tail = n
		q.head = n
	} else {
		q.tail.next = n
		q.tail = n
	}
	q.count++
}

func (q *JobQueue) Poll() jobInterface {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.head == nil {
		return nil
	}

	n := q.head
	q.head = n.next

	if q.head == nil {
		q.tail = nil
	}
	q.count--

	return n.data
}

func (q *JobQueue) Peek() jobInterface {
	q.lock.Lock()
	defer q.lock.Unlock()

	n := q.head
	if n == nil {
		return nil
	}

	return n.data
}

func (j deviceJob) commandSource() int {
    return j.source
}

func (j deviceJob) events() []portmidi.Event {
    return j.inputEvents
}
