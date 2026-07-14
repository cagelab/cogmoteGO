package backup

type Status string

const (
	Running            Status = "running"
	Succeeded          Status = "succeeded"
	PartiallySucceeded Status = "partially_succeeded"
	Failed             Status = "failed"
)

type Phase string

const (
	PhaseScanning   Phase = "scanning"
	PhaseUploading  Phase = "uploading"
	PhaseVerifying  Phase = "verifying"
	PhasePublishing Phase = "publishing"
	PhaseCompleted  Phase = "completed"
)

type EntryStatus string

const (
	EntryPending   EntryStatus = "pending"
	EntryRunning   EntryStatus = "running"
	EntrySucceeded EntryStatus = "succeeded"
	EntryFailed    EntryStatus = "failed"
)

type Source struct {
	RootID  string   `json:"root_id"`
	Entries []string `json:"entries"`
}

type Destination struct {
	Type   string `json:"type"`
	RootID string `json:"root_id"`
	Path   string `json:"path"`
}

type Entry struct {
	Path   string      `json:"path"`
	Type   string      `json:"type"`
	Status EntryStatus `json:"status"`
	Error  string      `json:"error,omitempty"`
}

type Task struct {
	ID               string      `json:"id"`
	Status           Status      `json:"status"`
	Phase            Phase       `json:"phase"`
	Source           Source      `json:"source"`
	Destination      Destination `json:"destination"`
	Entries          []Entry     `json:"entries"`
	CreatedAt        string      `json:"created_at"`
	StartedAt        string      `json:"started_at,omitempty"`
	FinishedAt       string      `json:"finished_at,omitempty"`
	FilesTotal       int64       `json:"files_total"`
	FilesCompleted   int64       `json:"files_completed"`
	BytesTotal       int64       `json:"bytes_total"`
	BytesTransferred int64       `json:"bytes_transferred"`
	CurrentPath      string      `json:"current_path,omitempty"`
	Error            string      `json:"error,omitempty"`
}

func cloneTask(task Task) Task {
	task.Source.Entries = append([]string(nil), task.Source.Entries...)
	task.Entries = append([]Entry(nil), task.Entries...)
	return task
}
